/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package share

import (
    "strings"

    "github.com/imind-lab/micro/v2/microctl/template"
)

const (
    _HelmPath     = "/deploy/helm/"
    _TemplatePath = "/templates/"
)

func CreateDeploy(data *template.Data) error {
    CreateDeployChart(data)
    CreateDeployValues(data)
    CreateDeployHelper(data)
    CreateDeployDeployment(data)
    CreateDeployHpa(data)
    CreateDeployConfigMap(data)
    CreateDeploySvc(data)
    CreateDeploySa(data)
    CreateDeployTraefik(data)
    return nil
}

// 生成docker
func CreateDeployChart(data *template.Data) error {
    // 生成Makefile
    var tpl = `apiVersion: v2
name: {{.Service}}{{.Suffix}}
description: A Helm chart for Kubernetes

# A chart can be either an 'application' or a 'library' chart.
#
# Application charts are a collection of templates that can be packaged into versioned archives
# to be deployed.
#
# Library charts provide useful utilities or functions for the chart developer. They're included as
# a dependency of application charts to inject those utilities and functions into the rendering
# pipeline. Library charts do not define any templates and therefore cannot be deployed.
type: application

# This is the chart version. This version number should be incremented each time you make changes
# to the chart and its templates, including the app version.
# Versions are expected to follow Semantic Versioning (https://semver.org/)
version: 0.1.0

# This is the version number of the application being deployed. This version number should be
# incremented each time you make changes to the application. Versions are not expected to
# follow Semantic Versioning. They should reflect the version the application is using.
# It is recommended to use it with quotes.
appVersion: "1.0.0"

icon: https://static.imind.tech/frontend/images/wechat/bj.png
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + "/"
    name := "Chart.yaml"

    return template.CreateFile(data, tpl, path, name)
}
func CreateDeployValues(data *template.Data) error {
    // 生成values.yaml
    tpl := `# Default values for imind.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 2

image:
  repository: registry.cn-beijing.aliyuncs.com/imind/{{.Service}}{{.Suffix}}
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: ""

imagePullSecrets:
  - name: regsecret

nameOverride: ""
fullnameOverride: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # Annotations to add to the service account
  annotations: {}
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name: ""

podAnnotations: {}

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  type: ClusterIP
  ports:
    - name: http
      port: 80
    - name: grpc
      port: 50051

traefik:
  enabled: true
  host: {{.Service}}{{.Suffix}}.imind.tech
  http:
    port: 80
  grpc:
    port: 50051
    tls: traefik-cert

ingress:
  enabled: false
  annotations:
    {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          backend:
            serviceName: chart-example.local
            servicePort: 80
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources:
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi

livenessProbe:
  exec:
    command:
      - /bin/grpc-health-probe
      - -addr=localhost:50051
      - -tls
      - -tls-ca-cert=/conf/ssl/tls.crt
      - -tls-server-name=www.imind.tech
  initialDelaySeconds: 10
  timeoutSeconds: 5
  periodSeconds: 10
  successThreshold: 1
  failureThreshold: 2

readinessProbe:
  exec:
    command:
      - /bin/grpc-health-probe
      - -addr=localhost:50051
      - -tls
      - -tls-ca-cert=/conf/ssl/tls.crt
      - -tls-server-name=www.imind.tech
  initialDelaySeconds: 10
  timeoutSeconds: 5
  periodSeconds: 10
  successThreshold: 1
  failureThreshold: 2

hostAliases:
  - hostnames:
      - kafka
    ip: 172.22.131.242

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}

config:
  global:
    rate:
      high:
        limit: 10
        capacity: 10
      low:
        limit: 10
        capacity: 10
    profile:
      rate: 1

  service:
    namespace: {{.Project}}
    name: {{.Service}}
    version: latest
    logLevel: -1
    logFormat: json
    port: #监听端口
      http: 80
      grpc: 50051

  db:
    logLevel: 4
    max:
      open: 10
      idle: 5
      life: 30
    timeout: 5s
    default:
      master:
        host: 127.0.0.1
        port: 3306
        user: root
        pass: imind
        name: imind
      replica:
        host: 127.0.0.1
        port: 3306
        user: root
        pass: imind
        name: imind

  redis:
    model: node
    timeout: 5s
    addr: '127.0.0.1:6379'
    #  pass: imind
    db: 0

  kafka:
    default:
      producer:
        - 'kafka.infra:9092'
      consumer:
        - 'kafka.infra:9092'
      topic:
        {{.Service}}Create: {{.Package}}_create
        {{.Service}}Update: {{.Package}}_update

  tracing:
    agent:
      host: '127.0.0.1'
      port: 6831

templates: {}
#   environment.tmpl: |-
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + "/"
    name := "values.yaml"

    return template.CreateFile(data, tpl, path, name)
}

func CreateDeployHelper(data *template.Data) error {
    // 生成helpers.tpl
    tpl := `{{/*
Expand the name of the chart.
*/}}
{{- define "imind.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "imind.fullname" -}}
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
{{- define "imind.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "imind.labels" -}}
helm-chart: {{ include "imind.chart" . }}
{{ include "imind.selectorLabels" . }}
version: {{ .Chart.AppVersion }}
managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "imind.selectorLabels" -}}
app: {{ include "imind.name" . }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "imind.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "imind.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "_helpers.tpl"
    return template.WriteFile(tpl, path, name)
}

func CreateDeployDeployment(data *template.Data) error {
    // 生成deployment.yaml
    tpl := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "imind.fullname" . }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "imind.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "imind.selectorLabels" . | nindent 8 }}
      annotations:
        checksum/config: {{ include (print $.Template.BasePath "/configmap.yaml") . | sha256sum }}
        {{- with .Values.podAnnotations }}
          {{- toYaml . | nindent 8 }}
        {{- end }}
    spec:
      {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "imind.serviceAccountName" . }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          command:
            - /bin/{{ .Chart.Name }}
            - server
          image: {{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            {{- range .Values.service.ports }}
            - name: {{ .name }}
              containerPort: {{ .port }}
              protocol: TCP
            {{- end }}
          livenessProbe:
            {{- toYaml .Values.livenessProbe | nindent 12 }}
          readinessProbe:
            {{- toYaml .Values.readinessProbe | nindent 12 }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
          volumeMounts:
            - name: conf
              mountPath: /conf/conf.yaml
              subPath: conf.yaml
      volumes:
        - name: conf
          configMap:
            name: {{ include "imind.fullname" . }}
            items:
              - key: conf.yaml
                path: conf.yaml
      {{- with .Values.hostAliases }}
      hostAliases:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
`
    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "deployment.yaml"

    return template.WriteFile(tpl, path, name)
}

func CreateDeployHpa(data *template.Data) error {
    // 生成hpa.yaml
    tpl := `{{- if .Values.autoscaling.enabled }}
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "imind.fullname" . }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "imind.fullname" . }}
  minReplicas: {{ .Values.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.autoscaling.maxReplicas }}
  metrics:
    {{- if .Values.autoscaling.targetCPUUtilizationPercentage }}
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
    {{- end }}
    {{- if .Values.autoscaling.targetMemoryUtilizationPercentage }}
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: {{ .Values.autoscaling.targetMemoryUtilizationPercentage }}
    {{- end }}
{{- end }}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "hpa.yaml"

    return template.WriteFile(tpl, path, name)
}

func CreateDeployConfigMap(data *template.Data) error {
    // 生成configmap.yaml
    tpl := `{{- if .Values.config }}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "imind.fullname" . }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
data:
  conf.yaml: |-
    {{- toYaml .Values.config | default "{}" | nindent 4 }}
  {{- range $key, $value := .Values.templates }}
  {{ $key }}: |-
    {{- $value | nindent 4 }}
  {{- end }}
{{- end }}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "configmap.yaml"

    return template.WriteFile(tpl, path, name)
}

func CreateDeploySvc(data *template.Data) error {
    // 生成service.yaml
    tpl := `apiVersion: v1
kind: Service
metadata:
  name: {{ include "imind.fullname" . }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    {{- toYaml .Values.service.ports | nindent 4 }}
  selector:
    {{- include "imind.selectorLabels" . | nindent 4 }}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "service.yaml"

    return template.WriteFile(tpl, path, name)
}

func CreateDeploySa(data *template.Data) error {
    // 生成serviceaccount.yaml
    tpl := `{{- if .Values.serviceAccount.create -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "imind.serviceAccountName" . }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
  {{- with .Values.serviceAccount.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
{{- end }}
`
    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "serviceaccount.yaml"

    return template.WriteFile(tpl, path, name)
}

func CreateDeployTraefik(data *template.Data) error {
    // 生成traefik.yaml
    tpl := `{{- if .Values.traefik.enabled -}}
{{- $fullName := include "imind.fullname" . -}}
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
spec:
  entryPoints:
    - web
  routes:
    - match: Host(^{{ .Values.traefik.host }}^)
      kind: Rule
      services:
        - name: {{ $fullName }}
          port: {{ .Values.traefik.http.port }}

---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: {{ $fullName }}-grpc
  labels:
    {{- include "imind.labels" . | nindent 4 }}
spec:
  entryPoints:
    - websecure
  routes:
    - match: Host(^{{ .Values.traefik.host }}^)
      kind: Rule
      services:
        - name: {{ $fullName }}
          port: {{ .Values.traefik.grpc.port }}
          kind: Service
          scheme: https
  tls:
    secretName: {{ .Values.traefik.grpc.tls }}
{{- end }}
`
    tpl = strings.Replace(tpl, "^", "`", -1)

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + data.Suffix + _HelmPath + data.Service + data.Suffix + _TemplatePath
    name := "traefik.yaml"

    return template.WriteFile(tpl, path, name)
}
