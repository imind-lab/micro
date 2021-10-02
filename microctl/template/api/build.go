/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package api

import (
	"os"
	"strings"
	"text/template"

	tp "github.com/imind-lab/micro/microctl/template"
)

// 生成docker
func CreateBuild(data *tp.Data) error {
	// 生成Makefile
	var tpl = `GOPATH := $(shell go env GOPATH)
VERSION := 0.0.1.0

gengo:
	protoc -I. --proto_path ../server/proto \
 --go_out ../server/proto --go_opt paths=source_relative --go-grpc_out ../server/proto --go-grpc_opt paths=source_relative \
 --grpc-gateway_out ../server/proto --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true {{.Service}}-api/{{.Service}}-api.proto
	protoc-go-inject-tag -input=../server/proto/{{.Service}}-api/{{.Service}}-api.pb.go

genphp:
	protoc -I. --proto_path ../server/proto \
 --php_out ../server/proto/{{.Service}}-api --grpc_out ../server/proto/{{.Service}}-api --plugin=protoc-gen-grpc=${GOPATH}/bin/grpc_php_plugin {{.Service}}-api/{{.Service}}-api.proto

depend: gengo
	go get ../...

build: depend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}} ../main.go

test:
	go test -v ../... -cover

docker: build health
	docker build -f ./Dockerfile -t registry.cn-beijing.aliyuncs.com/imind/{{.Service}}-api:$(VERSION) ../
	#docker push registry.cn-beijing.aliyuncs.com/imind/{{.Service}}-api:$(VERSION)
	rm -rf {{.Service}} grpc-health-probe

health:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o grpc-health-probe ../pkg/grpc-health-probe/main.go

helm:
	helm install {{.Service}} ./helm/{{.Service}} --set image.tag=$(VERSION)

clean:
	docker rmi registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION)

k8s:
	kubectl set image deployment/{{.Service}}-api {{.Service}}=registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION) -n imind-lab

.PHONY: gengo depend build test docker clean deploy k8s
`

	t, err := template.New("makefile").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/build/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "Makefile"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成Dockerfile
	tpl = `FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

WORKDIR .
ADD conf /conf
ADD build/{{.Service}} build/grpc-health-probe  /bin/
ENTRYPOINT [ "/bin/{{.Service}}" ]
`

	t, err = template.New("dockerfile").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/build/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "Dockerfile"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成ns-rbac.yaml
	tpl = `apiVersion: v1
kind: Namespace
metadata:
  name: imind-lab

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: julive-sa
  namespace: imind-lab

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: julive-registry
  namespace: imind-lab
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - configmaps
    verbs:
      - get
      - list
      - patch
      - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: julive-registry
  namespace: imind-lab
  labels:
    app: julive-rbac
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: julive-registry
subjects:
  - kind: ServiceAccount
    name: julive-sa
    namespace: imind-lab
`

	t, err = template.New("rbac").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "ns-rbac.yaml"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成deploy.yaml
	tpl = `apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Service}}-api
  namespace: imind-lab
data:
  conf.yaml: |-
    server:
      port:   #监听端口
        http: 80
        grpc: 50051
      profile:
        rate: 1

    db:
      hr:
        write:
          host: mysql.infra
          port: 3306
          user: root
          pass: 9WVULeeRPN
          name: hr
        read:
          - host: mysql.infra
            port: 3306
            user: root
            pass: 9WVULeeRPN
            name: hr
          - host: mysql.infra
            port: 3306
            user: root
            pass: 9WVULeeRPN
            name: hr

      bbs:
        write:
          host: mysql.infra
          port: 3306
          user: root
          pass: 9WVULeeRPN
          name: bbs
        read:
          - host: mysql.infra
            port: 3306
            user: root
            pass: 9WVULeeRPN
            name: bbs
          - host: mysql.infra
            port: 3306
            user: root
            pass: 9WVULeeRPN
            name: bbs

    redis:
      addr: 'redis-master.infra:6379'
      pass: '8l8GWyhWJx'
      db: 0

    kafka:
      business:
        producer:
          - 'kafka.infra:9092'
        consumer:
          - 'kafka.infra:9092'
        topic:
          commentAction: comment_action
          commonTask: common_task
          create{{.Svc}}: create_{{.Service}}
          update{{.Svc}}Count: update_{{.Service}}_count

      bigdata:
        producer:
          - 'kafka.infra:9092'
        consumer:
          - 'kafka.infra:9092'
        topic:
          commentAction: comment_action
          commonTask: common_task
          create{{.Svc}}: create_{{.Service}}
          update{{.Svc}}Count: update_{{.Service}}_count

    tracing:
      agent: '172.16.50.50:6831'
      type: const
      param: 1
      name:
        client: {{.Project}}-{{.Service}}-cli
        server: {{.Project}}-{{.Service}}-srv

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

---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: imind-lab
  name: {{.Service}}-api
spec:
  replicas: 1
  minReadySeconds: 30
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  selector:
    matchLabels:
      app: {{.Service}}-api
      version: latest
  template:
    metadata:
      labels:
        app: {{.Service}}-api
        version: latest
    spec:
      imagePullSecrets:
        - name: alisecret
      serviceAccountName: julive-sa
      containers:
        - name: {{.Service}}
          command:
            - /bin/{{.Service}}
            - server
          image: registry.cn-beijing.aliyuncs.com/julive/{{.Service}}-api:0.0.1.0
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
              name: http
            - containerPort: 50051
              name: grpc
          livenessProbe:
            exec:
              command:
                - /bin/grpc-health-probe
                - -addr=:50051
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
            initialDelaySeconds: 10
            timeoutSeconds: 5
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 2
          resources:
            limits:
              cpu: 100m
              memory: 100Mi
            requests:
              cpu: 100m
              memory: 100Mi
          volumeMounts:
            - name: conf
              mountPath: /conf/conf.yaml
              subPath: conf.yaml
      volumes:
        - name: conf
          configMap:
            name: {{.Service}}-api
            items:
              - key: conf.yaml
                path: conf.yaml
      hostAliases:
        - hostnames:
            - testkafka
          ip: 172.22.131.242

---
apiVersion: v1
kind: Service
metadata:
  name: {{.Service}}-api
  namespace: imind-lab
spec:
  ports:
    - name: http-80
      port: 80
    - name: grpc-50051
      port: 50051
  selector:
    app: {{.Service}}-api
    version: latest

---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: ingress-{{.Service}}-api
  namespace: imind-lab
spec:
  entryPoints:
    - web
  routes:
    - match: Host(${backtick}{{.Service}}-api.imind.tech${backtick})
      kind: Rule
      services:
        - name: {{.Service}}-api
          port: 80
`
	tpl = strings.Replace(tpl, "${backtick}", "`", -1)
	t, err = template.New("rbac").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + data.Service + "-api.yaml"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
