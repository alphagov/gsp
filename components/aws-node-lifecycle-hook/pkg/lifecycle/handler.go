package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/awsclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sdrainer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ASGScaleInEvent            = "EC2 Instance-terminate Lifecycle Action"
	ASGSource                  = "aws.autoscaling"
	ASGTransition              = "autoscaling:EC2_INSTANCE_TERMINATING"
	EC2SpotInterruptionWarning = "EC2 Spot Instance Interruption Warning"
	EC2SpotActionStop          = "stop"
	EC2SpotActionTerminate     = "terminate"
	EC2SpotActionHibernate     = "hibernate"
)

var (
	ASGActionContinue = aws.String("CONTINUE")
)

type EC2SpotInterruptionEventDetail struct {
	InstanceID     string `json:"instance-id"`
	InstanceAction string `json:"instance-action,omitempty"`
}

type ASGLifecycleEventDetail struct {
	LifecycleActionToken string
	AutoScalingGroupName string
	LifecycleHookName    string
	EC2InstanceId        string
	LifecycleTransition  string
}

type Handler struct {
	AWSClient        awsclient.Client
	KubernetesClient k8sclient.Client
	k8sdrainer.Drainer
	HeartbeatInterval time.Duration
}

// HandleCloudwatchEvent routes
func (h *Handler) HandleEvent(ctx context.Context, event *events.CloudWatchEvent) error {
	switch event.DetailType {
	case ASGScaleInEvent:
		var asgEvent ASGLifecycleEventDetail
		if err := json.Unmarshal(event.Detail, &asgEvent); err != nil {
			return fmt.Errorf("cannot decode ASGLifecycleEventDetail: %s", err)
		}
		if asgEvent.EC2InstanceId == "" {
			return fmt.Errorf("EC2InstanceId missing from ASGLifecycleEventDetail")
		}
		return h.ScaleIn(ctx, asgEvent)
	case EC2SpotInterruptionWarning:
		var spotTerminationEvent EC2SpotInterruptionEventDetail
		if err := json.Unmarshal(event.Detail, &spotTerminationEvent); err != nil {
			return fmt.Errorf("cannot decode EC2SpotInterruptionEventDetail: %s", err)
		}
		if spotTerminationEvent.InstanceID == "" {
			return fmt.Errorf("InstanceId missing from EC2SpotInterruptionEventDetail")
		}
		switch spotTerminationEvent.InstanceAction {
		case EC2SpotActionStop, EC2SpotActionTerminate, EC2SpotActionHibernate:
			return h.drain(ctx, spotTerminationEvent.InstanceID)
		default:
			return fmt.Errorf("cannot handle EC2SpotInterruptionWarning action: %v", spotTerminationEvent.InstanceAction)
		}
	default:
		return fmt.Errorf("cannot handle event of type: %s", event.DetailType)
	}
}

func (h *Handler) ScaleIn(ctx context.Context, asgEvent ASGLifecycleEventDetail) error {
	// cordon and drain node
	if err := h.scaleIn(ctx, asgEvent); err != nil {
		return err
	}
	// respond with CONTINUE to the lifecycle hook
	if _, err := h.AWSClient.CompleteLifecycleActionWithContext(ctx, &autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  &asgEvent.AutoScalingGroupName,
		InstanceId:            &asgEvent.EC2InstanceId,
		LifecycleHookName:     &asgEvent.LifecycleHookName,
		LifecycleActionToken:  &asgEvent.LifecycleActionToken,
		LifecycleActionResult: ASGActionContinue,
	}); err != nil {
		return fmt.Errorf("failed to report completion of lifecycle hook: %s", err)
	}
	return nil
}

func (h *Handler) scaleIn(ctx context.Context, asgEvent ASGLifecycleEventDetail) error {
	// wrap context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// start heartbeat to keep instance in wait state
	go heartbeat(ctx, h.AWSClient, asgEvent, h.HeartbeatInterval)
	// attempt cordon/drain
	return h.drain(ctx, asgEvent.EC2InstanceId)
}

func (h *Handler) drain(ctx context.Context, instanceID string) error {
	// find the kubernetes node by instance id
	nodes, err := h.KubernetesClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch nodes: %s", err)
	}
	var node *corev1.Node
	for _, n := range nodes.Items {
		if strings.HasSuffix(n.Spec.ProviderID, instanceID) {
			node = &n
			break
		}
	}
	if node == nil {
		log.Printf("no kubernetes node found (or node already drained/removed) for instance id %s", instanceID)
		return nil // nothing to do
	}
	// cordon the node so no pods can be bound
	if err := h.Cordon(h.KubernetesClient, node); err != nil {
		return fmt.Errorf("failed to cordon kubernetes node %s: %s", node.Spec.ProviderID, err)
	}
	// evict pods from the node
	if err := h.Drain(h.KubernetesClient, node); err != nil {
		return fmt.Errorf("failed to drain kubernetes node %s: %s", node.Spec.ProviderID, err)
	}
	// done
	return nil
}
