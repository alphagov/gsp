package controllers

import (
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	cache "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/aws/aws-sdk-go/aws"
)

// RedisCloudFormationController creates a Controller instance for provision
// an ElastiCache ReplicationGroup with cloudformation.
func RedisCloudFormationController(c sdk.Client) Controller {
	return &cloudformation.Controller{
		Kind:              &cache.Redis{},
		PrincipalListKind: &access.PrincipalList{},
		CloudFormationClient: &cloudformation.Client{
			Client: c,
		},
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:   aws.String(cache.VPCSecurityGroupIDParameterName),
				ParameterValue: aws.String(env.AWSRedisSecurityGroupID()),
			},
			{
				ParameterKey:   aws.String(cache.CacheSubnetGroupParameterName),
				ParameterValue: aws.String(env.AWSRedisSubnetGroupName()),
			},
		},
	}
}
