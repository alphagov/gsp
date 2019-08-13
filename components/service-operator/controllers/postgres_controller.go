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

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	database "github.com/alphagov/gsp/components/service-operator/api/v1beta1"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// PostgresReconciler reconciles a Postgres object
type PostgresReconciler struct {
	client.Client
	Log         logr.Logger
	ClusterName string
	secretName  string
	postgres    database.Postgres
}

const (
	finalizerName = "stack.aurora.postgres.database.gsp.k8s.io"
)

var (
	nonUpdatable = []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusReviewInProgress,
		cloudformation.StackStatusDeleteComplete,
	}
)

// +kubebuilder:rbac:groups=database.gsp.k8s.io,resources=postgres,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.gsp.k8s.io,resources=postgres/status,verbs=get;update;patch

func (r *PostgresReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("postgres", req.NamespacedName)

	if err := r.Get(ctx, req.NamespacedName, &r.postgres); err != nil {
		log.V(1).Info("unable to fetch Postgres - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, ignoreNotFound(err)
	}

	provisioner := os.Getenv("CLOUD_PROVIDER")
	switch provisioner {
	case "aws":
		return r.ReconcileAWS(ctx, req)
	default:
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 15}, fmt.Errorf("unsupported cloud provider: %s", provisioner)
	}
}

func (r *PostgresReconciler) ReconcileAWS(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	stackName := fmt.Sprintf("%s-%s-%s-%s", r.ClusterName, "gsp-service-operator-postgres", req.NamespacedName.Namespace, req.NamespacedName.Name)
	// secretName := coalesceString(postgres.Spec.Secret, postgres.Name)

	log := r.Log.WithValues("aws-postgres", req.NamespacedName)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// metadata := ec2metadata.New(sess)
	// awsRegion, err := metadata.Region()
	// if err != nil {
	// 	return ctrl.Result{}, fmt.Errorf("unable to get AWS region from metdata: %s", err)
	// }
	awsRegion := "eu-west-2"
	sess.Config.Region = aws.String(awsRegion)

	svc := cloudformation.New(sess, aws.NewConfig())

	stackID, stackStatus, stackStatusReason, stackExists := r.getCloudFormationStackStatus(svc, stackName)
	r.postgres.Status.ID = stackID
	r.postgres.Status.Status = stackStatus
	r.postgres.Status.Reason = stackStatusReason
	r.updatePostgres()

	if !r.postgres.ObjectMeta.DeletionTimestamp.IsZero() {
		// The resource needs deleting
		if !stackExists || stackStatus == cloudformation.StackStatusDeleteComplete {
			return r.markAsDeleted()
		}

		if stackStatus == cloudformation.StackStatusDeleteInProgress {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 2}, nil
		}

		return r.deleteCloudFormationStack(svc, stackName, log)
	}

	template := internalaws.AuroraPostgres(stackName, &r.postgres)
	yaml, err := template.YAML()
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, fmt.Errorf("error serialising template: %s", err)
	}

	if !stackExists { // create
		return r.createCloudFormationStack(yaml, svc, stackName, log)
	} else if !isInList(stackStatus, nonUpdatable...) { // update
		return r.updateCloudFormationStack(yaml, svc, stackName, log)
	}

	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
}

func (r *PostgresReconciler) updatePostgres() {
	if err := r.Update(context.Background(), &r.postgres); err != nil {
		r.Log.Error(err, "unable to update Postgres status")
	}
}

func (r *PostgresReconciler) getCloudFormationStackStatus(svc *cloudformation.CloudFormation, stackName string) (string, string, string, bool) {
	describeOutput, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)})
	if err != nil {
		return "", "", "", false
	}
	return *describeOutput.Stacks[0].StackId, *describeOutput.Stacks[0].StackStatus, "NoReasonGiven", true
}

func (r *PostgresReconciler) createCloudFormationStack(yaml []byte, svc *cloudformation.CloudFormation, stackName string, log logr.Logger) (ctrl.Result, error) {
	log.V(1).Info("creating stack...", "stackName", stackName)
	username, err := randomString(16, charactersUpper, charactersLower)
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error generating username: %s", err)
	}
	password, err := randomString(32, charactersUpper, charactersLower, charactersNumeric, charactersSpecial)
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error generating password: %s", err)
	}

	createOutput, err := svc.CreateStack(&cloudformation.CreateStackInput{
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
		Parameters: []*cloudformation.Parameter{
			&cloudformation.Parameter{
				ParameterKey:   aws.String("MasterUsername"),
				ParameterValue: aws.String(username),
			},
			&cloudformation.Parameter{
				ParameterKey:   aws.String("MasterPassword"),
				ParameterValue: aws.String(password),
			},
		},
	})
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error creating stack: %s", err)
	}

	// TODO: create Secret

	r.postgres.Status.ID = *createOutput.StackId
	r.postgres.Status.Status = cloudformation.StackStatusCreateInProgress
	r.postgres.ObjectMeta.Finalizers = append(r.postgres.ObjectMeta.Finalizers, finalizerName)
	r.updatePostgres()

	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
}

func (r *PostgresReconciler) updateCloudFormationStack(yaml []byte, svc *cloudformation.CloudFormation, stackName string, log logr.Logger) (ctrl.Result, error) {
	// TODO: get Secret

	log.V(1).Info("updating stack...", "stackName", stackName)

	_, err := svc.UpdateStack(&cloudformation.UpdateStackInput{
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
		Parameters: []*cloudformation.Parameter{
			&cloudformation.Parameter{
				ParameterKey:   aws.String("MasterUsername"),
				ParameterValue: aws.String("qwertyuiop"), // TODO: use Secret
			},
			&cloudformation.Parameter{
				ParameterKey:   aws.String("MasterPassword"),
				ParameterValue: aws.String("qwertyuiop1234567890"), // TODO: use Secret
			},
		},
	})
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error updating stack: %s", err)
	}

	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
}

func (r *PostgresReconciler) deleteCloudFormationStack(svc *cloudformation.CloudFormation, stackName string, log logr.Logger) (ctrl.Result, error) {
	if containsString(r.postgres.ObjectMeta.Finalizers, finalizerName) {
		log.V(1).Info("deleting stack...", "stackName", stackName)
		if _, err := svc.DeleteStack(&cloudformation.DeleteStackInput{StackName: aws.String(stackName)}); err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error deleting stack: %s", err)
		}
		// TODO: delete Secret
	}
	log.V(1).Info("no finalizers found for stack...", "stackName", stackName)
	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
}

func (r *PostgresReconciler) markAsDeleted() (ctrl.Result, error) {
	r.postgres.ObjectMeta.Finalizers = removeString(r.postgres.ObjectMeta.Finalizers, finalizerName)
	if err := r.Update(context.Background(), &r.postgres); err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
	}
	return ctrl.Result{}, nil
}

func (r *PostgresReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&database.Postgres{}).
		Complete(r)
}
