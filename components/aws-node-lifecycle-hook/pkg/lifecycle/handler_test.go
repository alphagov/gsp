package lifecycle_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/awsclient/fakeawsclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sclient/fakek8sclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sdrainer/fakek8sdrainer"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/lifecycle"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var origStderr = os.Stderr

var _ = Describe("LifecycleHandler", func() {

	var (
		ctx             = context.Background()
		handler         lifecycle.Handler
		fakeAWSClient   *fakeawsclient.FakeClient
		fakeK8sClient   *fakek8sclient.FakeClient
		fakeK8sDrainer  *fakek8sdrainer.FakeDrainer
		cloudwatchEvent *events.CloudWatchEvent
		corev1client    *fakek8sclient.FakeCoreV1Interface
	)

	var (
		node1 = v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-one",
			},
			Spec: v1.NodeSpec{
				ProviderID: "aws://eu-west-2a/i-111111111",
			},
		}
		node2 = v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-two",
			},
			Spec: v1.NodeSpec{
				ProviderID: "aws://eu-west-2b/i-222222222",
			},
		}
		node3 = v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-three",
			},
			Spec: v1.NodeSpec{
				ProviderID: "aws://eu-west-2c/i-333333333",
			},
		}
	)

	BeforeEach(func() {
		fakeAWSClient = &fakeawsclient.FakeClient{}
		fakeAWSClient.RecordLifecycleActionHeartbeatWithContextReturns(&autoscaling.RecordLifecycleActionHeartbeatOutput{}, nil)
		fakeAWSClient.CompleteLifecycleActionWithContextReturns(&autoscaling.CompleteLifecycleActionOutput{}, nil)

		corev1client = &fakek8sclient.FakeCoreV1Interface{}
		nodeLister := &fakek8sclient.FakeNodeInterface{}
		nodeLister.ListReturns(&v1.NodeList{
			Items: []v1.Node{node1, node2, node3},
		}, nil)
		corev1client.NodesReturns(nodeLister)
		fakeK8sClient = &fakek8sclient.FakeClient{}
		fakeK8sClient.CoreV1Returns(corev1client)

		fakeK8sDrainer = &fakek8sdrainer.FakeDrainer{}
		fakeK8sDrainer.CordonReturns(nil)
		fakeK8sDrainer.DrainReturns(nil)

		handler = lifecycle.Handler{
			AWSClient:        fakeAWSClient,
			KubernetesClient: fakeK8sClient,
			Drainer:          fakeK8sDrainer,
		}
	})

	Context("receives an event of unknown type", func() {
		JustBeforeEach(func() {
			cloudwatchEvent = &events.CloudWatchEvent{
				Version:    "0",
				ID:         "1",
				Source:     "aws.autoscaling",
				AccountID:  "123456789",
				Time:       time.Now(),
				Region:     "eu-west-2",
				Resources:  []string{},
				DetailType: "???",
			}
		})

		It("should return an error about unknown event type", func() {
			err := handler.HandleEvent(ctx, cloudwatchEvent)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot handle event of type"))
		})
	})

	Context("receives an event with malformed detail field", func() {
		JustBeforeEach(func() {
			cloudwatchEvent = &events.CloudWatchEvent{
				Version:    "xxxx",
				ID:         "1",
				Source:     "aws.autoscaling",
				AccountID:  "123456789",
				Time:       time.Now(),
				Region:     "eu-west-2",
				Resources:  []string{},
				DetailType: lifecycle.ASGScaleInEvent,
				Detail:     json.RawMessage(`{"LifecycleActionToken": BADJSON}`),
			}
		})

		It("should return an error about not being able to unmarshal payload", func() {
			err := handler.HandleEvent(ctx, cloudwatchEvent)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot decode event detail"))
		})
	})

	Context("receives an event with no EC2 instance ID", func() {
		JustBeforeEach(func() {
			cloudwatchEvent = &events.CloudWatchEvent{
				Version:    "xxxx",
				ID:         "1",
				Source:     "aws.autoscaling",
				AccountID:  "123456789",
				Time:       time.Now(),
				Region:     "eu-west-2",
				Resources:  []string{},
				DetailType: lifecycle.ASGScaleInEvent,
				Detail:     json.RawMessage(`{}`),
			}
		})

		It("should return an error about missing instance ID", func() {
			err := handler.HandleEvent(ctx, cloudwatchEvent)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("EC2InstanceId missing"))
		})
	})

	Context("receives a valid ASG instance terminating event", func() {

		var handlerErr error

		JustBeforeEach(func() {
			asgEvent := lifecycle.ASGLifecycleEventDetail{
				LifecycleActionToken: "token-0",
				AutoScalingGroupName: "asg-0",
				LifecycleHookName:    "hook-0",
				EC2InstanceId:        "i-222222222",
				LifecycleTransition:  "autoscaling:EC2_INSTANCE_TERMINATING",
			}
			asgEventBytes, err := json.Marshal(asgEvent)
			Expect(err).ToNot(HaveOccurred())
			cloudwatchEvent = &events.CloudWatchEvent{
				Version:    "xxxx",
				ID:         "1",
				Source:     "aws.autoscaling",
				AccountID:  "123456789",
				Time:       time.Now(),
				Region:     "eu-west-2",
				Resources:  []string{},
				DetailType: lifecycle.ASGScaleInEvent,
				Detail:     json.RawMessage(asgEventBytes),
			}
			handlerErr = handler.HandleEvent(ctx, cloudwatchEvent)
		})

		Context("and everything goes well", func() {

			It("should not error", func() {
				Expect(handlerErr).ToNot(HaveOccurred())
			})

			It("should cordon the node", func() {
				Expect(fakeK8sDrainer.CordonCallCount()).To(Equal(1))
				c, node := fakeK8sDrainer.CordonArgsForCall(0)
				Expect(c).To(Equal(fakeK8sClient))
				Expect(node.ObjectMeta.Name).To(Equal("node-two"))
			})

			It("should drain the node", func() {
				Expect(fakeK8sDrainer.DrainCallCount()).To(Equal(1))
				c, node := fakeK8sDrainer.DrainArgsForCall(0)
				Expect(c).To(Equal(fakeK8sClient))
				Expect(node.ObjectMeta.Name).To(Equal("node-two"))
			})

			It("should call the lifecycle action callback exactly once", func() {
				Expect(fakeAWSClient.CompleteLifecycleActionWithContextCallCount()).To(Equal(1))
			})

			It("should call the lifecycle action callback with the expected args", func() {
				callingContext, input, _ := fakeAWSClient.CompleteLifecycleActionWithContextArgsForCall(0)
				Expect(callingContext).ToNot(BeNil())
				Expect(input.AutoScalingGroupName).To(Equal(aws.String("asg-0")))
				Expect(input.InstanceId).To(Equal(aws.String("i-222222222")))
				Expect(input.LifecycleHookName).To(Equal(aws.String("hook-0")))
				Expect(input.LifecycleActionToken).To(Equal(aws.String("token-0")))
				Expect(input.LifecycleActionResult).To(Equal(aws.String("CONTINUE")))
			})
		})

		Context("but cordoning fails", func() {
			BeforeEach(func() {
				fakeK8sDrainer.CordonReturns(fmt.Errorf("CORDON_ERR"))
			})

			It("should return an error about cordon failure", func() {
				Expect(handlerErr).To(HaveOccurred())
				Expect(handlerErr.Error()).To(ContainSubstring("failed to cordon node"))
			})
		})

		Context("but draining fails", func() {
			BeforeEach(func() {
				fakeK8sDrainer.DrainReturns(fmt.Errorf("DRAIN_ERR"))
			})

			It("should return an error about drain failure", func() {
				Expect(handlerErr).To(HaveOccurred())
				Expect(handlerErr.Error()).To(ContainSubstring("failed to drain node"))
			})
		})

		Context("but it fails to call completion", func() {
			BeforeEach(func() {
				fakeAWSClient.CompleteLifecycleActionWithContextReturns(nil, fmt.Errorf("AWS_ERR"))
			})

			It("should return an error about missing node", func() {
				Expect(handlerErr).To(HaveOccurred())
				Expect(handlerErr.Error()).To(ContainSubstring("failed to report completion"))
			})
		})

		Context("but the node cannot be found", func() {
			BeforeEach(func() {
				nodeLister := &fakek8sclient.FakeNodeInterface{}
				nodeLister.ListReturns(&v1.NodeList{
					Items: []v1.Node{node1},
				}, nil)
				corev1client.NodesReturns(nodeLister)
			})

			It("should return an error about missing node", func() {
				Expect(handlerErr).To(HaveOccurred())
				Expect(handlerErr.Error()).To(ContainSubstring("failed to find node"))
			})
		})

		Context("but there is an error from the kubernetes api", func() {
			BeforeEach(func() {
				nodeLister := &fakek8sclient.FakeNodeInterface{}
				nodeLister.ListReturns(nil, fmt.Errorf("FAKE_K8S_ERROR"))
				corev1client.NodesReturns(nodeLister)
			})

			It("should return an error about fetching nodes", func() {
				Expect(handlerErr).To(HaveOccurred())
				Expect(handlerErr.Error()).To(ContainSubstring("failed to fetch nodes"))
			})
		})

		Context("and it takes a while to drain", func() {

			BeforeEach(func() {
				handler.HeartbeatInterval = time.Millisecond * 10
				fakeK8sDrainer.DrainStub = func(c k8sclient.Client, Node *v1.Node) error {
					time.Sleep(time.Second * 1)
					return nil
				}
			})

			It("should have periodically sent lifecycle action heartbeat", func() {
				Expect(fakeAWSClient.RecordLifecycleActionHeartbeatWithContextCallCount()).To(BeNumerically(">", 80))
			})

			It("should have called the heartbeat action with expected args", func() {
				callingContext, input, _ := fakeAWSClient.RecordLifecycleActionHeartbeatWithContextArgsForCall(0)
				Expect(callingContext).ToNot(BeNil())
				Expect(input.AutoScalingGroupName).To(Equal(aws.String("asg-0")))
				Expect(input.InstanceId).To(Equal(aws.String("i-222222222")))
				Expect(input.LifecycleHookName).To(Equal(aws.String("hook-0")))
				Expect(input.LifecycleActionToken).To(Equal(aws.String("token-0")))
			})

			It("should stop heartbeating after handler completes", func() {
				heartbeatCount := fakeAWSClient.RecordLifecycleActionHeartbeatWithContextCallCount()
				Consistently(fakeAWSClient.RecordLifecycleActionHeartbeatWithContextCallCount).Should(Equal(heartbeatCount))
			})

		})
	})
})
