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
	SchemeBuilder.Register(&ImageRepository{}, &ImageRepositoryList{})
}

const (
	ImageRepositoryResourceName      = "ImageRepository"
	ImageRepositoryName              = "ImageRepositoryName"
	ImageRepositoryURI               = "ImageRepositoryURI"
	ImageRepositoryRegion            = "ImageRepositoryRegion"
	ImageRepositoryResourceIAMPolicy = "ImageRepositoryIAMPolicy"
	AccountIdParameterName           = "AWS::AccountId"
	IAMRoleArnParameterName          = "IAMRoleArn"
)

// ensure implements required interfaces
var _ cloudformation.Stack = &ImageRepository{}
var _ cloudformation.StackPolicyAttacher = &ImageRepository{}
var _ object.SecretNamer = &ImageRepository{}
var _ cloudformation.StackSecretOutputter = &ImageRepository{}

// ImageRepositorySpec defines the desired state of ImageRepository
type ImageRepositorySpec struct {
	// AWS specific subsection of the resource.
	AWS AWS `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
}

// +kubebuilder:object:root=true

// ImageRepositoryList contains a list of ImageRepository
type ImageRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageRepository `json:"items"`
}

// +kubebuilder:object:root=true

// ImageRepository is the Schema for the ImageRepository API
type ImageRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec          ImageRepositorySpec `json:"spec,omitempty"`
	object.Status `json:"status,omitempty"`
}

// Name returns the name of the ImageRepository cloudformation stack
func (s *ImageRepository) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "ecr", s.Namespace, s.ObjectMeta.Name)
}

// SecretName returns the name of the secret that will be populated with data
func (s *ImageRepository) GetSecretName() string {
	if s.Spec.Secret == "" {
		return s.GetName()
	}
	return s.Spec.Secret
}

// Template returns a cloudformation Template for provisioning an ImageRepository
func (s *ImageRepository) GetStackTemplate() (*cloudformation.Template, error) {
	template := cloudformation.NewTemplate()

	template.Parameters[IAMRoleParameterName] = map[string]string{
		"Type": "String",
	}

	repositoryName := fmt.Sprintf("%s-%s-%s", env.ClusterName(), s.Namespace, s.ObjectMeta.Name)
	template.Resources[ImageRepositoryResourceName] = &cloudformation.AWSECRRepository{
		RepositoryName: repositoryName,
	}

	imageRepositoryArn := cloudformation.GetAtt(ImageRepositoryResourceName, "Arn")

	template.Resources[ImageRepositoryResourceIAMPolicy] = &cloudformation.AWSIAMPolicy{
		PolicyName:     cloudformation.Join("-", []string{"ecr", repositoryName}),
		PolicyDocument: cloudformation.NewRolePolicyDocument([]string{imageRepositoryArn}, []string{"ecr:*"}),
		Roles: []string{
			cloudformation.Ref(IAMRoleParameterName),
		},
	}

	template.Outputs[ImageRepositoryName] = map[string]interface{}{
		"Description": "Image repository name to be returned to the user.",
		"Value":       cloudformation.Ref(ImageRepositoryResourceName),
	}

	template.Outputs[ImageRepositoryURI] = map[string]interface{}{
		"Description": "Image repository URI to be returned to the user.",
		"Value":       cloudformation.Join("", []string{cloudformation.Ref(AccountIdParameterName), ".dkr.ecr.eu-west-2.amazonaws.com/", repositoryName}),
	}

	template.Outputs[ImageRepositoryRegion] = map[string]interface{}{
		"Description": "Region that the image repository lives in",
		"Value":       "eu-west-2",
	}

	template.Outputs[IAMRoleParameterName] = map[string]interface{}{
		"Description": "Name of the IAM role with access to the image repository",
		"Value":       cloudformation.Ref(IAMRoleParameterName),
	}

	template.Outputs[IAMRoleArnParameterName] = map[string]interface{}{
		"Description": "ARN of the IAM role with access to the image repository",
		"Value":       cloudformation.Join("", []string{"arn:aws:iam::", cloudformation.Ref(AccountIdParameterName), ":role/", cloudformation.Ref(IAMRoleParameterName)}),
	}

	return template, nil
}

// GetStackRoleParameters returns additional params based on a target principal resource
func (s *ImageRepository) GetStackRoleParameters(roleName string) ([]*cloudformation.Parameter, error) {
	params := []*cloudformation.Parameter{
		{
			ParameterKey:   aws.String(IAMRoleParameterName),
			ParameterValue: aws.String(roleName),
		},
	}
	return params, nil
}
