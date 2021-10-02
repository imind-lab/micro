/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
	"os"
	"strings"
	"text/template"
)

// 生成docker
func CreateDeploy(data *Data) error {
	// 生成Makefile
	var tpl = `apiVersion: v2
name: greet
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

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/deploy/helm/" + data.Service + "/"

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
  repository: registry.cn-beijing.aliyuncs.com/imind/greet
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
  http:
    host: test.imind.tech
    port: 80
  grpc:
    host: test-grpc.imind.tech
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

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/deploy/helm/" + data.Service + "/templates/"

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

	// 生成NOTES.text
	tpl = `1. Get the application URL by running these commands:
{{- if .Values.ingress.enabled }}
{{- range $host := .Values.ingress.hosts }}
  {{- range .paths }}
  http{{ if $.Values.ingress.tls }}s{{ end }}://{{ $host.host }}{{ .path }}
  {{- end }}
{{- end }}
{{- else if contains "NodePort" .Values.service.type }}
  export NODE_PORT=$(kubectl get --namespace {{ .Release.Namespace }} -o jsonpath="{.spec.ports[0].nodePort}" services {{ include "imind.fullname" . }})
  export NODE_IP=$(kubectl get nodes --namespace {{ .Release.Namespace }} -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
{{- else if contains "LoadBalancer" .Values.service.type }}
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace {{ .Release.Namespace }} svc -w {{ include "imind.fullname" . }}'
  export SERVICE_IP=$(kubectl get svc --namespace {{ .Release.Namespace }} {{ include "imind.fullname" . }} --template "{{"{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}"}}")
  echo http://$SERVICE_IP:{{ .Values.service.port }}
{{- else if contains "ClusterIP" .Values.service.type }}
  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app.kubernetes.io/name={{ include "imind.name" . }},app.kubernetes.io/instance={{ .Release.Name }}" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace {{ .Release.Namespace }} $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace {{ .Release.Namespace }} port-forward $POD_NAME 8080:$CONTAINER_PORT
{{- end }}
`

	fileName = dir + "NOTES.text"

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
	tpl = `{{- if .Values.traefik.enabled -}}
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
    - match: Host(^{{ .Values.traefik.http.host }}^)
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
    - match: Host(^{{ .Values.traefik.grpc.host }}^)
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

	fileName = dir + "traefik.yaml"

	err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
	if err != nil {
		return err
	}

	// 生成conf-test.yaml
	tpl = `server:
  port:   #监听端口
    http: 80
    grpc: 50051
  profile:
    rate: 1

db:
  logMode: 4
  hr:
    write:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: rStoAmBDJk
      name: hr
    read:
      - host: 127.0.0.1
        port: 3306
        user: root
        pass: rStoAmBDJk
        name: hr
      - host: 127.0.0.1
        port: 3306
        user: root
        pass: rStoAmBDJk
        name: hr

redis:
  addr: 'redis-master.infra:6379'
  pass: 'VrvwqhvvRz'
  db: 0

kafka:
  business:
    producer:
      - 'kafka:9092'
    consumer:
      - 'kafka:9092'
    topic:
      commentAction: comment_action
      commonTask: common_task

tracing:
  agent: '172.16.50.50:6831'
  type: const
  param: 1
  name:
    client: imind-greet-cli
    server: imind-greet-srv

log:
  path: './logs/ms.log'
  level: -1
  age: 7
  size: 128
  backup: 30
  compress: true
  format: json
`

	t, err = template.New("conf-test.yaml").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/deploy/helm/" + data.Service + "/conf/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "conf-test.yaml"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成conf-prod.yaml
	tpl = `server:
  port:   #监听端口
    http: 80
    grpc: 50051
  profile:
    rate: 1

db:
  logMode: 4
  hr:
    write:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: rStoAmBDJk
      name: hr
    read:
      - host: 127.0.0.1
        port: 3306
        user: root
        pass: rStoAmBDJk
        name: hr
      - host: 127.0.0.1
        port: 3306
        user: root
        pass: rStoAmBDJk
        name: hr

redis:
  addr: 'redis-master.infra:6379'
  pass: 'VrvwqhvvRz'
  db: 0

kafka:
  business:
    producer:
      - 'kafka:9092'
    consumer:
      - 'kafka:9092'
    topic:
      commentAction: comment_action
      commonTask: common_task

tracing:
  agent: '172.16.50.50:6831'
  type: const
  param: 1
  name:
    client: imind-greet-cli
    server: imind-greet-srv

log:
  path: './logs/ms.log'
  level: -1
  age: 7
  size: 128
  backup: 30
  compress: true
  format: json
`

	t, err = template.New("conf-prod.yaml").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "conf-prod.yaml"

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

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/deploy/helm/" + data.Service + "/templates/tests/"

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
