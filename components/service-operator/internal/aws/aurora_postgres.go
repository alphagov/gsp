package aws

import (
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

	charactersUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charactersLower   = "abcdefghijklmnopqrstuvwxyz"
	charactersNumeric = "0123456789"
	charactersSpecial = "~=+%^*()[]{}!#$?|"
)

type AuroraPostgres struct {
	PostgresConfig *database.Postgres
}

func (p *AuroraPostgres) Template(stackName string) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Parameters["MasterUsername"] = map[string]string{
		"Type": "String",
	}
	template.Parameters["MasterPassword"] = map[string]interface{}{
		"Type":   "String",
		"NoEcho": true,
	}

	template.Resources["RDSCluster"] = &resources.AWSRDSDBCluster{
		Engine:                      Engine,
		MasterUsername:              cloudformation.Ref("MasterUsername"),
		MasterUserPassword:          cloudformation.Ref("MasterPassword"),
		DBClusterParameterGroupName: cloudformation.Ref("RDSDBClusterParameterGroup"),
	}

	template.Resources["RDSDBInstance1"] = &resources.AWSRDSDBInstance{
		DBClusterIdentifier:  cloudformation.Ref("RDSCluster"),
		DBInstanceClass:      internal.CoalesceString(p.PostgresConfig.Spec.AWS.InstanceType, DefaultClass),
		Engine:               Engine,
		PubliclyAccessible:   false,
		DBParameterGroupName: cloudformation.Ref("RDSDBParameterGroup"),
	}

	template.Resources["RDSDBInstance2"] = &resources.AWSRDSDBInstance{
		DBClusterIdentifier:  cloudformation.Ref("RDSCluster"),
		DBInstanceClass:      internal.CoalesceString(p.PostgresConfig.Spec.AWS.InstanceType, DefaultClass),
		Engine:               Engine,
		PubliclyAccessible:   false,
		DBParameterGroupName: cloudformation.Ref("RDSDBParameterGroup"),
	}

	template.Resources["RDSDBClusterParameterGroup"] = &resources.AWSRDSDBClusterParameterGroup{
		Description: "GSP Service Operator Cluster Parameter Group",
		Family:      Family,
		Parameters: map[string]string{
			"timezone": "UTC",
		},
	}

	template.Resources["RDSDBParameterGroup"] = &resources.AWSRDSDBParameterGroup{
		Description: "GSP Service Operator Parameter Group",
		Family:      Family,
		Parameters: map[string]string{
			"application_name": stackName,
		},
	}

	return template
}

func (p *AuroraPostgres) Parameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{
		&awscloudformation.Parameter{
			ParameterKey:   aws.String("MasterUsername"),
			ParameterValue: aws.String("qwertyuiop"),
			//      ParameterValue: aws.String(username),
		},
		&awscloudformation.Parameter{
			ParameterKey:   aws.String("MasterPassword"),
			ParameterValue: aws.String("qwertyuiop1234567890"),
			//      ParameterValue: aws.String(password),
		},
	}, nil
}

func (p *AuroraPostgres) ResourceName() string {
	return "postgres"
}
