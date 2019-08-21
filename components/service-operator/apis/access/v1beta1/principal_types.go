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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
const AccessGroupLabel = "group.access.govsvc.uk"

// Event is a single action taken against the resource at any given time.
type Event struct {
	// Status of the currently running instance.
	Status string `json:"status"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
	// Time of the event cast.
	Time *metav1.Time `json:"time"`
}

// PrincipalStatus defines the observed state of Principal
type PrincipalStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// ID of an instance for a reference.
	ID string `json:"id"`
	// Status of the currently running instance.
	Status string `json:"status"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
	// Events will hold more in-depth details of the current state of the instance.
	Events []Event `json:"events,omitempty"`
	// ARN of the IAM Principal
	ARN string `json:"arn"`
}

// +kubebuilder:object:root=true

// Principal is the Schema for the Principal API
type Principal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status PrincipalStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PrincipalList contains a list of Principal
type PrincipalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Principal `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Principal{}, &PrincipalList{})
}
