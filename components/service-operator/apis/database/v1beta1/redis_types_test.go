package v1beta1_test

import (
	"os"
	"strconv"

	"github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Redis", func() {

	var o v1beta1.Redis

	BeforeEach(func() {
		os.Setenv("CLUSTER_NAME", "xxx") // required for env package
		o = v1beta1.Redis{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
				Labels: map[string]string{
					cloudformation.AccessGroupLabel: "test.access.group",
				},
			},
			Spec: v1beta1.RedisSpec{
				AWS: v1beta1.RedisAWSSpec{
					NodeType:         "cache.t3.micro",
					EngineVersion:    "5.0.6",
					NumCacheClusters: 2,
				},
			},
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
			v1beta1.RedisPrimaryHostnameOutputName: "test-endpoint.local.govsandbox.uk",
			v1beta1.RedisPrimaryPortOutputName:     "6379",
			v1beta1.RedisReadHostnamesOutputName:   "[test-endpoint-ro.local.govsandbox.uk]",
			v1beta1.RedisReadPortsOutputName:       "[6379]",
		}
		portnum, err := strconv.Atoi(outputs[v1beta1.RedisPrimaryPortOutputName])
		Expect(err).NotTo(HaveOccurred())

		specs, err := o.GetServiceEntrySpecs(outputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(specs).To(HaveLen(2))
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
				HaveKeyWithValue("hosts", ContainElement(outputs[v1beta1.RedisPrimaryHostnameOutputName])),
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
			And(
				HaveKeyWithValue("resolution", "STATIC"),
				HaveKeyWithValue("location", "MESH_EXTERNAL"),
				HaveKeyWithValue("addresses", ContainElement("127.0.0.1")),
				HaveKeyWithValue("endpoints", ContainElement(
					map[string]interface{}{
						"address": "127.0.0.1",
					},
				)),
				HaveKeyWithValue("hosts", ContainElement("test-endpoint-ro.local.govsandbox.uk")),
				HaveKeyWithValue("ports", ContainElement(
					map[string]interface{}{
						"name":     "redis",
						"number":   6379,
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
			v1beta1.RedisPrimaryHostnameOutputName: "test-endpoint",
			v1beta1.RedisPrimaryPortOutputName:     "asd",
			v1beta1.RedisReadHostnamesOutputName:   "[test-endpoint-ro]",
			v1beta1.RedisReadPortsOutputName:       "[6379]",
		}
		_, err := o.GetServiceEntrySpecs(outputs)
		Expect(err).To(HaveOccurred())

		outputs = cloudformation.Outputs{
			v1beta1.RedisPrimaryHostnameOutputName: "test-endpoint",
			v1beta1.RedisPrimaryPortOutputName:     "6379",
			v1beta1.RedisReadHostnamesOutputName:   "[test-endpoint-ro]",
			v1beta1.RedisReadPortsOutputName:       "[asd]",
		}
		_, err = o.GetServiceEntrySpecs(outputs)
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
				HaveKey("ClusterPrimaryRedisHostname"),
				HaveKey("ClusterPrimaryRedisPort"),
				HaveKey("SecretAuthToken"),
			))
		})

		Context("redis resource", func() {
			var cluster *cloudformation.AWSElastiCacheReplicationGroup

			JustBeforeEach(func() {
				t, err := o.GetStackTemplate()
				Expect(err).ToNot(HaveOccurred())
				Expect(t.Resources).To(ContainElement(BeAssignableToTypeOf(&cloudformation.AWSElastiCacheReplicationGroup{})))
				var ok bool
				cluster, ok = t.Resources[v1beta1.RedisResourceName].(*cloudformation.AWSElastiCacheReplicationGroup)
				Expect(ok).To(BeTrue())
			})

			It("should have a replication group ID prefixed with cluster and namespace name", func() {
				Expect(cluster.ReplicationGroupId).To(Equal("xxx-default-example"))
			})

			It("should be redis", func() {
				Expect(cluster.Engine).To(Equal("redis"))
			})

			It("should run on port 6379", func() {
				Expect(cluster.Port).To(Equal(6379))
			})

			It("should set our cache subnet group name", func() {
				Expect(cluster.CacheSubnetGroupName).To(Equal(cloudformation.Ref(v1beta1.CacheSubnetGroupParameterName)))
			})

			It("should have appropriate tags", func() {
				Expect(cluster.Tags).To(Equal(
					[]cloudformation.Tag{
						{Key: "Cluster", Value: env.ClusterName()},
						{Key: "Service", Value: "redis"},
						{Key: "Name", Value: "example"},
						{Key: "Namespace", Value: "default"},
						{Key: "Environment", Value: "default"},
					},
				))
			})

			It("should set our VPC subnet group ID", func() {
				Expect(cluster.SecurityGroupIds).To(ConsistOf(
					cloudformation.Ref(v1beta1.VPCSecurityGroupIDParameterName),
				))
			})

			It("should have an auth token set", func() {
				Expect(cluster.AuthToken).ToNot(BeEmpty())
			})

			It("should set node type appropriately", func() {
				Expect(cluster.CacheNodeType).To(Equal("cache.t3.micro"))
			})

			It("should set engine version appropriately", func() {
				Expect(cluster.EngineVersion).To(Equal("5.0.6"))
			})

			It("should set cache cluster count appropriately", func() {
				Expect(cluster.NumCacheClusters).To(Equal(2))
			})
		})
	})
})
