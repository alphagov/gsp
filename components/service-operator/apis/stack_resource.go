package apis

import (
	"github.com/alphagov/gsp/components/service-operator/internal/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type StackObject interface {
	runtime.Object
	metav1.Object
	aws.Stack
	GetSecretName() string
}
