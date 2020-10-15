package validating

import (
	"context"
	"fmt"
	"strings"
	"testing"

	concoursev1beta1 "github.com/alphagov/gsp/components/concourse-operator/pkg/apis/concourse/v1beta1"
	"sigs.k8s.io/yaml"
)

type ValidationTestCase struct {
	Name                 string
	Pipeline             string
	Allowed              bool
	Reason               string
	HandlerErrorContains string
}

var validationCases = []ValidationTestCase{
	{
		Name:    "valid-pipeline-task",
		Allowed: true,
		Reason:  "ok",
		Pipeline: `
apiVersion: concourse.govsvc.uk/v1beta1
kind: Pipeline
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: pipeline-sample
spec:
  config:
    jobs:
    - name: hello
      plan:
      - task: say
        config:
          platform: linux
          image_resource:
            type: docker-image
            source: {repository: busybox}
          run:
            path: echo
            args: [hello world]
`,
	},

	{
		Name:    "using-in_parallel-step",
		Allowed: true,
		Reason:  "ok",
		Pipeline: `
apiVersion: concourse.govsvc.uk/v1beta1
kind: Pipeline
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: pipeline-sample
spec:
  config:
    resources:
    - name: thing1
      type: docker-image
    - name: thing2
      type: docker-image
    jobs:
    - name: job-with-parallel
      plan:
      - in_parallel:
          limit: 2
          steps:
          - get: thing1
          - get: thing2
`,
	},

	{
		Name:    "non-existant-resource",
		Allowed: false,
		Reason:  "get(non-existant-resource): unknown resource",
		Pipeline: `
apiVersion: concourse.govsvc.uk/v1beta1
kind: Pipeline
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: pipeline-sample
spec:
  config:
    jobs:
    - name: bad-job
      plan:
      - get: non-existant-resource
      - task: say
        config:
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

func TestPipelineValidation(t *testing.T) {
	for _, tc := range validationCases {
		err := testPipelineValidation(tc)
		if err != nil {
			t.Fatalf("%v failed: %v", tc, err)
		}
	}
}

func testPipelineValidation(tc ValidationTestCase) error {
	h := &PipelineCreateUpdateHandler{}
	var pipeline concoursev1beta1.Pipeline
	unmarshalError := yaml.Unmarshal([]byte(tc.Pipeline), &pipeline)
	if unmarshalError != nil {
		return fmt.Errorf("did not expect unmarshalError but got: %v", unmarshalError)
	}

	res := h.handle(context.Background(), &pipeline)
	if res.Response.Allowed != tc.Allowed {
		return fmt.Errorf("expected handler to be allowed='%v' but got allowed='%v' with resaon='%v'", tc.Allowed, res.Response.Allowed, res.Response.Result.Reason)
	}
	if !res.Response.Allowed {
		if !strings.Contains(string(res.Response.Result.Reason), tc.Reason) {
			return fmt.Errorf("expected failure reason to contain '%v' but got '%v'", tc.Reason, res.Response.Result.Reason)
		}
	}
	return nil
}
