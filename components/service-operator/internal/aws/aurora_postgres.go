package aws

import (
	"fmt"

	database "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal"

	"github.com/aws/aws-sdk-go/aws"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	Engine       = "aurora-postgresql"
	Family       = "aurora-postgresql10"
	DefaultClass = "db.r5.large"

	InstanceCount = 2

	charactersUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charactersLower   = "abcdefghijklmnopqrstuvwxyz"
	charactersNumeric = "0123456789"
	charactersSpecial = "~=+%^*()[]{}!#$?|"

	PostgresResourceCluster               = "RDSCluster"
	PostgresResourceInstance              = "RDSDBInstance"
	PostgresResourceParameterGroup        = "RDSDBParameterGroup"
	PostgresResourceClusterParameterGroup = "RDSDBClusterParameterGroup"

	PostgresOutputEndpoint     = "Endpoint"
	PostgresOutputReadEndpoint = "ReadEndpoint"
	PostgresOutputPort         = "Port"
	PostgresOutputDBName       = "DBName"
	PostgresOutputUsername     = "DBUsername"
	PostgresOutputPassword     = "DBPassword"
	PostgresOutputEngine       = "Engine"
)

type AuroraPostgres struct {
	PostgresConfig *database.Postgres
	Credentials    internal.BasicAuth
}

func (p *AuroraPostgres) Template(stackName string, tags []resources.Tag) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Parameters[PostgresOutputUsername] = map[string]string{
		"Type": "String",
	}
	template.Parameters[PostgresOutputPassword] = map[string]interface{}{
		"Type":   "String",
		"NoEcho": true,
	}

	template.Resources[PostgresResourceCluster] = &resources.AWSRDSDBCluster{
		Engine:                      Engine,
		MasterUsername:              cloudformation.Ref(PostgresOutputUsername),
		MasterUserPassword:          cloudformation.Ref(PostgresOutputPassword),
		DBClusterParameterGroupName: cloudformation.Ref(PostgresResourceClusterParameterGroup),
		Tags:                        tags,
	}

	for i := 0; i < InstanceCount; i++ {
		template.Resources[fmt.Sprintf("%s%d", PostgresResourceInstance, i)] = &resources.AWSRDSDBInstance{
			DBClusterIdentifier:  cloudformation.Ref(PostgresResourceCluster),
			DBInstanceClass:      internal.CoalesceString(p.PostgresConfig.Spec.AWS.InstanceType, DefaultClass),
			Engine:               Engine,
			PubliclyAccessible:   false,
			DBParameterGroupName: cloudformation.Ref(PostgresResourceParameterGroup),
			Tags:                 tags,
		}
	}

	template.Resources[PostgresResourceClusterParameterGroup] = &resources.AWSRDSDBClusterParameterGroup{
		Description: "GSP Service Operator Cluster Parameter Group",
		Family:      Family,
		Parameters: map[string]string{
			"timezone": "UTC",
		},
		Tags: tags,
	}

	template.Resources[PostgresResourceParameterGroup] = &resources.AWSRDSDBParameterGroup{
		Description: "GSP Service Operator Parameter Group",
		Family:      Family,
		Parameters: map[string]string{
			"application_name": stackName,
		},
		Tags: tags,
	}

	template.Outputs[PostgresOutputEndpoint] = map[string]interface{}{
		"Description": "Postgres Endpoint used by the application to perform connection.",
		"Value":       cloudformation.GetAtt(PostgresResourceCluster, "Endpoint.Address"),
	}

	template.Outputs[PostgresOutputReadEndpoint] = map[string]interface{}{
		"Description": "Postgres reader Endpoint used by the application to perform connection.",
		"Value":       cloudformation.GetAtt(PostgresResourceCluster, "ReadEndpoint.Address"),
	}

	template.Outputs[PostgresOutputPort] = map[string]interface{}{
		"Description": "Postgres Port used by the application to perform connection.",
		"Value":       cloudformation.GetAtt(PostgresResourceCluster, "Endpoint.Port"),
	}

	template.Outputs[PostgresOutputDBName] = map[string]interface{}{
		"Description": "Postgres Database Name used by the application to perform connection.",
		"Value":       cloudformation.Ref(PostgresOutputDBName),
	}

	template.Outputs[PostgresOutputUsername] = map[string]interface{}{
		"Description": "Postgres Username used by the application to perform connection.",
		"Value":       cloudformation.Ref(PostgresOutputUsername),
	}

	template.Outputs[PostgresOutputPassword] = map[string]interface{}{
		"Description": "Postgres Password used by the application to perform connection.",
		"Value":       cloudformation.Ref(PostgresOutputPassword),
	}

	template.Outputs[PostgresOutputEngine] = map[string]interface{}{
		"Description": "Engine used by the application to perform connection.",
		"Value":       cloudformation.Ref(PostgresOutputEngine),
	}

	return template
}

func (p *AuroraPostgres) Parameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{
		&awscloudformation.Parameter{
			ParameterKey:   aws.String(PostgresOutputUsername),
			ParameterValue: aws.String(p.Credentials.Username),
		},
		&awscloudformation.Parameter{
			ParameterKey:   aws.String(PostgresOutputPassword),
			ParameterValue: aws.String(p.Credentials.Password),
		},
	}, nil
}

func (p *AuroraPostgres) ResourceType() string {
	return "postgres"
}
