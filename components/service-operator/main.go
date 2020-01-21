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

package main

import (
	"flag"
	"log"

	accessv1beta1 "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	databasev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	queuev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	storagev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/storage/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/controllers"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk"
	"github.com/alphagov/gsp/components/service-operator/internal/istio"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = istio.AddToScheme(scheme)

	_ = databasev1beta1.AddToScheme(scheme)
	_ = queuev1beta1.AddToScheme(scheme)
	_ = accessv1beta1.AddToScheme(scheme)
	_ = storagev1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func run() error {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		return err
	}

	c := sdk.NewClient()
	controllers := []controllers.Controller{
		controllers.PrincipalCloudFormationController(c),
		controllers.PostgresCloudFormationController(c),
		controllers.SQSCloudFormationController(c),
		controllers.S3CloudFormationController(c),
		controllers.ImageRepositoryCloudFormationController(c),
	}

	for _, c := range controllers {
		if err := c.SetupWithManager(mgr); err != nil {
			return err
		}
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
