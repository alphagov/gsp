package aws

import (
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	IAMRoleResourceName = "IAMRole"

	IAMRoleARN = "IAMRoleARN"
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

	template.Outputs[IAMRoleARN] = map[string]interface{}{
		"Description": "IAMRole ARN to be returned to the user.",
		"Value":       cloudformation.GetAtt(IAMRoleResourceName, "Arn"),
	}

	return template
}

func (s *IAMRole) Parameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{}, nil
}

func (p *IAMRole) ResourceType() string {
	return "principal"
}
