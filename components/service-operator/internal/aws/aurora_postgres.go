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
	PostgresResourceIAMPolicy             = "RDSIAMRole"

	PostgresEndpoint     = "Endpoint"
	PostgresReadEndpoint = "ReadEndpoint"
	PostgresPort         = "Port"
	PostgresUsername     = "DBUsername"
	PostgresPassword     = "DBPassword"
)

type AuroraPostgres struct {
	PostgresConfig *database.Postgres
	IAMRoleName    string
	SecurityGroup  string
	DBSubnetGroup  string
	MasterUsername string
	MasterPassword string
}

func (p *AuroraPostgres) Template(stackName string, tags []resources.Tag) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Parameters[PostgresUsername] = map[string]string{
		"Type": "String",
	}
	template.Parameters[PostgresPassword] = map[string]interface{}{
		"Type":   "String",
		"NoEcho": true,
	}

	template.Resources[PostgresResourceCluster] = &resources.AWSRDSDBCluster{
		Engine:                      Engine,
		MasterUsername:              cloudformation.Ref(PostgresUsername),
		MasterUserPassword:          cloudformation.Ref(PostgresPassword),
		DBClusterParameterGroupName: cloudformation.Ref(PostgresResourceClusterParameterGroup),
		Tags:                        tags,
		VpcSecurityGroupIds:         []string{p.SecurityGroup},
		DBSubnetGroupName:           p.DBSubnetGroup,
	}

	for i := 0; i < InstanceCount; i++ {
		template.Resources[fmt.Sprintf("%s%d", PostgresResourceInstance, i)] = &resources.AWSRDSDBInstance{
			DBClusterIdentifier:  cloudformation.Ref(PostgresResourceCluster),
			DBInstanceClass:      internal.CoalesceString(p.PostgresConfig.Spec.AWS.InstanceType, DefaultClass),
			Engine:               Engine,
			PubliclyAccessible:   false,
			DBParameterGroupName: cloudformation.Ref(PostgresResourceParameterGroup),
			Tags:                 tags,
			DBSubnetGroupName:    p.DBSubnetGroup,
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

	template.Resources[PostgresResourceIAMPolicy] = &resources.AWSIAMPolicy{
		PolicyName: cloudformation.Join("-", []string{"postgres", "access", cloudformation.Ref(PostgresResourceCluster)}),
		PolicyDocument: NewRolePolicyDocument(
			[]string{
				cloudformation.Join(
					":",
					[]string{
						"arn",
						"aws",
						"rds",
						cloudformation.Ref("AWS::Region"),
						cloudformation.Ref("AWS::AccountId"),
						"cluster",
						cloudformation.Ref(PostgresResourceCluster),
					},
				),
			},
			[]string{"rds-data:*"},
		),
		Roles: []string{p.IAMRoleName},
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

	return template
}

func (p *AuroraPostgres) CreateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{
		&awscloudformation.Parameter{
			ParameterKey:   aws.String(PostgresUsername),
			ParameterValue: aws.String(p.MasterUsername),
		},
		&awscloudformation.Parameter{
			ParameterKey:   aws.String(PostgresPassword),
			ParameterValue: aws.String(p.MasterPassword),
		},
	}, nil
}

func (p *AuroraPostgres) UpdateParameters() ([]*awscloudformation.Parameter, error) {
	return []*awscloudformation.Parameter{
		&awscloudformation.Parameter{
			ParameterKey:     aws.String(PostgresUsername),
			UsePreviousValue: aws.Bool(true),
		},
		&awscloudformation.Parameter{
			ParameterKey:     aws.String(PostgresPassword),
			UsePreviousValue: aws.Bool(true),
		},
	}, nil
}

func (p *AuroraPostgres) ResourceType() string {
	return "postgres"
}
