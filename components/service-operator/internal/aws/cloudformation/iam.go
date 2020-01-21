package cloudformation

// helpers for building iam documents in cloudformation

func NewAssumeRolePolicyDocument(principal, serviceOperatorRoleArn string) AssumeRolePolicyDocument {
	return AssumeRolePolicyDocument{
		Version: "2012-10-17",
		Statement: []AssumeRolePolicyStatement{
			{
				Effect: "Allow",
				Principal: PolicyPrincipal{
					AWS: []string{
						principal,
						serviceOperatorRoleArn,
					},
				},
				Action: []string{"sts:AssumeRole"},
			},
		},
	}
}

type PolicyDocument struct {
	Version   string
	Statement []PolicyStatement
}

type PolicyStatement struct {
	Effect   string
	Action   []string
	Resource []string
}

type AssumeRolePolicyDocument struct {
	Version   string
	Statement []AssumeRolePolicyStatement
}

type AssumeRolePolicyStatement struct {
	Effect    string
	Principal PolicyPrincipal
	Action    []string
}

type PolicyPrincipal struct {
	AWS []string
}

func NewRolePolicyDocument(resources, actions []string) PolicyDocument {
	return PolicyDocument{
		Version: "2012-10-17",
		Statement: []PolicyStatement{
			{
				Effect:   "Allow",
				Action:   actions,
				Resource: resources,
			},
		},
	}
}
