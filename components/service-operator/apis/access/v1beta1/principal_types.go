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
	"fmt"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ensure implements required interfaces
var _ cloudformation.Stack = &Principal{}
var _ cloudformation.StackOutputWhitelister = &Principal{}
var _ object.Principal = &Principal{}

func init() {
	SchemeBuilder.Register(&Principal{}, &PrincipalList{})
}

const (
	IAMRoleResourceName                 = "IAMRole"
	IAMRoleName                         = "IAMRoleName"
	IAMRolePrincipalParameterName       = "IAMRolePrincipal"
	IAMPermissionsBoundaryParameterName = "IAMPermissionsBoundary"
)

// ensure implements StackObject
// var _ apis.StackObject = &Principal{}

// +kubebuilder:object:root=true

// Principal is the Schema for the Principal API
type Principal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	object.Status `json:"status,omitempty"`
}

// GetStackName generates a unique name for the stack
func (s *Principal) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "principal", s.GetNamespace(), s.GetName())
}

// GetRoleName returns a generated unique name suitable for use as a role name
func (s *Principal) GetRoleName() string {
	return fmt.Sprintf("svcop-%s-%s-%s", env.ClusterName(), s.GetNamespace(), s.GetName())
}

// GetStackTemplate returns cloudformation to create an IAM role
func (s *Principal) GetStackTemplate() *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Parameters[IAMRolePrincipalParameterName] = map[string]string{
		"Type": "String",
	}
	template.Parameters[IAMPermissionsBoundaryParameterName] = map[string]string{
		"Type": "String",
	}

	template.Resources[IAMRoleResourceName] = &cloudformation.AWSIAMRole{
		RoleName:                 s.GetRoleName(),
		AssumeRolePolicyDocument: cloudformation.NewAssumeRolePolicyDocument(cloudformation.Ref(IAMRolePrincipalParameterName)),
		PermissionsBoundary:      cloudformation.Ref(IAMPermissionsBoundaryParameterName),
	}

	template.Outputs[IAMRoleName] = map[string]interface{}{
		"Description": "IAMRole ARN to be returned to the user.",
		"Value":       cloudformation.Ref(IAMRoleResourceName),
	}

	return template
}

// GetStackOutputWhitelist will whitelist any output keys from template that can be shown in resource Status
func (s *Principal) GetStackOutputWhitelist() []string {
	return []string{IAMRoleName}
}

var _ object.PrincipalLister = &PrincipalList{}

// +kubebuilder:object:root=true

// PrincipalList contains a list of Principal
type PrincipalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Principal `json:"items"`
}

// GetPrincipals implements object.PrincipalLister
func (p *PrincipalList) GetPrincipals() []object.Principal {
	ps := make([]object.Principal, len(p.Items))
	for i := range p.Items {
		ps[i] = &p.Items[i]
	}
	return ps
}
