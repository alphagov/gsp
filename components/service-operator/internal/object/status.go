// +kubebuilder:object:generate=true
package object

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type State string

var (
	ErrorState       State = "ERROR"
	DeletingState    State = "DELETE_IN_PROGRESS"
	ReadyState       State = "READY"
	ReconcilingState State = "RECONCILE_IN_PROGRESS"
)

// +kubebuilder:object:root=false

// Status is the type shared by most service resources
type Status struct {
	// Generic service state
	State State `json:"state,omitempty"`
	// AWS specific status
	AWS AWSStatus `json:"aws,omitempty"`
}

var _ StatusReader = &Status{}
var _ StatusWriter = &Status{}

// GetStatus returns the status
func (s *Status) GetStatus() Status {
	if s == nil {
		return Status{}
	}
	return *s
}

// SetStatus updates status fields
func (s *Status) SetStatus(status Status) {
	if s == nil {
		return
	}
	s.State = status.State
	s.AWS = status.AWS
}

// GetState returns the status
func (s *Status) GetState() State {
	if s == nil {
		return State("")
	}
	return s.State
}

// SetState returns the status
func (s *Status) SetState(state State) {
	if s == nil {
		return
	}
	s.State = state
}

// +kubebuilder:object:root=false

// AWSStatus a cloudformation Stack
type AWSStatus struct {
	// ID of an instance for a reference.
	ID string `json:"id,omitempty"`
	// Name of an instance for a reference.
	Name string `json:"name,omitempty"`
	// Status of the currently running instance.
	Status string `json:"status,omitempty"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
	// Events will hold more in-depth details of the current state of the instance.
	Events []AWSEvent `json:"events,omitempty"`
	// Info shows any outputs returned from GetStackOutputWhitelist
	Info map[string]string `json:"info,omitempty"`
}

// +kubebuilder:object:root=false

// AWSEvent is a single action taken against the resource at any given time.
type AWSEvent struct {
	// Status of the currently running instance.
	Status string `json:"status,omitempty"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
	// Time of the event cast.
	Time *metav1.Time `json:"time,omitempty"`
}
