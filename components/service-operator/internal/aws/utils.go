package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

func ValueFromOutputs(key string, outputs []*cloudformation.Output) []byte {
	for _, output := range outputs {
		if output.OutputKey != nil && *output.OutputKey == key {
			return []byte(*output.OutputValue)
		}
	}
	return nil
}

func DefineTags(clusterName, resourceName, namespace, resourceType string) []resources.Tag {
	return []resources.Tag{
		{
			Key:   "Cluster",
			Value: clusterName,
		},
		{
			Key:   "Name",
			Value: resourceName,
		},
		{
			Key:   "Service",
			Value: resourceType,
		},
		{
			Key:   "Namespace",
			Value: namespace,
		},
		{
			Key:   "Environment",
			Value: namespace,
		},
	}
}

type PolicyDocument struct {
	Version   string
	Statement []PolicyStatement
}

type PolicyStatement struct {
	Effect   string
	Action   []string
	Resource []string
}

type AssumeRolePolicyDocument struct {
	Version   string
	Statement []AssumeRolePolicyStatement
}

type AssumeRolePolicyStatement struct {
	Effect    string
	Principal PolicyPrincipal
	Action    []string
}

type PolicyPrincipal struct {
	AWS []string
}

func NewRolePolicyDocument(resources, actions []string) PolicyDocument {
	return PolicyDocument{
		Version: "2012-10-17",
		Statement: []PolicyStatement{
			{
				Effect:   "Allow",
				Action:   actions,
				Resource: resources,
			},
		},
	}
}
func NewAssumeRolePolicyDocument(principal string) AssumeRolePolicyDocument {
	return AssumeRolePolicyDocument{
		Version: "2012-10-17",
		Statement: []AssumeRolePolicyStatement{
			{
				Effect: "Allow",
				Principal: PolicyPrincipal{
					AWS: []string{principal},
				},
				Action: []string{"sts:AssumeRole"},
			},
		},
	}
}
