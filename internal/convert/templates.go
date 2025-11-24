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

package convert

const radiusContainersTemplate = `{{ $workloadName := .WorkloadName }}{{ $firstWorkload := true }}{{ $service := .Spec.Service }}{{ $resources := .Spec.Resources }}
extension radius

@description('The Radius Application ID. Injected automatically by the rad CLI.')
param application string

@description('The Radius Environment ID. Injected automatically by the rad CLI.')
param environment string
{{ range $containerName, $container := .Spec.Containers }}
{{- if $firstWorkload }}
resource {{ $workloadName }} 'Applications.Core/containers@2023-10-01-preview' = {
  name: '{{ $workloadName }}'
  properties: {
    application: application
    environment: environment
    container: {
      image: '{{ $container.Image }}'

      {{- if (gt (len $container.Command) 0) }}
      command: [
        {{- range $i, $cmd := $container.Command }}
        '{{ $cmd }}'
        {{- end }}
      ]{{- end }}

      {{- if (gt (len $container.Args) 0) }}
      args: [
        {{- range $i, $arg := $container.Args }}
        '{{ $arg }}'
        {{- end }}
      ]{{- end }}

      {{- if (gt (len $container.Variables) 0) }}
      env: {
        {{- range $variableName, $variableValue := $container.Variables }}
        {{ $variableName }}: {
          value: '{{ $variableValue }}'
        }{{- end }}
      }{{- end }}

      {{- if and (ne $service nil) (gt (len $service.Ports) 0) }}
      ports: {
        {{- range $portName, $port := $service.Ports }}
        '{{ $portName }}': {
          port: {{ $port.Port }}
          {{- if ne $port.Protocol nil }}
          protocol: '{{ $port.Protocol }}'
          {{- end }}
          {{- if ne $port.TargetPort nil }}
          containerPort: {{ $port.TargetPort }}
          {{- end }}
        }{{- end }}
      }{{- end }}

      {{- if (ne $container.LivenessProbe nil) }}
      livenessProbe: {
        {{- if (ne $container.LivenessProbe.Exec nil) }}
        kind: 'exec'
        command: {{ $container.LivenessProbe.Exec.Command }}
        {{- else if (ne $container.LivenessProbe.HttpGet nil) }}
        kind: 'httpGet'
        containerPort: {{ $container.LivenessProbe.HttpGet.Port }}
        {{- if (ne $container.LivenessProbe.HttpGet.Path "") }}
        path: '{{ $container.LivenessProbe.HttpGet.Path }}'
        {{- end }}
        {{- end }}
      }{{- end }}
      {{- if (ne $container.ReadinessProbe nil) }}
      readinessProbe: {
       {{- if (ne $container.ReadinessProbe.Exec nil) }}
        kind: 'exec'
        command: {{ $container.ReadinessProbe.Exec.Command }}
        {{- else if (ne $container.ReadinessProbe.HttpGet nil) }}
        kind: 'httpGet'
        containerPort: {{ $container.ReadinessProbe.HttpGet.Port }}
        {{- if (ne $container.ReadinessProbe.HttpGet.Path "") }}
        path: '{{ $container.ReadinessProbe.HttpGet.Path }}'
        {{- end }}
        {{- end }}
      }{{- end }}
    }

    {{- if gt (len $resources) 0 }}
    connections: {
      {{- range $resource := $resources }}
      {{- $resourceId := splitList "." $resource.Id | last }}
      {{ $resourceId }}: {
        source: {{ $resourceId }}.id
        disableDefaultEnvVars: {{ and (ne $resource.Params.disableDefaultEnvVars nil) $resource.Params.disableDefaultEnvVars }}
      }
      {{- end }}
    }
    {{- end }}
  }
}{{ $firstWorkload := false }}{{- end }}{{ end }}
`
