/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
)

// var (
// 	nonUpdatable = []string{
// 		awscloudformation.StackStatusCreateInProgress,
// 		awscloudformation.StackStatusRollbackInProgress,
// 		awscloudformation.StackStatusDeleteInProgress,
// 		awscloudformation.StackStatusUpdateInProgress,
// 		awscloudformation.StackStatusUpdateCompleteCleanupInProgress,
// 		awscloudformation.StackStatusUpdateRollbackInProgress,
// 		awscloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
// 		awscloudformation.StackStatusReviewInProgress,
// 		awscloudformation.StackStatusDeleteComplete,
// 	}
// 	// Update cannot be called with the stack in it's current state
// 	ErrNonUpdatableState = fmt.Errorf("NON_UPDATABLE_STATE")
// )

var (
	capabilities = []*string{
		aws.String("CAPABILITY_NAMED_IAM"),
	}
)

var (
	// The stack does not exist, or has been deleted
	ErrStackNotFound = fmt.Errorf("STACK_NOT_FOUND")
)

var (
	// String to match in error from aws to detect if nothing to update
	NoUpdatesErrMatch = "No updates"
)

type CloudFormationClient struct {
	// ClusterName is used to prefix any generated names to avoid clashes
	ClusterName string
	// Client is the AWS Client implementation to use see NewAWSClient()
	Client AWSClient
	// PollingInterval is the duration between calls to check state when waiting for apply/destroy to complete
	PollingInterval time.Duration
}

// Apply provisions and reconciles the given cloudformation stack and blocks
// until the stack is no longer in an creating/applying state or a ctx timeout is hit
// Calls should be retried in DeadlineExceeded errors are hit
// Returns any outputs on successful apply.
// Will update stack with current status
func (r *CloudFormationClient) Apply(ctx context.Context, stack Stack, params ...*awscloudformation.Parameter) ([]*awscloudformation.Output, error) {
	// always update stack status
	defer func() {
		_ = r.updateStatus(ctx, stack)
	}()
	// check if exists
	exists, err := r.exists(ctx, stack)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := r.create(ctx, stack, params...)
		if err != nil {
			return nil, err
		}
	}
	_, err = r.waitUntilReadyState(ctx, stack)
	if err != nil {
		return nil, err
	}
	err = r.update(ctx, stack, params...)
	if err != nil {
		return nil, err
	}
	state, err := r.waitUntilReadyState(ctx, stack)
	if err != nil {
		return nil, err
	}
	return state.Outputs, nil
}

func (r *CloudFormationClient) create(ctx context.Context, stack Stack, params ...*awscloudformation.Parameter) error {
	yaml, err := stack.GetStackTemplate().YAML()
	if err != nil {
		return err
	}
	stackParams, err := stack.GetStackCreateParameters()
	if err != nil {
		return err
	}
	_, err = r.Client.CreateStackWithContext(ctx, &awscloudformation.CreateStackInput{
		Capabilities: capabilities,
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stack.GetStackName()),
		Parameters:   append(params, stackParams...),
	})
	if err != nil {
		return err
	}
	return nil
}

// Update the stack and wait for update to complete.
func (r *CloudFormationClient) update(ctx context.Context, stack Stack, params ...*awscloudformation.Parameter) error {
	yaml, err := stack.GetStackTemplate().YAML()
	if err != nil {
		return err
	}
	stackParams, err := stack.GetStackUpdateParameters()
	if err != nil {
		return err
	}
	_, err = r.Client.UpdateStackWithContext(ctx, &awscloudformation.UpdateStackInput{
		Capabilities: capabilities,
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stack.GetStackName()),
		Parameters:   append(params, stackParams...),
	})
	if err != nil && !IsNoUpdateError(err) {
		return err
	}
	return nil
}

// Destroy will attempt to deprovision the cloudformation stack and block until complete or the ctx Deadline expires
// Calls should be retried if DeadlineExceeded errors are hit
// Will update stack with current status
func (r *CloudFormationClient) Destroy(ctx context.Context, stack Stack) error {
	// always update stack status
	defer func() {
		_ = r.updateStatus(ctx, stack)
	}()
	// fetch current state
	state, err := r.get(ctx, stack)
	if err == ErrStackNotFound {
		// resource is already deleted (or never existsed)
		// so we're done here
		return nil
	} else if err != nil {
		// failed to get stack status
		return err
	}
	if *state.StackStatus == awscloudformation.StackStatusDeleteComplete {
		// resource already deleted
		return nil
	}
	// trigger a delete unless we're already in a deleting state
	if *state.StackStatus != awscloudformation.StackStatusDeleteInProgress {
		_, err := r.Client.DeleteStackWithContext(ctx, &awscloudformation.DeleteStackInput{
			StackName: aws.String(stack.GetStackName()),
		})
		if err != nil {
			return err
		}
	}
	_, err = r.waitUntilDestroyedState(ctx, stack)
	if err != nil {
		return err
	}
	return nil
}

// Get fetches the cloudformation stack state
// Returns ErrStackNotFound if stack does not exist
func (r *CloudFormationClient) get(ctx context.Context, stack Stack) (*awscloudformation.Stack, error) {
	describeOutput, err := r.Client.DescribeStacksWithContext(ctx, &awscloudformation.DescribeStacksInput{
		StackName: aws.String(stack.GetStackName()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return nil, ErrStackNotFound
			}
		}
		return nil, err
	}
	if describeOutput == nil {
		return nil, fmt.Errorf("describeOutput was nil, potential issue with AWSClient")
	}
	if len(describeOutput.Stacks) == 0 {
		return nil, fmt.Errorf("describeOutput contained no Stacks, potential issue with AWSClient")
	}
	state := describeOutput.Stacks[0]
	if state.StackStatus == nil {
		return nil, fmt.Errorf("describeOutput contained a nil StackStatus, potential issue with AWSClient")
	}
	return state, nil
}

// update stack with current state and events
func (r *CloudFormationClient) updateStatus(ctx context.Context, stack Stack) error {
	state, _ := r.get(ctx, stack)
	events, _ := r.events(ctx, stack)
	return stack.SetStackStatus(state, events)
}

func (r *CloudFormationClient) events(ctx context.Context, stack Stack) ([]*awscloudformation.StackEvent, error) {
	eventsOutput, err := r.Client.DescribeStackEventsWithContext(ctx, &awscloudformation.DescribeStackEventsInput{
		StackName: aws.String(stack.GetStackName()),
	})
	if err != nil {
		return nil, err
	}
	if eventsOutput == nil {
		return []*awscloudformation.StackEvent{}, nil
	}
	return eventsOutput.StackEvents, nil
}

// Exists checks if the stack has been provisioned
func (r *CloudFormationClient) exists(ctx context.Context, stack Stack) (bool, error) {
	_, err := r.get(ctx, stack)
	if err == ErrStackNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CloudFormationClient) waitUntilReadyState(ctx context.Context, stack Stack) (*awscloudformation.Stack, error) {
	return r.waitUntilState(ctx, stack, []string{
		awscloudformation.StackStatusCreateComplete,
		awscloudformation.StackStatusUpdateComplete,
		awscloudformation.StackStatusUpdateRollbackComplete,
		awscloudformation.StackStatusRollbackComplete,
	})
}

func (r *CloudFormationClient) waitUntilDestroyedState(ctx context.Context, stack Stack) (*awscloudformation.Stack, error) {
	return r.waitUntilState(ctx, stack, []string{
		awscloudformation.StackStatusDeleteComplete,
	})
}

func (r *CloudFormationClient) waitUntilState(ctx context.Context, stack Stack, desiredStates []string) (*awscloudformation.Stack, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		default:
			state, err := r.get(ctx, stack)
			if err != nil {
				return nil, err
			}
			if in(*state.StackStatus, desiredStates) {
				return state, nil
			}
		}
		time.Sleep(r.PollingInterval)
	}
}

func in(needle string, haystack []string) bool {
	for _, s := range haystack {
		if needle == s {
			return true
		}
	}
	return false
}

func IsNoUpdateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), NoUpdatesErrMatch)
}
