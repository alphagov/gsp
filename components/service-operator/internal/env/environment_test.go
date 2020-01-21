package env_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/gsp/components/service-operator/internal/env"
)

var _ = Describe("Environment", func() {

	var knownValue = "value-set-in-test"

	Context("MustGet", func() {
		It("should read value from environment", func() {
			os.Setenv("TEST_EXAMPLE_VAR", knownValue)
			Expect(env.MustGet("TEST_EXAMPLE_VAR")).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("TEST_EXAMPLE_VAR")
			Expect(func() { env.MustGet("TEST_EXAMPLE_VAR") }).To(Panic())
		})
	})

	Context("ClusterName", func() {
		It("should read value from environment", func() {
			os.Setenv("CLUSTER_NAME", knownValue)
			Expect(env.ClusterName()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("CLUSTER_NAME")
			Expect(func() { env.ClusterName() }).To(Panic())
		})
	})

	Context("AWSRDSSecurityGroupID", func() {
		It("should read value from environment", func() {
			os.Setenv("AWS_RDS_SECURITY_GROUP_ID", knownValue)
			Expect(env.AWSRDSSecurityGroupID()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("AWS_RDS_SECURITY_GROUP_ID")
			Expect(func() { env.AWSRDSSecurityGroupID() }).To(Panic())
		})
	})

	Context("AWSRDSSubnetGroupName", func() {
		It("should read value from environment", func() {
			os.Setenv("AWS_RDS_SUBNET_GROUP_NAME", knownValue)
			Expect(env.AWSRDSSubnetGroupName()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("AWS_RDS_SUBNET_GROUP_NAME")
			Expect(func() { env.AWSRDSSubnetGroupName() }).To(Panic())
		})
	})

	Context("AWSPrincipalServerRoleARN", func() {
		It("should read value from environment", func() {
			os.Setenv("AWS_PRINCIPAL_SERVER_ROLE_ARN", knownValue)
			Expect(env.AWSPrincipalServerRoleARN()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("AWS_PRINCIPAL_SERVER_ROLE_ARN")
			Expect(func() { env.AWSPrincipalServerRoleARN() }).To(Panic())
		})
	})

	Context("AWSOIDCProviderURL", func() {
		It("should read value from environment", func() {
			os.Setenv("AWS_OIDC_PROVIDER_URL", knownValue)
			Expect(env.AWSOIDCProviderURL()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("AWS_OIDC_PROVIDER_URL")
			Expect(func() { env.AWSOIDCProviderURL() }).To(Panic())
		})
	})

	Context("AWSOIDCProviderARN", func() {
		It("should read value from environment", func() {
			os.Setenv("AWS_OIDC_PROVIDER_ARN", knownValue)
			Expect(env.AWSOIDCProviderARN()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("AWS_OIDC_PROVIDER_ARN")
			Expect(func() { env.AWSOIDCProviderARN() }).To(Panic())
		})
	})

	Context("AWSPrincipalPermissionsBoundaryARN", func() {
		It("should read value from environment", func() {
			os.Setenv("AWS_PRINCIPAL_PERMISSIONS_BOUNDARY_ARN", knownValue)
			Expect(env.AWSPrincipalPermissionsBoundaryARN()).To(Equal(knownValue))
		})
		It("should panic if not set", func() {
			os.Unsetenv("AWS_PRINCIPAL_PERMISSIONS_BOUNDARY_ARN")
			Expect(func() { env.AWSPrincipalPermissionsBoundaryARN() }).To(Panic())
		})
	})

	Context("AWSIntegrationTestEnabled", func() {
		It("should return true if environment set to true", func() {
			os.Setenv("AWS_INTEGRATION", "true")
			Expect(env.AWSIntegrationTestEnabled()).To(BeTrue())
		})
		It("should return false if environment not set", func() {
			os.Setenv("AWS_INTEGRATION", "")
			Expect(env.AWSIntegrationTestEnabled()).To(BeFalse())
		})
	})

})
