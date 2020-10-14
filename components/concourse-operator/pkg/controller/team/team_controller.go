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

package team

import (
	"context"
	"fmt"
	"os"
	"strings"

	concoursev1beta1 "github.com/alphagov/gsp/components/concourse-operator/pkg/apis/concourse/v1beta1"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

var MainTeam = "main"

const (
	deleteFinalizer = "pipeline.finalizers.concourse-ci-org"
)

// Add creates a new Team Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, clientFn func(team string) (concourse.Client, error)) error {
	r := &ReconcileTeam{
		Client:    mgr.GetClient(),
		scheme:    mgr.GetScheme(),
		newClient: clientFn,
	}
	return add(mgr, r)
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("team-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to Team
	err = c.Watch(&source.Kind{Type: &concoursev1beta1.Team{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileTeam{}

// ReconcileTeam reconciles a Teams object
type ReconcileTeam struct {
	client.Client
	scheme    *runtime.Scheme
	newClient func(team string) (concourse.Client, error)
}

// Reconcile reads that state of the cluster for a Team object and makes changes based on the state read
// Automatically generate RBAC rules to allow the Controller to read and write Team resources
// +kubebuilder:rbac:groups=concourse.govsvc.uk,resources=teams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=concourse.govsvc.uk,resources=teams/status,verbs=get;update;patch
func (r *ReconcileTeam) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	if r.newClient == nil {
		return reconcile.Result{}, fmt.Errorf("newClient is undefined")
	}

	instance := &concoursev1beta1.Team{} // Fetch the Team instance
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err // Error reading the object - requeue the request
	}

	namespacePrefix := os.Getenv("CONCOURSE_NAMESPACE_PREFIX")
	nameFromNamespace := strings.TrimPrefix(instance.ObjectMeta.Namespace, namespacePrefix)
	if nameFromNamespace != instance.ObjectMeta.Name {
		return reconcile.Result{}, fmt.Errorf("Team name %s does not match namespace name %s (full namespace is %s)", instance.ObjectMeta.Name, nameFromNamespace, instance.ObjectMeta.Namespace)
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object.
		if !containsString(instance.ObjectMeta.Finalizers, deleteFinalizer) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, deleteFinalizer)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
		if err := r.update(instance); err != nil {
			fmt.Println("UPDATE FAILED:", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, err)
			return reconcile.Result{}, err
		}
	} else {
		// The object is being deleted
		if containsString(instance.ObjectMeta.Finalizers, deleteFinalizer) {
			// our finalizer is present, so lets handle our external dependency
			if err := r.destroy(instance); err != nil {
				// we only log the error here not block deleteion of the resource as there are some
				// limitations around deleting teams that are hard to reconcile (must always be one admin team for example)
				fmt.Println("DESTROY FAILED:", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, err)
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

func (r *ReconcileTeam) update(instance *concoursev1beta1.Team) error {

	team := atc.Team{
		Name: instance.ObjectMeta.Name,
		Auth: atc.TeamAuth{},
	}

	for _, role := range instance.Spec.Roles {
		users := []string{}
		groups := []string{}

		for _, user := range role.Github.Users {
			if user != "" {
				users = append(users, "github:"+strings.ToLower(user))
			}
		}
		for _, group := range role.Github.Teams {
			if group != "" {
				groups = append(groups, "github:"+strings.ToLower(group))
			}
		}
		for _, user := range role.Local.Users {
			if user != "" {
				users = append(users, "local:"+strings.ToLower(user))
			}
		}

		team.Auth[role.Name] = map[string][]string{
			"users":  users,
			"groups": groups,
		}
	}

	concourseClient, err := r.newClient(MainTeam)
	if err != nil {
		return err
	}
	_, _, _, _, err = concourseClient.Team(team.Name).CreateOrUpdate(team)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileTeam) destroy(instance *concoursev1beta1.Team) error {
	if instance.ObjectMeta.Name == MainTeam {
		return nil // deleting the main team is not possible
	}
	concourseClient, err := r.newClient(MainTeam)
	if err != nil {
		return err
	}
	err = concourseClient.Team(instance.ObjectMeta.Name).DestroyTeam(instance.ObjectMeta.Name)
	if err != nil {
		return err
	}
	fmt.Println("DESTROYED", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name)
	return nil
}

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
