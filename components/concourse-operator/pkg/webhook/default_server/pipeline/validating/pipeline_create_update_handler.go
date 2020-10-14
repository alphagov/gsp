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

package validating

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	concoursev1beta1 "github.com/alphagov/gsp/components/concourse-operator/pkg/apis/concourse/v1beta1"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/configvalidate"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validating-create-update-pipeline"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &PipelineCreateUpdateHandler{})
}

// PipelineCreateUpdateHandler handles Pipeline
type PipelineCreateUpdateHandler struct {
	// To use the client, you need to do the following:
	// - uncomment it
	// - import sigs.k8s.io/controller-runtime/pkg/client
	// - uncomment the InjectClient method at the bottom of this file.
	// Client  client.Client

	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *PipelineCreateUpdateHandler) Validate(config *atc.Config) (bool, string, error) {
	warnings, err := h.validationWarnings(config)
	if err != nil {
		msg := fmt.Sprintf("unable to parse pipeline: %s", err.Error())
		return false, msg, nil
	}
	if len(warnings) > 0 {
		msg := fmt.Sprintf("pipeline validation failed: %s", strings.Join(warnings, ", "))
		return false, msg, nil
	}
	return true, "ok", nil
}

func (h *PipelineCreateUpdateHandler) validationWarnings(config *atc.Config) ([]string, error) {
	warnings := []string{}

	warningMessages, errorMessages := configvalidate.Validate(*config)

	if len(warningMessages) > 0 {
		for _, warning := range warningMessages {
			warnings = append(warnings, warning.Message)
		}
	}

	if len(errorMessages) > 0 {
		for _, err := range errorMessages {
			warnings = append(warnings, err)
		}
	}

	return warnings, nil
}

var _ admission.Handler = &PipelineCreateUpdateHandler{}

// Handle handles admission requests.
func (h *PipelineCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &concoursev1beta1.Pipeline{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	var config atc.Config

	if len(obj.Spec.Config.Jobs) > 0 {
		config = obj.Spec.Config
	} else if obj.Spec.PipelineString != "" {
		err := atc.UnmarshalConfig([]byte(obj.Spec.PipelineString), &config)
		if err != nil {
			return admission.ErrorResponse(http.StatusInternalServerError, err)
		}
	} else {
		return admission.ValidationResponse(false, "need to define `config` or `pipelineString`")
	}

	allowed, reason, err := h.Validate(&obj.Spec.Config)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

//var _ inject.Client = &PipelineCreateUpdateHandler{}
//
//// InjectClient injects the client into the PipelineCreateUpdateHandler
//func (h *PipelineCreateUpdateHandler) InjectClient(c client.Client) error {
//	h.Client = c
//	return nil
//}

var _ inject.Decoder = &PipelineCreateUpdateHandler{}

// InjectDecoder injects the decoder into the PipelineCreateUpdateHandler
func (h *PipelineCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
