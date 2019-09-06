package cloudformation

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	goformation "github.com/awslabs/goformation/cloudformation"
	goresources "github.com/awslabs/goformation/cloudformation/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Alias types from the various cloudformation packages so we can access
// relevant parts via this package for convinience
type State = cloudformation.Stack
type StateEvent = cloudformation.StackEvent
type Output = cloudformation.Output
type DescribeStacksInput = cloudformation.DescribeStacksInput
type CreateStackInput = cloudformation.CreateStackInput
type UpdateStackInput = cloudformation.UpdateStackInput
type DeleteStackInput = cloudformation.DeleteStackInput
type DescribeStackEventsInput = cloudformation.DescribeStackEventsInput
type DescribeStacksOutput = cloudformation.DescribeStacksOutput
type GetSecretValueInput = secretsmanager.GetSecretValueInput
type Parameter = cloudformation.Parameter
type Template = goformation.Template
type Tag = goresources.Tag
type AWSRDSDBCluster = goresources.AWSRDSDBCluster
type AWSRDSDBInstance = goresources.AWSRDSDBInstance
type AWSRDSDBClusterParameterGroup = goresources.AWSRDSDBClusterParameterGroup
type AWSRDSDBParameterGroup = goresources.AWSRDSDBParameterGroup
type AWSIAMPolicy = goresources.AWSIAMPolicy
type AWSIAMRole = goresources.AWSIAMRole
type AWSSecretsManagerSecret = goresources.AWSSecretsManagerSecret
type AWSSecretsManagerSecretTargetAttachment = goresources.AWSSecretsManagerSecretTargetAttachment
type GenerateSecretString = goresources.AWSSecretsManagerSecret_GenerateSecretString

var NewTemplate = goformation.NewTemplate
var Join = goformation.Join
var GetAtt = goformation.GetAtt
var Ref = goformation.Ref

const CreateInProgress = cloudformation.StackStatusCreateInProgress
const DeleteInProgress = cloudformation.StackStatusDeleteInProgress
const UpdateInProgress = cloudformation.StackStatusUpdateInProgress
const CreateComplete = cloudformation.StackStatusCreateComplete
const DeleteComplete = cloudformation.StackStatusDeleteComplete
const UpdateComplete = cloudformation.StackStatusUpdateComplete
const CreateFailed = cloudformation.StackStatusCreateFailed
const DeleteFailed = cloudformation.StackStatusDeleteFailed
const RollbackComplete = cloudformation.StackStatusRollbackComplete
const UpdateRollbackComplete = cloudformation.StackStatusUpdateRollbackComplete

var (
	// capabilities required by cloudformation
	capabilities = []*string{
		aws.String("CAPABILITY_NAMED_IAM"),
	}
	// ErrStackNotFound returned when stack does not exist, or has been deleted
	ErrStackNotFound = fmt.Errorf("STACK_NOT_FOUND")
	// NoUpdatesErrMatch is string to match in error from aws to detect if nothing to update
	NoUpdatesErrMatch = "No updates"
	// NoExistErrMatch is a string to match if stack does not exist
	NoExistErrMatch = "does not exist"
)

// Outputs is used as a more friendly version of cloudformation.Output
type Outputs map[string]string

// Client performs cloudformation operations on objects that implement the Stack interface
type Client struct {
	// ClusterName is used to prefix any generated names to avoid clashes
	ClusterName string
	// Client is the AWS SDK Client implementation to use
	Client sdk.Client
	// PollingInterval is the duration between calls to check state when waiting for apply/destroy to complete
	PollingInterval time.Duration
}

// Apply reconciles the state of the remote cloudformation stack and blocks
// until the stack is no longer in an creating/applying state or a ctx timeout is hit
// Calls should be retried if DeadlineExceeded errors are hit
// Returns any outputs on successful apply.
// Will update stack with current status
func (r *Client) Apply(ctx context.Context, stack Stack, params ...*Parameter) (Outputs, error) {
	// always update stack status
	defer r.updateStatus(stack)
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
	_, err = r.waitUntilCompleteState(ctx, stack)
	if err != nil {
		return nil, err
	}
	err = r.update(ctx, stack, params...)
	if err != nil {
		return nil, err
	}
	state, err := r.waitUntilCompleteState(ctx, stack)
	if err != nil {
		return nil, err
	}
	return r.resolveOutputs(ctx, state.Outputs)
}

// validateParams checks for any unset template parameters
func (r *Client) validateTemplateParams(t *Template, params []*Parameter) error {
	missing := map[string]interface{}{}
	// copy all wanted params into missing
	for k, v := range t.Parameters {
		missing[k] = v
	}
	// remove items from missing list as found
	for wantedKey := range t.Parameters {
		for _, param := range params {
			if param.ParameterKey == nil {
				continue
			}
			// phew found it
			if *param.ParameterKey == wantedKey {
				delete(missing, wantedKey)
			}
		}
	}
	// if any left, then we have an issue
	if len(missing) > 0 {
		keys := []string{}
		for k := range missing {
			keys = append(keys, k)
		}
		keysCSV := strings.Join(keys, ",")
		return fmt.Errorf("missing required input parameters: [%s]", keysCSV)
	}
	return nil
}

// create initiates a cloudformation create passing in the given params
func (r *Client) create(ctx context.Context, stack Stack, params ...*Parameter) error {
	// fetch and validate template
	t := stack.GetStackTemplate()
	err := r.validateTemplateParams(t, params)
	if err != nil {
		return err
	}
	yaml, err := t.YAML()
	if err != nil {
		return err
	}
	_, err = r.Client.CreateStackWithContext(ctx, &CreateStackInput{
		Capabilities: capabilities,
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stack.GetStackName()),
		Parameters:   params,
	})
	if err != nil {
		return err
	}
	return nil
}

// Update the stack and wait for update to complete.
func (r *Client) update(ctx context.Context, stack Stack, params ...*Parameter) error {
	// fetch and validate template params
	t := stack.GetStackTemplate()
	err := r.validateTemplateParams(t, params)
	if err != nil {
		return err
	}
	yaml, err := t.YAML()
	if err != nil {
		return err
	}
	_, err = r.Client.UpdateStackWithContext(ctx, &UpdateStackInput{
		Capabilities: capabilities,
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stack.GetStackName()),
		Parameters:   params,
	})
	if err != nil && !IsNoUpdateError(err) {
		return err
	}
	return nil
}

// Destroy will attempt to deprovision the cloudformation stack and block until complete or the ctx Deadline expires
// Calls should be retried if DeadlineExceeded errors are hit
// Will update stack with current status
func (r *Client) Destroy(ctx context.Context, stack Stack) error {
	// always update stack status
	defer r.updateStatus(stack)
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
	if *state.StackStatus == DeleteComplete {
		// resource already deleted
		return nil
	}
	// trigger a delete unless we're already in a deleting state
	if *state.StackStatus != DeleteInProgress {
		_, err := r.Client.DeleteStackWithContext(ctx, &DeleteStackInput{
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

// Outputs fetches the cloudformation outputs for the given stack
// Returns ErrStackNotFound if stack does not exist
func (r *Client) Outputs(ctx context.Context, stack Stack) (Outputs, error) {
	state, err := r.get(ctx, stack)
	if err != nil {
		return nil, err
	}
	return r.resolveOutputs(ctx, state.Outputs)
}

// resolveOutputs returns cloudformation outputs in a map format and resolves
// any values that are stored in AWS Secrets Manager
func (r *Client) resolveOutputs(ctx context.Context, list []*Output) (Outputs, error) {
	outputs := Outputs{}
	for _, item := range list {
		if item.OutputKey == nil || item.OutputValue == nil {
			continue
		}
		key := *item.OutputKey
		value := *item.OutputValue
		// we automatically resolve references to AWS Secrets Manager
		// secrets here, so that we are able to make use of encrypted
		// sensitive values in cloudformation templates
		if strings.HasPrefix(value, "{{resolve:secretsmanager:") {
			// extract ARN and key name from reference
			secretARNMatcher := regexp.MustCompile(`{{resolve:secretsmanager:(.*):SecretString:(.*)}}`)
			matches := secretARNMatcher.FindStringSubmatch(value)
			if len(matches) == 0 {
				return nil, fmt.Errorf("failed to extract ARN and key name from secretsmanager value: %s", secretARNMatcher)
			}
			arn := matches[1]
			subkey := matches[2]
			v, err := r.Client.GetSecretValueWithContext(ctx, &GetSecretValueInput{
				SecretId: aws.String(arn),
			})
			if err != nil {
				return nil, err
			}
			if v.SecretString == nil {
				return nil, fmt.Errorf("unexpected nil value in SecretString of %s", arn)
			}
			secrets := map[string]interface{}{}
			err = json.Unmarshal([]byte(*v.SecretString), &secrets)
			if err != nil {
				return nil, err
			}
			subval, haveSubkey := secrets[subkey]
			if !haveSubkey {
				return nil, fmt.Errorf("could not find subkey %s in SecretString of %s", subkey, arn)
			}
			subvalString, ok := subval.(string)
			if !ok {
				return nil, fmt.Errorf("subval at subkey %s in SecretString of %s is not a string", subkey, arn)
			}
			value = subvalString
		}
		outputs[key] = value
	}
	return outputs, nil
}

// get fetches the cloudformation stack state
// Returns ErrStackNotFound if stack does not exist
func (r *Client) get(ctx context.Context, stack Stack) (*State, error) {
	describeOutput, err := r.Client.DescribeStacksWithContext(ctx, &DescribeStacksInput{
		StackName: aws.String(stack.GetStackName()),
	})
	if err != nil {
		if IsNotFoundError(err) {
			return nil, ErrStackNotFound
		}
		return nil, err
	}
	if describeOutput == nil {
		return nil, fmt.Errorf("describeOutput was nil, potential issue with AWS Client")
	}
	if len(describeOutput.Stacks) == 0 {
		return nil, fmt.Errorf("describeOutput contained no Stacks, potential issue with AWS Client")
	}
	state := describeOutput.Stacks[0]
	if state.StackStatus == nil {
		return nil, fmt.Errorf("describeOutput contained a nil StackStatus, potential issue with AWS Client")
	}
	return state, nil
}

// update mutates the stack's status with current state, events and any
// whitelisted outputs. ignores any errors encountered and just updates
// whatever it can with the intension of getting as much info visible as
// possible even under error conditions.
func (r *Client) updateStatus(stack Stack) {
	// use a fresh context or we might not be able to update status after
	// deadline is hit, but this feels a little wrong
	ctx := context.Background()
	state, _ := r.get(ctx, stack)
	events, _ := r.events(ctx, stack)
	s := stack.GetStatus()
	// update aws specific state
	if state != nil {
		if state.StackId != nil {
			s.AWS.ID = *state.StackId
		}
		if state.StackName != nil {
			s.AWS.Name = *state.StackName
		}
		if state.StackStatus != nil {
			s.AWS.Status = *state.StackStatus
		}
		if state.StackStatusReason != nil {
			s.AWS.Reason = *state.StackStatusReason
		}
	}
	// add any event details
	if events != nil {
		s.AWS.Events = []object.AWSEvent{}
		for _, event := range events {
			reason := "-"
			if event.ResourceStatusReason != nil {
				reason = *event.ResourceStatusReason
			}
			s.AWS.Events = append(s.AWS.Events, object.AWSEvent{
				Status: *event.ResourceStatus,
				Reason: reason,
				Time:   &metav1.Time{Time: *event.Timestamp},
			})
		}
	}
	// update generic state
	switch s.AWS.Status {
	case DeleteFailed, CreateFailed, RollbackComplete, UpdateRollbackComplete:
		s.State = object.ErrorState
	case DeleteInProgress, DeleteComplete:
		s.State = object.DeletingState
	case CreateComplete, UpdateComplete:
		s.State = object.ReadyState
	default:
		s.State = object.ReconcilingState
	}
	// if object implements whitelisting of output keys, then update info
	if w, ok := stack.(StackOutputWhitelister); ok {
		if s.AWS.Info == nil {
			s.AWS.Info = map[string]string{}
		}
		outputs, _ := r.Outputs(ctx, stack)
		for _, whitelistedKey := range w.GetStackOutputWhitelist() {
			s.AWS.Info[whitelistedKey] = outputs[whitelistedKey]
		}
	}
	stack.SetStatus(s)
}

func (r *Client) events(ctx context.Context, stack Stack) ([]*StateEvent, error) {
	eventsOutput, err := r.Client.DescribeStackEventsWithContext(ctx, &DescribeStackEventsInput{
		StackName: aws.String(stack.GetStackName()),
	})
	if err != nil {
		return nil, err
	}
	if eventsOutput == nil {
		return []*StateEvent{}, nil
	}
	return eventsOutput.StackEvents, nil
}

// Exists checks if the stack has been provisioned
func (r *Client) exists(ctx context.Context, stack Stack) (bool, error) {
	_, err := r.get(ctx, stack)
	if err == ErrStackNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Client) waitUntilCompleteState(ctx context.Context, stack Stack) (*State, error) {
	return r.waitUntilState(ctx, stack, []string{
		CreateComplete,
		UpdateComplete,
		UpdateRollbackComplete,
		RollbackComplete,
	})
}

func (r *Client) waitUntilDestroyedState(ctx context.Context, stack Stack) (*State, error) {
	return r.waitUntilState(ctx, stack, []string{
		DeleteComplete,
	})
}

func (r *Client) waitUntilState(ctx context.Context, stack Stack, desiredStates []string) (*State, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		default:
			state, err := r.get(ctx, stack)
			if IsNotFoundError(err) && in(DeleteComplete, desiredStates) {
				// If we are waiting for DeleteComplete state and the
				// stack has gone missing, consider this DeleteComplete
				return &State{}, nil
			} else if err != nil {
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

func IsNotFoundError(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "ResourceNotFoundException" {
			return true
		} else if awsErr.Code() == "ValidationError" && strings.Contains(awsErr.Message(), NoExistErrMatch) {
			return true
		}
	}
	return false
}
