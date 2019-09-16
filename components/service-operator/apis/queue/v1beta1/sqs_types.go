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
	SchemeBuilder.Register(&SQS{}, &SQSList{})
}

const (
	SQSResourceName      = "SQSQueue"
	SQSOutputURL         = "QueueURL"
	SQSResourceIAMPolicy = "SQSSIAMPolicy"
	IAMRoleParameterName = "IAMRoleName"
)

var (
	allowedActions = []string{
		"sqs:SendMessage",
		"sqs:ReceiveMessage",
		"sqs:DeleteMessage",
		"sqs:GetQueueAttributes",
	}
)

// ensure implements required interfaces
var _ cloudformation.Stack = &SQS{}
var _ cloudformation.StackPolicyAttacher = &SQS{}
var _ object.SecretNamer = &SQS{}

// AWS allows specifying configuration for the SQS queue
type AWS struct {
	ContentBasedDeduplication     bool   `json:"contentBasedDeduplication,omitempty"`
	DelaySeconds                  int    `json:"delaySeconds,omitempty"`
	FifoQueue                     bool   `json:"fifoQueue,omitempty"`
	MaximumMessageSize            int    `json:"maximumMessageSize,omitempty"`
	MessageRetentionPeriod        int    `json:"messageRetentionPeriod,omitempty"`
	ReceiveMessageWaitTimeSeconds int    `json:"receiveMessageWaitTimeSeconds,omitempty"`
	RedrivePolicy                 string `json:"redrivePolicy,omitempty"`
	VisibilityTimeout             int    `json:"visibilityTimeout,omitempty"`
}

// SQSSpec defines the desired state of SQS
type SQSSpec struct {
	// AWS specific subsection of the resource.
	AWS AWS `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
}

// +kubebuilder:object:root=true

// SQSList contains a list of SQS
type SQSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQS `json:"items"`
}

// +kubebuilder:object:root=true

// SQS is the Schema for the SQS API
type SQS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec          SQSSpec `json:"spec,omitempty"`
	object.Status `json:"status,omitempty"`
}

// Name returns the name of the SQS cloudformation stack
func (s *SQS) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "sqs", s.Namespace, s.ObjectMeta.Name)
}

// SecretName returns the name of the secret that will be populated with data
func (s *SQS) GetSecretName() string {
	if s.Spec.Secret == "" {
		return s.GetName()
	}
	return s.Spec.Secret
}

// Template returns a cloudformation Template for provisioning an SQS queue
func (s *SQS) GetStackTemplate() *cloudformation.Template {
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
			Value: "sqs",
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

	queueName := fmt.Sprintf("%s-%s-%s", env.ClusterName(), s.Namespace, s.ObjectMeta.Name)
	template.Resources[SQSResourceName] = &resources.AWSSQSQueue{
		QueueName:                     queueName,
		Tags:                          tags,
		ContentBasedDeduplication:     s.Spec.AWS.ContentBasedDeduplication,
		DelaySeconds:                  s.Spec.AWS.DelaySeconds,
		FifoQueue:                     s.Spec.AWS.FifoQueue,
		MaximumMessageSize:            s.Spec.AWS.MaximumMessageSize,
		MessageRetentionPeriod:        s.Spec.AWS.MessageRetentionPeriod,
		ReceiveMessageWaitTimeSeconds: s.Spec.AWS.ReceiveMessageWaitTimeSeconds,
		RedrivePolicy:                 s.Spec.AWS.RedrivePolicy,
		VisibilityTimeout:             s.Spec.AWS.VisibilityTimeout,
	}

	template.Resources[SQSResourceIAMPolicy] = &resources.AWSIAMPolicy{
		PolicyName:     cloudformation.Join("-", []string{"sqs", "access", cloudformation.GetAtt(SQSResourceName, "QueueName")}),
		PolicyDocument: cloudformation.NewRolePolicyDocument([]string{cloudformation.GetAtt(SQSResourceName, "Arn")}, allowedActions),
		Roles: []string{
			cloudformation.Ref(IAMRoleParameterName),
		},
	}

	template.Outputs[SQSOutputURL] = map[string]interface{}{
		"Description": "SQSQueue URL to be returned to the user.",
		"Value":       cloudformation.Ref(SQSResourceName),
	}

	template.Outputs[IAMRoleParameterName] = map[string]interface{}{
		"Description": "Name of the IAM role with access to queue",
		"Value":       cloudformation.Ref(IAMRoleParameterName),
	}

	return template
}

// GetStackRoleParameters returns additional params based on a target principal resource
func (s *SQS) GetStackRoleParameters(roleName string) ([]*cloudformation.Parameter, error) {
	params := []*cloudformation.Parameter{
		{
			ParameterKey:   aws.String(IAMRoleParameterName),
			ParameterValue: aws.String(roleName),
		},
	}
	return params, nil
}
