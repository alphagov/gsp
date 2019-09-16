package controllers

import (
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	database "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/aws/aws-sdk-go/aws"
)

// PostgresCloudFormationController creates a Controller instance for provision
// Postgres with cloudformation.
func PostgresCloudFormationController(c sdk.Client) Controller {
	return &cloudformation.Controller{
		Kind:              &database.Postgres{},
		PrincipalListKind: &access.PrincipalList{},
		CloudFormationClient: &cloudformation.Client{
			Client: c,
		},
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:   aws.String(database.VPCSecurityGroupIDParameterName),
				ParameterValue: aws.String(env.AWSRDSSecurityGroupID()),
			},
			{
				ParameterKey:   aws.String(database.DBSubnetGroupNameParameterName),
				ParameterValue: aws.String(env.AWSRDSSubnetGroupName()),
			},
		},
	}
}
