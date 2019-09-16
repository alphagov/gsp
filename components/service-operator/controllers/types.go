package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Controller is the interface for all "reconcilers". It adds a consistent way
// to to register with managers
type Controller interface {
	reconcile.Reconciler
	SetupWithManager(ctrl.Manager) error
}

// ControllerWrapper allows us to wrap a reconcile.Reconciler to collect and
// inspect any errors. This is useful for tests
//
// Example:
//
//     myWrappedReconciler := &ControllerWrapper{
//         Reconciler: &myReconcilerUnderTest{},
//     }
//
//     go myWrappeedReconciler.Reconcile(...) // Reconcile called aysync elsewhere
//
//     Consistently(myWrappedReconciler.Err, time.Second*2).ShouldNot(HaveOccurred())
//
//
type ControllerWrapper struct {
	Reconciler Controller
	errs       []error
}

var _ Controller = &ControllerWrapper{}

// Reconcile forwards the Reconcile call to the real Controller and collects any errors
func (r *ControllerWrapper) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	res, err := r.Reconciler.Reconcile(req)
	if err != nil {
		r.errs = append(r.errs, err)
	}
	return res, err
}

// SetupWithManager emulates the SetupWithManager call
func (r *ControllerWrapper) SetupWithManager(mgr ctrl.Manager) error {
	return r.Reconciler.SetupWithManager(mgr)
}

// Err shifts the first error off the list of collected errors, returns it and flushes the error queue
func (r *ControllerWrapper) Err() error {
	errs := r.Errs()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Errs returns all errors collected and flushes the error queue
func (r *ControllerWrapper) Errs() []error {
	defer func() {
		r.errs = []error{}
	}()
	return r.errs
}
