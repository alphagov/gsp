package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
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
	ASGScaleInEvent = "EC2 Instance-terminate Lifecycle Action"
	ASGSource       = "aws.autoscaling"
	ASGTransition   = "autoscaling:EC2_INSTANCE_TERMINATING"
)

var (
	ASGActionContinue = aws.String("CONTINUE")
)

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
			return fmt.Errorf("cannot decode event detail")
		}
		if asgEvent.EC2InstanceId == "" {
			return fmt.Errorf("EC2InstanceId missing from event detail")
		}
		return h.ScaleIn(ctx, asgEvent)
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
	// find the kubernetes node by instance id
	nodes, err := h.KubernetesClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch nodes: %s", err)
	}
	var node *corev1.Node
	for _, n := range nodes.Items {
		if strings.HasSuffix(n.Spec.ProviderID, asgEvent.EC2InstanceId) {
			node = &n
			break
		}
	}
	if node == nil {
		return fmt.Errorf("failed to find node with id: %s", asgEvent.EC2InstanceId)
	}
	// cordon the node so no pods can be bound
	if err := h.Cordon(h.KubernetesClient, node); err != nil {
		return fmt.Errorf("failed to cordon node: %s", err)
	}
	// evict pods from the node
	if err := h.Drain(h.KubernetesClient, node); err != nil {
		return fmt.Errorf("failed to drain node: %s", err)
	}
	// done
	return nil
}
