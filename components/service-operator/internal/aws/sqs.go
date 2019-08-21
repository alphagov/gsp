package aws

import (
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	SQSResourceName = "SQSQueue"

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

type SQS struct {
	SQSConfig   *queue.SQS
	QueueName   string
	IAMRoleName string
}

func (s *SQS) Template(stackName string, tags []resources.Tag) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources[SQSResourceName] = &resources.AWSSQSQueue{
		QueueName: s.QueueName,
		Tags:      tags,
	}

	template.Resources[SQSResourceIAMPolicy] = &resources.AWSIAMPolicy{
		PolicyName:     cloudformation.Join("-", []string{"sqs", "access", cloudformation.GetAtt(SQSResourceName, "QueueName")}),
		PolicyDocument: NewRolePolicyDocument([]string{cloudformation.GetAtt(SQSResourceName, "Arn")}, allowedActions),
		Roles:          []string{s.IAMRoleName},
	}

	template.Outputs[SQSOutputURL] = map[string]interface{}{
		"Description": "SQSQueue URL to be returned to the user.",
		"Value":       cloudformation.Ref(SQSResourceName),
	}

	return template
}

func (s *SQS) CreateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

func (s *SQS) UpdateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

func (p *SQS) ResourceType() string {
	return "sqs"
}
