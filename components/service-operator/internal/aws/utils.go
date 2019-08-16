package aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func ValueFromOutputs(key string, outputs []*cloudformation.Output) []byte {
	for _, output := range outputs {
		if output.OutputKey != nil && *output.OutputKey == key {
			return []byte(base64.StdEncoding.EncodeToString([]byte(*output.OutputValue)))
		}
	}
	return nil
}
