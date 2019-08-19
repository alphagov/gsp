package aws_test

import (
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	It("Should return nil if key not found", func() {
		outputs := []*cloudformation.Output{
			&cloudformation.Output{
				OutputKey:   aws.String("test-key"),
				OutputValue: aws.String("test-value"),
			},
		}

		Expect(internalaws.ValueFromOutputs("test-key-2", outputs)).To(BeNil())
	})

	It("Should return base64 encoded byte array of value of output that matches key", func() {
		outputs := []*cloudformation.Output{
			&cloudformation.Output{
				OutputKey:   aws.String("test-key"),
				OutputValue: aws.String("test-value"),
			},
		}

		expected := []byte([]byte("test-value"))
		Expect(internalaws.ValueFromOutputs("test-key", outputs)).To(Equal(expected))
	})
})
