package v1beta1_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
)

var _ = Describe("Principal", func() {

	var principal v1beta1.Principal

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		principal = v1beta1.Principal{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
		}
	})

	It("should return a unique role name", func() {
		Expect(principal.GetRoleName()).To(Equal("svcop-xxx-default-example"))
	})

	It("should implement runtime.Object", func() {
		o2 := principal.DeepCopyObject()
		Expect(o2).ToNot(BeZero())
	})

	Context("cloudformation", func() {

		It("should generate a unique stack name prefixed with cluster name", func() {
			Expect(principal.GetStackName()).To(HavePrefix("xxx-principal-default-example"))
		})

		It("should have outputs for role name", func() {
			t := principal.GetStackTemplate()
			Expect(t.Outputs).To(HaveKey("IAMRoleName"))
		})

		It("should whitelist the IAMRoleName output", func() {
			whitelist := principal.GetStackOutputWhitelist()
			Expect(whitelist).To(ContainElement("IAMRoleName"))
		})

		Context("role resource", func() {

			var role *cloudformation.AWSIAMRole

			JustBeforeEach(func() {
				t := principal.GetStackTemplate()
				Expect(t.Resources[v1beta1.IAMRoleResourceName]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMRole{}))
				role = t.Resources[v1beta1.IAMRoleResourceName].(*cloudformation.AWSIAMRole)
			})

			It("should set unique role name", func() {
				Expect(role.RoleName).To(Equal(principal.GetRoleName()))
			})

			It("should set a permissions boundary", func() {
				Expect(role.PermissionsBoundary).ToNot(BeEmpty())
			})

		})

	})

})
