package controllers_test

import (
	"context"
	"fmt"
	"time"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("SQS Cloudormation Controller", func() {

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

	It("Should create and destroy an SQS queue", func() {

		var (
			name                   = fmt.Sprintf("test-queue-%s", time.Now().Format("20060102150405"))
			secretName             = "test-secret"
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
			sqs = queue.SQS{
				TypeMeta: metav1.TypeMeta{
					APIVersion: queue.GroupVersion.Group,
					Kind:       "SQS",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						cloudformation.AccessGroupLabel: "test.access.group",
					},
				},
				Spec: queue.SQSSpec{
					Secret: secretName,
				},
			}
			secret core.Secret
		)

		By("creating a prequisite Principal resource with kubernetes api", func() {
			Expect(client.Create(ctx, &principal)).To(Succeed())
		})

		By("creating an SQS resource with kubernetes api", func() {
			Expect(client.Create(ctx, &sqs)).To(Succeed())
		})

		By("displaying a READY resource status after initial creation", func() {
			Eventually(func() object.State {
				_ = client.Get(ctx, resourceNamespacedName, &sqs)
				return sqs.GetState()
			}, timeout).Should(Equal(object.ReadyState))
		})

		By("displaying an AWS CREATE_COMPLETE resource status after initial creation", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &sqs)
				return sqs.Status.AWS.Status
			}, timeout).Should(Equal(cloudformation.CreateComplete))
		})

		By("displaying an AWS stack id in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &sqs)
				return sqs.Status.AWS.ID
			}).ShouldNot(BeEmpty())
		})

		By("displaying a stack name prefixed with cluster name in resource status", func() {
			Eventually(func() string {
				_ = client.Get(ctx, resourceNamespacedName, &sqs)
				return sqs.Status.AWS.Name
			}).Should(ContainSubstring("xxx-sqs-test-test-queue"))
		})

		By("ensuring a finalizaer is present on resource to prevent deletion", func() {
			Eventually(func() []string {
				_ = client.Get(ctx, resourceNamespacedName, &sqs)
				return sqs.Finalizers
			}).Should(ContainElement(cloudformation.Finalizer))
		})

		By("ensuring no DeletionTimestamp exists", func() {
			Eventually(func() bool {
				_ = client.Get(ctx, resourceNamespacedName, &sqs)
				return sqs.ObjectMeta.DeletionTimestamp == nil
			}).Should(BeTrue())
		})

		By("creating a secret with queue connection details", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("QueueURL"))
		})

		By("creating a secret with the prinical role name", func() {
			Eventually(func() map[string][]byte {
				_ = client.Get(ctx, secretNamespacedName, &secret)
				return secret.Data
			}).Should(HaveKey("IAMRoleName"))
		})

		By("deleteing SQS resource with kubernetes api", func() {
			err := client.Get(ctx, resourceNamespacedName, &sqs)
			Expect(err).ToNot(HaveOccurred())
			Expect(client.Delete(ctx, &sqs)).To(Succeed())
		})

		By("ensuring the SQS resources have been removed", func() {
			Eventually(func() int {
				var list queue.SQSList
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
