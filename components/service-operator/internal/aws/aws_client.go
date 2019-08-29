package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . AWSClient

type AWSClient interface {
	DescribeStacksWithContext(aws.Context, *cloudformation.DescribeStacksInput, ...request.Option) (*cloudformation.DescribeStacksOutput, error)
	DescribeStackEventsWithContext(aws.Context, *cloudformation.DescribeStackEventsInput, ...request.Option) (*cloudformation.DescribeStackEventsOutput, error)
	CreateStackWithContext(aws.Context, *cloudformation.CreateStackInput, ...request.Option) (*cloudformation.CreateStackOutput, error)
	UpdateStackWithContext(aws.Context, *cloudformation.UpdateStackInput, ...request.Option) (*cloudformation.UpdateStackOutput, error)
	DeleteStackWithContext(aws.Context, *cloudformation.DeleteStackInput, ...request.Option) (*cloudformation.DeleteStackOutput, error)
}

func NewAWSClient() (AWSClient, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// FIXME: We could/should obtain this automatically... At the moment, we have no plans to run this outside of london.
	awsRegion := "eu-west-2"
	sess.Config.Region = aws.String(awsRegion)
	svc := cloudformation.New(sess, aws.NewConfig())
	return svc, nil
}
