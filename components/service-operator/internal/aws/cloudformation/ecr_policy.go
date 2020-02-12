package cloudformation

const (
	ECRLifecycleMoreThan     = "imageCountMoreThan"
	ECRLifecyclePolicyExpire = "expire"
)

type ECRLifecyclePolicySelection struct {
	TagStatus   string `json:"tagStatus,omitempty"`
	CountType   string `json:"countType,omitempty"`
	CountNumber int    `json:"countNumber,omitempty"`
}

type ECRLifecyclePolicyAction struct {
	Type string `json:"type,omitempty"`
}

type ECRLifecyclePolicyRule struct {
	RulePriority int64                       `json:"rulePriority"`
	Description  string                      `json:"description,omitempty"`
	Selection    ECRLifecyclePolicySelection `json:"selection,omitempty"`
	Action       ECRLifecyclePolicyAction    `json:"action,omitempty"`
}

type ECRLifecyclePolicy struct {
	Rules []ECRLifecyclePolicyRule `json:"rules,omitempty"`
}
