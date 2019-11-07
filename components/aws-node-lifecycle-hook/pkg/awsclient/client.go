package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeawsclient/fake_awsclient.go . Client

type Client interface {
	CompleteLifecycleActionWithContext(context.Context, *autoscaling.CompleteLifecycleActionInput, ...request.Option) (*autoscaling.CompleteLifecycleActionOutput, error)
	RecordLifecycleActionHeartbeatWithContext(context.Context, *autoscaling.RecordLifecycleActionHeartbeatInput, ...request.Option) (*autoscaling.RecordLifecycleActionHeartbeatOutput, error)
}

func New() (Client, error) {
	return autoscaling.New(session.New()), nil
}
