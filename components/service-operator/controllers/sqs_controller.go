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

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
)

// SQSReconciler reconciles a SQS object
type SQSReconciler struct {
	client.Client
	Log                      logr.Logger
	CloudFormationReconciler internalaws.CloudFormationReconciler
	sqs                      queue.SQS
}

const (
	SQSFinalizerName = "stack.sqs.queue.queue.gsp.k8s.io"
)

// +kubebuilder:rbac:groups=queue.gsp.k8s.io,resources=sqs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=queue.gsp.k8s.io,resources=sqs/status,verbs=get;update;patch

func (r *SQSReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("sqs", req.NamespacedName)

	var sqs queue.SQS
	if err := r.Get(ctx, req.NamespacedName, &sqs); err != nil {
		log.V(1).Info("unable to fetch SQS Resource - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, internal.IgnoreNotFound(err)
	}

	var secret core.Secret
	secretName := internal.CoalesceString(sqs.Spec.Secret, sqs.Name)
	if err := r.Get(ctx, k8stypes.NamespacedName{Name: secretName, Namespace: req.Namespace}, &secret); internal.IgnoreNotFound(err) != nil {
		log.V(1).Info("unable to fetch SQS Secret - waiting 5 minutes")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, err
	}

	provisioner := os.Getenv("CLOUD_PROVIDER")
	switch provisioner {
	case "aws":
		sqsCloudFormationTemplate := internalaws.SQS{SQSConfig: &sqs}
		action, stackData, err := r.CloudFormationReconciler.Reconcile(ctx, log, req, &sqsCloudFormationTemplate, !sqs.ObjectMeta.DeletionTimestamp.IsZero())
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 2}, err
		}
		newSecret := sqsOutputsToSecret(secretName, req.Namespace, stackData.Outputs)
		sqs.Status.ID = stackData.ID
		sqs.Status.Status = stackData.Status
		sqs.Status.Reason = stackData.Reason

		for _, event := range stackData.Events {
			sqs.Status.Events = append(sqs.Status.Events, queue.Event{
				Status: *event.ResourceStatus,
				Reason: *event.ResourceStatusReason,
				Time:   event.Timestamp,
			})
		}

		backoff := ctrl.Result{Requeue: true, RequeueAfter: time.Minute}

		switch action {
		case internal.Create:
			sqs.ObjectMeta.Finalizers = append(sqs.ObjectMeta.Finalizers, SQSFinalizerName)
			err := r.Update(ctx, &sqs)
			if err != nil {
				return backoff, err
			}

			return backoff, r.Create(ctx, &newSecret)
		case internal.Update:
			err := r.Update(ctx, &newSecret)
			if err != nil {
				return backoff, err
			}

			return backoff, r.Update(ctx, &sqs)
		case internal.Delete:
			sqs.ObjectMeta.Finalizers = internal.RemoveString(sqs.ObjectMeta.Finalizers, SQSFinalizerName)
			err := r.Update(ctx, &sqs)
			if err != nil {
				return backoff, err
			}

			return backoff, r.Delete(ctx, &secret)
		default:
			return backoff, r.Update(ctx, &sqs)
		}

	default:
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 15}, fmt.Errorf("unsupported cloud provider: %s", provisioner)
	}
}

func (r *SQSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&queue.SQS{}).
		Complete(r)
}

func sqsOutputsToSecret(secretName, namespace string, outputs []*cloudformation.Output) core.Secret {
	return core.Secret{
		Type: core.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Annotations: map[string]string{
				"operator": "gsp-service-operator",
				"group":    queue.GroupVersion.Group,
				"version":  queue.GroupVersion.Version,
			},
		},
		Data: map[string][]byte{
			"QueueURL": internalaws.ValueFromOutputs(internalaws.SQSOutputURL, outputs),
		},
	}
}
