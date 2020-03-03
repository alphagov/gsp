package v1beta1_test

import (
	"context"
	"encoding/base64"
	"os"
	"time"

	"github.com/alphagov/gsp/components/service-operator/apis/storage/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk/sdkfakes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"

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

	It("should call the AWS API when emptied", func() {
		client := &sdkfakes.FakeClient{}
		client.DescribeImagesPagesWithContextStub = func(_ context.Context, input *ecr.DescribeImagesInput, fn func(page *ecr.DescribeImagesOutput, lastPage bool) bool, o ...request.Option) error {
			fn(
				&ecr.DescribeImagesOutput{
					ImageDetails: []*ecr.ImageDetail{
						&ecr.ImageDetail{
							ImageDigest:      aws.String("sha256:some long sha256 sum"),
							ImagePushedAt:    aws.Time(time.Now()),
							ImageScanStatus:  &ecr.ImageScanStatus{
								Description:  aws.String("not done"),
								Status:       aws.String("fake client happy"),
							},
							ImageSizeInBytes: aws.Int64(42),
							ImageTags:        []*string{
								aws.String("latest"),
							},
							RegistryId:       input.RegistryId,
							RepositoryName:   input.RepositoryName,
						},
					},
				},
				false,
			)
			fn(
				&ecr.DescribeImagesOutput{
					ImageDetails: []*ecr.ImageDetail{
						&ecr.ImageDetail{
							ImageDigest:      aws.String("sha256:another long sha256 sum"),
							ImagePushedAt:    aws.Time(time.Now()),
							ImageScanStatus:  &ecr.ImageScanStatus{
								Description:  aws.String("not done"),
								Status:       aws.String("fake client happy"),
							},
							ImageSizeInBytes: aws.Int64(42),
							ImageTags:        []*string{
								aws.String("latest"),
							},
							RegistryId:       input.RegistryId,
							RepositoryName:   input.RepositoryName,
						},
					},
				},
				true,
			)
			return nil
		}
		client.BatchDeleteImageWithContextStub = func(_ context.Context, input *ecr.BatchDeleteImageInput, o ...request.Option) (*ecr.BatchDeleteImageOutput, error) {
			return &ecr.BatchDeleteImageOutput{
				Failures: []*ecr.ImageFailure{},
				ImageIds: input.ImageIds,
			}, nil
		}
		err := o.Empty(context.Background(), client)
		Expect(err).ToNot(HaveOccurred())
		Expect(client.BatchDeleteImageWithContextCallCount()).To(Equal(2))

		_, input, _ := client.BatchDeleteImageWithContextArgsForCall(0)
		Expect(input.RepositoryName).To(Equal(aws.String(o.GetAWSName())))
		Expect(input.ImageIds).To(ConsistOf(&ecr.ImageIdentifier{
			ImageDigest: aws.String("sha256:some long sha256 sum"),
		}))

		_, input, _ = client.BatchDeleteImageWithContextArgsForCall(1)
		Expect(input.RepositoryName).To(Equal(aws.String(o.GetAWSName())))
		Expect(input.ImageIds).To(ConsistOf(&ecr.ImageIdentifier{
			ImageDigest: aws.String("sha256:another long sha256 sum"),
		}))
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

			It("should have a lifecycle policy than only keeps last 100 images", func() {
				Expect(repository.LifecyclePolicy.LifecyclePolicyText).To(MatchJSON(`{
					"rules": [{
						"rulePriority": 1,
						"description": "only keep 100 images",
						"selection": {
							"tagStatus": "any",
							"countType": "imageCountMoreThan",
							"countNumber": 100
						},
						"action": {
							"type": "expire"
						}
					}]
				}`))
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
