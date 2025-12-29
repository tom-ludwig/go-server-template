{{/*
Expand the name of the chart.
*/}}
{{- define "go-server.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "go-server.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "go-server.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "go-server.labels" -}}
helm.sh/chart: {{ include "go-server.chart" . }}
{{ include "go-server.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "go-server.selectorLabels" -}}
app.kubernetes.io/name: {{ include "go-server.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create image name with tag
*/}}
{{- define "go-server.image" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag }}
{{- printf "%s:%s" .Values.image.repository $tag }}
{{- end }}

{{/*
Render environment variable from flexible format
Supports: value, secretKeyRef, configMapKeyRef
*/}}
{{- define "go-server.envVar" -}}
{{- $name := .name -}}
{{- $config := .config -}}
{{- if $config.value }}
- name: {{ $name }}
  value: {{ $config.value | quote }}
{{- else if $config.secretKeyRef }}
- name: {{ $name }}
  valueFrom:
    secretKeyRef:
      name: {{ $config.secretKeyRef.name }}
      key: {{ $config.secretKeyRef.key }}
{{- else if $config.configMapKeyRef }}
- name: {{ $name }}
  valueFrom:
    configMapKeyRef:
      name: {{ $config.configMapKeyRef.name }}
      key: {{ $config.configMapKeyRef.key }}
{{- end }}
{{- end }}

{{/*
Render all environment variables from the env map
*/}}
{{- define "go-server.envVars" -}}
{{- range $name, $config := .Values.env }}
{{- include "go-server.envVar" (dict "name" $name "config" $config) }}
{{- end }}
{{- with .Values.extraEnv }}
{{- toYaml . }}
{{- end }}
{{- end }}

{{/*
Migration job name with release revision for uniqueness
*/}}
{{- define "go-server.migrationName" -}}
{{- printf "%s-migrate-%d" (include "go-server.fullname" .) .Release.Revision }}
{{- end }}

