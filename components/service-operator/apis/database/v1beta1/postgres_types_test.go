package v1beta1_test

import (
	"os"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws"
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

	It("should default service entry name to object name", func() {
		Expect(postgres.GetServiceEntryName()).To(Equal(postgres.GetName()))
	})

	It("should use service entry name from spec.ServiceEntry if set", func() {
		postgres.Spec.ServiceEntry = "my-target-service-entry"
		Expect(postgres.GetServiceEntryName()).To(Equal("my-target-service-entry"))
	})

	It("should produce the correct service entry", func() {
		outputs := cloudformation.Outputs{
			v1beta1.PostgresEndpoint:     "test-endpoint",
			v1beta1.PostgresReadEndpoint: "test-read-endpoint",
			v1beta1.PostgresPort:         "3306",
		}
		portnum, err := strconv.Atoi(outputs[v1beta1.PostgresPort])
		Expect(err).NotTo(HaveOccurred())

		specs, err := postgres.GetServiceEntrySpecs(outputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(specs).To(HaveLen(2))
		Expect(specs).To(ConsistOf(
			And(
				HaveKeyWithValue("resolution", "DNS"),
				HaveKeyWithValue("location", "MESH_EXTERNAL"),
				HaveKeyWithValue("hosts", ContainElement(outputs[v1beta1.PostgresEndpoint])),
				HaveKeyWithValue("ports", ContainElement(
					map[string]interface{}{
						"name":     "aurora",
						"number":   portnum,
						"protocol": "TCP",
					},
				)),
				HaveKeyWithValue("exportTo", And(
					HaveLen(1),
					ContainElement("."),
				)),
			),
			And(
				HaveKeyWithValue("resolution", "DNS"),
				HaveKeyWithValue("location", "MESH_EXTERNAL"),
				HaveKeyWithValue("hosts", ContainElement(outputs[v1beta1.PostgresReadEndpoint])),
				HaveKeyWithValue("ports", ContainElement(
					map[string]interface{}{
						"name":     "aurora",
						"number":   portnum,
						"protocol": "TCP",
					},
				)),
				HaveKeyWithValue("exportTo", And(
					HaveLen(1),
					ContainElement("."),
				)),
			),
		))

	})

	It("should error if port is not numeric", func() {
		outputs := cloudformation.Outputs{
			v1beta1.PostgresEndpoint:     "test-endpoint",
			v1beta1.PostgresReadEndpoint: "test-read-endpoint",
			v1beta1.PostgresPort:         "asd",
		}
		_, err := postgres.GetServiceEntrySpecs(outputs)
		Expect(err).To(HaveOccurred())
	})

	It("should generate a unique stack name prefixed with cluster name", func() {
		Expect(postgres.GetStackName()).To(HavePrefix("xxx-postgres-default-example"))
	})

	It("should have a sensible stack policy", func() {
		expectedStackPolicyDocument := aws.StackPolicyDocument{
			Statement: []aws.StatementEntry{
				{
					Effect:    "Deny",
					Action:    []string{"Update:Replace", "Update:Delete"},
					Principal: "*",
					Resource:  "LogicalResourceId/RDSCluster",
				},
				{
					Effect:    "Allow",
					Action:    []string{"Update:Modify"},
					Principal: "*",
					Resource:  "LogicalResourceId/RDSCluster",
				},
				{
					Effect:    "Deny",
					Action:    []string{"Update:Replace", "Update:Delete"},
					Principal: "*",
					Resource:  "LogicalResourceId/RDSDBInstance0",
				},
				{
					Effect:    "Allow",
					Action:    []string{"Update:Modify"},
					Principal: "*",
					Resource:  "LogicalResourceId/RDSDBInstance0",
				},
				{
					Effect:    "Deny",
					Action:    []string{"Update:Replace", "Update:Delete"},
					Principal: "*",
					Resource:  "LogicalResourceId/RDSDBInstance1",
				},
				{
					Effect:    "Allow",
					Action:    []string{"Update:Modify"},
					Principal: "*",
					Resource:  "LogicalResourceId/RDSDBInstance1",
				},
			},
		}

		actualStackPolicyDocument := postgres.GetStackPolicy()
		Expect(actualStackPolicyDocument).To(Equal(expectedStackPolicyDocument))
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
				Expect(cluster.BackupRetentionPeriod).To(Equal(7))
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
					Expect(instance.DeleteAutomatedBackups).To(Equal(false))
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
