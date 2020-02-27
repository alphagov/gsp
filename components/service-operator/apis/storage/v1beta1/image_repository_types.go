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
	"context"
	"encoding/json"
	"fmt"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
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
var _ cloudformation.StackObjectEmptier = &ImageRepository{}

// DefaultLifecyclePolicy is the default policy assigned to ECR repositories
var DefaultLifecyclePolicy = cloudformation.ECRLifecyclePolicy{
	Rules: []cloudformation.ECRLifecyclePolicyRule{
		{
			RulePriority: 1,
			Description:  "only keep 100 images",
			Selection: cloudformation.ECRLifecyclePolicySelection{
				TagStatus:   "any",
				CountType:   cloudformation.ECRLifecycleMoreThan,
				CountNumber: 100,
			},
			Action: cloudformation.ECRLifecyclePolicyAction{
				Type: cloudformation.ECRLifecyclePolicyExpire,
			},
		},
	},
}

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

	lifecyclePolicyText, err := json.Marshal(DefaultLifecyclePolicy)
	if err != nil {
		return nil, err
	}

	repositoryName := s.GetAWSName()
	template.Resources[ImageRepositoryResourceName] = &cloudformation.AWSECRRepository{
		RepositoryName: repositoryName,
		LifecyclePolicy: &cloudformation.AWSECRRepository_LifecyclePolicy{
			LifecyclePolicyText: string(lifecyclePolicyText),
		},
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

func (s *ImageRepository) GetAWSName() string {
	return fmt.Sprintf("%s-%s-%s", env.ClusterName(), s.Namespace, s.ObjectMeta.Name)
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

func (s *ImageRepository) Empty(ctx context.Context, client sdk.Client) error {
	var imageIdBatches [][]*ecr.ImageIdentifier
	err := client.DescribeImagesPagesWithContext(
		ctx,
		&ecr.DescribeImagesInput{
			RepositoryName: aws.String(s.GetAWSName()),
		},
		func(images *ecr.DescribeImagesOutput, _ bool) bool {
			var imageIds []*ecr.ImageIdentifier
			for _, image := range images.ImageDetails {
				imageIds = append(imageIds, &ecr.ImageIdentifier{
					ImageDigest: image.ImageDigest,
				})
			}
			if len(imageIds) > 0 {
				imageIdBatches = append(imageIdBatches, imageIds)
			}
			return true
		},
	)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "RepositoryNotFoundException" {
				return nil
			}
		}

		return err
	}

	for _, batch := range imageIdBatches {
		output, err := client.BatchDeleteImageWithContext(
			ctx,
			&ecr.BatchDeleteImageInput{
				RepositoryName: aws.String(s.GetAWSName()),
				ImageIds:       batch,
			},
		)
		if err != nil {
			return err
		}
		for _, failure := range output.Failures {
			return fmt.Errorf("%s", *failure.FailureReason)
		}
	}

	return nil
}
