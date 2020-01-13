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

package pipeline

import (
	"context"
	"fmt"
	"os"
	"strings"

	concoursev1beta1 "github.com/alphagov/gsp/components/concourse-operator/pkg/apis/concourse/v1beta1"
	"github.com/concourse/concourse/go-concourse/concourse"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	deleteFinalizer = "pipeline.finalizers.concourse-ci-org"
)

// Add creates a new Pipeline Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, clientFn func(team string) (concourse.Client, error)) error {
	r := &ReconcilePipeline{
		Client:    mgr.GetClient(),
		scheme:    mgr.GetScheme(),
		newClient: clientFn,
	}
	return add(mgr, r)
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("pipeline-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to Pipeline
	err = c.Watch(&source.Kind{Type: &concoursev1beta1.Pipeline{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcilePipeline{}

// ReconcilePipeline reconciles a Pipeline object
type ReconcilePipeline struct {
	client.Client
	scheme    *runtime.Scheme
	newClient func(team string) (concourse.Client, error)
}

// Reconcile reads that state of the cluster for a Pipeline object and makes changes based on the state read
// and what is in the Pipeline.Spec
// +kubebuilder:rbac:groups=concourse.govsvc.uk,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcilePipeline) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Pipeline instance
	instance := &concoursev1beta1.Pipeline{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	namespacePrefix := os.Getenv("CONCOURSE_NAMESPACE_PREFIX")
	teamName := strings.TrimPrefix(instance.ObjectMeta.Namespace, namespacePrefix)
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object.
		if !containsString(instance.ObjectMeta.Finalizers, deleteFinalizer) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, deleteFinalizer)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
		if err := r.update(teamName, instance); err != nil {
			fmt.Println("UPDATE FAILED:", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, err)
			return reconcile.Result{}, err
		}
	} else {
		// The object is being deleted
		if containsString(instance.ObjectMeta.Finalizers, deleteFinalizer) {
			// our finalizer is present, so lets handle our external dependency
			if err := r.destroy(teamName, instance); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				fmt.Println("DESTROY FAILED:", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, err)
				return reconcile.Result{}, err
			}
			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, deleteFinalizer)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	}
	return reconcile.Result{}, nil
}

func (r *ReconcilePipeline) update(teamName string, instance *concoursev1beta1.Pipeline) error {
	pipelineName := instance.ObjectMeta.Name
	var pipelineYAML []byte
	var err error

	if len(instance.Spec.Config.Jobs) > 0 {
		pipelineYAML, err = yaml.Marshal(instance.Spec.Config)
		if err != nil {
			return err
		}
	} else if instance.Spec.PipelineString != "" {
		pipelineYAML = []byte(instance.Spec.PipelineString)
	} else {
		return fmt.Errorf("need to define `config` or `pipelineString` for pipeline '%s' in team '%s'", pipelineName, teamName)
	}

	// create a token client
	// fetch the pipeline yaml version if it exists
	concourseClient, err := r.newClient(teamName)
	if err != nil {
		return err
	}
	_, existingConfigVersion, found, err := concourseClient.Team(teamName).PipelineConfig(pipelineName)
	if found && err != nil {
		return fmt.Errorf("couldn't obtain existing pipeline config: %s", err)
	}
	// set pipeline
	_, _, _, err = concourseClient.Team(teamName).CreateOrUpdatePipelineConfig(pipelineName, existingConfigVersion, pipelineYAML, true)
	if err != nil {
		return fmt.Errorf("couldn't CreateOrUpdatePipelineConfig '%s' for team '%s': %s", pipelineName, teamName, err)
	}

	// set publicity
	if instance.Spec.Paused {
		if _, err = concourseClient.Team(teamName).PausePipeline(pipelineName); err != nil {
			return fmt.Errorf("couldn't pause pipeline '%s' in team '%s': %s", pipelineName, teamName, err)
		}
	} else {
		if _, err := concourseClient.Team(teamName).UnpausePipeline(pipelineName); err != nil {
			return fmt.Errorf("couldn't unpause pipeline '%s' in team '%s': %s", pipelineName, teamName, err)
		}
	}

	// set exposure
	if instance.Spec.Exposed {
		if _, err := concourseClient.Team(teamName).ExposePipeline(pipelineName); err != nil {
			return fmt.Errorf("couldn't expose pipeline '%s' in team '%s': %s", pipelineName, teamName, err)
		}
	} else {
		if _, err := concourseClient.Team(teamName).HidePipeline(pipelineName); err != nil {
			return fmt.Errorf("couldn't hide pipeline '%s' in team '%s': %s", pipelineName, teamName, err)
		}
	}

	fmt.Println("UPDATED", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, string(pipelineYAML))
	return nil
}

func (r *ReconcilePipeline) destroy(teamName string, instance *concoursev1beta1.Pipeline) error {
	concourseClient, err := r.newClient(teamName)
	if err != nil {
		return err
	}
	_, err = concourseClient.Team(teamName).DeletePipeline(instance.ObjectMeta.Name)
	if err != nil {
		return err
	}
	fmt.Println("DESTROYED", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name)
	return nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
