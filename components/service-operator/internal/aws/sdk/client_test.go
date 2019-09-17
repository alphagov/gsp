package sdk_test

import (
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Client", func() {

	It("should return a valid aws client", func() {
		var _ sdk.Client = sdk.NewClient()
	})

})
