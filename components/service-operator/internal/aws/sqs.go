package aws

import (
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	SQSResourceName = "SQSQueue"

	SQSOutputURL = "QueueURL"
)

type SQS struct {
	SQSConfig *queue.SQS
}

func (s *SQS) Template(stackName string) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources[SQSResourceName] = &resources.AWSSQSQueue{
		QueueName: s.SQSConfig.Spec.AWS.QueueName,
	}

	template.Outputs[SQSOutputURL] = map[string]interface{}{
		"Description": "SQSQueue URL to be returned to the user.",
		"Value":       cloudformation.Ref(SQSResourceName),
	}

	return template
}

func (s *SQS) Parameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

func (p *SQS) ResourceType() string {
	return "sqs"
}
