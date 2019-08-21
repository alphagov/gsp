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
	"os"
	"time"

	"github.com/alphagov/gsp/components/service-operator/internal"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
)

// PrincipalReconciler reconciles a Principal object
type PrincipalReconciler struct {
	client.Client
	Log                      logr.Logger
	CloudFormationReconciler internalaws.CloudFormationReconciler
	ClusterName              string
	RolePrincipal            string
	PermissionsBoundary      string
	principal                access.Principal
}

const (
	PrincipalFinalizerName = "stack.principal.access.govsvc.uk"
)

// +kubebuilder:rbac:groups=access.govsvc.uk,resources=principals,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=access.govsvc.uk,resources=principals/status,verbs=get;update;patch

func (r *PrincipalReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("access", req.NamespacedName)

	var principal access.Principal
	if err := r.Get(ctx, req.NamespacedName, &principal); err != nil {
		log.V(1).Info("unable to fetch Principal Resource - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, internal.IgnoreNotFound(err)
	}

	provisioner := os.Getenv("CLOUD_PROVIDER")
	switch provisioner {
	case "aws":
		roleName := fmt.Sprintf("%s-%s-%s", r.ClusterName, req.Namespace, principal.ObjectMeta.Name)

		principalCloudFormationTemplate := internalaws.IAMRole{
			RoleConfig:          &principal,
			RoleName:            roleName,
			RolePrincipal:       r.RolePrincipal,
			PermissionsBoundary: r.PermissionsBoundary,
		}
		action, stackData, err := r.CloudFormationReconciler.Reconcile(ctx, log, req, &principalCloudFormationTemplate, !principal.ObjectMeta.DeletionTimestamp.IsZero())
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 2}, err
		}
		principal.Status.ID = stackData.ID
		principal.Status.Status = stackData.Status
		principal.Status.Reason = stackData.Reason
		principal.Status.ARN = string(internalaws.ValueFromOutputs(internalaws.IAMRoleARN, stackData.Outputs))

		for _, event := range stackData.Events {
			principal.Status.Events = append(principal.Status.Events, access.Event{
				Status: *event.ResourceStatus,
				Reason: *event.ResourceStatusReason,
				Time:   &metav1.Time{Time: *event.Timestamp},
			})
		}

		backoff := ctrl.Result{Requeue: true, RequeueAfter: time.Minute}

		switch action {
		case internal.Create:
			principal.ObjectMeta.Finalizers = append(principal.ObjectMeta.Finalizers, PrincipalFinalizerName)
			return backoff, r.Update(ctx, &principal)
		case internal.Delete:
			principal.ObjectMeta.Finalizers = internal.RemoveString(principal.ObjectMeta.Finalizers, PrincipalFinalizerName)
			return backoff, r.Update(ctx, &principal)
		default:
			return backoff, r.Update(ctx, &principal)
		}

	default:
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 15}, fmt.Errorf("unsupported cloud provider: %s", provisioner)
	}
}

func (r *PrincipalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&access.Principal{}).
		Complete(r)
}
