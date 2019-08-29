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

package controllers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	accessv1beta1 "github.com/alphagov/gsp/components/service-operator/apis/access/v1beta1"
	databasev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/database/v1beta1"
	queue "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	queuev1beta1 "github.com/alphagov/gsp/components/service-operator/apis/queue/v1beta1"
	"github.com/alphagov/gsp/components/service-operator/controllers"
	"github.com/alphagov/gsp/components/service-operator/internal/aws"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/awsfakes"
	"github.com/alphagov/gsp/components/service-operator/internal/controllertest"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var log logr.Logger

// var principalReconciler *controllertest.ReconcilerWrapper
var sqsReconciler *controllertest.ReconcilerWrapper

// var postgresReconciler *controllertest.ReconcilerWrapper
var ctx context.Context
var mgrStopChan = make(chan struct{})
var fakeAWSClient *awsfakes.FakeAWSClient

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	ctx = context.Background()

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{envtest.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	os.Setenv("CLOUD_PROVIDER", "aws")
	os.Setenv("CLUSTER_NAME", "xxx")

	log = zap.LoggerTo(GinkgoWriter, true)
	logf.SetLogger(log)

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = databasev1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = queuev1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = accessv1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	By("starting control plane")
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	k8sClient = mgr.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	// wait for control plane to be happy
	Eventually(func() error {
		return k8sClient.List(ctx, &core.SecretList{})
	}, time.Second*20).Should(Succeed())

	// postgresCloudformationReconciler = internalawsmocks.NewMockCloudFormationReconciler(mockCtrl)
	// postgresReconciler = &controllertest.ReconcilerWrapper{
	// 	Reconciler: &controllers.PostgresReconciler{
	// 		Client:                   k8sClient,
	// 		Log:                      ctrl.Log.WithName("controllers").WithName("Postgres"),
	// 		CloudFormationReconciler: postgresCloudformationReconciler,
	// 	},
	// }
	// err = postgresReconciler.SetupWithManager(mgr, &database.Postgres{})
	// Expect(err).ToNot(HaveOccurred())

	fakeAWSClient = awsfakes.NewFakeAWSClient(nil)

	logger := log.WithName("controller-runtime").WithName("controller")

	sqsReconciler = &controllertest.ReconcilerWrapper{
		Reconciler: &controllers.SQSReconciler{
			Kind:        &queue.SQS{},
			Client:      k8sClient,
			Scheme:      mgr.GetScheme(),
			ClusterName: os.Getenv("CLUSTER_NAME"),
			CloudFormationClient: &aws.CloudFormationClient{
				Client:          fakeAWSClient,
				PollingInterval: time.Second * 1,
			},
			Log:              logger.WithName("sqs"),
			ReconcileTimeout: time.Second * 1,
			RequeueTimeout:   time.Millisecond * 100,
		},
	}
	err = sqsReconciler.SetupWithManager(mgr, &queue.SQS{})
	Expect(err).ToNot(HaveOccurred())

	// principalCloudformationReconciler = internalawsmocks.NewMockCloudFormationReconciler(mockCtrl)
	// principalReconciler = &controllertest.ReconcilerWrapper{
	// 	Reconciler: &controllers.PrincipalReconciler{
	// 		Client:                   k8sClient,
	// 		Log:                      ctrl.Log.WithName("controllers").WithName("Principal"),
	// 		CloudFormationReconciler: principalCloudformationReconciler,
	// 		RolePrincipal:            "arn:aws:iam::123456789012:role/kiam",
	// 		PermissionsBoundary:      "arn:aws:iam::123456789012:policy/permissions-boundary",
	// 	},
	// }
	// err = principalReconciler.SetupWithManager(mgr, &access.Principal{})
	// Expect(err).ToNot(HaveOccurred())

	By("starting controller manager")
	go func() {
		<-ctrl.SetupSignalHandler()
		close(mgrStopChan)
	}()
	go func() {
		err = mgr.Start(mgrStopChan)
		Expect(err).ToNot(HaveOccurred())
	}()
}, 60)

var _ = AfterSuite(func() {
	By("stopping controller manager")
	close(mgrStopChan)
	By("stopping control plane")
	Expect(testEnv.Stop()).To(Succeed())
})
