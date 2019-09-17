package controllers

import (
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	storage "github.com/alphagov/gsp/components/service-operator/apis/storage/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
)

func S3CloudFormationController(c sdk.Client) Controller {
	return &cloudformation.Controller{
		Kind:              &storage.S3Bucket{},
		PrincipalListKind: &access.PrincipalList{},
		CloudFormationClient: &cloudformation.Client{
			Client: c,
		},
	}
}
