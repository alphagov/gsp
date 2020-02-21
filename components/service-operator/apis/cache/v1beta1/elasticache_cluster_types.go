/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
	"net"
	"strconv"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&ElasticacheCluster{}, &ElasticacheClusterList{})
}

const (
	ElasticacheClusterResourceName            = "ElasticacheCluster"
	CacheSubnetGroupParameterName             = "CacheSubnetGroup"
	VPCSecurityGroupIDParameterName           = "VPCSecurityGroupID"
	ElasticacheClusterRedisHostnameOutputName = "ClusterRedisHostname"
	ElasticacheClusterRedisPortOutputName     = "ClusterRedisPort"
)

// ensure implements required interfaces
var _ cloudformation.Stack = &ElasticacheCluster{}
var _ object.SecretNamer = &ElasticacheCluster{}
var _ cloudformation.StackSecretOutputter = &ElasticacheCluster{}
var _ cloudformation.ServiceEntryCreator = &ElasticacheCluster{}

// AWS allows specifying configuration for the elasticache cluster
type ElasticacheClusterAWSSpec struct {
	// InstanceType essentially defines the amount of memory and cpus on the database.
	//InstanceType string `json:"instanceType,omitempty"`
	// InstanceCount is the number of database instances in the cluster (defaults to 2 if not set)
	//InstanceCount int `json:"instanceCount,omitempty"`
}

// ElasticacheClusterSpec defines the desired state of ElasticacheCluster
type ElasticacheClusterSpec struct {
	// AWS specific subsection of the resource.
	AWS ElasticacheClusterAWSSpec `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
	// ServiceEntry name to be used for storing the egress firewall rule to allow tenant access to the cluster
	ServiceEntry string `json:"serviceEntry,omitempty"`
}

// +kubebuilder:object:root=true

// ElasticacheClusterList contains a list of ElasticacheCluster
type ElasticacheClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticacheCluster `json:"items"`
}

// +kubebuilder:object:root=true

// ElasticacheCluster is the Schema for the ElasticacheCluster API
type ElasticacheCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec          ElasticacheClusterSpec `json:"spec,omitempty"`
	object.Status `json:"status,omitempty"`
}

// Name returns the name of the ElasticacheCluster cloudformation stack
func (s *ElasticacheCluster) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "elasticache", s.Namespace, s.ObjectMeta.Name)
}

// SecretName returns the name of the secret that will be populated with data
func (s *ElasticacheCluster) GetSecretName() string {
	if s.Spec.Secret == "" {
		return s.GetName()
	}
	return s.Spec.Secret
}

// Template returns a cloudformation Template for provisioning an ElasticacheCluster
func (s *ElasticacheCluster) GetStackTemplate() (*cloudformation.Template, error) {
	template := cloudformation.NewTemplate()

	template.Parameters[VPCSecurityGroupIDParameterName] = map[string]interface{}{
		"Type": "String",
	}
	template.Parameters[CacheSubnetGroupParameterName] = map[string]interface{}{
		"Type": "String",
	}

	clusterName := fmt.Sprintf("%s-%s-%s", env.ClusterName(), s.Namespace, s.ObjectMeta.Name)
	template.Resources[ElasticacheClusterResourceName] = &cloudformation.AWSElastiCacheReplicationGroup{
		// TODO: make PreferredMaintenanceWindow configurable?
		// TODO: add Tags?

		Engine:                      "redis",
		AutomaticFailoverEnabled:    true,
		ReplicationGroupDescription: "", // TODO
		ReplicationGroupId:          clusterName,
		CacheNodeType:               "cache.t3.micro", // TODO: make configurable
		EngineVersion:               "1.4.24", // TODO: make configurable
		NumCacheClusters:            1, // TODO: make configurable
		Port:                        6379,
		CacheSubnetGroupName:        cloudformation.Ref(CacheSubnetGroupParameterName),
		SecurityGroupIds:            []string{
			cloudformation.Ref(VPCSecurityGroupIDParameterName),
		},
		TransitEncryptionEnabled:    true,
		AuthToken:                   "hunter2hunter2hunter2", // TODO
	}

	template.Outputs[ElasticacheClusterRedisHostnameOutputName] = map[string]interface{}{
		"Description": "Elasticache Cluster Redis hostname to be returned to the user.",
		"Value":       cloudformation.GetAtt(ElasticacheClusterResourceName, "PrimaryEndPoint.Address"),
	}
	template.Outputs[ElasticacheClusterRedisPortOutputName] = map[string]interface{}{
		"Description": "Elasticache Cluster Redis port to be returned to the user.",
		"Value":       cloudformation.GetAtt(ElasticacheClusterResourceName, "PrimaryEndPoint.Port"),
	}

	return template, nil
}

func (s *ElasticacheCluster) GetServiceEntryName() string {
	if s.Spec.ServiceEntry == "" {
		return s.GetName()
	}
	return s.Spec.ServiceEntry
}

// ServiceEntry to whitelist egress access to cluster port and hosts.
func (s *ElasticacheCluster) GetServiceEntrySpecs(outputs cloudformation.Outputs) ([]map[string]interface{}, error) {
	port, err := strconv.Atoi(outputs[ElasticacheClusterRedisPortOutputName])
	if err != nil {
		return nil, err
	}

	rwAddresses, err := net.LookupIP(outputs[ElasticacheClusterRedisHostnameOutputName])
	if err != nil {
		return nil, err
	}
	if len(rwAddresses) < 1 {
		return nil, fmt.Errorf("list of endpoint IPs was empty - failed to resolve?")
	}
	rwAddress := rwAddresses[0].String()
	if rwAddress == "<nil>" {
		return nil, fmt.Errorf("unexpected nil returned for endpoint IP")
	}

	specs := []map[string]interface{}{
		{
			"addresses": []string{
				rwAddress,
			},
			"endpoints": []map[string]interface{}{
				{
					"address": rwAddress,
				},
			},
			"hosts": []string{
				outputs[ElasticacheClusterRedisHostnameOutputName],
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "redis",
					"number":   port,
					"protocol": "TCP",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "STATIC",
			"exportTo":   []string{"."},
		},
	}
	return specs, nil
}