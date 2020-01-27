package controllers_test

import (
	"context"
	"fmt"
	"time"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	storage "github.com/alphagov/gsp/components/service-operator/apis/storage/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ImageRepositoryCloudFormationController", func() {

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

	It("Should create and destroy an ECR repository", func() {

		var (
			name                   = fmt.Sprintf("test-image-%s", time.Now().Format("20060102150405"))
			secretName             = "test-ecr-secret"
			principalName          = "test-role"
			namespace              = "test"
			resourceNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			}
			secretNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      secretName,
			}
			principal = access.Principal{
				TypeMeta: metav1.TypeMeta{
					APIVersion: access.GroupVersion.Group,
					Kind:       "Principal",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      principalName,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
			}
			repository = storage.ImageRepository{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ImageRepository",
					APIVersion: storage.GroupVersion.Group,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:         name,
					GenerateName: "",
					Namespace:    namespace,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
				Spec: storage.ImageRepositorySpec{
					Secret: secretName,
				},
			}
			secret core.Secret
		)

		By("creating a prequisite Principal resource with kubernetes api", func() {
			Expect(client.Create(ctx, &principal)).To(Succeed())
		})

		By("creating an ImageRepository resource with kubernetes api", func() {
			Expect(client.Create(ctx, &repository)).To(Succeed())
		})

		By("displaying a READY resource status after initial creation", func() {
			Eventually(func() object.State {
				_ = client.Get(ctx, resourceNamespacedName, &repository)
				return repository.GetState()
			}, timeout).Should(Equal(object.ReadyState))
		})

		By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &repository)
				return repository.Status.AWS.Status
			}, timeout).Should(Equal(cloudformation.CreateComplete))
		})

		By("displaying an AWS stack id in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &repository)
				return repository.Status.AWS.ID
			}).ShouldNot(BeEmpty())
		})

		By("displaying a stack name prefixed with cluster name in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &repository)
				return repository.Status.AWS.Name
			}).Should(ContainSubstring("xxx-ecr-test-test-image"))
		})

		By("ensuring a finalizer is present on resource to prevent deletion", func() {
			Eventually(func() []string {
				_ = client.Get(ctx, resourceNamespacedName, &repository)
				return repository.Finalizers
			}).Should(ContainElement(cloudformation.Finalizer))
		})

		By("ensuring no DeletionTimestamp exists", func() {
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &repository)
				return repository.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a secret with repository details", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(And(
				HaveKeyWithValue("ImageRepositoryName", ContainSubstring("xxx-test-test-image")),
				HaveKeyWithValue("ImageRepositoryURI", ContainSubstring("011571571136.dkr.ecr.eu-west-2.amazonaws.com/xxx-test-test-image")),
				HaveKeyWithValue("IAMRoleName", BeEquivalentTo("svcop-xxx-test-test-role")),
			))
		})

		By("deleting ImageRepository resource with Kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &repository)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &repository)).To(Succeed())
		})

		By("ensuring the ImageRepository resources have been removed", func() {
			Eventually(func() int {
				var list storage.ImageRepositoryList
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				return len(list.Items)
			}, timeout).Should(Equal(0))
		})
	})
})
