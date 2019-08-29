package apis

import (
	"fmt"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrorState       = "ERROR"
	DeletingState    = "DELETE_IN_PROGRESS"
	ReadyState       = "READY"
	ReconcilingState = "RECONCILE_IN_PROGRESS"
)

type ServiceStatus struct {
	// one of CREATING, RECONCILING, DELETEING or READY
	State string `json:"state"`
}

// SetFromStack updates this ServiceStatus based on cloudformation stack data
// If state is nil then this will have no effect
func (s *ServiceStatus) SetFromStack(state *awscloudformation.Stack, events []*awscloudformation.StackEvent) error {
	if state != nil {
		if state.StackStatus == nil {
			s.State = ErrorState
			return fmt.Errorf("StackStatus should not be nil")
		}
		switch *state.StackStatus {
		case awscloudformation.StackStatusDeleteFailed, awscloudformation.StackStatusCreateFailed:
			s.State = ErrorState
		case awscloudformation.StackStatusDeleteInProgress, awscloudformation.StackStatusDeleteComplete:
			s.State = DeletingState
		case awscloudformation.StackStatusCreateComplete, awscloudformation.StackStatusUpdateComplete, awscloudformation.StackStatusUpdateRollbackComplete, awscloudformation.StackStatusRollbackComplete:
			s.State = ReadyState
		default:
			s.State = ReconcilingState
		}
	}
	return nil
}

// AWSEvent is a single action taken against the resource at any given time.
type AWSEvent struct {
	// Status of the currently running instance.
	Status string `json:"status"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
	// Time of the event cast.
	Time *metav1.Time `json:"time"`
}

type AWSServiceStatus struct {
	// ID of an instance for a reference.
	ID string `json:"id"`
	// Name of an instance for a reference.
	Name string `json:"name"`
	// Status of the currently running instance.
	Status string `json:"status"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
	// Events will hold more in-depth details of the current state of the instance.
	Events []AWSEvent `json:"events,omitempty"`
}

// SetFromStack updates the status from cloudformation stack state
// If state is nil then this will have no effect
func (s *AWSServiceStatus) SetFromStack(state *awscloudformation.Stack, events []*awscloudformation.StackEvent) error {
	if state != nil {
		if state.StackId != nil {
			s.ID = *state.StackId
		}
		if state.StackName != nil {
			s.Name = *state.StackName
		}
		if state.StackStatus != nil {
			s.Status = *state.StackStatus
		}
		if state.StackStatusReason != nil {
			s.Reason = *state.StackStatusReason
		}
	}
	if events != nil {
		s.Events = []AWSEvent{}
		for _, event := range events {
			reason := "-"
			if event.ResourceStatusReason != nil {
				reason = *event.ResourceStatusReason
			}
			s.Events = append(s.Events, AWSEvent{
				Status: *event.ResourceStatus,
				Reason: reason,
				Time:   &metav1.Time{Time: *event.Timestamp},
			})
		}
	}
	return nil
}
