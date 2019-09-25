package v1beta1_test

import (
	"fmt"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
)

var _ = Describe("Postgres", func() {

	var postgres v1beta1.Postgres
	var tags []cloudformation.Tag

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		postgres = v1beta1.Postgres{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.PostgresSpec{},
		}
		tags = []cloudformation.Tag{
			{Key: "Cluster", Value: env.ClusterName()},
			{Key: "Service", Value: "postgres"},
			{Key: "Name", Value: "example"},
			{Key: "Namespace", Value: "default"},
			{Key: "Environment", Value: "default"},
		}
	})

	It("should default secret name to object name", func() {
		Expect(postgres.GetSecretName()).To(Equal("example"))
	})

	It("should use secret name from spec.Secret if set ", func() {
		postgres.Spec.Secret = "my-target-secret"
		Expect(postgres.GetSecretName()).To(Equal("my-target-secret"))
	})

	It("should base egress whitelisted host name off object name", func() {
		outputs := cloudformation.Outputs {
			v1beta1.PostgresEndpoint: "test-endpoint",
			v1beta1.PostgresReadEndpoint: "test-read-endpoint",
			v1beta1.PostgresPort: "3306",
		}

		ret, err := postgres.GetServiceEntry(outputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret.GetObjectMeta().Name).To(Equal(fmt.Sprintf("svcop-postgres-%s", postgres.GetName())))
		Expect(ret.GetObjectMeta().Namespace).To(Equal(postgres.GetNamespace()))
		Expect(ret.GetSpec()["resolution"]).To(Equal("DNS"))
		Expect(ret.GetSpec()["location"]).To(Equal("MESH_EXTERNAL"))
		Expect(ret.GetSpec()["hosts"]).To(ContainElement(outputs[v1beta1.PostgresEndpoint]))
		Expect(ret.GetSpec()["hosts"]).To(ContainElement(outputs[v1beta1.PostgresReadEndpoint]))
		ports, ok := ret.GetSpec()["ports"].([]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(ports)).To(BeNumerically(">", 0))
		port, ok := ports[0].(map[string]interface{})
		Expect(port["name"]).To(Equal("aurora"))
		portnum, err := strconv.Atoi(outputs[v1beta1.PostgresPort])
		Expect(err).NotTo(HaveOccurred())
		Expect(port["number"]).To(Equal(portnum))
		Expect(port["protocol"]).To(Equal("TLS"))
	})

	It("should error if port is not numeric", func() {
		outputs := cloudformation.Outputs {
			v1beta1.PostgresEndpoint: "test-endpoint",
			v1beta1.PostgresReadEndpoint: "test-read-endpoint",
			v1beta1.PostgresPort: "asd",
		}
		_, err := postgres.GetServiceEntry(outputs)
		Expect(err).To(HaveOccurred())
	})

	It("should generate a unique stack name prefixed with cluster name", func() {
		Expect(postgres.GetStackName()).To(HavePrefix("xxx-postgres-default-example"))
	})

	Context("cloudformation", func() {

		It("should have inputs for vpc config", func() {
			t := postgres.GetStackTemplate()
			Expect(t.Parameters).To(HaveKey("DBSubnetGroup"))
			Expect(t.Parameters).To(HaveKey("VPCSecurityGroupID"))
		})

		It("should have outputs for connection details", func() {
			t := postgres.GetStackTemplate()
			Expect(t.Outputs).To(HaveKey("Endpoint"))
			Expect(t.Outputs).To(HaveKey("ReadEndpoint"))
			Expect(t.Outputs).To(HaveKey("Port"))
			Expect(t.Outputs).To(HaveKey("Username"))
			Expect(t.Outputs).To(HaveKey("Password"))
		})

		Context("cluster resource", func() {

			var cluster *cloudformation.AWSRDSDBCluster

			JustBeforeEach(func() {
				t := postgres.GetStackTemplate()
				Expect(t.Resources[v1beta1.PostgresResourceCluster]).To(BeAssignableToTypeOf(&cloudformation.AWSRDSDBCluster{}))
				cluster = t.Resources[v1beta1.PostgresResourceCluster].(*cloudformation.AWSRDSDBCluster)
			})

			It("should have an RDS cluster resource with sensible defaults", func() {
				Expect(cluster.Engine).To(Equal("aurora-postgresql"))
				Expect(cluster.DBClusterParameterGroupName).ToNot(BeEmpty())
				Expect(cluster.VpcSecurityGroupIds).ToNot(BeNil())
				Expect(cluster.MasterUsername).ToNot(BeEmpty())
				Expect(cluster.MasterUserPassword).ToNot(BeEmpty())
				Expect(cluster.Tags).To(Equal(tags))
			})

		})

		Context("instance resources", func() {

			var instances []*cloudformation.AWSRDSDBInstance

			JustBeforeEach(func() {
				t := postgres.GetStackTemplate()
				instances = []*cloudformation.AWSRDSDBInstance{}
				for _, r := range t.Resources {
					inst, ok := r.(*cloudformation.AWSRDSDBInstance)
					if !ok {
						continue
					}
					instances = append(instances, inst)
				}
			})

			It("should default to 2 instances", func() {
				Expect(instances).To(HaveLen(2))
			})

			It("should have RDS instance resources with sensible defaults", func() {
				for _, instance := range instances {
					Expect(instance.PubliclyAccessible).To(BeFalse())
					Expect(instance.DBInstanceClass).To(Equal("db.r5.large"))
					Expect(instance.Engine).To(Equal("aurora-postgresql"))
					Expect(instance.Tags).To(Equal(tags))
				}
			})

			Context("when spec.aws.instanceCount is set", func() {
				BeforeEach(func() {
					postgres.Spec.AWS.InstanceCount = 3
				})
				It("should set number of instances from spec", func() {
					Expect(instances).To(HaveLen(3))
				})
			})

			Context("when spec.aws.instanceType is set", func() {
				BeforeEach(func() {
					postgres.Spec.AWS.InstanceType = "db.t3.medium"
				})
				It("should set instances from spec", func() {
					for _, instance := range instances {
						Expect(instance.DBInstanceClass).To(Equal("db.t3.medium"))
					}
				})
			})

		})

	})

})
