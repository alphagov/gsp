package k8sclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLifecycleHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KubernetesClientSuite")
}
