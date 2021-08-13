{{/*
Common labels
*/}}
{{- define "chart.labels" -}}
{{ include "chart.selectorLabels" . }}
project: {{ .Values.project }}
component: {{ .Values.component }}
release: {{ .Release.Name }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "chart.selectorLabels" -}}
name: {{ .Release.Name }}
instance: {{ .Release.Name }}
{{- end }}
