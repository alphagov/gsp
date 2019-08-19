package controllers_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	. "github.com/alphagov/gsp/components/service-operator/controllers"
	"github.com/alphagov/gsp/components/service-operator/internal"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	internalawsmocks "github.com/alphagov/gsp/components/service-operator/internal/aws/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("PrincipalController", func() {
	const roleName = "test-role"
	var (
		principal    access.Principal
		request      reconcile.Request
		reconciler   PrincipalReconciler
		cfReconciler *internalawsmocks.MockCloudFormationReconciler
	)

	BeforeEach(func() {
		request = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "test",
				Name:      roleName,
			},
		}
		principal = access.Principal{
			TypeMeta: metav1.TypeMeta{
				APIVersion: access.GroupVersion.Group,
				Kind:       "Principal",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      roleName,
			},
		}
		k8sClient.Create(context.TODO(), &principal)
		cfReconciler = internalawsmocks.NewMockCloudFormationReconciler(mockCtrl)
		reconciler = PrincipalReconciler{
			Client:                   k8sClient,
			Log:                      log,
			CloudFormationReconciler: cfReconciler,
			ClusterName:              "test-cluster",
			RolePrincipal:            "arn:aws:iam::123456789012:role/kiam",
			PermissionsBoundary:      "arn:aws:iam::123456789012:policy/permissions-boundary",
		}
	})

	AfterEach(func() {
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

				var updatedPrincipal access.Principal
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      roleName,
				}, &updatedPrincipal)

				checkPrincipalStatusUpdates(stackData, updatedPrincipal)
				Expect(updatedPrincipal.Finalizers).To(ContainElement(PrincipalFinalizerName))
				Expect(updatedPrincipal.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})

		Context("When updating a resource", func() {
			It("Should update the kubernetes resource", func() {
				stackData := internalaws.StackData{
					ID:     "test-id",
					Status: "updated",
					Reason: "because-of-update",
					Outputs: []*cloudformation.Output{
						&cloudformation.Output{
							OutputKey:   aws.String(internalaws.IAMRoleARN),
							OutputValue: aws.String("arn:aws:iam::123456789012:role/test-cluster-test-test-role"),
						},
					},
				}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), false).
					Return(internal.Update, stackData, nil).
					Times(1)

				principal.ObjectMeta.Finalizers = append(principal.Finalizers, PrincipalFinalizerName)
				k8sClient.Update(context.TODO(), &principal)

				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedPrincipal access.Principal
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      roleName,
				}, &updatedPrincipal)
				checkPrincipalStatusUpdates(stackData, updatedPrincipal)
				Expect(updatedPrincipal.Status.ARN).To(Equal("arn:aws:iam::123456789012:role/test-cluster-test-test-role"))
				Expect(updatedPrincipal.ObjectMeta.Finalizers).To(ContainElement(PrincipalFinalizerName))
				Expect(updatedPrincipal.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})

		Context("When deleting a resource", func() {
			It("Should delete the kubernetes resource", func() {
				stackData := internalaws.StackData{}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), true).
					Return(internal.Delete, stackData, nil).
					Times(1)

				principal.ObjectMeta.Finalizers = append(principal.Finalizers, PrincipalFinalizerName)
				k8sClient.Update(context.TODO(), &principal)

				k8sClient.Delete(context.TODO(), &principal)
				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedPrincipal access.Principal
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      roleName,
				}, &updatedPrincipal)
				Expect(updatedPrincipal.ObjectMeta.Finalizers).ToNot(ContainElement(PrincipalFinalizerName))
				Expect(updatedPrincipal.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})
	})
})

func checkPrincipalStatusUpdates(stackData internalaws.StackData, principal access.Principal) {
	Expect(principal.Status.ID).To(Equal(stackData.ID))
	Expect(principal.Status.Status).To(Equal(stackData.Status))
	Expect(principal.Status.Reason).To(Equal(stackData.Reason))
}
