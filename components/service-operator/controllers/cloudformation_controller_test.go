package controllers_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/gsp/components/service-operator/apis"
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	. "github.com/alphagov/gsp/components/service-operator/controllers"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("SQSController", func() {
	var (
		queueName = "test-queue"
		// queueURL            = "https://sqs.eu-west-2.amazonaws.com/1234567890/test-queue"
		secretName          = "test-secret"
		principalName       = "test-role"
		namespace           = "test"
		queueNamespacedName = types.NamespacedName{
			Namespace: namespace,
			Name:      queueName,
		}
	)

	Context("aws", func() {

		BeforeEach(func() {
		})

		AfterEach(func() {
			Consistently(sqsReconciler.Err, time.Second*2).ShouldNot(HaveOccurred())
			// Consistently(principalReconciler.Err, time.Second*2).ShouldNot(HaveOccurred())
		})

		It("Should create and destroy an SQS queue", func() {

			By("creating a prequisite Principal resource with kubernetes api", func() {
				principal := access.Principal{
					TypeMeta: metav1.TypeMeta{
						APIVersion: access.GroupVersion.Group,
						Kind:       "Principal",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      principalName,
						Labels: map[string]string{
							access.AccessGroupLabel: "test.access.group",
						},
					},
				}
				Expect(k8sClient.Create(ctx, &principal)).To(Succeed())
			})

			By("creating an SQS resource with kubernetes api", func() {
				sqs := queue.SQS{
					TypeMeta: metav1.TypeMeta{
						APIVersion: queue.GroupVersion.Group,
						Kind:       "SQS",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      queueName,
						Namespace: namespace,
						Labels: map[string]string{
							access.AccessGroupLabel: "test.access.group",
						},
					},
					Spec: queue.SQSSpec{
						Secret: secretName,
					},
				}
				Expect(k8sClient.Create(ctx, &sqs)).To(Succeed())
			})

			By("ensuring no reconcile errors during create", func() {
				Consistently(sqsReconciler.Err, time.Second*2).ShouldNot(HaveOccurred())
			})

			By("displaying a READY resource status after initial creation", func() {
				Eventually(func() string {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.Status.Service.State
				}).Should(Equal(apis.ReadyState))
			})

			By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
				Eventually(func() string {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.Status.AWS.Status
				}, time.Second*5).Should(Equal(awscloudformation.StackStatusCreateComplete))
			})

			By("displaying an AWS stack id in resource status", func() {
				Eventually(func() string {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.Status.AWS.ID
				}).ShouldNot(BeEmpty())
			})

			By("displaying a stack name prefixed with cluster name in resource status", func() {
				Eventually(func() string {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.Status.AWS.Name
				}).Should(Equal("xxx-sqs-test-test-queue"))
			})

			By("displaying an AWS state reason in resource status", func() {
				Eventually(func() string {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.Status.AWS.Reason
				}).ShouldNot(BeEmpty())
			})

			By("ensuring a finalizaer is present on resource to prevent deletion", func() {
				Eventually(func() []string {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.Finalizers
				}).Should(ContainElement(SQSFinalizerName))
			})

			By("ensuring no DeletionTimestamp exists", func() {
				Eventually(func() bool {
					var sqs queue.SQS
					_ = k8sClient.Get(ctx, queueNamespacedName, &sqs)
					return sqs.ObjectMeta.DeletionTimestamp == nil
				}).Should(BeTrue())
			})

			By("creating a secret with a secret with queue connection details", func() {
				Eventually(func() map[string][]byte {
					var secret core.Secret
					_ = k8sClient.Get(ctx, types.NamespacedName{
						Namespace: namespace,
						Name:      secretName,
					}, &secret)
					return secret.Data
				}).Should(HaveKey("QueueURL"))
			})

			// By("updating kubernetes resource should update the stack", func() {
			// 	sqsCloudformationReconciler.
			// 		EXPECT().
			// 		Reconcile(ctx, gomock.Any(), gomock.Eq(sqsRequest), gomock.Any(), false).
			// 		Return(internal.Create, sqsStackData, nil).
			// 		Times(1)
			// 	var sqs queue.SQS
			// 	err := k8sClient.Get(ctx, queueNamespacedName, &sqs)
			// 	Expect(err).ToNot(HaveOccurred())
			// 	sqs.Spec.AWS.MaximumMessageSize = 2 // arbitary change
			// 	Expect(k8sClient.Update(ctx, &sqs)).To(Succeed())

			// 	// TODO: check the reequeue stuff
			// 	// result, err := reconciler.Reconcile(request)
			// 	// Expect(err).ToNot(HaveOccurred())
			// 	// Expect(result.Requeue).To(BeTrue())
			// 	// Expect(result.RequeueAfter).To(Equal(time.Minute))

			// 	var updatedSQS queue.SQS
			// 	Expect(k8sClient.Get(ctx, types.NamespacedName{
			// 		Namespace: "test",
			// 		Name:      queueName,
			// 	}, &updatedSQS)).To(Succeed())
			// 	checkSQSStatusUpdates(stackData, updatedSQS)
			// 	Expect(updatedSQS.ObjectMeta.Finalizers).To(ContainElement(SQSFinalizerName))
			// 	Expect(updatedSQS.ObjectMeta.DeletionTimestamp).To(BeNil())
			// })

			// BeforeEach(func() {
			// 	sqs.ObjectMeta.Finalizers = append(sqs.Finalizers, SQSFinalizerName)
			// 	Expect(k8sClient.Update(ctx, &sqs)).To(Succeed())
			// 	Expect(k8sClient.Create(ctx, &secret)).To(Succeed())
			// })

			By("deleteing SQS resource with kubernetes api", func() {
				var sqs queue.SQS
				err := k8sClient.Get(ctx, queueNamespacedName, &sqs)
				Expect(err).ToNot(HaveOccurred())
				Expect(k8sClient.Delete(ctx, &sqs)).To(Succeed())
			})

			By("ensuring no reconcile errors during delete", func() {
				Consistently(sqsReconciler.Err, time.Second*2).ShouldNot(HaveOccurred())
			})

			By("ensuring the SQS resources have been removed", func() {
				var sqsList queue.SQSList
				Eventually(func() int {
					err := k8sClient.List(ctx, &sqsList)
					Expect(err).ToNot(HaveOccurred())
					return len(sqsList.Items)
				}, time.Second*10).Should(Equal(0))
			})

			// By("deleting the queue resource should delete the Secret", func() {
			// 	var secretList core.SecretList
			// 	Eventually(func() int {
			// 		err := k8sClient.List(ctx, &secretList)
			// 		Expect(err).ToNot(HaveOccurred())
			// 		return len(secretList.Items)
			// 	}, time.Second*10).Should(Equal(0))
			// })
		})

		// It("Should fail if there is no matching IAM role", func() {
		// 	stackData := internalaws.StackData{
		// 		ID:     "test-stack",
		// 		Status: "created",
		// 		Reason: "because-of-create",
		// 	}
		// 	cfReconciler.
		// 		EXPECT().
		// 		Reconcile(ctx, gomock.Any(), request, gomock.Any(), false).
		// 		Return(internal.Create, stackData, nil).
		// 		Times(1)
		// 	Expect(k8sClient.Delete(ctx, &principal)).To(Succeed())

		// 	result, _ := reconciler.Reconcile(request)
		// 	Expect(result.Requeue).To(BeTrue())
		// 	Expect(result.RequeueAfter).To(Equal(time.Minute * 2))
		// })
	})
})
