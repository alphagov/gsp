package aws_test

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	"github.com/aws/aws-sdk-go/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
  It("Should return empty string if key not found", func() {
    outputs := []*cloudformation.Output{
      &cloudformation.Output{
        OutputKey: aws.String("test-key"),
        OutputValue: aws.String("test-value"),
      },
    }

    Expect(internalaws.ValueFromOutputs("test-key-2", outputs)).To(Equal(""))
  })

  It("Should return value of output that matches key", func() {
    outputs := []*cloudformation.Output{
      &cloudformation.Output{
        OutputKey: aws.String("test-key"),
        OutputValue: aws.String("test-value"),
      },
    }

    Expect(internalaws.ValueFromOutputs("test-key", outputs)).To(Equal("test-value"))
  })
})
