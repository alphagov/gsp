package cloudformation

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/alphagov/gsp/components/service-operator/internal/object"
)

const (
	// Finalizer is assigned to objects that have cloudformation state
	Finalizer = "cloudformatiton.finalizers.govsvc.uk"
	// AccessGroupLabel is a label used reference a target IAM role for StackPolicyAttacher
	AccessGroupLabel = "group.access.govsvc.uk"
)

var (
	// DefaultReconcileDeadline is the default timeout for reconcile
	// context. DeadlineExceeded errors will be retried later.
	DefaultReconcileDeadline = time.Minute * 5
	// DefaultRequeueTimeout is the default time when a reconcile needs
	// requeuing after deadline is hit
	DefaultRequeueTimeout = time.Second * 1
	// DefaultPollingInterval is the frequency that cloudformation client
	// polls for state changes
	DefaultPollingInterval = time.Second * 5
	// ErrPrincipalNotFound is returned if no Principal (role record) can
	// be found to attach a policy to
	ErrPrincipalNotFound = fmt.Errorf("PRINCIPAL_NOT_FOUND")
	// ErrPrincipalMultipleMatches is returned if a label selector matches
	// multiple principals, which is not currently supported
	ErrPrincipalMultipleMatches = fmt.Errorf("PRINCIPAL_MULTIPLE_MATCHES")
	// ErrMissingKind returned on config error
	ErrMissingKind = fmt.Errorf("MISSING_CONTROLLER_RESOURCE_KIND")
	// ErrMissingPrincipalKind returned on config error
	ErrMissingPrincipalKind = fmt.Errorf("MISSING_CONTROLLER_PRINCIPAL_KIND")
	// ErrMissingLogger returned is logger not present during setup
	ErrMissingLogger = fmt.Errorf("MISSING_CONTROLLER_LOGGER")
	// ErrMissingCloudformationClient returned if no cloudformation client setup
	ErrMissingCloudformationClient = fmt.Errorf("MISSING_CLOUDFORMATION_CLIENT")
	// ErrMissingAWSClient return on config error
	ErrMissingAWSClient = fmt.Errorf("MISSING_AWS_CLIENT")
)

// Controller implements kubernetes controller-runtime reconcile.Reconciler to
// reconcile a kubernetes resource type that implements the Stack interface
// using cloudformation.
//
// It should be initialized with the Kind that it will reconcile (ie
// v1beta1.Postgres{}) any parameters that should always be passed to
// cloudformation during reconciliation.  Setting parameters here is useful for
// variables that are based on the environment rather than the resource options,
// for example which region or vpc to deploy to may not be something you wish
// to make configurable via the resource, but rather globally on the manager
//
// TODO: add example
//
type Controller struct {
	Scheme               *runtime.Scheme        // Scheme is required for operations like gc
	Log                  logr.Logger            // Log will be used to report each reconcile
	KubernetesClient     client.Client          // KubernetesClient is required to talk to api
	CloudFormationClient *Client                // CloudFormationClient is required to talk to aws
	Parameters           []*Parameter           // Parameters are default params always passed to Apply
	Kind                 Stack                  // Kind is the kubernetes resource type to reconcile
	PrincipalListKind    object.PrincipalLister // PrincipalListKind is the type that will be used to lookup role data for StackPolicyAttacher
	ReconcileTimeout     time.Duration          // ReconcileTimeout is the max execution time on Reconcile before requeuing
	RequeueTimeout       time.Duration          // RequeueTimeout is the delay before trying again after ReconcileTimeout is hit
}

// SetupWithManager validates and registers this controller with the manager and api
func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	r.KubernetesClient = mgr.GetClient()
	r.Scheme = mgr.GetScheme()
	// validate and defaults
	if r.Kind == nil {
		return ErrMissingKind
	}
	if r.PrincipalListKind == nil {
		return ErrMissingPrincipalKind
	}
	if r.CloudFormationClient == nil {
		return ErrMissingCloudformationClient
	}
	if r.CloudFormationClient.Client == nil {
		return ErrMissingAWSClient
	}
	if r.CloudFormationClient.PollingInterval == 0 {
		r.CloudFormationClient.PollingInterval = DefaultPollingInterval
	}
	if r.ReconcileTimeout == 0 {
		r.ReconcileTimeout = DefaultReconcileDeadline
	}
	if r.RequeueTimeout == 0 {
		r.RequeueTimeout = DefaultRequeueTimeout
	}
	// ensure that any controller params are set
	// this means we can fail early on missing config
	for _, p := range r.Parameters {
		if p.ParameterKey == nil || *p.ParameterKey == "" {
			return fmt.Errorf("invalid controller parameter: ParameterKey must be set")
		}
		if p.ParameterValue == nil || *p.ParameterValue == "" {
			return fmt.Errorf("missing required controller parameter: %s", *p.ParameterKey)
		}
	}
	// setup logger
	r.Log = ctrl.Log.WithName("controllers").WithName(r.Kind.GetResourceVersion())
	// register with manager
	return ctrl.NewControllerManagedBy(mgr).
		For(r.Kind).
		Complete(r)
}

// +kubebuilder:rbac:groups=queue.govsvc.uk,resources=sqs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=queue.govsvc.uk,resources=sqs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=database.govsvc.uk,resources=postgres,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.govsvc.uk,resources=postgres/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=access.govsvc.uk,resources=principals,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=access.govsvc.uk,resources=principals/status,verbs=get;update;patch

// Reconcile syncronizes state between the resource and a cloudformation stack
func (r *Controller) Reconcile(req ctrl.Request) (res ctrl.Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.ReconcileTimeout)
	defer cancel()
	// execute reconciliation and log changes
	op, err := r.reconcileWithContext(ctx, req)
	if err == context.DeadlineExceeded {
		// ran out of time, most likely waiting on
		// a long running provisioning, come back a bit later
		res.Requeue = true
		res.RequeueAfter = r.RequeueTimeout
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

// reconcileWithContext fetches the resource to reconcile and executes reconcileObjectWithContext and returns if any changes were made
func (r *Controller) reconcileWithContext(ctx context.Context, req ctrl.Request) (controllerutil.OperationResult, error) {
	bg := context.Background()
	o := r.Kind.DeepCopyObject().(Stack)
	if err := r.KubernetesClient.Get(bg, req.NamespacedName, o); apierrs.IsNotFound(err) {
		// nothing we can do if the resource has gone missing, so
		// ignore any not found errors and let the api carry on
		return controllerutil.OperationResultNone, nil
	} else if err != nil {
		// issue communicating with the api
		// return err and we'll retry later
		return controllerutil.OperationResultNone, err
	}
	// track changes to our object resource and call the main reconcile func
	return controllerutil.CreateOrUpdate(bg, r.KubernetesClient, o, func() error {
		return r.reconcileObjectWithContext(ctx, req, o)
	})
}

// reconcileObjectWithContext is the main loop, it will mutate "o" with any changes required
func (r *Controller) reconcileObjectWithContext(ctx context.Context, req ctrl.Request, o Stack) error {
	defer r.Log.Info("reconcileObjectWithContext",
		"o", o,
	)
	// examine DeletionTimestamp to determine if object is under deletion
	if !o.GetDeletionTimestamp().IsZero() {
		// The object is being deleted
		return r.destroyObjectWithContext(ctx, req, o)
	}

	// The object is not being deleted, so we ensure that our finalizer is present
	object.SetFinalizer(o, Finalizer)

	// lookup the iam role name from a Principal resource with labels
	roleParams, err := r.getRoleParams(ctx, o)
	if err != nil {
		return err
	}

	// append any default params
	params := append(r.Parameters, roleParams...)

	// create or update stack as required
	outputs, err := r.CloudFormationClient.Apply(ctx, o, params...)
	if err != nil {
		return err
	}

	// create or update secret
	err = r.updateCredentialsSecret(ctx, o, outputs)
	if err != nil {
		return err
	}

	return nil
}

// updateCredentialsSecret will write any cloudformation outputs to a secret so
// it can be consumed by other kubernetes resources like Pods. Not all stacks
// have outputs that need sharing, only Stacks that have both cloudformation
// outputs defined in their template AND implement object.SecretNamer will
// result in a secret being created.
func (r *Controller) updateCredentialsSecret(ctx context.Context, o Stack, outputs Outputs) error {
	if len(outputs) == 0 {
		// no outputs to write to secret
		return nil
	}
	s, ok := o.(object.SecretNamer)
	if !ok {
		// this objects is not a SecretNamer and does not
		// require writing credentials to a secret
		return nil
	}
	secret := core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.GetSecretName(),
			Namespace: o.GetNamespace(),
		},
		Data: map[string][]byte{},
	}
	secretKey, err := client.ObjectKeyFromObject(&secret)
	if err != nil {
		return err
	}
	err = r.KubernetesClient.Get(ctx, secretKey, &secret)
	if err != nil && !apierrs.IsNotFound(err) {
		return err
	}
	op, err := controllerutil.CreateOrUpdate(ctx, r.KubernetesClient, &secret, func() error {
		secret.Type = core.SecretTypeOpaque
		secret.Annotations = map[string]string{
			"operator": "gsp-service-operator",
			"group":    o.GroupVersionKind().Group,
			"version":  o.GroupVersionKind().Version,
		}
		for key, value := range outputs {
			secret.Data[key] = []byte(value)
		}
		return nil
	})
	r.Log.Info("update-secret",
		"secret", secretKey,
		"op", op,
		"err", err,
	)
	if err != nil {
		return err
	}
	// mark the secret as owned by the o resource so it gets gc'd
	if err := controllerutil.SetControllerReference(o, &secret, r.Scheme); err != nil {
		return err
	}
	return nil
}

// getRoleParams fetches additional IAM role params from any Stack that also
// implements StackPolicyAttacher.  Implementing StackPolicyAttacher tells the
// controller that this Stack wants to attach a policy to an IAM role, and
// provides the methods to locate the object that represents the role and
// extract the required role name.
// FIXME: this method needs to know too much about "Principal". The presence of
//        the PrincipalListKind var , PrincipalLister and Principal interface
//        types hints that this abstraction is failing. which leads to some
//        over complicated implementation of this controller. We should
//        consider either: (1) dropping the "Prinipal" concept and letting this
//        controller target roles based on annotation, (simple but may not work
//        well with multiple CLOUD_PROVIDER values in the future OR (2)
//        targeting "Pods" with a label selector instead of "Principal", and
//        letting a controller manage both the creation of IAMRoles and the
//        assignment of kiam annotations which would mean this controller would
//        only need to know instead an IAMPolicy kubernetes object and have a
//        policy OR (3) something else
func (r *Controller) getRoleParams(ctx context.Context, o Stack) ([]*Parameter, error) {
	serviceThatRequiresIAM, ok := o.(StackPolicyAttacher)
	if !ok {
		return nil, nil
	}
	list := r.PrincipalListKind.DeepCopyObject().(object.PrincipalLister)
	listOptsFunc := func(opts *client.ListOptions) {
		opts.Namespace = serviceThatRequiresIAM.GetNamespace()
		opts.LabelSelector = labels.SelectorFromSet(map[string]string{
			AccessGroupLabel: serviceThatRequiresIAM.GetLabels()[AccessGroupLabel],
		})
	}
	err := r.KubernetesClient.List(ctx, list, listOptsFunc)
	if err != nil {
		return nil, err
	}
	principals := list.GetPrincipals()
	if len(principals) == 0 {
		return nil, ErrPrincipalNotFound
	}
	// FIXME: allow assigning policy to multiple records
	if len(principals) > 1 {
		return nil, ErrPrincipalMultipleMatches
	}
	principal := principals[0]
	// fetch the outputs from the principal's stack
	params, err := serviceThatRequiresIAM.GetStackRoleParameters(principal.GetRoleName())
	if err != nil {
		return nil, err
	}
	return params, nil
}

// destroyObjectWithContext triggers the stack destroy and removes finalizer once done
func (r *Controller) destroyObjectWithContext(ctx context.Context, _ ctrl.Request, o Stack) error {
	if object.HasFinalizer(o, Finalizer) {
		// our finalizer is present, so lets attempt deletion
		err := r.CloudFormationClient.Destroy(ctx, o)
		if err != nil {
			return err
		}
		// delete succeeded so remove finalizer and update
		object.RemoveFinalizer(o, Finalizer)
	}
	return nil
}
