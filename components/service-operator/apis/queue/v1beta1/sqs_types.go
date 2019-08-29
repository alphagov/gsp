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
	"os"

	"github.com/alphagov/gsp/components/service-operator/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	SQSResourceName      = "SQSQueue"
	SQSOutputURL         = "QueueURL"
	SQSResourceIAMPolicy = "SQSSIAMPolicy"
)

var (
	allowedActions = []string{
		"sqs:SendMessage",
		"sqs:ReceiveMessage",
		"sqs:DeleteMessage",
		"sqs:GetQueueAttributes",
	}
)

var IAMRoleParameterName = "IAMRoleName"

// ensure implements StackObject
var _ apis.StackObject = &SQS{}

// +kubebuilder:object:root=true

// SQS is the Schema for the SQS API
type SQS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SQSSpec   `json:"spec,omitempty"`
	Status SQSStatus `json:"status,omitempty"`
}

// Name returns the name of the SQS cloudformation stack
func (s *SQS) GetStackName() string {
	CLUSTER_NAME := os.Getenv("CLUSTER_NAME") // FIXME: this is not right
	return fmt.Sprintf("%s-%s-%s-%s", CLUSTER_NAME, "sqs", s.Namespace, s.ObjectMeta.Name)
}

// SecretName returns the name of the secret that will be populated with data
func (s *SQS) GetSecretName() string {
	return s.Spec.Secret
}

// Template returns a cloudformation Template for provisioning an SQS queue
func (s *SQS) GetStackTemplate() *cloudformation.Template {
	CLUSTER_NAME := os.Getenv("CLUSTER_NAME") // FIXME: this is not right
	template := cloudformation.NewTemplate()

	template.Parameters[IAMRoleParameterName] = map[string]string{
		"Type": "String",
	}

	tags := []resources.Tag{
		{
			Key:   "Cluster",
			Value: CLUSTER_NAME,
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

	queueName := fmt.Sprintf("%s-%s-%s", CLUSTER_NAME, s.Namespace, s.ObjectMeta.Name)
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
		PolicyDocument: internalaws.NewRolePolicyDocument([]string{cloudformation.GetAtt(SQSResourceName, "Arn")}, allowedActions),
		Roles: []string{
			cloudformation.Ref(IAMRoleParameterName),
		},
	}

	template.Outputs[SQSOutputURL] = map[string]interface{}{
		"Description": "SQSQueue URL to be returned to the user.",
		"Value":       cloudformation.Ref(SQSResourceName),
	}

	return template
}

// CreateParameters returns any params used during stack creation
func (s *SQS) GetStackCreateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

// UpdateParameters returns any params used during stack update
func (s *SQS) GetStackUpdateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

// SetFromStack updates status fields based on given cloudformation stack
func (s *SQS) SetStackStatus(state *awscloudformation.Stack, events []*awscloudformation.StackEvent) error {
	if err := s.Status.Service.SetFromStack(state, events); err != nil {
		return err
	}
	if err := s.Status.AWS.SetFromStack(state, events); err != nil {
		return err
	}
	return nil
}

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
	// Important: Run "make" to regenerate code after modifying this file

	// AWS specific subsection of the resource.
	AWS AWS `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
}

// SQSStatus defines the observed state of SQS
type SQSStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Generic service status
	Service apis.ServiceStatus `json:"service"`
	// AWS specific status
	AWS apis.AWSServiceStatus `json:"aws,omitempty"`
}

// +kubebuilder:object:root=true

// SQSList contains a list of SQS
type SQSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SQS{}, &SQSList{})
}
