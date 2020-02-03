package v1beta1_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
)

var _ = Describe("Principal", func() {

	var principal v1beta1.Principal
	var principalWithSvcAccTrust v1beta1.Principal

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
		principalWithSvcAccTrust = v1beta1.Principal{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.PrincipalSpec{
				TrustServiceAccount: "testacc",
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

		It("should have expected output keys", func() {
			t, err := principal.GetStackTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Outputs).To(And(
				HaveKey("IAMRoleName"),
				HaveKey("IAMRoleArn"),
			))
		})

		It("should safelist the IAMRoleName output", func() {
			safelist := principal.GetStackOutputWhitelist()
			Expect(safelist).To(ContainElement("IAMRoleName"))
		})

		Context("role resource", func() {
			var role *cloudformation.AWSIAMRole

			JustBeforeEach(func() {
				t, err := principal.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources[v1beta1.IAMRoleResourceName]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMRole{}))
				role = t.Resources[v1beta1.IAMRoleResourceName].(*cloudformation.AWSIAMRole)
			})

			It("should set unique role name", func() {
				Expect(role.RoleName).To(Equal(principal.GetRoleName()))
			})

			It("should set a permissions boundary", func() {
				Expect(role.PermissionsBoundary).ToNot(BeEmpty())
			})

			It("should set an assume role policy", func() {
				subEncoded, ok := role.AssumeRolePolicyDocument.(string)
				Expect(ok).To(Equal(true))

				subJson, err := base64.StdEncoding.DecodeString(subEncoded)
				Expect(err).ToNot(HaveOccurred())
				type SubStruct struct {
					FnSub string `json:"Fn::Sub"`
				}
				var sub SubStruct
				err = json.Unmarshal(subJson, &sub)
				Expect(err).ToNot(HaveOccurred())

				var arpd cloudformation.AssumeRolePolicyDocument
				err = json.Unmarshal([]byte(sub.FnSub), &arpd)
				Expect(err).ToNot(HaveOccurred())

				Expect(arpd.Version).To(Equal("2012-10-17"))
				Expect(arpd.Statement).To(HaveLen(1))
				statement := arpd.Statement[0]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Principal.Federated).To(BeEmpty())
				Expect(statement.Principal.AWS).To(ConsistOf(
					fmt.Sprintf("${%s}", v1beta1.IAMRolePrincipalParameterName),
					fmt.Sprintf("${%s}", v1beta1.ServiceOperatorIAMRoleArn),
				))
				Expect(statement.Action).To(ConsistOf("sts:AssumeRole"))
			})

		})

		Context("role resource with service account trust", func() {
			var role *cloudformation.AWSIAMRole

			JustBeforeEach(func() {
				t, err := principalWithSvcAccTrust.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources[v1beta1.IAMRoleResourceName]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMRole{}))
				role = t.Resources[v1beta1.IAMRoleResourceName].(*cloudformation.AWSIAMRole)
			})

			It("should set an assume role policy with references to a service account if configured", func() {
				subEncoded, ok := role.AssumeRolePolicyDocument.(string)
				Expect(ok).To(Equal(true))

				subJson, err := base64.StdEncoding.DecodeString(subEncoded)
				Expect(err).ToNot(HaveOccurred())
				type SubStruct struct {
					FnSub string `json:"Fn::Sub"`
				}
				var sub SubStruct
				err = json.Unmarshal(subJson, &sub)
				Expect(err).ToNot(HaveOccurred())

				var arpd cloudformation.AssumeRolePolicyDocument
				err = json.Unmarshal([]byte(sub.FnSub), &arpd)
				Expect(err).ToNot(HaveOccurred())

				Expect(arpd.Version).To(Equal("2012-10-17"))
				Expect(arpd.Statement).To(HaveLen(2))

				statement := arpd.Statement[0]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Principal.Federated).To(BeEmpty())
				Expect(statement.Principal.AWS).To(ConsistOf(
					fmt.Sprintf("${%s}", v1beta1.IAMRolePrincipalParameterName),
					fmt.Sprintf("${%s}", v1beta1.ServiceOperatorIAMRoleArn),
				))
				Expect(statement.Action).To(ConsistOf("sts:AssumeRole"))

				statement = arpd.Statement[1]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Principal.AWS).To(BeEmpty())
				Expect(statement.Principal.Federated).To(ConsistOf(fmt.Sprintf("${%s}", v1beta1.IAMOIDCProviderARNParameterName)))
				Expect(statement.Condition.StringEquals).To(And(
					HaveLen(1),
					HaveKeyWithValue(
						fmt.Sprintf("${%s}:sub", v1beta1.IAMOIDCProviderURLParameterName),
						fmt.Sprintf("system:serviceaccount:%s:%s", principalWithSvcAccTrust.GetNamespace(), principalWithSvcAccTrust.Spec.TrustServiceAccount),
					),
				))
			})
		})

		Context("policy resource", func() {
			var policy *cloudformation.AWSIAMPolicy
			var doc cloudformation.PolicyDocument

			JustBeforeEach(func() {
				t, err := principal.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources[v1beta1.SharedPolicyResourceName]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMPolicy{}))
				policy = t.Resources[v1beta1.SharedPolicyResourceName].(*cloudformation.AWSIAMPolicy)
				Expect(policy.PolicyDocument).To(BeAssignableToTypeOf(cloudformation.PolicyDocument{}))
				doc = policy.PolicyDocument.(cloudformation.PolicyDocument)
			})

			It("should have a policy name", func() {
				Expect(policy.PolicyName).To(Equal(principal.GetRoleName()))
			})

			It("should assign policy to the given role name", func() {
				Expect(policy.Roles).To(ContainElement(cloudformation.Ref(v1beta1.IAMRoleResourceName)))
			})

			It("should have a policy document with relevant actions", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Action).To(ConsistOf(
					"ecr:GetAuthorizationToken",
				))
			})

			It("is generally scoped", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]
				Expect(statement.Resource[0]).To(Equal("*"))
			})
		})
	})

})
