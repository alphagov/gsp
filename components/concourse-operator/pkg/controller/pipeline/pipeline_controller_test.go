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
	"os"
	"testing"
	"time"

	concoursev1beta1 "github.com/alphagov/gsp/components/concourse-operator/pkg/apis/concourse/v1beta1"
	"github.com/concourse/concourse/go-concourse/concourse"
	fakes "github.com/concourse/concourse/go-concourse/concourse/concoursefakes"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const timeout = time.Second * 5

type pipelineArgs struct {
	Name       string
	Version    string
	Pipeline   []byte
	checkCreds bool
}

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// set a namespace prefix the namespace prefix is
	// stripped off namespace name to determin team name)
	os.Setenv("CONCOURSE_NAMESPACE_PREFIX", "xxxx-")

	// setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// set the fake concourse client
	concourseClient := &fakes.FakeClient{}
	teamClient := &fakes.FakeTeam{}
	concourseClient.TeamReturns(teamClient)

	// setup the test env
	kubeClient := mgr.GetClient()
	ctx := context.TODO()

	// create a manager that uses the mock TeamClient
	r := &ReconcilePipeline{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		newClient: func(team string) (concourse.Client, error) {
			if team != "myteam" {
				t.Fatalf("modifying pipelines must be done from target team got: %s expected: myteam", team)
			}
			return concourseClient, nil
		},
	}
	rWrapped, requests := SetupTestReconcile(r, t)
	g.Expect(add(mgr, rWrapped)).NotTo(gomega.HaveOccurred())

	// start the manager
	stopMgr, mgrStopped := StartTestManager(mgr, g)
	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// create a Pipeline resource
	pipelineResource := &concoursev1beta1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foopipeline",
			Namespace: "xxxx-myteam",
		},
		Spec: concoursev1beta1.PipelineSpec{
			PipelineString: `
				---
				platform: linux

				image_resource:
				  type: docker-image
				  source: {repository: busybox}

				run:
				  path: echo
				  args: [hello world]
			`,
		},
	}
	err = kubeClient.Create(ctx, pipelineResource)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// reconcile should be called
	expectedRequest := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "foopipeline",
			Namespace: "xxxx-myteam",
		},
	}
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	// Team should have been called with the target team
	g.Eventually(func() string {
		if concourseClient.TeamCallCount() < 1 {
			return ""
		}
		return concourseClient.TeamArgsForCall(0)
	}, timeout).Should(gomega.Equal("myteam"))

	// CreateOrUpdatePipeline should be called with the expected atc.Team
	expectedPipelineArgs := &pipelineArgs{
		Name:       pipelineResource.ObjectMeta.Name,
		Version:    "",
		Pipeline:   []byte(pipelineResource.Spec.PipelineString),
		checkCreds: true,
	}
	g.Eventually(func() *pipelineArgs {
		if teamClient.CreateOrUpdatePipelineConfigCallCount() < 1 {
			return nil
		}
		pipelineName, configVersion, pipelineBytes, checkCreds := teamClient.CreateOrUpdatePipelineConfigArgsForCall(0)
		return &pipelineArgs{
			Name:       pipelineName,
			Version:    configVersion,
			Pipeline:   pipelineBytes,
			checkCreds: checkCreds,
		}
	}, timeout).Should(gomega.Equal(expectedPipelineArgs))

	// delete the Pipeline resource
	err = kubeClient.Delete(ctx, pipelineResource)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// reconcile should be called again
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	// DeletePipeline should be called with the pipeline name
	g.Eventually(func() string {
		if teamClient.DeletePipelineCallCount() < 1 {
			return ""
		}
		return teamClient.DeletePipelineArgsForCall(0)
	}, timeout).Should(gomega.Equal(pipelineResource.ObjectMeta.Name))
}
