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
}

// +kubebuilder:rbac:groups=database.gsp.k8s.io,resources=postgres,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.gsp.k8s.io,resources=postgres/status,verbs=get;update;patch

func (r *PostgresReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("postgres", req.NamespacedName)

	var postgres database.Postgres
	if err := r.Get(ctx, req.NamespacedName, &postgres); err != nil {
		log.V(1).Info("unable to fetch Postgres - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, ignoreNotFound(err)
	}

	provisioner := os.Getenv("CLOUD_PROVIDER")
	switch provisioner {
	case "aws":
		return r.ReconcileAWS(ctx, req, &postgres)
	default:
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 15}, fmt.Errorf("unsupported cloud provider: %s", provisioner)
	}
}

func (r *PostgresReconciler) ReconcileAWS(ctx context.Context, req ctrl.Request, postgres *database.Postgres) (ctrl.Result, error) {
	finalizerName := "stack.aurora.postgres.database.gsp.k8s.io"
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

	if postgres.ObjectMeta.DeletionTimestamp.IsZero() { // create or update
		template := internalaws.AuroraPostgres(stackName, postgres)
		yaml, err := template.YAML()
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, fmt.Errorf("error serialising template: %s", err)
		}

		describeOutput, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)})
		if err != nil { // create
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

			postgres.Status.ID = createOutput.String()
			postgres.Status.Status = cloudformation.StackStatusCreateInProgress
			if err := r.Status().Update(ctx, postgres); err != nil {
				log.Error(err, "unable to update Postgres status")
				return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
			}

			// if err := svc.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}); err != nil {
			// 	return ctrl.Result{}, fmt.Errorf("error waiting for stack to be created: %s", err)
			// }

			// postgres.Status.Status = fmt.Sprintf("stack-%s", StateAvailable)
			// if err := r.Status().Update(ctx, postgres); err != nil {
			// 	log.Error(err, "unable to update Postgres status")
			// 	return ctrl.Result{}, err
			// }
		} else if *describeOutput.Stacks[0].StackStatus == cloudformation.StackStatusCreateComplete { // update
			// TODO: get Secret

			log.V(1).Info("updating stack...", "stackName", stackName)

			_, err = svc.UpdateStack(&cloudformation.UpdateStackInput{
				TemplateBody: aws.String(string(yaml)),
				StackName:   aws.String(stackName),
				Parameters: []*cloudformation.Parameter{
					&cloudformation.Parameter{
						ParameterKey:   aws.String("MasterUsername"),
						ParameterValue: aws.String(""), // TODO: use Secret
					},
					&cloudformation.Parameter{
						ParameterKey:   aws.String("MasterPassword"),
						ParameterValue: aws.String(""), // TODO: use Secret
					},
				},
			})
			if err != nil {
				return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error updating stack: %s", err)
			}

			// postgres.Status.ID = createOutput.String()
			// postgres.Status.Status = fmt.Sprintf("stack-%s", StateCreating)
			// if err := r.Status().Update(ctx, postgres); err != nil {
			// 	log.Error(err, "unable to update Postgres status")
			// 	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
			// }
		}

		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
	} else { // delete
		if containsString(postgres.ObjectMeta.Finalizers, finalizerName) {
			log.V(1).Info("deleting stack...", "stackName", stackName)
			if _, err := svc.DeleteStack(&cloudformation.DeleteStackInput{StackName: aws.String(stackName)}); err != nil {
				return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, fmt.Errorf("error deleting stack: %s", err)
			}

			// TODO: delete Secret

			postgres.ObjectMeta.Finalizers = removeString(postgres.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), postgres); err != nil {
				return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
			}
		}

		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
	}

	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, nil
}

func (r *PostgresReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&database.Postgres{}).
		Complete(r)
}
