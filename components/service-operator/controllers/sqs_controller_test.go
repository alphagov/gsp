package controllers_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	. "github.com/alphagov/gsp/components/service-operator/controllers"
	"github.com/alphagov/gsp/components/service-operator/internal"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	internalawsmocks "github.com/alphagov/gsp/components/service-operator/internal/aws/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("SQSController", func() {
	var (
		sqs          queue.SQS
		queueName    string
		secret       core.Secret
		secretName   string
		principal    access.Principal
		request      reconcile.Request
		reconciler   SQSReconciler
		cfReconciler *internalawsmocks.MockCloudFormationReconciler
	)

	BeforeEach(func() {
		queueName, _ = internal.RandomString(8, internal.CharactersLower)
		secretName, _ = internal.RandomString(8, internal.CharactersLower)
		request = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "test",
				Name:      queueName,
			},
		}
		sqs = queue.SQS{
			TypeMeta: metav1.TypeMeta{
				APIVersion: queue.GroupVersion.Group,
				Kind:       "SQS",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      queueName,
				Namespace: "test",
				Labels: map[string]string{
					access.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: queue.SQSSpec{
				Secret: secretName,
			},
		}
		secret = core.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: "test",
			},
			Data: map[string][]byte{
				"QueueURL": []byte("test-queue-url"),
			},
		}
		principal = access.Principal{
			TypeMeta: metav1.TypeMeta{
				APIVersion: access.GroupVersion.Group,
				Kind:       "Principal",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      "test-role",
				Labels: map[string]string{
					access.AccessGroupLabel: "test.access.group",
				},
			},
		}
		k8sClient.Create(context.TODO(), &principal)
		k8sClient.Create(context.TODO(), &sqs)
		cfReconciler = internalawsmocks.NewMockCloudFormationReconciler(mockCtrl)
		reconciler = SQSReconciler{
			Client:                   k8sClient,
			Log:                      log,
			CloudFormationReconciler: cfReconciler,
		}
	})

	AfterEach(func() {
		k8sClient.Delete(context.TODO(), &sqs)
		k8sClient.Delete(context.TODO(), &secret)
		k8sClient.Delete(context.TODO(), &principal)
	})

	Context("When using an undefined provisioner", func() {
		BeforeEach(func() {
			os.Setenv("CLOUD_PROVIDER", "undefined")
		})

		It("Should backoff for 15 minutes", func() {

			result, err := reconciler.Reconcile(request)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported cloud provider: undefined"))
			Expect(result.Requeue).To(BeTrue())
			Expect(result.RequeueAfter).To(Equal(time.Minute * 15))
		})
	})

	Context("When using the AWS provisioner", func() {
		BeforeEach(func() {
			os.Setenv("CLOUD_PROVIDER", "aws")
		})

		Context("When creating a new resource", func() {
			It("Should update the kubernetes resource", func() {
				stackData := internalaws.StackData{
					ID:     "test-id",
					Status: "created",
					Reason: "because-of-create",
				}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), false).
					Return(internal.Create, stackData, nil).
					Times(1)

				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedSQS queue.SQS
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      queueName,
				}, &updatedSQS)

				checkSQSStatusUpdates(stackData, updatedSQS)
				Expect(updatedSQS.Finalizers).To(ContainElement(SQSFinalizerName))
				Expect(updatedSQS.ObjectMeta.DeletionTimestamp).To(BeNil())
			})

			It("Should create a secret with the queue name", func() {
				url := "https://sqs.eu-west-2.amazonaws.com/1234567890/test-queue"
				stackData := internalaws.StackData{
					Outputs: []*cloudformation.Output{
						&cloudformation.Output{
							OutputKey:   aws.String("QueueURL"),
							OutputValue: aws.String(url),
						},
					},
				}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), false).
					Return(internal.Create, stackData, nil).
					Times(1)

				_, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())

				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      secretName,
				}, &secret)
				Expect(string(secret.Data["QueueURL"])).To(Equal(url))
			})

			It("Should fail if there is no matching IAM role", func() {
				stackData := internalaws.StackData{
					ID:     "test-id",
					Status: "created",
					Reason: "because-of-create",
				}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), false).
					Return(internal.Create, stackData, nil).
					Times(1)
				k8sClient.Delete(context.TODO(), &principal)

				result, _ := reconciler.Reconcile(request)
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute * 2))
			})
		})

		Context("When updating a resource", func() {
			It("Should update the kubernetes resource", func() {
				stackData := internalaws.StackData{
					ID:     "test-id",
					Status: "updated",
					Reason: "because-of-update",
				}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), false).
					Return(internal.Update, stackData, nil).
					Times(1)

				sqs.ObjectMeta.Finalizers = append(sqs.Finalizers, SQSFinalizerName)
				k8sClient.Update(context.TODO(), &sqs)
				k8sClient.Create(context.TODO(), &secret)

				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedSQS queue.SQS
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      queueName,
				}, &updatedSQS)
				checkSQSStatusUpdates(stackData, updatedSQS)
				Expect(updatedSQS.ObjectMeta.Finalizers).To(ContainElement(SQSFinalizerName))
				Expect(updatedSQS.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})

		Context("When deleting a resource", func() {
			BeforeEach(func() {
				sqs.ObjectMeta.Finalizers = append(sqs.Finalizers, SQSFinalizerName)
				k8sClient.Update(context.TODO(), &sqs)
				k8sClient.Create(context.TODO(), &secret)
			})

			It("Should delete the kubernetes resource", func() {
				stackData := internalaws.StackData{}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), true).
					Return(internal.Delete, stackData, nil).
					Times(1)

				k8sClient.Delete(context.TODO(), &sqs)
				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedSQS queue.SQS
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      queueName,
				}, &updatedSQS)
				Expect(updatedSQS.ObjectMeta.Finalizers).ToNot(ContainElement(SQSFinalizerName))
				Expect(updatedSQS.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})
	})
})

func checkSQSStatusUpdates(stackData internalaws.StackData, sqs queue.SQS) {
	Expect(sqs.Status.ID).To(Equal(stackData.ID))
	Expect(sqs.Status.Status).To(Equal(stackData.Status))
	Expect(sqs.Status.Reason).To(Equal(stackData.Reason))
}
