package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/labels"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/object"
)

type ServiceAccountController struct {
	Scheme               *runtime.Scheme        // Scheme is required for operations like gc
	Log                  logr.Logger            // Log will be used to report each reconcile
	KubernetesClient     client.Client          // KubernetesClient is required to talk to api
}

// SetupWithManager validates and registers this controller with the manager and api
func (r *ServiceAccountController) SetupWithManager(mgr ctrl.Manager) error {
	r.KubernetesClient = mgr.GetClient()
	r.Scheme = mgr.GetScheme()
	// setup logger
	r.Log = ctrl.Log.WithName("controllers").WithName((&core.ServiceAccount{}).GetResourceVersion())
	// register with manager
	return ctrl.NewControllerManagedBy(mgr).
		For(&core.ServiceAccount{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=,resources=serviceaccount,verbs=get;list;watch;create;update;patch;delete

// Reconcile synchronises state between the resource and a cloudformation stack
func (r *ServiceAccountController) Reconcile(req ctrl.Request) (res ctrl.Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()
	// execute reconciliation and log changes
	op, err := r.reconcileWithContext(ctx, req)
	if err == context.DeadlineExceeded {
		// ran out of time, most likely waiting on
		// a long running provisioning, come back a bit later
		res.Requeue = true
		res.RequeueAfter = time.Second * 1
		err = nil
	}
	r.Log.Info("reconciled",
		"service", req.NamespacedName,
		"requeue", res.Requeue,
		"after", res.RequeueAfter,
		"op", op,
		"err", err,
	)
	return res, err
}

// reconcileWithContext fetches the resource to reconcile and executes reconcileServiceAccountWithContext and returns if any changes were made
func (r *ServiceAccountController) reconcileWithContext(ctx context.Context, req ctrl.Request) (controllerutil.OperationResult, error) {
	bg := context.Background()
	o := &core.ServiceAccount{}
	if err := r.KubernetesClient.Get(bg, req.NamespacedName, o); apierrs.IsNotFound(err) {
		// nothing we can do if the resource has gone missing, so
		// ignore any not found errors and let the api carry on
		return controllerutil.OperationResultNone, nil
	} else if err != nil {
		// issue communicating with the api
		// return err and we'll retry later
		return controllerutil.OperationResultNone, err
	}

	if _, ok := o.ObjectMeta.Labels[cloudformation.AccessGroupLabel]; !ok {
		return controllerutil.OperationResultNone, nil
	}

	// track changes to our object resource and call the main reconcile func
	var reconcileErr error
	op, updateErr := controllerutil.CreateOrUpdate(bg, r.KubernetesClient, o, func() error {
		reconcileErr = r.reconcileServiceAccountWithContext(ctx, req, o)
		return nil // always try to update
	})
	if reconcileErr != nil {
		return op, reconcileErr
	} else if updateErr != nil {
		return op, updateErr
	}
	return op, nil
}

func (r *ServiceAccountController) updatePrincipal(ctx context.Context, o *core.ServiceAccount) (*access.Principal, error) {
	// List principals with labels set by o.ObjectMeta.Labels
	list := &access.PrincipalList{}
	listOptsFunc := func(opts *client.ListOptions) {
		opts.Namespace = o.GetNamespace()
		opts.LabelSelector = labels.SelectorFromSet(o.ObjectMeta.Labels)
	}
	err := r.KubernetesClient.List(ctx, list, listOptsFunc)
	if err != nil {
		return nil, err
	}
	principals := list.GetPrincipals()
	var principal *access.Principal
	if len(principals) == 0 {
		// if none are found, make a new one with GenerateName
		principal = &access.Principal{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: o.GetName(),
				Namespace:    o.GetNamespace(),
			},
		}
	} else if len(principals) == 1 {
		// if one is found, use it
		var ok bool
		principal, ok = principals[0].(*access.Principal)
		if !ok {
			return nil, fmt.Errorf("found principal but could not cast to *access.Principal")
		}
	} else if len(principals) > 1 {
		return nil, fmt.Errorf("multiple principals found with service operator's labels")
	}

	op, err := controllerutil.CreateOrUpdate(ctx, r.KubernetesClient, principal, func() error {
		principal.ObjectMeta.Labels = o.ObjectMeta.Labels
		principal.Spec = access.PrincipalSpec{
			TrustServiceAccount: o.GetName(),
		}
		// mark the principal as owned by the o resource so it gets gc'd
		if err := controllerutil.SetControllerReference(o, principal, r.Scheme); err != nil {
			return err
		}
		return nil
	})
	r.Log.Info("update-principal",
		"namespace", o.GetNamespace(),
		"svcacc", o.GetName(),
		"op", op,
		"err", err,
	)
	if err != nil {
		return nil, err
	}
	return principal, nil
}

// reconcileServiceAccountWithContext is the main loop, it will mutate "o" with any changes required
func (r *ServiceAccountController) reconcileServiceAccountWithContext(ctx context.Context, req ctrl.Request, sa *core.ServiceAccount) error {
	defer r.Log.Info("reconcileServiceAccountWithContext",
		"sa", sa,
	)
	// examine DeletionTimestamp to determine if object is under deletion
	if !sa.GetDeletionTimestamp().IsZero() {
		// The object is being deleted
		return nil
	}

	principal, err := r.updatePrincipal(ctx, sa)
	if err != nil {
		return err
	}

	if principal.Status.State != object.ReadyState {
		return fmt.Errorf("principal not ready")
	}

	sa.Annotations = map[string]string{
		"eks.amazonaws.com/role-arn": principal.Status.AWS.Info[access.IAMRoleArnOutputName],
	}

	return nil
}
