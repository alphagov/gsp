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
	"os"

	accessv1beta1 "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	databasev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	queuev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/controllers"
	internalaws "github.com/alphagov/gsp/components/service-operator/internal/aws"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = databasev1beta1.AddToScheme(scheme)
	_ = queuev1beta1.AddToScheme(scheme)
	_ = accessv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var clusterName string
	var kiamServerRole string
	var rolePermissionsBoundary string
	var rdsFromWorkerSecurityGroup string
	var dbSubnetGroup string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&clusterName, "cluster", "", "The name of the k8s cluster")
	flag.StringVar(&kiamServerRole, "kiam-server-role-arn", "", "The ARN of the kiam server role")
	flag.StringVar(&rolePermissionsBoundary, "role-permissions-boundary-arn", "", "The ARN of the permissions boundary to apply to created IAM roles")
	flag.StringVar(&rdsFromWorkerSecurityGroup, "rds-from-worker-security-group", "", "The name of the security group to apply to created RDS databases")
	flag.StringVar(&dbSubnetGroup, "db-subnet-group", "", "The name of the DB Subnet Group to apply to created RDS instances")
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
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// FIXME: how to handle cluster name, shouldn't be both flag and envvar
	if clusterName != "" {
		os.Setenv("CLUSTER_NAME", clusterName)
	}

	cloudFormationController := internalaws.CloudFormationController{
		ClusterName: clusterName,
	}

	if err = (&controllers.PostgresReconciler{
		Client:                   mgr.GetClient(),
		Log:                      ctrl.Log.WithName("controllers").WithName("Postgres"),
		CloudFormationReconciler: &cloudFormationController,
		SecurityGroup:            rdsFromWorkerSecurityGroup,
		DBSubnetGroup:            dbSubnetGroup,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Postgres")
		os.Exit(1)
	}
	if err = (&controllers.SQSReconciler{
		Client:                   mgr.GetClient(),
		Log:                      ctrl.Log.WithName("controllers").WithName("SQS"),
		CloudFormationReconciler: &cloudFormationController,
		ClusterName:              clusterName,
		Provisioner:              os.Getenv("CLOUD_PROVIDER"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SQS")
		os.Exit(1)
	}
	if err = (&controllers.PrincipalReconciler{
		Client:                   mgr.GetClient(),
		Log:                      ctrl.Log.WithName("controllers").WithName("Principal"),
		CloudFormationReconciler: &cloudFormationController,
		ClusterName:              clusterName,
		RolePrincipal:            kiamServerRole,
		PermissionsBoundary:      rolePermissionsBoundary,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SQS")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
