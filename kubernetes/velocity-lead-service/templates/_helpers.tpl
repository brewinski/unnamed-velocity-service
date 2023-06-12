{{/*
Expand the name of the chart.
*/}}
{{- define "velocity-lead-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "velocity-lead-service.fullname" -}}
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
{{- define "velocity-lead-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "velocity-lead-service.labels" -}}
helm.sh/chart: {{ include "velocity-lead-service.chart" . }}
{{ include "velocity-lead-service.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "velocity-lead-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "velocity-lead-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "velocity-lead-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "velocity-lead-service.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Create a fully qualified DNS name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "velocity-lead-service.dnsname" -}}
{{- $name := .Release.Name -}}
{{- $domain := default ".local" .Values.environment.SITE_ENV_DOMAIN_NAME -}}
{{- if and .Values.deploymentEnvironment (eq .Values.deploymentEnvironment "production") -}}
    {{- if eq (.Values.deploymentEnv | toString ) (.Values.activeProdEnv | toString) -}}
        {{- printf "%s" $domain | trunc 63 | trimSuffix "-" -}}
    {{- else }}
        {{- printf "%s.api.canstar.com.au" $name | trunc 63 | trimSuffix "-" -}}
    {{- end -}}
{{- else if and .Values.deploymentEnvironment (eq .Values.deploymentEnvironment "staging") -}}
    {{- printf "%s.uat.api.canstar.com.au" $name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- printf "%s.local" $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}