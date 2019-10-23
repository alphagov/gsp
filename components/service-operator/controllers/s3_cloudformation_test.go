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
	istio "istio.io/istio/pilot/pkg/config/kube/crd"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("S3CloudFormationController", func() {

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

	It("Should create and destroy an S3Bucket bucket", func() {

		var (
			name                   = fmt.Sprintf("test-bucket-%s", time.Now().Format("20060102150405"))
			secretName             = "test-s3-secret"
			serviceEntryName       = "test-s3-service-entry"
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
			serviceEntryNamespacedName0 = types.NamespacedName{
				Namespace: namespace,
				Name:      fmt.Sprintf("%s-0", serviceEntryName),
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
			bucket = storage.S3Bucket{
				TypeMeta: metav1.TypeMeta{
					APIVersion: storage.GroupVersion.Group,
					Kind:       "S3Bucket",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
				Spec: storage.S3BucketSpec{
					Secret:       secretName,
					ServiceEntry: serviceEntryName,
				},
			}
			secret       core.Secret
			serviceEntry istio.ServiceEntry
		)

		By("creating a prequisite Principal resource with kubernetes api", func() {
			Expect(client.Create(ctx, &principal)).To(Succeed())
		})

		By("creating an S3Bucket resource with kubernetes api", func() {
			Expect(client.Create(ctx, &bucket)).To(Succeed())
		})

		By("displaying a READY resource status after initial creation", func() {
			Eventually(func() object.State {
				_ = client.Get(ctx, resourceNamespacedName, &bucket)
				return bucket.GetState()
			}, timeout).Should(Equal(object.ReadyState))
		})

		By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &bucket)
				return bucket.Status.AWS.Status
			}, timeout).Should(Equal(cloudformation.CreateComplete))
		})

		By("displaying an AWS stack id in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &bucket)
				return bucket.Status.AWS.ID
			}).ShouldNot(BeEmpty())
		})

		By("displaying a stack name prefixed with cluster name in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &bucket)
				return bucket.Status.AWS.Name
			}).Should(ContainSubstring("xxx-s3-test-test-bucket"))
		})

		By("ensuring a finalizer is present on resource to prevent deletion", func() {
			Eventually(func() []string {
				_ = client.Get(ctx, resourceNamespacedName, &bucket)
				return bucket.Finalizers
			}).Should(ContainElement(cloudformation.Finalizer))
		})

		By("ensuring no DeletionTimestamp exists", func() {
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &bucket)
				return bucket.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a secret with bucket details", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("S3BucketName"))
		})

		By("creating a secret with the principal role name", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("IAMRoleName"))
		})

		By("creating a service entry with the endpoints", func() {
			Eventually(func() map[string]interface{} {
				_ = client.Get(ctx, serviceEntryNamespacedName0, &serviceEntry)
				return serviceEntry.Spec
			}).Should(And(
				HaveKey("hosts"),
				HaveKey("ports"),
				HaveKey("location"),
				HaveKey("resolution"),
				HaveKey("exportTo"),
			))
		})

		By("creating a service entry with an owner reference", func() {
			Eventually(func() []metav1.OwnerReference {
				_ = client.Get(ctx, serviceEntryNamespacedName0, &serviceEntry)
				return serviceEntry.ObjectMeta.OwnerReferences
			}).Should(HaveLen(1))
		})

		By("deleting S3Bucket resource with Kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &bucket)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &bucket)).To(Succeed())
		})

		By("ensuring the S3Bucket resources have been removed", func() {
			Eventually(func() int {
				var list storage.S3BucketList
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				return len(list.Items)
			}, timeout).Should(Equal(0))
		})

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
