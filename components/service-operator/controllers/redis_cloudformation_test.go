package controllers_test

import (
	"context"
	"fmt"
	"time"

	cache "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
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

var _ = Describe("RedisCloudFormationController", func() {

	var within30mins time.Duration = time.Minute * 30
	var within5mins = time.Minute * 5
	var client client.Client
	var ctx context.Context = context.Background()
	var teardown func()

	BeforeEach(func() {
		client, teardown = SetupControllerEnv()
	})

	AfterEach(func() {
		teardown()
	})

	It("Should create and destroy a Redis", func() {

		var (
			name                   = fmt.Sprintf("test-cluster-%s", time.Now().Format("20060102150405"))
			secretName             = "test-redis-secret"
			serviceEntryName       = "test-redis-service-entry"
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
			ecc = cache.Redis{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
				Spec: cache.RedisSpec{
					Secret:       secretName,
					ServiceEntry: serviceEntryName,
					AWS:          cache.RedisAWSSpec{
						NodeType:         "cache.t3.micro",
						EngineVersion:    "5.0.6",
						NumCacheClusters: 2,
					},
				},
			}
			secret       core.Secret
			serviceEntry istio.ServiceEntry
		)

		By("creating an resource with kubernetes api", func() {
			Expect(client.Create(ctx, &ecc)).To(Succeed())
		})

		By("displaying a READY resource status after initial creation", func() {
			Eventually(func() object.State {
				_ = client.Get(ctx, resourceNamespacedName, &ecc)
				return ecc.GetState()
			}, within30mins).Should(Equal(object.ReadyState))
		})

		By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &ecc)
				return ecc.Status.AWS.Status
			}, within5mins).Should(Equal(cloudformation.CreateComplete))
		})

		By("displaying an AWS stack id in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &ecc)
				return ecc.Status.AWS.ID
			}).ShouldNot(BeEmpty())
		})

		By("displaying a stack name prefixed with cluster name in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &ecc)
				return ecc.Status.AWS.Name
			}).Should(ContainSubstring("xxx-redis-test-test-cluster"))
		})

		By("ensuring a finalizer is present on resource to prevent deletion", func() { // TODO: move to cloudformation.Controller unit test
			Eventually(func() []string {
				_ = client.Get(ctx, resourceNamespacedName, &ecc)
				return ecc.Finalizers
			}).Should(ContainElement(cloudformation.Finalizer))
		})

		By("ensuring no DeletionTimestamp exists", func() { // TODO: move to cloudformation.Controller unit test
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &ecc)
				return ecc.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a secret with credentials", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}, within5mins).Should(And(
				HaveKey("ClusterPrimaryRedisHostname"),
				HaveKey("ClusterPrimaryRedisPort"),
				HaveKey("SecretAuthToken"),
			))
		})

		By("creating a service entry with the endpoint", func() {
			Eventually(func() map[string]interface{} {
				_ = client.Get(ctx, serviceEntryNamespacedName0, &serviceEntry)
				return serviceEntry.Spec
			}, within5mins).Should(And(
				HaveKey("hosts"),
				HaveKey("ports"),
				HaveKey("addresses"),
				HaveKey("endpoints"),
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

		By("connecting to resource", func() {
			// TODO
		})

		By("deleting resource with kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &ecc)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &ecc)).To(Succeed())
		})

		By("ensuring the resources have been removed", func() {
			var list cache.RedisList
			Eventually(func() int {
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				return len(list.Items)
			}, within30mins).Should(Equal(0))
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
