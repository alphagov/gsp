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
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	database "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
)

// PostgresReconciler reconciles a Postgres object
type PostgresReconciler struct {
	client.Client
	Log                      logr.Logger
	CloudFormationController internalaws.CloudFormationController
	postgres                 database.Postgres
}

// +kubebuilder:rbac:groups=database.gsp.k8s.io,resources=postgres,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.gsp.k8s.io,resources=postgres/status,verbs=get;update;patch

func (r *PostgresReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	finalizerName := "stack.aurora.postgres.database.gsp.k8s.io"
	ctx := context.Background()
	log := r.Log.WithValues("postgres", req.NamespacedName)

	var postgres database.Postgres
	if err := r.Get(ctx, req.NamespacedName, &postgres); err != nil {
		log.V(1).Info("unable to fetch Postgres - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, internal.IgnoreNotFound(err)
	}

	provisioner := os.Getenv("CLOUD_PROVIDER")
	switch provisioner {
	case "aws":
		postgresCloudFormationTemplate := internalaws.AuroraPostgres{PostgresConfig: &postgres}
		action, stackData, err := r.CloudFormationController.Reconcile(log, ctx, req, &postgresCloudFormationTemplate, !postgres.ObjectMeta.DeletionTimestamp.IsZero())
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 2}, err
		}
		postgres.Status.ID = stackData.ID
		postgres.Status.Status = stackData.Status
		postgres.Status.Reason = stackData.Reason

		backoff := ctrl.Result{Requeue: true, RequeueAfter: time.Minute}

		switch action {
		case internal.Create:
			postgres.ObjectMeta.Finalizers = append(postgres.ObjectMeta.Finalizers, finalizerName)
			return backoff, r.Update(context.Background(), &postgres)
		case internal.Delete:
			postgres.ObjectMeta.Finalizers = internal.RemoveString(postgres.ObjectMeta.Finalizers, finalizerName)
			return backoff, r.Update(context.Background(), &postgres)
		default:
			return backoff, r.Update(context.Background(), &postgres)
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
