package sdk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSDK(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aws Suite")
}
