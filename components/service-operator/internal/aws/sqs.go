package aws

import (
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

type SQS struct {
	SQSConfig *queue.SQS
}

func (s *SQS) Template(stackName string) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources["SQSQueue"] = &resources.AWSSQSQueue{
		QueueName: s.SQSConfig.Spec.AWS.QueueName,
	}

	return template
}

func (s *SQS) Parameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}
