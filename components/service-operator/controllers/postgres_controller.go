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
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	access "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	database "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
)

// PostgresReconciler reconciles a Postgres object
type PostgresReconciler struct {
	client.Client
	Log                      logr.Logger
	CloudFormationReconciler internalaws.CloudFormationReconciler
}

const (
	PostgresFinalizerName = "stack.aurora.postgres.database.govsvc.uk"
)

// +kubebuilder:rbac:groups=database.govsvc.uk,resources=postgres,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.govsvc.uk,resources=postgres/status,verbs=get;update;patch

func (r *PostgresReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("postgres", req.NamespacedName)

	var postgres database.Postgres
	if err := r.Get(ctx, req.NamespacedName, &postgres); err != nil {
		log.V(1).Info("unable to fetch Postgres Resource - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, internal.IgnoreNotFound(err)
	}

	var secret core.Secret
	secretName := internal.CoalesceString(postgres.Spec.Secret, postgres.Name)
	if err := r.Get(ctx, k8stypes.NamespacedName{Name: secretName, Namespace: req.Namespace}, &secret); internal.IgnoreNotFound(err) != nil {
		log.V(1).Info("unable to fetch Postgres Secret - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, err
	}

	provisioner := os.Getenv("CLOUD_PROVIDER")
	switch provisioner {
	case "aws":
		var roles access.PrincipalList
		err := r.List(ctx, &roles, client.MatchingLabels(map[string]string{access.AccessGroupLabel: postgres.Labels[access.AccessGroupLabel]}))
		if err != nil || len(roles.Items) != 1 {
			log.V(1).Info("unable to find unique IAM Role in same gsp-access-group - waiting 5 minutes", "gsp-access-group", postgres.Labels[access.AccessGroupLabel])
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 2}, err
		}

		postgresCloudFormationTemplate := internalaws.AuroraPostgres{PostgresConfig: &postgres, IAMRoleName: roles.Items[0].Status.Name}
		action, stackData, err := r.CloudFormationReconciler.Reconcile(ctx, log, req, &postgresCloudFormationTemplate, !postgres.ObjectMeta.DeletionTimestamp.IsZero())
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 2}, err
		}
		newSecret := postgresOutputsToSecret(secretName, req.Namespace, stackData.Outputs)

		postgres.Status.ID = stackData.ID
		postgres.Status.Status = stackData.Status
		postgres.Status.Reason = stackData.Reason

		events := []database.Event{}
		for _, event := range stackData.Events {
			reason := "-"
			if event.ResourceStatusReason != nil {
				reason = *event.ResourceStatusReason
			}
			events = append(events, database.Event{
				Status: *event.ResourceStatus,
				Reason: reason,
				Time:   &metav1.Time{Time: *event.Timestamp},
			})
		}
		postgres.Status.Events = events

		backoff := ctrl.Result{Requeue: true, RequeueAfter: time.Minute}

		switch action {
		case internal.Create:
			postgres.ObjectMeta.Finalizers = append(postgres.ObjectMeta.Finalizers, PostgresFinalizerName)
			err := r.Update(ctx, &postgres)
			if err != nil {
				return backoff, err
			}

			return backoff, r.Create(ctx, &newSecret)
		case internal.Update:
			err := r.Update(ctx, &newSecret)
			if err != nil {
				return backoff, err
			}

			return backoff, r.Update(ctx, &postgres)
		case internal.Delete:
			postgres.ObjectMeta.Finalizers = internal.RemoveString(postgres.ObjectMeta.Finalizers, PostgresFinalizerName)
			err := r.Update(ctx, &postgres)
			if err != nil {
				return backoff, err
			}

			return backoff, r.Delete(ctx, &secret)
		default:
			return backoff, r.Update(ctx, &postgres)
		}

	default:
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 15}, fmt.Errorf("unsupported cloud provider: %s", provisioner)
	}
}

func (r *PostgresReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&database.Postgres{}).
		Complete(r)
}

func postgresOutputsToSecret(secretName, namespace string, outputs []*cloudformation.Output) core.Secret {
	return core.Secret{
		Type: core.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Annotations: map[string]string{
				"operator": "gsp-service-operator",
				"group":    database.GroupVersion.Group,
				"version":  database.GroupVersion.Version,
			},
		},
		Data: map[string][]byte{
			"Endpoint": internalaws.ValueFromOutputs(internalaws.PostgresEndpoint, outputs),
			"Port":     internalaws.ValueFromOutputs(internalaws.PostgresPort, outputs),
			"DBName":   internalaws.ValueFromOutputs(internalaws.PostgresDBName, outputs),
			"Engine":   internalaws.ValueFromOutputs(internalaws.PostgresEngine, outputs),
		},
	}
}
