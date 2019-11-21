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
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alphagov/gsp/components/service-operator/internal/aws"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/env"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
)

func init() {
	SchemeBuilder.Register(&Postgres{}, &PostgresList{})
}

const (
	Engine                       = "aurora-postgresql"
	Family                       = "aurora-postgresql10"
	DefaultClass                 = "db.r5.large"
	DefaultInstanceCount         = 2
	DefaultBackupRetentionPeriod = 7

	PostgresResourceMasterCredentials           = "MasterCredentials"
	PostgresResourceMasterCredentialsAttachment = "MasterCredentialsAttachment"
	PostgresResourceCluster                     = "RDSCluster"
	PostgresResourceInstance                    = "RDSDBInstance"
	PostgresResourceParameterGroup              = "RDSDBParameterGroup"
	PostgresResourceClusterParameterGroup       = "RDSDBClusterParameterGroup"

	VPCSecurityGroupIDParameterName = "VPCSecurityGroupID"
	DBSubnetGroupNameParameterName  = "DBSubnetGroup"

	PostgresEndpoint     = "Endpoint"
	PostgresReadEndpoint = "ReadEndpoint"
	PostgresPort         = "Port"
	PostgresUsername     = "Username"
	PostgresPassword     = "Password"
)

var _ cloudformation.Stack = &Postgres{}
var _ object.SecretNamer = &Postgres{}
var _ cloudformation.ServiceEntryCreator = &Postgres{}

// AWS allows specifying configuration for the Postgres RDS instance
type PostgresAWSSpec struct {
	// InstanceType essentially defines the amount of memory and cpus on the database.
	InstanceType string `json:"instanceType,omitempty"`
	// InstanceCount is the number of database instances in the cluster (defaults to 2 if not set)
	InstanceCount int `json:"instanceCount,omitempty"`
}

// PostgresSpec defines the desired state of Postgres
type PostgresSpec struct {
	// AWS specific subsection of the resource.
	AWS PostgresAWSSpec `json:"aws,omitempty"`
	// Secret name to be used for storing relevant instance secrets for further use.
	Secret string `json:"secret,omitempty"`
	// ServiceEntry name to be used for storing the egress firewall rule to allow tenant access to the database
	ServiceEntry string `json:"serviceEntry,omitempty"`
}

// +kubebuilder:object:root=true

// Postgres is the Schema for the postgres API
type Postgres struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	object.Status     `json:"status,omitempty"`

	Spec PostgresSpec `json:"spec,omitempty"`
}

// GetSecretName returns the name of a secret that will get populated with
// connection details
func (p *Postgres) GetSecretName() string {
	if p.Spec.Secret == "" {
		return p.GetName()
	}
	return p.Spec.Secret
}

// Name returns the name of the SQS cloudformation stack
func (p *Postgres) GetStackName() string {
	return fmt.Sprintf("%s-%s-%s-%s", env.ClusterName(), "postgres", p.GetNamespace(), p.GetName())
}

// GetStackTemplate implements cloudformation.Stack to serialize the configuration as a cloudformaiton template
func (p *Postgres) GetStackTemplate() *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Parameters[VPCSecurityGroupIDParameterName] = map[string]interface{}{
		"Type": "String",
	}
	template.Parameters[DBSubnetGroupNameParameterName] = map[string]interface{}{
		"Type": "String",
	}

	tags := []cloudformation.Tag{
		{
			Key:   "Cluster",
			Value: env.ClusterName(),
		},
		{
			Key:   "Service",
			Value: "postgres",
		},
		{
			Key:   "Name",
			Value: p.GetName(),
		},
		{
			Key:   "Namespace",
			Value: p.GetNamespace(),
		},
		{
			Key:   "Environment",
			Value: p.GetNamespace(),
		},
	}

	// generate secret in cloudformation not in operator (keeps state in aws)
	template.Resources[PostgresResourceMasterCredentials] = &cloudformation.AWSSecretsManagerSecret{
		Description: "Master Credentials for postgres instance",
		GenerateSecretString: &cloudformation.GenerateSecretString{
			SecretStringTemplate: `{"username": "master"}`,
			ExcludeCharacters:    `"@/\'`,
			GenerateStringKey:    "password",
			PasswordLength:       32,
		},
	}

	// create a reference to the values created in secrets manager
	masterUsernameSecretRef := cloudformation.Join("", []string{"{{resolve:secretsmanager:", cloudformation.Ref(PostgresResourceMasterCredentials), ":SecretString:username}}"})
	masterPasswordSecretRef := cloudformation.Join("", []string{"{{resolve:secretsmanager:", cloudformation.Ref(PostgresResourceMasterCredentials), ":SecretString:password}}"})

	template.Resources[PostgresResourceCluster] = &cloudformation.AWSRDSDBCluster{
		Engine:                      Engine,
		MasterUsername:              masterUsernameSecretRef,
		MasterUserPassword:          masterPasswordSecretRef,
		DBClusterParameterGroupName: cloudformation.Ref(PostgresResourceClusterParameterGroup),
		Tags:                        tags,
		VpcSecurityGroupIds: []string{
			cloudformation.Ref(VPCSecurityGroupIDParameterName),
		},
		DBSubnetGroupName:     cloudformation.Ref(DBSubnetGroupNameParameterName),
		BackupRetentionPeriod: DefaultBackupRetentionPeriod,
	}

	template.Resources[PostgresResourceMasterCredentialsAttachment] = &cloudformation.AWSSecretsManagerSecretTargetAttachment{
		SecretId:   cloudformation.Ref(PostgresResourceMasterCredentials),
		TargetId:   cloudformation.Ref(PostgresResourceCluster),
		TargetType: "AWS::RDS::DBCluster",
	}

	instanceCount := p.Spec.AWS.InstanceCount
	if instanceCount < 1 {
		instanceCount = DefaultInstanceCount
	}
	for i := 0; i < instanceCount; i++ {
		template.Resources[fmt.Sprintf("%s%d", PostgresResourceInstance, i)] = &cloudformation.AWSRDSDBInstance{
			DBClusterIdentifier:    cloudformation.Ref(PostgresResourceCluster),
			DBInstanceClass:        coalesce(p.Spec.AWS.InstanceType, DefaultClass),
			Engine:                 Engine,
			PubliclyAccessible:     false,
			DBParameterGroupName:   cloudformation.Ref(PostgresResourceParameterGroup),
			Tags:                   tags,
			DBSubnetGroupName:      cloudformation.Ref(DBSubnetGroupNameParameterName),
			DeleteAutomatedBackups: false,
		}
	}

	template.Resources[PostgresResourceClusterParameterGroup] = &cloudformation.AWSRDSDBClusterParameterGroup{
		Description: "GSP Service Operator Cluster Parameter Group",
		Family:      Family,
		Parameters: map[string]string{
			"timezone": "UTC",
		},
		Tags: tags,
	}

	template.Resources[PostgresResourceParameterGroup] = &cloudformation.AWSRDSDBParameterGroup{
		Description: "GSP Service Operator Parameter Group",
		Family:      Family,
		Parameters: map[string]string{
			"application_name": p.GetStackName(),
		},
		Tags: tags,
	}

	template.Outputs[PostgresEndpoint] = map[string]interface{}{
		"Description": "Postgres Endpoint used by the application to perform connection.",
		"Value":       cloudformation.GetAtt(PostgresResourceCluster, "Endpoint.Address"),
	}

	template.Outputs[PostgresReadEndpoint] = map[string]interface{}{
		"Description": "Postgres reader Endpoint used by the application to perform connection.",
		"Value":       cloudformation.GetAtt(PostgresResourceCluster, "ReadEndpoint.Address"),
	}

	template.Outputs[PostgresPort] = map[string]interface{}{
		"Description": "Postgres Port used by the application to perform connection.",
		"Value":       cloudformation.GetAtt(PostgresResourceCluster, "Endpoint.Port"),
	}

	template.Outputs[PostgresUsername] = map[string]interface{}{
		"Description": "Postgres master username",
		"Value":       masterUsernameSecretRef,
	}

	template.Outputs[PostgresPassword] = map[string]interface{}{
		"Description": "Postgres master password",
		"Value":       masterPasswordSecretRef,
	}

	return template
}

func (p *Postgres) GetServiceEntryName() string {
	if p.Spec.ServiceEntry == "" {
		return p.GetName()
	}
	return p.Spec.ServiceEntry
}

// ServiceEntry to whitelist egress access to Postgres port and hosts.
func (p *Postgres) GetServiceEntrySpecs(outputs cloudformation.Outputs) ([]map[string]interface{}, error) {
	port, err := strconv.Atoi(outputs[PostgresPort])
	if err != nil {
		return nil, err
	}
	specs := []map[string]interface{}{
		{
			"hosts": []string{
				outputs[PostgresEndpoint],
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "aurora",
					"number":   port,
					"protocol": "TCP",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "DNS",
			"exportTo":   []string{"."},
		}, {
			"hosts": []string{
				outputs[PostgresReadEndpoint],
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":     "aurora",
					"number":   port,
					"protocol": "TCP",
				},
			},
			"location":   "MESH_EXTERNAL",
			"resolution": "DNS",
			"exportTo":   []string{"."},
		},
	}
	return specs, nil
}

// GetStackPolicy implements cloudformation.Stack to return a serialised form of the stack policy, or nil if one is
// not needed
func (p *Postgres) GetStackPolicy() aws.StackPolicyDocument {
	statements := []aws.StatementEntry{
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
	}

	instanceCount := p.Spec.AWS.InstanceCount
	if instanceCount < 1 {
		instanceCount = DefaultInstanceCount
	}

	for i := 0; i < instanceCount; i++ {
		statements = append(statements, aws.StatementEntry{
			Effect:    "Deny",
			Action:    []string{"Update:Replace", "Update:Delete"},
			Principal: "*",
			Resource:  fmt.Sprintf("LogicalResourceId/%s%d", PostgresResourceInstance , i),
		})

		statements = append(statements, aws.StatementEntry{
			Effect:    "Allow",
			Action:    []string{"Update:Modify"},
			Principal: "*",
			Resource:  fmt.Sprintf("LogicalResourceId/%s%d", PostgresResourceInstance, i),
		})
	}

	stackPolicy := aws.StackPolicyDocument{
		Statement: statements,
	}

	return stackPolicy
}

// +kubebuilder:object:root=true

// PostgresList contains a list of Postgres
type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Postgres `json:"items,omitempty"`
}

func coalesce(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
