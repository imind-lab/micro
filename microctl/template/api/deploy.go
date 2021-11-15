/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package api

import (
	tp "github.com/imind-lab/micro/microctl/template"
	"os"
	"strings"
	"text/template"
)

// 生成docker
func CreateDeploy(data *tp.Data) error {
	// 生成Makefile
	var tpl = `apiVersion: v2
name: {{.Service}}-api
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

	t, err := template.New("chart").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/deploy/helm/" + data.Service + "-api/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "Chart.yaml"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成values.yaml
	tpl = `# Default values for imind.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 2

image:
  repository: 348681422678.dkr.ecr.ap-southeast-1.amazonaws.com/{{.Project}}/{{.Service}}-api
  pullPolicy: Always
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

istio:
  enabled: true
  http:
    host: {{.Service}}-api.chope.co
    port: 80

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
      - -addr=:50051
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
      - -addr=:50051
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
`

	t, err = template.New("values").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "values.yaml"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成helpers.tpl
	tpl = `{{/*
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

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/deploy/helm/" + data.Service + "-api/templates/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "_helpers.tpl"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成deployment.yaml
	tpl = `apiVersion: apps/v1
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
        checksum/config: {{ include (print $.Template.BasePath "/secret.yaml") . | sha256sum }}
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
          secret:
            secretName: {{ include "imind.fullname" . }}
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

	fileName = dir + "deployment.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成hpa.yaml
	tpl = `{{- if .Values.autoscaling.enabled }}
apiVersion: autoscaling/v2beta1
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
        targetAverageUtilization: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
    {{- end }}
    {{- if .Values.autoscaling.targetMemoryUtilizationPercentage }}
    - type: Resource
      resource:
        name: memory
        targetAverageUtilization: {{ .Values.autoscaling.targetMemoryUtilizationPercentage }}
    {{- end }}
{{- end }}
`

	fileName = dir + "hpa.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成secret.yaml
	tpl = `{{- $fullName := include "imind.fullname" . -}}
apiVersion: v1
kind: Secret
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "imind.labels" . | nindent 4 }}
type: Opaque
data:
  {{ (.Files.Glob "conf/conf.yaml").AsSecrets | indent 2 }}
`

	fileName = dir + "secret.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成service.yaml
	tpl = `apiVersion: v1
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

	fileName = dir + "service.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成serviceaccount.yaml
	tpl = `{{- if .Values.serviceAccount.create -}}
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

	fileName = dir + "serviceaccount.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成traefik.yaml
	tpl = `{{- if .Values.istio.enabled -}}
{{- $fullName := include "imind.fullname" . -}}
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: {{ $fullName }}-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      hosts:
        - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ $fullName }}
spec:
  hosts:
    - "{{ .Values.istio.http.host }}"
  gateways:
    - {{ $fullName }}-gateway
  http:
    - match:
        - uri:
            prefix: /
      route:
        - destination:
            host: {{ $fullName }}
            port:
              number: {{ .Values.istio.http.port }}
{{- end }}
`
	tpl = strings.Replace(tpl, "^", "`", -1)

	fileName = dir + "istio.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成conf.yaml
	tpl = `service:
  project: {{.Project}}
  name: {{.Service}}-api
  port: #监听端口
    http: 80
    grpc: 50051
  profile:
    rate: 1

db:
  logMode: 4
  chope:
    write:
      host: uat-daniel-ms-instance-1.ctn3ok75ssrh.ap-southeast-1.rds.amazonaws.com
      port: 3306
      user: uatdanielms
      pass: 9xBf6cdlVlxwiJeN
      name: chope
    read:
      - host: uat-daniel-ms-instance-1.ctn3ok75ssrh.ap-southeast-1.rds.amazonaws.com
        port: 3306
        user: uatdanielms
        pass: 9xBf6cdlVlxwiJeN
        name: chope
      - host: uat-daniel-ms-instance-1.ctn3ok75ssrh.ap-southeast-1.rds.amazonaws.com
        port: 3306
        user: uatdanielms
        pass: 9xBf6cdlVlxwiJeN
        name: chope

redis:
  addr: '127.0.0.1:6379'
  db: 0

kafka:
  business:
    producer:
      - '127.0.0.1:9092'
    consumer:
      - '127.0.0.1:9092'
    topic:
      {{.Service}}Create: {{.Service}}_create
      {{.Service}}Update: {{.Service}}_update

tracing:
  agent: '172.16.50.50:6831'
  type: const
  param: 1
  name:
    client: {{.Project}}-{{.Service}}-api-cli
    server: {{.Project}}-{{.Service}}-api-srv

log:
  path: './logs/ms.log'
  level: -1
  age: 7
  size: 128
  backup: 30
  compress: true
  format: json
rpc:
  {{.Service}}:
    service: {{.Service}}
    port: 50051
`

	t, err = template.New("conf.yaml").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/deploy/helm/" + data.Service + "-api/conf/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "conf.yaml"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成_helpers.tpl
	tpl = `apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "imind.fullname" . }}-test-connection"
  labels:
    {{- include "imind.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args: ['{{ include "imind.fullname" . }}:{{ .Values.service.port }}']
  restartPolicy: Never
`

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/deploy/helm/" + data.Service + "-api/templates/tests/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "test-connection.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
