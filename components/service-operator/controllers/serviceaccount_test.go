package controllers_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
)

var _ = Describe("ServiceAccountController", func() {

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

	It("Should create and destroy a ServiceAccount", func() {
		var (
			name                   = fmt.Sprintf("test-svcacc-%s", time.Now().Format("20060102150405"))
			namespace              = "test"
			resourceNamespacedName = types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			}
			svcacc = core.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
			}
		)

		By("creating an resource with kubernetes api", func() {
			Expect(client.Create(ctx, &svcacc)).To(Succeed())
		})

		By("ensuring no DeletionTimestamp exists", func() {
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &svcacc)
				return svcacc.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a principal with labels", func() {
			Eventually(func() map[string]string {
				var list access.PrincipalList
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				if len(list.Items) == 1 {
					return list.Items[0].ObjectMeta.Labels
				}
				return map[string]string{}
			}, time.Minute*5).Should(HaveKeyWithValue(cloudformation.AccessGroupLabel, "test.access.group"))
		})

		By("creating a principal with an owner reference", func() {
			Eventually(func() []metav1.OwnerReference {
				var list access.PrincipalList
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				if len(list.Items) == 1 {
					return list.Items[0].ObjectMeta.OwnerReferences
				}
				return []metav1.OwnerReference{}
			}).Should(HaveLen(1))
		})

		By("ensuring the EKS Role ARN annotation has been applied", func() {
			Eventually(func() map[string]string {
				_ = client.Get(ctx, resourceNamespacedName, &svcacc)
				return svcacc.ObjectMeta.Annotations
			}, time.Minute*1).Should(HaveKey("eks.amazonaws.com/role-arn"))
		})

		By("deleting resource with kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &svcacc)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &svcacc)).To(Succeed())
		})

		By("ensuring the resources have been removed", func() {
			var list core.ServiceAccountList
			Eventually(func() int {
				err := client.List(ctx, &list)
				Expect(err).ToNot(HaveOccurred())
				return len(list.Items)
			}, timeout).Should(Equal(0))
		})

		// GC will remove this in a real cluster, but we don't have the hooks installed in our tests :(
		// By("ensuring principal has been removed", func() {
		// 	var list access.PrincipalList
		// 	Eventually(func() int {
		// 		err := client.List(ctx, &list)
		// 		Expect(err).ToNot(HaveOccurred())
		// 		return len(list.Items)
		// 	}, time.Second*10).Should(Equal(0))
		// })
	})
})
