/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&S3Bucket{}, &S3BucketList{})
}

const (
	S3BucketResourceName      = "S3Bucket"
	S3BucketName              = "S3BucketName"
	S3BucketURL               = "S3BucketURL"
	S3BucketRegion            = "S3BucketRegion"
	S3BucketResourceIAMPolicy = "S3BucketIAMPolicy"
	IAMRoleParameterName      = "IAMRoleName"
)

var (
	// allowedActions defines the set of actions that get attached to the role, for the rationale see:
	// https://github.com/alphagov/gsp/blob/master/docs/architecture/adr/ADR041-service-operated-policies.md
	allowedActions = []string{
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
	}
)

// ensure implements required interfaces
var _ cloudformation.Stack = &S3Bucket{}
var _ cloudformation.StackPolicyAttacher = &S3Bucket{}
var _ object.SecretNamer = &S3Bucket{}
var _ cloudformation.ServiceEntryCreator = &S3Bucket{}

// AWS allows specifying configuration for the S3Bucket
type AWS struct {
}

// S3BucketSpec defines the desired state of S3Bucket
type S3BucketSpec struct {
	// AWS specific subsection of the resource.
	AWS AWS `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
	// ServiceEntry name to be used for storing the egress firewall rule to allow tenant access to the bucket
	ServiceEntry string `json:"serviceEntry,omitempty"`
}

// +kubebuilder:object:root=true

// S3BucketList contains a list of S3Bucket
type S3BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []S3Bucket `json:"items"`
}

// +kubebuilder:object:root=true

// S3Bucket is the Schema for the S3Bucket API
type S3Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec          S3BucketSpec `json:"spec,omitempty"`
	object.Status `json:"status,omitempty"`
}

// Name returns the name of the S3Bucket cloudformation stack
func (s *S3Bucket) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "s3", s.Namespace, s.ObjectMeta.Name)
}

// SecretName returns the name of the secret that will be populated with data
func (s *S3Bucket) GetSecretName() string {
	if s.Spec.Secret == "" {
		return s.GetName()
	}
	return s.Spec.Secret
}

// Template returns a cloudformation Template for provisioning an S3Bucket
func (s *S3Bucket) GetStackTemplate() (*cloudformation.Template, error) {
	template := cloudformation.NewTemplate()

	template.Parameters[IAMRoleParameterName] = map[string]string{
		"Type": "String",
	}

	tags := []cloudformation.Tag{
		{
			Key:   "Cluster",
			Value: env.ClusterName(),
		},
		{
			Key:   "Service",
			Value: "s3",
		},
		{
			Key:   "Name",
			Value: s.ObjectMeta.Name,
		},
		{
			Key:   "Namespace",
			Value: s.Namespace,
		},
		{
			Key:   "Environment",
			Value: s.Namespace,
		},
	}

	bucketName := fmt.Sprintf("%s-%s-%s", env.ClusterName(), s.Namespace, s.ObjectMeta.Name)
	template.Resources[S3BucketResourceName] = &cloudformation.AWSS3Bucket{
		BucketName: bucketName,
		Tags:       tags,
	}

	s3BucketArn := cloudformation.Join("", []string{cloudformation.GetAtt(S3BucketResourceName, "Arn")})
	s3ResourceArn := cloudformation.Join("", []string{cloudformation.GetAtt(S3BucketResourceName, "Arn"), "/*"})

	template.Resources[S3BucketResourceIAMPolicy] = &cloudformation.AWSIAMPolicy{
		PolicyName:     cloudformation.Join("-", []string{"s3", "access", bucketName}),
		PolicyDocument: cloudformation.NewRolePolicyDocument([]string{s3BucketArn, s3ResourceArn}, allowedActions),
		Roles: []string{
			cloudformation.Ref(IAMRoleParameterName),
		},
	}

	template.Outputs[S3BucketName] = map[string]interface{}{
		"Description": "S3Bucket name to be returned to the user.",
		"Value":       cloudformation.Ref(S3BucketResourceName),
	}

	template.Outputs[S3BucketURL] = map[string]interface{}{
		"Description": "Bucket URL to be returned to the user.",
		"Value":       fmt.Sprintf("https://%s.s3.eu-west-2.amazonaws.com", bucketName),
	}

	template.Outputs[S3BucketRegion] = map[string]interface{}{
		"Description": "Region that the bucket lives in",
		"Value":       "eu-west-2",
	}

	template.Outputs[IAMRoleParameterName] = map[string]interface{}{
		"Description": "Name of the IAM role with access to bucket",
		"Value":       cloudformation.Ref(IAMRoleParameterName),
	}

	return template, nil
}

func (s *S3Bucket) GetServiceEntryName() string {
	if s.Spec.ServiceEntry == "" {
		return s.GetName()
	}
	return s.Spec.ServiceEntry
}

// ServiceEntry to whitelist egress access to S3 hostname.
func (s *S3Bucket) GetServiceEntrySpecs(outputs cloudformation.Outputs) ([]map[string]interface{}, error) {
	specs := []map[string]interface{}{
		{
			"hosts": []string{
				fmt.Sprintf("%s.s3.%s.amazonaws.com", outputs[S3BucketName], outputs[S3BucketRegion]),
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "https",
					"number":   443,
					"protocol": "TLS",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "DNS",
			"exportTo":   []string{"."},
		},
		{
			"hosts": []string{
				fmt.Sprintf("%s.s3.amazonaws.com", outputs[S3BucketName]),
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "https",
					"number":   443,
					"protocol": "TLS",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "DNS",
			"exportTo":   []string{"."},
		},
	}
	return specs, nil
}

// GetStackRoleParameters returns additional params based on a target principal resource
func (s *S3Bucket) GetStackRoleParameters(roleName string) ([]*cloudformation.Parameter, error) {
	params := []*cloudformation.Parameter{
		{
			ParameterKey:   aws.String(IAMRoleParameterName),
			ParameterValue: aws.String(roleName),
		},
	}
	return params, nil
}
