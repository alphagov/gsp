package controllers_test

import (
	"context"
	"fmt"
	"time"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("PrincipalCloudFormationController", func() {

	var timeout time.Duration = time.Minute * 15
	var client client.Client
	var ctx context.Context = context.Background()
	var teardown func()

	BeforeEach(func() {
		client, teardown = SetupControllerEnv()
	})

	AfterEach(func() {
		teardown()
	})

	It("Should create and destroy an IAM role", func() {

		var (
			name                   = fmt.Sprintf("test-role-%s", time.Now().Format("20060102150405"))
			namespace              = "test"
			secretName             = "test-ecr-secret"
			resourceNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			}
			principal = access.Principal{
				TypeMeta: metav1.TypeMeta{
					APIVersion: access.GroupVersion.Group,
					Kind:       "Principal",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      name,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
				Spec: access.PrincipalSpec{
					Secret: secretName,
				},
			}
			secretNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      secretName,
			}
			secret core.Secret
		)

		By("creating a Principal resource with kubernetes api", func() {
			Expect(client.Create(ctx, &principal)).To(Succeed())
		})

		By("displaying a READY resource status after initial creation", func() {
			Eventually(func() object.State {
				_ = client.Get(ctx, resourceNamespacedName, &principal)
				return principal.GetState()
			}, timeout).Should(Equal(object.ReadyState))
		})

		By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &principal)
				return principal.Status.AWS.Status
			}, timeout).Should(Equal(cloudformation.CreateComplete))
		})

		By("displaying an the created role name in status info", func() {
			Eventually(func() map[string]string {
				_ = client.Get(ctx, resourceNamespacedName, &principal)
				return principal.Status.AWS.Info
			}).Should(HaveKey(access.IAMRoleName))
		})

		By("displaying a stack name prefixed with cluster name in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &principal)
				return principal.Status.AWS.Name
			}).Should(ContainSubstring("xxx-principal-test-test-role"))
		})

		By("ensuring a finalizer is present on resource to prevent deletion", func() {
			Eventually(func() []string {
				_ = client.Get(ctx, resourceNamespacedName, &principal)
				return principal.Finalizers
			}).Should(ContainElement(cloudformation.Finalizer))
		})

		By("ensuring no DeletionTimestamp exists", func() {
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &principal)
				return principal.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a secret with registry details", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(And(
				HaveKeyWithValue("ImageRegistryUsername", BeEquivalentTo("AWS")),
				HaveKey("ImageRegistryPassword"),
				HaveKeyWithValue("ImageRegistryEndpoint", BeEquivalentTo("https://011571571136.dkr.ecr.eu-west-2.amazonaws.com")),
			))
		})

		By("ensuring the secret is updated on credentials renewal", func() {
			originalPassword := secret.Data["ImageRegistryPassword"]
			Eventually(func() []byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				Expect(secret.Data["ImageRegistryPassword"]).ToNot(BeEmpty())
				return secret.Data["ImageRegistryPassword"]
			}, time.Minute*5).ShouldNot(BeEquivalentTo(originalPassword))
		})

		By("deleting resource with kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &principal)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &principal)).To(Succeed())
		})

		By("ensuring the resources have been removed", func() {
			var list access.PrincipalList
			Eventually(func() int {
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				return len(list.Items)
			}, timeout).Should(Equal(0))
		})

	})
})
