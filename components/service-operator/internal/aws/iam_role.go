package aws

import (
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	IAMRoleResourceName = "IAMRole"

	IAMRoleName = "IAMRoleName"
)

type IAMRole struct {
	RoleConfig          *access.Principal
	RoleName            string
	RolePrincipal       string
	PermissionsBoundary string
}

func (s *IAMRole) Template(stackName string, tags []resources.Tag) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources[IAMRoleResourceName] = &resources.AWSIAMRole{
		RoleName:                 s.RoleName,
		AssumeRolePolicyDocument: NewAssumeRolePolicyDocument(s.RolePrincipal),
		PermissionsBoundary:      s.PermissionsBoundary,
	}

	template.Outputs[IAMRoleName] = map[string]interface{}{
		"Description": "IAMRole ARN to be returned to the user.",
		"Value":       cloudformation.Ref(IAMRoleResourceName),
	}

	return template
}

func (s *IAMRole) CreateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

func (s *IAMRole) UpdateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

func (p *IAMRole) ResourceType() string {
	return "principal"
}
