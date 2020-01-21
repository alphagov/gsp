package v1beta1_test

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
)

var _ = Describe("SQS", func() {

	var sqs v1beta1.SQS
	var tags []cloudformation.Tag

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		sqs = v1beta1.SQS{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.SQSSpec{},
		}
		tags = []cloudformation.Tag{
			{Key: "Cluster", Value: env.ClusterName()},
			{Key: "Service", Value: "sqs"},
			{Key: "Name", Value: "example"},
			{Key: "Namespace", Value: "default"},
			{Key: "Environment", Value: "default"},
			{Key: "QueueType", Value: "Main"},
		}
	})

	It("should default secret name to object name", func() {
		Expect(sqs.GetSecretName()).To(Equal("example"))
	})

	It("should use secret name from spec.Secret if set ", func() {
		sqs.Spec.Secret = "my-target-secret"
		Expect(sqs.GetSecretName()).To(Equal("my-target-secret"))
	})

	It("implements runtime.Object", func() {
		o2 := sqs.DeepCopyObject()
		Expect(o2).ToNot(BeZero())
	})

	Context("cloudformation", func() {

		It("should generate a unique stack name prefixed with cluster name", func() {
			Expect(sqs.GetStackName()).To(HavePrefix("xxx-sqs-default-example"))
		})

		It("should require an IAM role input", func() {
			t, err := sqs.GetStackTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Parameters).To(HaveKey("IAMRoleName"))
		})

		It("should have outputs for connection details", func() {
			t, err := sqs.GetStackTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Outputs).To(And(
				HaveKey("QueueURL"),
				HaveKey("DLQueueURL"),
				HaveKey("IAMRoleName"),
			))
		})

		It("should map role name to role parameter", func() {
			params, err := sqs.GetStackRoleParameters("fake-role")
			Expect(err).ToNot(HaveOccurred())
			Expect(params).To(ContainElement(&cloudformation.Parameter{
				ParameterKey:   aws.String("IAMRoleName"),
				ParameterValue: aws.String("fake-role"),
			}))
		})

		Context("queue resource", func() {

			var queue *cloudformation.AWSSQSQueue

			JustBeforeEach(func() {
				t, err := sqs.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSSQSQueue{})))
				var ok bool
				queue, ok = t.Resources[v1beta1.SQSResourceName].(*cloudformation.AWSSQSQueue)
				Expect(ok).To(BeTrue())
			})

			It("should have a queue name prefixed with cluster and namespace name", func() {
				Expect(queue.QueueName).To(Equal("xxx-default-example"))
			})

			It("should have suitable tags set", func() {
				Expect(queue.Tags).To(Equal(tags))
			})

			Context("defaults", func() {
				It("should have sensible default values", func() {
					Expect(queue.ContentBasedDeduplication).To(BeFalse())
					Expect(queue.DelaySeconds).To(BeZero())
					Expect(queue.FifoQueue).To(BeFalse())
					Expect(queue.MaximumMessageSize).To(BeZero())
					Expect(queue.MessageRetentionPeriod).To(BeZero())
					Expect(queue.ReceiveMessageWaitTimeSeconds).To(BeZero())
					Expect(queue.RedrivePolicy).To(BeEmpty())
					Expect(queue.VisibilityTimeout).To(BeZero())
				})
			})

			Context("when spec.aws.contentBasedDeduplication is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.ContentBasedDeduplication = true
				})
				It("should set queue ContentBasedDeduplication from spec", func() {
					Expect(queue.ContentBasedDeduplication).To(BeTrue())
				})
			})

			Context("when spec.aws.delaySeconds is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.DelaySeconds = 5
				})
				It("should set queue DelaySeconds from spec", func() {
					Expect(queue.DelaySeconds).To(Equal(5))
				})
			})

			Context("when spec.aws.fifoQueue is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.FifoQueue = true
				})
				It("should set queue DelaySeconds from spec", func() {
					Expect(queue.FifoQueue).To(BeTrue())
				})
			})

			Context("when spec.aws.maximumMessageSize is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.MaximumMessageSize = 101
				})
				It("should set queue MaximumMessageSize from spec", func() {
					Expect(queue.MaximumMessageSize).To(Equal(101))
				})
			})

			Context("when spec.aws.messageRetentionPeriod is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.MessageRetentionPeriod = 3600
				})
				It("should set queue MessageRetentionPeriod from spec", func() {
					Expect(queue.MessageRetentionPeriod).To(Equal(3600))
				})
			})

			Context("when spec.aws.receiveMessageWaitTimeSeconds is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.ReceiveMessageWaitTimeSeconds = 60
				})
				It("should set queue ReceiveMessageWaitTimeSeconds from spec", func() {
					Expect(queue.ReceiveMessageWaitTimeSeconds).To(Equal(60))
				})
			})

			Context("when spec.aws.redriveMaxReceiveCount is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.RedriveMaxReceiveCount = 10
				})
				It("should set queue RedrivePolicy from spec", func() {
					policy, ok := queue.RedrivePolicy.(map[string]interface{})
					Expect(ok).To(BeTrue())
					Expect(policy).To(And(
						HaveKeyWithValue("maxReceiveCount", 10),
						HaveKey("deadLetterTargetArn"),
					))
				})
			})

			Context("when spec.aws.visibilityTimeout is set", func() {
				BeforeEach(func() {
					sqs.Spec.AWS.VisibilityTimeout = 120
				})
				It("should set queue VisibilityTimeout from spec", func() {
					Expect(queue.VisibilityTimeout).To(Equal(120))
				})
			})

		})

		Context("policy resource", func() {
			var policy *cloudformation.AWSIAMPolicy
			var doc cloudformation.PolicyDocument

			JustBeforeEach(func() {
				t, err := sqs.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources[v1beta1.SQSResourceIAMPolicy]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMPolicy{}))
				policy = t.Resources[v1beta1.SQSResourceIAMPolicy].(*cloudformation.AWSIAMPolicy)
				Expect(policy.PolicyDocument).To(BeAssignableToTypeOf(cloudformation.PolicyDocument{}))
				doc = policy.PolicyDocument.(cloudformation.PolicyDocument)
			})

			It("should have a policy name", func() {
				Expect(policy.PolicyName).ToNot(BeEmpty())
			})

			It("should asign policy to the given role name", func() {
				Expect(policy.Roles).To(ContainElement(cloudformation.Ref("IAMRoleName")))
			})

			It("should have a policy document that allows access to resource", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Action).To(ConsistOf(
					"sqs:ChangeMessageVisibility",
					"sqs:DeleteMessage",
					"sqs:GetQueueAttributes",
					"sqs:GetQueueUrl",
					"sqs:ListDeadLetterSourceQueues",
					"sqs:ListQueueTags",
					"sqs:PurgeQueue",
					"sqs:ReceiveMessage",
					"sqs:SendMessage",
				))
				Expect(statement.Resource).To(ContainElement(
					cloudformation.GetAtt(v1beta1.SQSResourceName, "Arn"),
				))
			})

		})

	})

})
