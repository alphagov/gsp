package v1beta1_test

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/storage/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
)

var _ = Describe("S3Bucket", func() {

	var o v1beta1.S3Bucket
	var tags []cloudformation.Tag

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		o = v1beta1.S3Bucket{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.S3BucketSpec{},
		}
		tags = []cloudformation.Tag{
			{Key: "Cluster", Value: env.ClusterName()},
			{Key: "Service", Value: "s3"},
			{Key: "Name", Value: "example"},
			{Key: "Namespace", Value: "default"},
			{Key: "Environment", Value: "default"},
		}
	})

	It("should default secret name to object name", func() {
		Expect(o.GetSecretName()).To(Equal("example"))
	})

	It("should use secret name from spec.Secret if set ", func() {
		o.Spec.Secret = "my-target-secret"
		Expect(o.GetSecretName()).To(Equal("my-target-secret"))
	})

	It("should default service entry name to object name", func() {
		Expect(o.GetServiceEntryName()).To(Equal(o.GetName()))
	})

	It("should use service entry name from spec.ServiceEntry if set", func() {
		o.Spec.ServiceEntry = "my-target-service-entry"
		Expect(o.GetServiceEntryName()).To(Equal("my-target-service-entry"))
	})

	It("should produce the correct service entry", func() {
		outputs := cloudformation.Outputs{
			v1beta1.S3BucketName: "test",
			v1beta1.S3BucketURL:  "testing",
		}

		specs, err := o.GetServiceEntrySpecs(outputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(specs).To(HaveLen(1))
		spec := specs[0]
		Expect(spec).To(And(
			HaveKeyWithValue("resolution", "DNS"),
			HaveKeyWithValue("location", "MESH_EXTERNAL"),
			HaveKey("hosts"),
			HaveKey("ports"),
			HaveKey("exportTo"),
		))
		Expect(spec["hosts"]).To(ContainElement(fmt.Sprintf("%s.s3.eu-west-2.amazonaws.com", outputs[v1beta1.S3BucketName])))
		Expect(spec["ports"]).To(ContainElement(
			map[string]interface{}{
				"name":     "https",
				"number":   443,
				"protocol": "TLS",
			},
		))
		Expect(spec["exportTo"]).To(And(
			HaveLen(1),
			ContainElement("."),
		))
	})

	It("implements runtime.Object", func() {
		o2 := o.DeepCopyObject()
		Expect(o2).ToNot(BeZero())
	})

	Context("cloudformation", func() {

		It("should generate a unique stack name prefixed with cluster name", func() {
			Expect(o.GetStackName()).To(HavePrefix("xxx-s3-default-example"))
		})

		It("should require an IAM role input", func() {
			t := o.GetStackTemplate()
			Expect(t.Parameters).To(HaveKey("IAMRoleName"))
		})

		It("should have outputs for connection details", func() {
			t := o.GetStackTemplate()
			Expect(t.Outputs).To(And(
				HaveKey("S3BucketName"),
				HaveKey("S3BucketURL"),
				HaveKey("IAMRoleName"),
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

		Context("bucket resource", func() {

			var bucket *cloudformation.AWSS3Bucket

			JustBeforeEach(func() {
				t := o.GetStackTemplate()
				Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSS3Bucket{})))
				var ok bool
				bucket, ok = t.Resources[v1beta1.S3BucketResourceName].(*cloudformation.AWSS3Bucket)
				Expect(ok).To(BeTrue())
			})

			It("should have a bucket name prefixed with cluster and namespace name", func() {
				Expect(bucket.BucketName).To(Equal("xxx-default-example"))
			})

			It("should have suitable tags set", func() {
				Expect(bucket.Tags).To(Equal(tags))
			})
		})

		Context("policy resource", func() {
			var policy *cloudformation.AWSIAMPolicy
			var doc cloudformation.PolicyDocument

			JustBeforeEach(func() {
				t := o.GetStackTemplate()
				Expect(t.Resources[v1beta1.S3BucketResourceIAMPolicy]).To(BeAssignableToTypeOf(&cloudformation.AWSIAMPolicy{}))
				policy = t.Resources[v1beta1.S3BucketResourceIAMPolicy].(*cloudformation.AWSIAMPolicy)
				Expect(policy.PolicyDocument).To(BeAssignableToTypeOf(cloudformation.PolicyDocument{}))
				doc = policy.PolicyDocument.(cloudformation.PolicyDocument)
			})

			It("should have a policy name", func() {
				Expect(policy.PolicyName).ToNot(BeEmpty())
			})

			It("should asign policy to the given role name", func() {
				Expect(policy.Roles).To(ContainElement(cloudformation.Ref("IAMRoleName")))
			})

			It("should have a policy document with relevant actions", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]
				Expect(statement.Effect).To(Equal("Allow"))
				Expect(statement.Action).To(ConsistOf(
					"s3:DeleteObject",
					"s3:DeleteObjectVersion",
					"s3:GetAccelerateConfiguration",
					"s3:GetAnalyticsConfiguration",
					"s3:GetBucketAcl",
					"s3:GetBucketCORS",
					"s3:GetBucketLocation",
					"s3:GetBucketLogging",
					"s3:GetBucketNotification",
					"s3:GetBucketObjectLockConfiguration",
					"s3:GetBucketPolicy",
					"s3:GetBucketPolicyStatus",
					"s3:GetBucketPublicAccessBlock",
					"s3:GetBucketRequestPayment",
					"s3:GetBucketTagging",
					"s3:GetBucketVersioning",
					"s3:GetBucketWebsite",
					"s3:GetEncryptionConfiguration",
					"s3:GetInventoryConfiguration",
					"s3:GetLifecycleConfiguration",
					"s3:GetMetricsConfiguration",
					"s3:GetObject",
					"s3:GetObjectAcl",
					"s3:GetObjectLegalHold",
					"s3:GetObjectRetention",
					"s3:GetObjectTagging",
					"s3:GetObjectTorrent",
					"s3:GetObjectVersion",
					"s3:GetObjectVersionAcl",
					"s3:GetObjectVersionForReplication",
					"s3:GetObjectVersionTagging",
					"s3:GetObjectVersionTorrent",
					"s3:GetReplicationConfiguration",
					"s3:ListBucket",
					"s3:ListBucketByTags",
					"s3:ListBucketMultipartUploads",
					"s3:ListBucketVersions",
					"s3:ListMultipartUploadParts",
					"s3:PutBucketObjectLockConfiguration",
					"s3:PutObject",
					"s3:PutObjectLegalHold",
					"s3:PutObjectRetention",
					"s3:PutObjectVersionAcl",
					"s3:ReplicateObject",
					"s3:RestoreObject",
				))
			})

			It("is scoped to the correct resources (a single S3 bucket and its objects)", func() {
				Expect(doc.Statement).To(HaveLen(1))
				statement := doc.Statement[0]

				wantedBucketArn, err := base64.StdEncoding.DecodeString(cloudformation.Join("", []string{
					cloudformation.GetAtt(v1beta1.S3BucketResourceName, "Arn"),
				}))
				Expect(err).ToNot(HaveOccurred())

				wantedResourceArn, err := base64.StdEncoding.DecodeString(cloudformation.Join("", []string{
					cloudformation.GetAtt(v1beta1.S3BucketResourceName, "Arn"),
					"/*",
				}))
				Expect(err).ToNot(HaveOccurred())

				s3BucketArn, err := base64.StdEncoding.DecodeString(statement.Resource[0])
				Expect(err).ToNot(HaveOccurred())
				Expect(string(s3BucketArn)).To(Equal(string(wantedBucketArn)))

				s3ResourceArn, err := base64.StdEncoding.DecodeString(statement.Resource[1])
				Expect(err).ToNot(HaveOccurred())
				Expect(string(s3ResourceArn)).To(Equal(string(wantedResourceArn)))
			})
		})
	})
})
