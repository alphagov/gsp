package cloudformation

import (
	"context"
	"github.com/alphagov/gsp/components/service-operator/internal/aws"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Stack

// Stack represents a type that can marshall itself to a cloudformation config for use with Client
type Stack interface {
	GetStackName() string
	GetStackTemplate() (*Template, error)
	object.Service
}

// StackPolicyAttacher adds an additional method to Stack for injecting a
// target role name into stack parameters that will be the receiver of policies
// provisioned by this stack. Role parameters are added at both Update and
// Create time. The role provided will already exist and may contain other
// existing policies which should not be altered.
type StackPolicyAttacher interface {
	Stack
	GetStackRoleParameters(role string) ([]*Parameter, error)
}

// StackOutputWhitelister allows a type to return a list of Output keys whose values are safe to be
// displayed in the Status of the resource. This is useful for when you want to display some outputs
// without having them embedded into a Secret
type StackOutputWhitelister interface {
	GetStackOutputWhitelist() []string
}

// StackSecretOutputter allows a type to return the name of a kubernetes Secret
// that will be populated with any cloudformation outputs. This is useful when
// the cloudformation stack returns sensitive information that must be consumed
// as configuration, for example a username, password and connection string for
// a database.
type StackSecretOutputter interface {
	Stack
	object.SecretNamer
}

// ServiceEntryCreator names a ServiceEntry to hold what are essentially egress
// firewall rules.
type ServiceEntryCreator interface {
	Stack
	GetServiceEntryName() string
	GetServiceEntrySpecs(outputs Outputs) ([]map[string]interface{}, error)
}

// ServiceAccountCreator names a ServiceAccount to be used to access created resources
type ServiceAccountCreator interface {
	Stack
	GetServiceAccountName() string
	GetServiceAccountAnnotations(outputs Outputs) map[string]string
}

// StackPolicyProvider allows a type to return a string representation of a
// CloudFormationStack policy to prevent unauthorised operations on a given
// stack.
type StackPolicyProvider interface {
	GetStackPolicy() aws.StackPolicyDocument
}

type StackSecretContributor interface {
	Stack
	GetTemplateSecrets(ctx context.Context, client sdk.Client, outputs Outputs) (map[string]string, error)
}
