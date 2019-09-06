package sdk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client

const (
	DefaultAWSRegion = "eu-west-2"
)

type Client interface {
	DescribeStacksWithContext(aws.Context, *cloudformation.DescribeStacksInput, ...request.Option) (*cloudformation.DescribeStacksOutput, error)
	DescribeStackEventsWithContext(aws.Context, *cloudformation.DescribeStackEventsInput, ...request.Option) (*cloudformation.DescribeStackEventsOutput, error)
	CreateStackWithContext(aws.Context, *cloudformation.CreateStackInput, ...request.Option) (*cloudformation.CreateStackOutput, error)
	UpdateStackWithContext(aws.Context, *cloudformation.UpdateStackInput, ...request.Option) (*cloudformation.UpdateStackOutput, error)
	DeleteStackWithContext(aws.Context, *cloudformation.DeleteStackInput, ...request.Option) (*cloudformation.DeleteStackOutput, error)
	GetSecretValueWithContext(aws.Context, *secretsmanager.GetSecretValueInput, ...request.Option) (*secretsmanager.GetSecretValueOutput, error)
}

// NewClient creates a new AWS client that implements the Client interface.
func NewClient() Client {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sess.Config.Region = aws.String(DefaultAWSRegion)
	return &client{
		SecretsManager: secretsmanager.New(sess, aws.NewConfig().WithRegion(DefaultAWSRegion)),
		CloudFormation: cloudformation.New(sess, aws.NewConfig()),
	}
}

// client combines multiple required aws sdk service clients into a
// single kind that share a single session. This makes it easier to configure
// and mock. Use NewClient to create a client.
type client struct {
	*secretsmanager.SecretsManager
	*cloudformation.CloudFormation
}
