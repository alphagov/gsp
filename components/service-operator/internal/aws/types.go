package aws

import (
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
)

type StackData struct {
	ID     string
	Name   string
	Status string
	Reason string

	Events  []*awscloudformation.StackEvent
	Outputs []*awscloudformation.Output
}
