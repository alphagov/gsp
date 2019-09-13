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
	"github.com/awslabs/goformation/cloudformation/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&S3Bucket{}, &S3BucketList{})
}

const (
	S3BucketResourceName      = "S3Bucket"
	S3BucketName              = "S3BucketName"
	S3BucketResourceIAMPolicy = "S3BucketIAMPolicy"
	IAMRoleParameterName      = "IAMRoleName"
)

var (
	allowedActions = []string{
		"s3:Get*",
		"s3:Put*",
		"s3:Delete*",
	}
)

// ensure implements required interfaces
var _ cloudformation.Stack = &S3Bucket{}
var _ cloudformation.StackPolicyAttacher = &S3Bucket{}
var _ object.SecretNamer = &S3Bucket{}

// AWS allows specifying configuration for the S3Bucket
type AWS struct {
}

// S3BucketSpec defines the desired state of S3Bucket
type S3BucketSpec struct {
	// AWS specific subsection of the resource.
	AWS AWS `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
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
func (s *S3Bucket) GetStackTemplate() *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Parameters[IAMRoleParameterName] = map[string]string{
		"Type": "String",
	}

	tags := []resources.Tag{
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
	template.Resources[S3BucketResourceName] = &resources.AWSS3Bucket{
		BucketName: bucketName,
		Tags:       tags,
	}

	template.Resources[S3BucketResourceIAMPolicy] = &resources.AWSIAMPolicy{
		PolicyName:     cloudformation.Join("-", []string{"s3", "access", bucketName}),
		PolicyDocument: cloudformation.NewRolePolicyDocument([]string{cloudformation.GetAtt(S3BucketResourceName, "Arn")}, allowedActions),
		Roles: []string{
			cloudformation.Ref(IAMRoleParameterName),
		},
	}

	template.Outputs[S3BucketName] = map[string]interface{}{
		"Description": "S3Bucket name to be returned to the user.",
		"Value":       cloudformation.Ref(S3BucketResourceName),
	}

	template.Outputs[IAMRoleParameterName] = map[string]interface{}{
		"Description": "Name of the IAM role with access to bucket",
		"Value":       cloudformation.Ref(IAMRoleParameterName),
	}

	return template
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
