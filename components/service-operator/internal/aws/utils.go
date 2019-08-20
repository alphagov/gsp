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
	Effect    string
	Principal PolicyPrincipal
	Action    []string
	Resources []string
}

type PolicyPrincipal struct {
	AWS []string
}

func NewRolePolicyDocument(principal, resource string, actions []string) PolicyDocument {
	return PolicyDocument{
		Version: "2012-10-17",
		Statement: []PolicyStatement{
			PolicyStatement{
				Effect: "Allow",
				Principal: PolicyPrincipal{
					AWS: []string{principal},
				},
				Action:    actions,
				Resources: []string{resource},
			},
		},
	}
}
