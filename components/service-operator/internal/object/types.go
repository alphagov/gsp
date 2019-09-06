package object

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SecretNamer names a Secret to hold sensitive details
type SecretNamer interface {
	GetSecretName() string
}

// Service is the interface shared by all service resources
type Service interface {
	runtime.Object
	metav1.Object
	schema.ObjectKind
	StatusReader
	StatusWriter
}

// StatusReader can fetch a status
type StatusReader interface {
	GetStatus() Status
	GetState() State
}

// StatusWriter can set status fields
type StatusWriter interface {
	SetStatus(Status)
	SetState(State)
}

// PrincipalLister declares that a type can return a list of principals
type PrincipalLister interface {
	runtime.Object
	GetPrincipals() []Principal
}

// Principal is the interface shared by all principal types
type Principal interface {
	GetRoleName() string
}
