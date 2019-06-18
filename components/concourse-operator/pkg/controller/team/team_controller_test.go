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
	"os"
	"testing"
	"time"

	concoursev1beta1 "github.com/alphagov/gsp-concourse-pipeline-controller/pkg/apis/concourse/v1beta1"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
	fakes "github.com/concourse/concourse/go-concourse/concourse/concoursefakes"
	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const timeout = time.Second * 5

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
	r := &ReconcileTeam{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		newClient: func(team string) (concourse.Client, error) {
			if team != "main" {
				t.Fatalf("modifying teams must be done from the 'main' team")
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

	// create a Team object
	teamResource := &concoursev1beta1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "footeam",
			Namespace: "xxxx-footeam",
		},
		Spec: concoursev1beta1.TeamSpec{
			Roles: []concoursev1beta1.RoleSpec{
				{
					Name: "owner",
					Github: concoursev1beta1.GithubAuth{
						Users: []string{"jeff"},
					},
				},
			},
		},
	}
	err = kubeClient.Create(ctx, teamResource)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// reconcile should be called
	expectedRequest := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "footeam",
			Namespace: "xxxx-footeam",
		},
	}
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	// CreateOrUpdate should be called with the expected atc.Team
	expectedTeamConfig := &atc.Team{
		ID:   0,
		Name: "footeam",
		Auth: atc.TeamAuth{
			"owner": map[string][]string{
				"groups": []string{},
				"users":  []string{"github:jeff"},
			},
		},
	}
	g.Eventually(func() *atc.Team {
		if teamClient.CreateOrUpdateCallCount() < 1 {
			return nil
		}
		t := teamClient.CreateOrUpdateArgsForCall(0)
		return &t
	}, timeout).Should(gomega.Equal(expectedTeamConfig))

	// delete the Team resource
	err = kubeClient.Delete(ctx, teamResource)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// reconcile should be called again
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	// DestroyTeam should be called with the team name
	g.Eventually(func() string {
		if teamClient.DestroyTeamCallCount() < 1 {
			return ""
		}
		return teamClient.DestroyTeamArgsForCall(0)
	}, timeout).Should(gomega.Equal("footeam"))
}
