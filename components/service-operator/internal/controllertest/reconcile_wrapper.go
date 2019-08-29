package controllertest

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcilerWrapper allows us to wrap a reconcile.Reconciler to collect and
// inspect any errors returned from within tests.
//
// Example:
//
//     myWrappedReconciler := &ReconcilerWrapper{
//         Reconciler: &myReconcilerUnderTest{},
//     }
//
//     go myWrappeedReconciler.Reconcile(...) // Reconcile called aysync elsewhere
//
//     Consistently(myWrappedReconciler.Err, time.Second*2).ShouldNot(HaveOccurred())
//
//
type ReconcilerWrapper struct {
	Reconciler reconcile.Reconciler
	errs       []error
}

// Reconcile forwards the Reconcile call to the real Reconciler and collects any errors
func (r *ReconcilerWrapper) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	res, err := r.Reconciler.Reconcile(req)
	if err != nil {
		r.errs = append(r.errs, err)
	}
	return res, err
}

// SetupWithManager emulates the SetupWithManager call
func (r *ReconcilerWrapper) SetupWithManager(mgr ctrl.Manager, t runtime.Object) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(t).
		Complete(r)
}

// Err shifts the first error off the list of collected errors, returns it and flushes the error queue
func (r *ReconcilerWrapper) Err() error {
	errs := r.Errs()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Errs returns all errors collected and flushes the error queue
func (r *ReconcilerWrapper) Errs() []error {
	defer func() {
		r.errs = []error{}
	}()
	return r.errs
}
