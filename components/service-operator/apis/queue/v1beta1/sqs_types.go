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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AWS allows specifying configuration for the Postgres RDS instance
type AWS struct {
	// QueueName will define given queue on the provider.
	QueueName string `json:"queueName"`
}

// SQSSpec defines the desired state of SQS
type SQSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AWS specific subsection of the resource.
	AWS AWS `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
}

// SQSStatus defines the observed state of SQS
type SQSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ID of an instance for a reference.
	ID string `json:"id"`
	// Status of the currently running instance.
	Status string `json:"status"`
	// Reason for the current status of the instance.
	Reason string `json:"reason,omitempty"`
}

// +kubebuilder:object:root=true

// SQS is the Schema for the SQS API
type SQS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SQSSpec   `json:"spec,omitempty"`
	Status SQSStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SQSList contains a list of SQS
type SQSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SQS{}, &SQSList{})
}
