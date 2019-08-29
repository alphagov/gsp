package awsfakes

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	. "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sanathkr/yaml"
)

type FakeOutput struct {
	Description string      `json:"Description"`
	Value       interface{} `json:"Value"`
	Ref         string      `json:"Ref"`
}

type FakeTemplate struct {
	Outputs map[string]FakeOutput `json:"Outputs"`
}

// NewFakeAWSClient creates a fake AWSClient that attempts to stub the state
// transitions of a cloudformation stack for use in tests without real AWSClient
func NewFakeAWSClient(initialStackState *Stack) *FakeAWSClient {

	client := &FakeAWSClient{}
	stack := initialStackState

	setStackState := func(s, reason string) {
		// initial/uncreated state
		if s == "" {
			stack = nil
			return
		}
		// set state
		stack.StackStatus = aws.String(s)
		stack.StackStatusReason = aws.String(reason)
	}

	setStackOutputsFromTemplate := func(templateYAML string) {
		var t FakeTemplate
		err := yaml.Unmarshal([]byte(templateYAML), &t)
		if err != nil {
			panic(err)
		}
		stack.Outputs = []*Output{}
		for k, v := range t.Outputs {
			stack.Outputs = append(stack.Outputs, &Output{
				Description: aws.String(v.Description),
				OutputKey:   aws.String(k),
				OutputValue: aws.String("FAKED_VALUE"),
			})
		}
	}

	client.DescribeStacksWithContextStub = func(context.Context, *DescribeStacksInput, ...request.Option) (*DescribeStacksOutput, error) {
		if stack == nil {
			return nil, ResourceNotFoundException
		}
		return &DescribeStacksOutput{
			Stacks: []*Stack{stack},
		}, nil
	}

	client.CreateStackWithContextStub = func(_ context.Context, input *CreateStackInput, o ...request.Option) (*CreateStackOutput, error) {
		if stack == nil {
			stack = &Stack{
				StackId:   aws.String(fmt.Sprintf("stack-%d", rand.Intn(10000))),
				StackName: input.StackName,
			}
			setStackState(StackStatusCreateInProgress, "fake-create-stack-called")
			// extract the cloudformation outputs from the given template
			if input.TemplateBody == nil {
				return nil, fmt.Errorf("TemplateBody is required")
			}
			// start timer to swtich to CREATE_COMPLETE state and add Outputs
			go func() {
				time.Sleep(time.Second * 2)
				setStackOutputsFromTemplate(*input.TemplateBody)
				setStackState(StackStatusCreateComplete, "fake-creation-timer-completed")
			}()
			return &CreateStackOutput{
				StackId: stack.StackId,
			}, nil
		}
		return nil, fmt.Errorf("CANNOT_CREATE_ALREADY_CREATED")
	}

	client.DeleteStackWithContextStub = func(context.Context, *DeleteStackInput, ...request.Option) (*DeleteStackOutput, error) {
		if stack == nil {
			return nil, fmt.Errorf("CANNOT_DELETE_BEFORE_CREATE")
		}
		switch *stack.StackStatus {
		case StackStatusCreateComplete, StackStatusUpdateComplete, StackStatusUpdateRollbackComplete, StackStatusRollbackComplete:
			go func() {
				// after a while transition to DELETE_COMPLETE state
				time.Sleep(time.Second * 2)
				setStackState(StackStatusDeleteComplete, "fake-deletion-timer-completed")
			}()
			return &DeleteStackOutput{}, nil
		default:
			return nil, fmt.Errorf("CANNOT_DELETE_FROM_CURRENT_STATE: %s", *stack.StackStatus)
		}
	}

	client.UpdateStackWithContextStub = func(context.Context, *UpdateStackInput, ...request.Option) (*UpdateStackOutput, error) {
		return nil, NoUpdateRequiredException
	}

	return client
}
