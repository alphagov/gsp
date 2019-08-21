package aws

import (
	"fmt"

	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	SQSResourceName = "SQSQueue"

	SQSOutputURL = "QueueURL"
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
	SQSConfig  *queue.SQS
	IAMRoleARN string
}

func (s *SQS) Template(stackName string, tags []resources.Tag) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources[SQSResourceName] = &resources.AWSSQSQueue{
		QueueName: fmt.Sprintf("%s-%s-%s", s.SQSConfig.ClusterName, s.SQSConfig.Namespace, s.SQSConfig.Name),
		Tags:      tags,
	}

	template.Resources[PostgresResourceIAMPolicy] = &resources.AWSIAMPolicy{
		PolicyName:     cloudformation.Join("-", []string{"sqs", "access", cloudformation.GetAtt(SQSResourceName, "QueueName")}),
		PolicyDocument: NewRolePolicyDocument(s.IAMRoleARN, []string{cloudformation.Ref(PostgresResourceCluster)}, allowedActions),
		Roles:          []string{s.IAMRoleARN},
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
