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

	"github.com/sanathkr/yaml"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&Redis{}, &RedisList{})
}

const (
	RedisResourceName                    = "Redis"
	CacheSubnetGroupParameterName        = "CacheSubnetGroup"
	RedisVPCSecurityGroupIDParameterName = "VPCSecurityGroupID"
	RedisPrimaryHostnameOutputName       = "ClusterPrimaryRedisHostname"
	RedisPrimaryPortOutputName           = "ClusterPrimaryRedisPort"
	RedisReadHostnamesOutputName         = "ClusterReadRedisHostnames"
	RedisReadPortsOutputName             = "ClusterReadRedisPorts"
	AuthTokenSecretResourceName          = "AuthTokenSecret"
	RedisAuthTokenOutputName             = "SecretAuthToken"
)

// ensure implements required interfaces
var _ cloudformation.Stack = &Redis{}
var _ object.SecretNamer = &Redis{}
var _ cloudformation.StackSecretOutputter = &Redis{}
var _ cloudformation.ServiceEntryCreator = &Redis{}

// AWS allows specifying configuration for the redis
type RedisAWSSpec struct {
	// NodeType defines the amount of RAM and CPUs nodes in the cluster have as well as their network performance
	NodeType string `json:"nodeType"`

	// EngineVersion defines the version of Redis running in the cluster.
	EngineVersion string `json:"engineVersion"`

	// NumCacheClusters defines the number of clusters that belong to our replication group. A number between 2 and 6 inclusive.
	NumCacheClusters int `json:"numCacheClusters"`

	// PreferredMaintenanceWindow defines the weekly window during which maintenance is performed on the cluster. The minimum period is 60 minutes.
	PreferredMaintenanceWindow string `json:"preferredMaintenanceWindow,omitempty"`
}

// RedisSpec defines the desired state of Redis
type RedisSpec struct {
	// AWS specific subsection of the resource.
	AWS RedisAWSSpec `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
	// ServiceEntry name to be used for storing the egress firewall rule to allow tenant access to the cluster
	ServiceEntry string `json:"serviceEntry,omitempty"`
}

// +kubebuilder:object:root=true

// RedisList contains a list of Redis
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

// +kubebuilder:object:root=true

// Redis is the Schema for the Redis API
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec          RedisSpec `json:"spec,omitempty"`
	object.Status `json:"status,omitempty"`
}

// Name returns the name of the Redis cloudformation stack
func (s *Redis) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "redis", s.Namespace, s.ObjectMeta.Name)
}

// SecretName returns the name of the secret that will be populated with data
func (s *Redis) GetSecretName() string {
	if s.Spec.Secret == "" {
		return s.GetName()
	}
	return s.Spec.Secret
}

// Template returns a cloudformation Template for provisioning an Redis
func (s *Redis) GetStackTemplate() (*cloudformation.Template, error) {
	template := cloudformation.NewTemplate()

	template.Parameters[RedisVPCSecurityGroupIDParameterName] = map[string]interface{}{
		"Type": "String",
	}
	template.Parameters[CacheSubnetGroupParameterName] = map[string]interface{}{
		"Type": "String",
	}

	// generate secret in cloudformation not in operator (keeps state in aws)
	template.Resources[AuthTokenSecretResourceName] = &cloudformation.AWSSecretsManagerSecret{
		Description: "Auth token for the redis",
		GenerateSecretString: &cloudformation.GenerateSecretString{
			ExcludeCharacters:    "\"%'()*+,./:;=?@[\\]_`{|}~",
			PasswordLength:       128,
			SecretStringTemplate: `{}`,
			GenerateStringKey:    "AuthToken",
		},
	}

	clusterName := fmt.Sprintf("%s-%s-%s", env.ClusterName(), s.Namespace, s.ObjectMeta.Name)
	authTokenRef := cloudformation.Join(":", []string{
		"{{resolve",
		"secretsmanager",
		cloudformation.Ref(AuthTokenSecretResourceName),
		"SecretString",
		"AuthToken}}",
	})
	template.Resources[RedisResourceName] = &cloudformation.AWSElastiCacheReplicationGroup{
		Engine:                      "redis",
		AutomaticFailoverEnabled:    true,
		ReplicationGroupDescription: clusterName,
		ReplicationGroupId:          clusterName,
		CacheNodeType:               s.Spec.AWS.NodeType,
		EngineVersion:               s.Spec.AWS.EngineVersion,
		NumCacheClusters:            s.Spec.AWS.NumCacheClusters,
		Port:                        6379,
		CacheSubnetGroupName:        cloudformation.Ref(CacheSubnetGroupParameterName),
		SecurityGroupIds:            []string{
			cloudformation.Ref(RedisVPCSecurityGroupIDParameterName),
		},
		TransitEncryptionEnabled:    true,
		AuthToken:                   authTokenRef,
		Tags:                        []cloudformation.Tag{
			{
				Key:   "Cluster",
				Value: env.ClusterName(),
			},
			{
				Key:   "Service",
				Value: "redis",
			},
			{
				Key:   "Name",
				Value: s.GetName(),
			},
			{
				Key:   "Namespace",
				Value: s.GetNamespace(),
			},
			{
				Key:   "Environment",
				Value: s.GetNamespace(),
			},
		},
		PreferredMaintenanceWindow: s.Spec.AWS.PreferredMaintenanceWindow,
	}

	template.Outputs[RedisPrimaryHostnameOutputName] = map[string]interface{}{
		"Description": "Redis primary hostname to be returned to the user.",
		"Value":       cloudformation.GetAtt(RedisResourceName, "PrimaryEndPoint.Address"),
	}
	template.Outputs[RedisPrimaryPortOutputName] = map[string]interface{}{
		"Description": "Redis primary port to be returned to the user.",
		"Value":       cloudformation.GetAtt(RedisResourceName, "PrimaryEndPoint.Port"),
	}

	template.Outputs[RedisReadHostnamesOutputName] = map[string]interface{}{
		"Description": "Redis read hostnames to be returned to the user.",
		"Value":       cloudformation.GetAtt(RedisResourceName, "ReadEndPoint.Addresses"),
	}
	template.Outputs[RedisReadPortsOutputName] = map[string]interface{}{
		"Description": "Redis read ports to be returned to the user.",
		"Value":       cloudformation.GetAtt(RedisResourceName, "ReadEndPoint.Ports"),
	}

	template.Outputs[RedisAuthTokenOutputName] = map[string]interface{}{
		"Description": "Redis authentication token to be returned to the user.",
		"Value":       authTokenRef,
	}

	return template, nil
}

func (s *Redis) GetServiceEntryName() string {
	if s.Spec.ServiceEntry == "" {
		return s.GetName()
	}
	return s.Spec.ServiceEntry
}

// ServiceEntry to whitelist egress access to cluster port and hosts.
func (s *Redis) GetServiceEntrySpecs(outputs cloudformation.Outputs) ([]map[string]interface{}, error) {
	// primary
	rwPort, err := strconv.Atoi(outputs[RedisPrimaryPortOutputName])
	if err != nil {
		return nil, err
	}

	rwAddresses, err := net.LookupIP(outputs[RedisPrimaryHostnameOutputName])
	if err != nil {
		return nil, err
	}
	if len(rwAddresses) < 1 {
		return nil, fmt.Errorf("list of rw endpoint IPs was empty - failed to resolve?")
	}
	rwAddress := rwAddresses[0].String()
	if rwAddress == "<nil>" {
		return nil, fmt.Errorf("unexpected nil returned for rw endpoint IP")
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
				outputs[RedisPrimaryHostnameOutputName],
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "redis",
					"number":   rwPort,
					"protocol": "TCP",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "STATIC",
			"exportTo":   []string{"."},
		},
	}

	// read-only endpoints
	var roHostnames []string
	err = yaml.Unmarshal([]byte(outputs[RedisReadHostnamesOutputName]), &roHostnames)
	if err != nil {
		return nil, fmt.Errorf("error YAML-unmarshalling read hostnames")
	}

	var roPorts []int
	err = yaml.Unmarshal([]byte(outputs[RedisReadPortsOutputName]), &roPorts)
	if err != nil {
		return nil, fmt.Errorf("error YAML-unmarshalling read ports")
	}

	if len(roHostnames) != len(roPorts) {
		return nil, fmt.Errorf("read hostnames and read ports lists have different lengths")
	}

	for i := 0; i < len(roHostnames); i++ {
		adresses, err := net.LookupIP(roHostnames[i])
		if err != nil {
			return nil, err
		}
		if len(adresses) < 1 {
			return nil, fmt.Errorf("list of ro endpoint IPs was empty - failed to resolve?")
		}
		address := adresses[0].String()
		if address == "<nil>" {
			return nil, fmt.Errorf("unexpected nil returned for ro endpoint IP")
		}

		specs = append(specs, map[string]interface{}{
			"addresses": []string{
				address,
			},
			"endpoints": []map[string]interface{}{
				{
					"address": address,
				},
			},
			"hosts": []string{
				roHostnames[i],
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "redis",
					"number":   roPorts[i],
					"protocol": "TCP",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "STATIC",
			"exportTo":   []string{"."},
		})
	}

	return specs, nil
}
