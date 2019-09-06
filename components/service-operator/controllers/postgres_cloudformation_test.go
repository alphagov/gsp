package controllers_test

import (
	"context"
	"fmt"
	"time"

	database "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("PostgresCloudFormationController", func() {

	var timeout time.Duration = time.Minute * 30
	var client client.Client
	var ctx context.Context = context.Background()
	var teardown func()

	BeforeEach(func() {
		client, teardown = SetupControllerEnv()
	})

	AfterEach(func() {
		teardown()
	})

	It("Should create and destroy an Postgres database", func() {

		var (
			name                   = fmt.Sprintf("test-db-%s", time.Now().Format("20060102150405"))
			secretName             = "test-secret"
			namespace              = "test"
			resourceNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			}
			secretNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      secretName,
			}
			pg = database.Postgres{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
				Spec: database.PostgresSpec{
					Secret: secretName,
					AWS: database.PostgresAWSSpec{
						InstanceCount: 1,
						InstanceType:  "db.t3.medium",
					},
				},
			}
			secret core.Secret
		)

		By("creating an resource with kubernetes api", func() {
			Expect(client.Create(ctx, &pg)).To(Succeed())
		})

		By("displaying a READY resource status after initial creation", func() {
			Eventually(func() object.State {
				_ = client.Get(ctx, resourceNamespacedName, &pg)
				return pg.GetState()
			}, timeout).Should(Equal(object.ReadyState))
		})

		By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &pg)
				return pg.Status.AWS.Status
			}, timeout).Should(Equal(cloudformation.CreateComplete))
		})

		By("displaying an AWS stack id in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &pg)
				return pg.Status.AWS.ID
			}).ShouldNot(BeEmpty())
		})

		By("displaying a stack name prefixed with cluster name in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &pg)
				return pg.Status.AWS.Name
			}).Should(ContainSubstring("xxx-postgres-test-test-db"))
		})

		By("ensuring a finalizaer is present on resource to prevent deletion", func() { // TODO: move to cloudformation.Controller unit test
			Eventually(func() []string {
				_ = client.Get(ctx, resourceNamespacedName, &pg)
				return pg.Finalizers
			}).Should(ContainElement(cloudformation.Finalizer))
		})

		By("ensuring no DeletionTimestamp exists", func() { // TODO: move to cloudformation.Controller unit test
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &pg)
				return pg.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a secret with username", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("Username"))
		})

		By("creating a secret with password", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("Password"))
		})

		By("creating a secret with endpoint", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("Endpoint"))
		})

		By("connecting to resource", func() {
			// TODO
		})

		By("deleteing resource with kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &pg)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &pg)).To(Succeed())
		})

		By("ensuring the resources have been removed", func() {
			var list database.PostgresList
			Eventually(func() int {
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				return len(list.Items)
			}, timeout).Should(Equal(0))
		})

		// GC will remove this in a real cluster, but we don't have the hooks installed in our tests :(
		// By("ensuring secret has been removed", func() {
		// 	var secretList core.SecretList
		// 	Eventually(func() int {
		// 		err := client.List(ctx, &secretList)
		// 		Expect(err).ToNot(HaveOccurred())
		// 		return len(secretList.Items)
		// 	}, time.Second*10).Should(Equal(0))
		// })
	})
})
