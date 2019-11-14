package aws

type StackPolicyDocument struct {
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect    string
	Action    []string
	Principal string
	Resource  string
}
