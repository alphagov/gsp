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

// TeamSpec defines the desired state of Team
// +k8s:openapi-gen=true
type TeamSpec struct {
	Roles []RoleSpec `json:"roles"`
}

// GithubAuth allows configuring github users
// +k8s:openapi-gen=true
type GithubAuth struct {
	Users []string `json:"users,omitempty"`
	Teams []string `json:"teams,omitempty"`
}

// LocalAuth allows defining hardcoded local users
// +k8s:openapi-gen=true
type LocalAuth struct {
	Users []string `json:"users,omitempty"`
}

// RoleSpec maps concourse roles to various auth backends
// +k8s:openapi-gen=true
type RoleSpec struct {
	Name   string     `json:"name"`
	Github GithubAuth `json:"github,omitempty"`
	Local  LocalAuth  `json:"local,omitempty"`
}

// TeamStatus define the status
type TeamStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Team is the Schema for the teams API
// +k8s:openapi-gen=true
type Team struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamSpec   `json:"spec,omitempty"`
	Status TeamStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeamList contains a list of Team
type TeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Team `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Team{}, &TeamList{})
}
