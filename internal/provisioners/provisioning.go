// Copyright 2024 The Score Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provisioners

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/score-spec/score-go/framework"

	"github.com/score-spec/score-radius/internal/state"
)

type Provisioner struct {
	Uri         string `yaml:"uri"`
	ResType     string `yaml:"type"`
	Format      string `yaml:"format"`
	Class       string `yaml:"class"`
	Description string `yaml:"description,omitempty"`
	// The InitTemplate is always evaluated first, it is used as temporary or working set data that may be needed in the
	// later templates. It has access to the resource inputs and previous state.
	InitTemplate      string `yaml:"init,omitempty"`
	ManifestsTemplate string `yaml:"manifests,omitempty"`
	// Params is a list of inputs that the provisioner expects to be passed in.
	Params []string `yaml:"params,omitempty"`
	// Outputs is a list of outputs that the provisioner should return.
	Outputs []string `yaml:"expected_outputs,omitempty"`
	// Outputs is a list of actual outputs evaluated from the template.
	OutputsTemplate string `yaml:"outputs,omitempty"`
}

type Data struct {
	Id           string
	Init         map[string]interface{}
	WorkloadName string
}

func ProvisionResources(currentState *state.State, provisioners []Provisioner) (string, *state.State, error) {
	out := currentState
	manifests := ""

	// provision in sorted order
	orderedResources, err := currentState.GetSortedResourceUids()
	if err != nil {
		return "", nil, fmt.Errorf("failed to determine sort order for provisioning: %w", err)
	}

	out.Resources = maps.Clone(out.Resources)
	for _, resUid := range orderedResources {
		resState := out.Resources[resUid]

		provisionerIndex := slices.IndexFunc(provisioners, func(provisioner Provisioner) bool {
			return provisioner.ResType == resUid.Type() && provisioner.Class == resUid.Class()
		})
		if provisionerIndex < 0 {
			return "", nil, fmt.Errorf("resource '%s' is not supported by any provisioner. Please implement a custom resource provisioner to support this resource type '%s' with class '%s'", resUid, resUid.Type(), resUid.Class())
		}

		var params map[string]interface{}
		if len(resState.Params) > 0 {
			resOutputs, err := out.GetResourceOutputForWorkload(resState.SourceWorkload)
			if err != nil {
				return "", nil, fmt.Errorf("%s: failed to find resource params for resource: %w", resUid, err)
			}
			sf := framework.BuildSubstitutionFunction(out.Workloads[resState.SourceWorkload].Spec.Metadata, resOutputs)
			rawParams, err := framework.Substitute(resState.Params, sf)
			if err != nil {
				return "", nil, fmt.Errorf("%s: failed to substitute params for resource: %w", resUid, err)
			}
			params = rawParams.(map[string]interface{})
		}
		resState.Params = params
		provisioner := provisioners[provisionerIndex]
		resState.ProvisionerUri = provisioner.Uri

		init := make(map[string]interface{})
		data := Data{
			Id:           resState.Id,
			Init:         init,
			WorkloadName: resState.SourceWorkload,
		}

		if err := renderTemplateAndDecode(provisioner.InitTemplate, &data, &data.Init); err != nil {
			return "", nil, fmt.Errorf("init template failed: %w", err)
		}

		resState.Outputs = make(map[string]interface{})
		if err := renderTemplateAndDecode(provisioner.OutputsTemplate, &data, &resState.Outputs); err != nil {
			return "", nil, fmt.Errorf("outputs template failed: %w", err)
		}

		var resourceManifest string
		resourceManifest, err = generateResourceManifest(provisioner.ManifestsTemplate, data)
		if err != nil {
			return "", nil, fmt.Errorf("failed to generate resource manifest %s: %w", resUid.Type(), err)
		}
		slog.Info(fmt.Sprintf("Resource %s's manifests generated", resUid.Type()))

		out.Resources[resUid] = resState
		resourceManifest = "\n" + resourceManifest
		manifests = manifests + resourceManifest
	}

	return manifests, out, nil
}

func renderTemplateAndDecode(raw string, data interface{}, out interface{}) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	prepared, err := template.New("").Funcs(sprig.FuncMap()).Parse(raw)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	buff := new(bytes.Buffer)
	if err := prepared.Execute(buff, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	buffContents := buff.String()
	if strings.TrimSpace(buff.String()) == "" {
		return nil
	}
	var intermediate interface{}
	if err := yaml.Unmarshal([]byte(buffContents), &intermediate); err != nil {
		slog.Debug(fmt.Sprintf("template output was '%s' from template '%s'", buffContents, raw))
		return fmt.Errorf("failed to decode output: %w", err)
	}
	err = mapstructure.Decode(intermediate, &out)
	if err != nil {
		slog.Debug(fmt.Sprintf("template output was '%s' from template '%s'", intermediate, raw))
		return fmt.Errorf("failed to decode output: %w", err)
	}
	return nil
}

func generateResourceManifest(resourceTypeTemplate string, data Data) (string, error) {
	t, err := template.New("").Funcs(sprig.FuncMap()).Parse(resourceTypeTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
