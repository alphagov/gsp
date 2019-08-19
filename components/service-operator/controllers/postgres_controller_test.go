package controllers_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	database "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	. "github.com/alphagov/gsp/components/service-operator/controllers"
	"github.com/alphagov/gsp/components/service-operator/internal"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	internalawsmocks "github.com/alphagov/gsp/components/service-operator/internal/aws/mocks"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("PostgresController", func() {
	var (
		postgres     database.Postgres
		request      reconcile.Request
		reconciler   PostgresReconciler
		cfReconciler *internalawsmocks.MockCloudFormationReconciler
	)

	BeforeEach(func() {
		request = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "test",
				Name:      "test-postgres",
			},
		}
		postgres = database.Postgres{
			TypeMeta: metav1.TypeMeta{
				APIVersion: database.GroupVersion.Group,
				Kind:       "Postgres",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-postgres",
				Namespace: "test",
			},
			Spec: database.PostgresSpec{
				AWS: database.AWS{
					InstanceType: "db.t3.medium",
				},
			},
		}
		k8sClient.Create(context.TODO(), &postgres)
		cfReconciler = internalawsmocks.NewMockCloudFormationReconciler(mockCtrl)
		reconciler = PostgresReconciler{
			Client:                   k8sClient,
			Log:                      log,
			CloudFormationReconciler: cfReconciler,
		}
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

				var updatedPostgres database.Postgres
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      "test-postgres",
				}, &updatedPostgres)

				checkPostgresStatusUpdates(stackData, updatedPostgres)
				Expect(updatedPostgres.Finalizers).To(ContainElement(PostgresFinalizerName))
				Expect(updatedPostgres.ObjectMeta.DeletionTimestamp).To(BeNil())
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

				postgres.Spec.AWS.InstanceType = "db.m5.large"
				postgres.ObjectMeta.Finalizers = append(postgres.Finalizers, PostgresFinalizerName)
				k8sClient.Update(context.TODO(), &postgres)

				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedPostgres database.Postgres
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      "test-postgres",
				}, &updatedPostgres)
				checkPostgresStatusUpdates(stackData, updatedPostgres)
				Expect(updatedPostgres.ObjectMeta.Finalizers).To(ContainElement(PostgresFinalizerName))
				Expect(updatedPostgres.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})

		Context("When deleting a resource", func() {
			BeforeEach(func() {
				postgres.ObjectMeta.Finalizers = append(postgres.Finalizers, PostgresFinalizerName)
				k8sClient.Update(context.TODO(), &postgres)
			})

			It("Should delete the kubernetes resource", func() {
				stackData := internalaws.StackData{}
				cfReconciler.
					EXPECT().
					Reconcile(context.TODO(), gomock.Any(), request, gomock.Any(), true).
					Return(internal.Delete, stackData, nil).
					Times(1)

				k8sClient.Delete(context.TODO(), &postgres)
				result, err := reconciler.Reconcile(request)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Requeue).To(BeTrue())
				Expect(result.RequeueAfter).To(Equal(time.Minute))

				var updatedPostgres database.Postgres
				k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: "test",
					Name:      "test-postgres",
				}, &updatedPostgres)
				Expect(updatedPostgres.ObjectMeta.Finalizers).ToNot(ContainElement(PostgresFinalizerName))
				Expect(updatedPostgres.ObjectMeta.DeletionTimestamp).To(BeNil())
			})
		})
	})
})

func checkPostgresStatusUpdates(stackData internalaws.StackData, postgres database.Postgres) {
	Expect(postgres.Status.ID).To(Equal(stackData.ID))
	Expect(postgres.Status.Status).To(Equal(stackData.Status))
	Expect(postgres.Status.Reason).To(Equal(stackData.Reason))
}
