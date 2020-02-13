package sdk

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sts"
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
	GetAuthorizationTokenWithContext(aws.Context, *ecr.GetAuthorizationTokenInput, ...request.Option) (*ecr.GetAuthorizationTokenOutput, error)
	AssumeRole(roleArn string) Client
	GetRoleCredentials(roleARN string, sessionDuration time.Duration) *credentials.Credentials
}

// NewClient creates a new AWS client that implements the Client interface.
func NewClient(optionalConfig ...*aws.Config) Client {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sess.Config.Region = aws.String(DefaultAWSRegion)
	cfg := aws.NewConfig()
	for _, providedConfig := range optionalConfig {
		cfg = providedConfig
	}
	cfg = cfg.WithRegion(DefaultAWSRegion)
	return &client{
		SecretsManager: secretsmanager.New(sess, cfg),
		CloudFormation: cloudformation.New(sess, cfg),
		ECR:            ecr.New(sess, cfg),
	}
}

// client combines multiple required aws sdk service clients into a
// single kind that share a single session. This makes it easier to configure
// and mock. Use NewClient to create a client.
type client struct {
	*secretsmanager.SecretsManager
	*cloudformation.CloudFormation
	*ecr.ECR
}

func (c *client) GetRoleCredentials(roleARN string, sessionDuration time.Duration) *credentials.Credentials {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return credentials.NewCredentials(&stscreds.AssumeRoleProvider{
		Client:   sts.New(sess),
		RoleARN:  roleARN,
		Duration: sessionDuration,
	})
}

func (c *client) AssumeRole(roleARN string) Client {
	return NewClient(&aws.Config{Credentials: c.GetRoleCredentials(roleARN, stscreds.DefaultDuration)})
}
