/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/alphagov/gsp/components/service-operator/internal"

	core "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/alphagov/gsp/components/service-operator/apis"
	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
)

var (
	DefaultReconcileDeadline = time.Minute * 15
)

// CloudFormationReconciler reconciles resources that implement the StackResource interface
type CloudFormationReconciler struct {
	Scheme               *runtime.Scheme
	Log                  logr.Logger
	KubernetesClient     client.Client
	CloudFormationClient *internalaws.CloudFormationClient
	ClusterName          string
	Kind                 apis.StackObject
	ReconcileTimeout     time.Duration
	RequeueTimeout       time.Duration
}

const (
	Finalizer = "cloudformatiton.finalizers.govsvc.uk"
)

func (r *CloudFormationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.Kind).
		Complete(r)
}

// +kubebuilder:rbac:groups=queue.govsvc.uk,resources=sqs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=queue.govsvc.uk,resources=sqs/status,verbs=get;update;patch

func (r *CloudFormationReconciler) Reconcile(req ctrl.Request) (res ctrl.Result, err error) {
	timeout := DefaultReconcileDeadline
	if r.ReconcileTimeout != 0 {
		timeout = r.ReconcileTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	o := r.Kind.DeepCopyObject().(apis.StackObject)
	if err := r.KubernetesClient.Get(context.Background(), req.NamespacedName, o); err != nil {
		// nothing we can do if the resource has gone missing, so
		// ignore any not found errors and let the api carry on
		return ctrl.Result{}, internal.IgnoreNotFound(err)
	}
	op, err := controllerutil.CreateOrUpdate(context.Background(), r.KubernetesClient, o, func() error {
		err := r.reconcile(ctx, req, o)
		if err == context.DeadlineExceeded {
			// ran out of time, most likely waiting on
			// a long running provisioning, come back a bit later
			res.Requeue = true
			res.RequeueAfter = r.RequeueTimeout
			return nil
		}
		return err
	})
	r.Log.Info("reconciled",
		"service", req.NamespacedName,
		"requeue", res.Requeue,
		"after", res.RequeueAfter,
		"op", op,
		"err", err,
	)
	return res, err
}

func (r *CloudFormationReconciler) reconcile(ctx context.Context, req ctrl.Request, o apis.StackObject) error {
	// examine DeletionTimestamp to determine if object is under deletion
	finalizers := o.GetFinalizers()
	if !o.GetDeletionTimestamp().IsZero() {
		// The object is being deleted
		if internal.ContainsString(finalizers, Finalizer) {
			// our finalizer is present, so lets attempt deletion
			err := r.CloudFormationClient.Destroy(ctx, o)
			if err != nil {
				return err
			}
			// delete succeeded so remove finalizer and update
			o.SetFinalizers(internal.RemoveString(finalizers, Finalizer))
		}
		return nil
	}

	// lookup the target iam role name from the target principal
	roleName, err := r.getRoleName(ctx, o)
	if err != nil {
		return err
	}

	// The object is not being deleted, so if it does not have our finalizer,
	// then lets register our finalizer and update the object immediately.
	if !internal.ContainsString(finalizers, Finalizer) {
		o.SetFinalizers(append(finalizers, Finalizer))
	}

	// create or update stack as required
	outputs, err := r.CloudFormationClient.Apply(ctx, o, &awscloudformation.Parameter{
		ParameterKey:   aws.String(queue.IAMRoleParameterName),
		ParameterValue: aws.String(roleName),
	})
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

func (r *CloudFormationReconciler) updateCredentialsSecret(ctx context.Context, o apis.StackObject, outputs []*awscloudformation.Output) error {
	secret := core.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      o.GetSecretName(),
		Namespace: o.GetNamespace(),
	}}
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
			"group":    queue.GroupVersion.Group,
			"version":  queue.GroupVersion.Version,
		}
		secret.Data = map[string][]byte{ //TODO this should be from a generic o.StackOutputs()
			"QueueURL": internalaws.ValueFromOutputs(queue.SQSOutputURL, outputs),
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

func (r *CloudFormationReconciler) getRoleName(ctx context.Context, m metav1.Object) (string, error) {
	var roles access.PrincipalList
	listOptsFunc := func(opts *client.ListOptions) {
		opts.Namespace = m.GetNamespace()
		opts.LabelSelector = labels.SelectorFromSet(map[string]string{
			access.AccessGroupLabel: m.GetLabels()[access.AccessGroupLabel],
		})
	}
	err := r.KubernetesClient.List(ctx, &roles, listOptsFunc)
	if err != nil {
		return "", err
	}
	if len(roles.Items) != 1 {
		return "", fmt.Errorf("PRINCIPAL_NOT_FOUND")
	}
	roleName := roles.Items[0].Status.Name
	if roleName == "" {
		return "", fmt.Errorf("PRINCIPAL_HAS_NO_ROLE_NAME_YET")
	}
	return roleName, nil
}
