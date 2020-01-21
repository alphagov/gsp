package cloudformation

// helpers for building iam documents in cloudformation

func NewAssumeRolePolicyDocument(awsPrincipal, serviceOperatorRoleArn string) AssumeRolePolicyDocument {
	return AssumeRolePolicyDocument{
		Version: "2012-10-17",
		Statement: []AssumeRolePolicyStatement{
			{
				Effect: "Allow",
				Principal: PolicyPrincipal{
					AWS: []string{
						awsPrincipal,
						serviceOperatorRoleArn,
					},
				},
				Action: []string{"sts:AssumeRole"},
			},
		},
	}
}

func NewAssumeRolePolicyDocumentWithServiceAccount(awsPrincipal string, serviceOperatorRoleArn string, federatedPrincipal string, federatedConditionKey string, federatedConditionValue string) AssumeRolePolicyDocument {
	return AssumeRolePolicyDocument{
		Version: "2012-10-17",
		Statement: []AssumeRolePolicyStatement{
			{
				Effect: "Allow",
				Principal: PolicyPrincipal{
					AWS: []string{
						awsPrincipal,
						serviceOperatorRoleArn,
					},
				},
				Action: []string{"sts:AssumeRole"},
			},
			{
				Effect: "Allow",
				Action: []string{"sts:AssumeRoleWithWebIdentity"},
				Principal: PolicyPrincipal{
					Federated: []string{federatedPrincipal},
				},
				Condition: PolicyCondition{
					StringEquals: map[string]string{
						federatedConditionKey: federatedConditionValue,
					},
				},
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
	Condition PolicyCondition `json:"Condition,omitempty"`
}

type PolicyPrincipal struct {
	AWS       []string `json:"AWS,omitempty"`
	Federated []string `json:"Federated,omitempty"`
}

type PolicyCondition struct {
	StringEquals map[string]string `json:"StringEquals,omitempty"`
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
