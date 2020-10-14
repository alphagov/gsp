package validating

import (
	"fmt"
	"strings"
	"testing"

	"github.com/concourse/concourse/atc"
)

type ValidationTestCase struct {
	Name                      string
	Pipeline                  string
	Valid                     bool
	ValidationMessageContains string
	HandlerErrorContains      string
}

var validationCases = []ValidationTestCase{
	{
		Name:  "valid-pipeline-task",
		Valid: true,
		Pipeline: `
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
		Name:  "using-in_parallel-step",
		Valid: true,
		Pipeline: `
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
		Name:                      "non-existant-resource",
		Valid:                     false,
		ValidationMessageContains: "jobs.bad-job.plan.do[0].get(non-existant-resource): unknown resource",
		Pipeline: `
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
	var config atc.Config
	unmarshalError := atc.UnmarshalConfig([]byte(tc.Pipeline), &config)
	if unmarshalError != nil {
		return fmt.Errorf("did not expect unmarshalError but got: %v", unmarshalError)
	}

	valid, validationMessage, handlerError := h.Validate(&config)
	if handlerError != nil {
		if tc.HandlerErrorContains == "" {
			return fmt.Errorf("did not expect handlerError but got: %v", handlerError)
		}
		if !strings.Contains(handlerError.Error(), tc.HandlerErrorContains) {
			return fmt.Errorf("expected handlerError to contain '%v' but got: %v", tc.HandlerErrorContains, handlerError)
		}
		return nil
	}
	if validationMessage != "ok" {
		if tc.ValidationMessageContains == "" {
			return fmt.Errorf("did not expect validationMessage but got: %v", validationMessage)
		}
		if !strings.Contains(validationMessage, tc.ValidationMessageContains) {
			return fmt.Errorf("expected validationMessage to contain '%v' but got: %v", tc.ValidationMessageContains, validationMessage)
		}
	}
	if tc.Valid != valid {
		return fmt.Errorf("expected valid=%v got: %v", tc.Valid, valid)
	}
	return nil
}
