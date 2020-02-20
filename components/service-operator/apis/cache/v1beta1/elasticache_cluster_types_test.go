package v1beta1_test

import (
	"encoding/base64"
	"os"

	"github.com/alphagov/gsp/components/service-operator/apis/cache/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ElasticacheCluster", func() {

	var o v1beta1.ElasticacheCluster

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		o = v1beta1.ElasticacheCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.ElasticacheClusterSpec{},
		}
	})

	It("should default secret name to object name", func() {
		Expect(o.GetSecretName()).To(Equal("example"))
	})

	It("should use secret name from spec.Secret if set ", func() {
		o.Spec.Secret = "my-target-secret"
		Expect(o.GetSecretName()).To(Equal("my-target-secret"))
	})


	It("should default service entry name to object name", func() {
		Expect(o.GetServiceEntryName()).To(Equal(o.GetName()))
	})

	It("should use service entry name from spec.ServiceEntry if set", func() {
		o.Spec.ServiceEntry = "my-target-service-entry"
		Expect(o.GetServiceEntryName()).To(Equal("my-target-service-entry"))
	})

	It("should produce the correct service entry", func() {
		outputs := cloudformation.Outputs{
			v1beta1.ElasticacheClusterRedisHostnameOutputName: "test-endpoint.local.govsandbox.uk",
			v1beta1.ElasticacheClusterRedisPortOutputName:     "6379",
		}
		portnum, err := strconv.Atoi(outputs[v1beta1.ElasticacheClusterRedisPortOutputName])
		Expect(err).NotTo(HaveOccurred())

		specs, err := o.GetServiceEntrySpecs(outputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(specs).To(HaveLen(1))
		Expect(specs).To(ConsistOf(
			And(
				HaveKeyWithValue("resolution", "STATIC"),
				HaveKeyWithValue("location", "MESH_EXTERNAL"),
				HaveKeyWithValue("addresses", ContainElement("127.0.0.1")),
				HaveKeyWithValue("endpoints", ContainElement(
					map[string]interface{}{
						"address": "127.0.0.1",
					},
				)),
				HaveKeyWithValue("hosts", ContainElement(outputs[v1beta1.ElasticacheClusterRedisHostnameOutputName])),
				HaveKeyWithValue("ports", ContainElement(
					map[string]interface{}{
						"name":     "redis",
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
			v1beta1.ElasticacheClusterRedisHostnameOutputName: "test-endpoint",
			v1beta1.ElasticacheClusterRedisPortOutputName:     "asd",
		}
		_, err := o.GetServiceEntrySpecs(outputs)
		Expect(err).To(HaveOccurred())
	})

	It("implements runtime.Object", func() {
		o2 := o.DeepCopyObject()
		Expect(o2).ToNot(BeZero())
	})

	Context("cloudformation", func() {

		It("should generate a unique stack name prefixed with cluster name", func() {
			Expect(o.GetStackName()).To(HavePrefix("xxx-ecr-default-example"))
		})

		It("should have outputs for connection details", func() {
			t, err := o.GetStackTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Outputs).To(And(
				HaveKey("ClusterRedisHostname"),
				HaveKey("ClusterRedisPort"),
			))
		})

		Context("elasticache cluster resource", func() {
			var cluster *cloudformation.AWSElastiCacheCluster

			JustBeforeEach(func() {
				t, err := o.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSElastiCacheCluster{})))
				var ok bool
				cluster, ok = t.Resources[v1beta1.ElasticacheClusterResourceName].(*cloudformation.AWSElastiCacheCluster)
				Expect(ok).To(BeTrue())
			})

			It("should have a cluster name prefixed with cluster and namespace name", func() {
				Expect(cluster.ClusterName).To(Equal("xxx-default-example"))
			})

			It("should be redis", func() {
				Expect(cluster.Engine).To(Equal("redis"))
			})

			It("should run on port 6379", func() {
				Expect(cluster.Port).To(Equal(6379))
			})

			It("should set our cache subnet group name", func() {
				Expect(cluster.CacheSubnetGroupName).To(Equal(cloudformation.Ref(CacheSubnetGroupParameterName)))
			})

			It("should set our VPC subnet group ID", func() {
				Expect(cluster.VpcSecurityGroupIds).To(ConsistOf(
					cloudformation.Ref(VPCSecurityGroupIDParameterName),
				))
			})
/*
TODO: configurable CacheNodeType
TODO: configurable EngineVersion
TODO: configurable NumCacheNodes
TODO: configurable PreferredMaintenanceWindow
TODO: Tags
*/
		})
	})
})