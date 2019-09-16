package object

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HasFinalizer is a helper for checking if finalizer exists
func HasFinalizer(o metav1.Object, finalizer string) bool {
	finalizers := o.GetFinalizers()
	return contains(finalizers, finalizer)
}

// SetFinalizer adds finalizer to object if not exists
func SetFinalizer(o metav1.Object, finalizer string) {
	finalizers := o.GetFinalizers()
	if !contains(finalizers, finalizer) {
		o.SetFinalizers(append(finalizers, finalizer))
	}
}

// RemoveFinalizer removes finalizer from object if it exists
func RemoveFinalizer(o metav1.Object, finalizer string) {
	finalizers := o.GetFinalizers()
	o.SetFinalizers(remove(finalizers, finalizer))
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func remove(slice []string, s string) []string {
	result := []string{}
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return result
}
