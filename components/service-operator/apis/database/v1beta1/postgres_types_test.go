package v1beta1_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
)

var _ = Describe("Postgres", func() {

	var o v1beta1.Postgres
	var tags []cloudformation.Tag

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		o = v1beta1.Postgres{
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
		Expect(o.GetSecretName()).To(Equal("example"))
	})

	It("should use secret name from spec.Secret if set ", func() {
		o.Spec.Secret = "my-target-secret"
		Expect(o.GetSecretName()).To(Equal("my-target-secret"))
	})

	It("should generate a unique stack name prefixed with cluster name", func() {
		Expect(o.GetStackName()).To(HavePrefix("xxx-postgres-default-example"))
	})

	It("should have inputs for vpc config", func() {
		t := o.GetStackTemplate()
		Expect(t.Parameters).To(HaveKey("DBSubnetGroup"))
		Expect(t.Parameters).To(HaveKey("VPCSecurityGroupID"))
	})

	It("should have outputs for connection details", func() {
		t := o.GetStackTemplate()
		Expect(t.Outputs).To(HaveKey("Endpoint"))
		Expect(t.Outputs).To(HaveKey("ReadEndpoint"))
		Expect(t.Outputs).To(HaveKey("Port"))
		Expect(t.Outputs).To(HaveKey("Username"))
		Expect(t.Outputs).To(HaveKey("Password"))
	})

	It("should have an RDS cluster resource with sensible defaults", func() {
		t := o.GetStackTemplate()
		Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSRDSDBCluster{})))
		cluster, ok := t.Resources[v1beta1.PostgresResourceCluster].(*cloudformation.AWSRDSDBCluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.Engine).To(Equal("aurora-postgresql"))
		Expect(cluster.DBClusterParameterGroupName).ToNot(BeEmpty())
		Expect(cluster.VpcSecurityGroupIds).ToNot(BeNil())
		Expect(cluster.MasterUsername).ToNot(BeEmpty())
		Expect(cluster.MasterUserPassword).ToNot(BeEmpty())
		Expect(cluster.Tags).To(Equal(tags))
	})

	It("should have RDS instance resources with sensible defaults", func() {
		t := o.GetStackTemplate()
		Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSRDSDBInstance{})))
		count := 0
		for _, r := range t.Resources {
			instance, ok := r.(*cloudformation.AWSRDSDBInstance)
			if !ok {
				continue
			}
			count++
			Expect(instance.PubliclyAccessible).To(BeFalse())
			Expect(instance.DBInstanceClass).To(Equal("db.r5.large"))
			Expect(instance.Engine).To(Equal("aurora-postgresql"))
			Expect(instance.Tags).To(Equal(tags))
		}
		Expect(count).To(BeNumerically(">", 0))
		Expect(count).To(Equal(v1beta1.DefaultInstanceCount))
	})

	It("should get number of instances from spec.AWS.InstanceCount", func() {
		o.Spec.AWS.InstanceCount = 3
		t := o.GetStackTemplate()
		count := 0
		for _, r := range t.Resources {
			_, ok := r.(*cloudformation.AWSRDSDBInstance)
			if !ok {
				continue
			}
			count++
		}
		Expect(count).To(Equal(o.Spec.AWS.InstanceCount))
	})

	It("should get instance size from spec.AWS.InstanceType", func() {
		o.Spec.AWS.InstanceType = "db.t3.medium"
		t := o.GetStackTemplate()
		count := 0
		for _, r := range t.Resources {
			instance, ok := r.(*cloudformation.AWSRDSDBInstance)
			if !ok {
				continue
			}
			count++
			Expect(instance.DBInstanceClass).To(Equal(o.Spec.AWS.InstanceType))
		}
		Expect(count).To(BeNumerically(">", 0))
	})

})
