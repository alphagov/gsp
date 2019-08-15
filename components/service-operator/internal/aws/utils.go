package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func ValueFromOutputs(key string, outputs []*cloudformation.Output) string {
	for _, output := range outputs {
		if output.OutputKey != nil && *output.OutputKey == key {
			return *output.OutputValue
		}
	}
	return ""
}
