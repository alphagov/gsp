package aws

import (
	database "github.com/alphagov/gsp/components/service-operator/api/v1beta1"

	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

const (
	Engine       = "aurora-postgresql"
	Family       = "aurora-postgresql10"
	DefaultClass = "db.r5.large"
)

func AuroraPostgres(stackName string, postgresConfig *database.Postgres) *cloudformation.Template {
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
		DBInstanceClass:      coalesceString(postgresConfig.Spec.AWS.InstanceType, DefaultClass),
		Engine:               Engine,
		PubliclyAccessible:   false,
		DBParameterGroupName: cloudformation.Ref("RDSDBParameterGroup"),
	}

	template.Resources["RDSDBInstance2"] = &resources.AWSRDSDBInstance{
		DBClusterIdentifier:  cloudformation.Ref("RDSCluster"),
		DBInstanceClass:      coalesceString(postgresConfig.Spec.AWS.InstanceType, DefaultClass),
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

func coalesceString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
