package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

func ValueFromOutputs(key string, outputs []*cloudformation.Output) []byte {
	for _, output := range outputs {
		if output.OutputKey != nil && *output.OutputKey == key {
			return []byte([]byte(*output.OutputValue))
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
