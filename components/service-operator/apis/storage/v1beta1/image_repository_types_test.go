package v1beta1_test

import (
	"encoding/base64"
	"os"

	"github.com/alphagov/gsp/components/service-operator/apis/storage/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ImageRepository", func() {

	var o v1beta1.ImageRepository

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		o = v1beta1.ImageRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.ImageRepositorySpec{},
		}
	})

	It("should default secret name to object name", func() {
		Expect(o.GetSecretName()).To(Equal("example"))
	})

	It("should use secret name from spec.Secret if set ", func() {
		o.Spec.Secret = "my-target-secret"
		Expect(o.GetSecretName()).To(Equal("my-target-secret"))
	})

	It("implements runtime.Object", func() {
		o2 := o.DeepCopyObject()
		Expect(o2).ToNot(BeZero())
	})

	Context("cloudformation", func() {

		It("should generate a unique stack name prefixed with cluster name", func() {
			Expect(o.GetStackName()).To(HavePrefix("xxx-ecr-default-example"))
		})

		It("should require an IAM role input", func() {
			t, err := o.GetStackTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Parameters).To(HaveKey("IAMRoleName"))
		})

		It("should have outputs for connection details", func() {
			t, err := o.GetStackTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Outputs).To(And(
				HaveKey("ImageRepositoryName"),
				HaveKey("ImageRepositoryURI"),
				HaveKey("ImageRepositoryRegion"),
				HaveKey("IAMRoleArn"),
			))
		})

		It("should map role name to role parameter", func() {
			params, err := o.GetStackRoleParameters("fake-role")
			Expect(err).ToNot(HaveOccurred())
			Expect(params).To(ContainElement(&cloudformation.Parameter{
				ParameterKey:   aws.String("IAMRoleName"),
				ParameterValue: aws.String("fake-role"),
			}))
		})

		Context("ecr repository resource", func() {

			var repository *cloudformation.AWSECRRepository

			JustBeforeEach(func() {
				t, err := o.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSECRRepository{})))
				var ok bool
				repository, ok = t.Resources[v1beta1.ImageRepositoryResourceName].(*cloudformation.AWSECRRepository)
				Expect(ok).To(BeTrue())
			})

			It("should have a repository name prefixed with cluster and namespace name", func() {
				Expect(repository.RepositoryName).To(Equal("xxx-default-example"))
			})
		})

		Context("policy resource", func() {
			var policy *cloudformation.AWSIAMPolicy
			var doc cloudformation.PolicyDocument

			JustBeforeEach(func() {
				t, err := o.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources[v1beta1.ImageRepositoryResourceIAMPolicy]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMPolicy{}))
				policy = t.Resources[v1beta1.ImageRepositoryResourceIAMPolicy].(*cloudformation.AWSIAMPolicy)
				Expect(policy.PolicyDocument).To(BeAssignableToTypeOf(cloudformation.PolicyDocument{}))
				doc = policy.PolicyDocument.(cloudformation.PolicyDocument)
			})

			It("should have a policy name", func() {
				Expect(policy.PolicyName).ToNot(BeEmpty())
			})

			It("should assign policy to the given role name", func() {
				Expect(policy.Roles).To(ContainElement(cloudformation.Ref("IAMRoleName")))
			})

			It("should have a policy document with relevant actions", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Action).To(ConsistOf(
					"ecr:*",
				))
			})

			It("is scoped to a single repository", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]

				wantedRepositoryArn, err := base64.StdEncoding.DecodeString(cloudformation.GetAtt(v1beta1.ImageRepositoryResourceName, "Arn"))
				Expect(err).ToNot(HaveOccurred())

				repositoryArn, err := base64.StdEncoding.DecodeString(statement.Resource[0])
				Expect(err).ToNot(HaveOccurred())
				Expect(string(repositoryArn)).To(Equal(string(wantedRepositoryArn)))
			})
		})
	})
})
