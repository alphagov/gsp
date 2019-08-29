package aws

import (
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	goformation "github.com/awslabs/goformation/cloudformation"
)

type Stack interface {
	GetStackName() string
	GetStackTemplate() *goformation.Template
	GetStackCreateParameters() ([]*awscloudformation.Parameter, error)
	GetStackUpdateParameters() ([]*awscloudformation.Parameter, error)
	SetStackStatus(state *awscloudformation.Stack, events []*awscloudformation.StackEvent) error
}
