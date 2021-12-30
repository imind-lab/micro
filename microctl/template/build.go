/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
	"os"
	"text/template"
)

// 生成docker
func CreateBuild(data *Data) error {
	// 生成Makefile
	var tpl = `GOPATH := $(shell go env GOPATH)
VERSION := 0.0.1.0

gengo:
	protoc -I. --proto_path ../application/{{.Service}}/proto:../pkg/proto \
 --go_out ../application/{{.Service}}/proto --go_opt paths=source_relative --go-grpc_out ../application/{{.Service}}/proto --go-grpc_opt paths=source_relative \
 --grpc-gateway_out ../application/{{.Service}}/proto --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=false {{.Service}}.proto
	protoc-go-inject-tag -input=../application/{{.Service}}/proto/{{.Service}}.pb.go

depend:
	go get ../...

build: depend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}} ../main.go

test:
	go test -v ../... -cover

docker: build health
	docker build -f ./Dockerfile -t registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION) ../
	#docker push registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION)
	rm -rf {{.Service}} grpc-health-probe

health:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o grpc-health-probe ../pkg/grpc-health-probe/main.go

helm:
	helm upgrade --install {{.Service}} ./helm/{{.Service}} --set image.tag=$(VERSION)

clean:
	docker rmi registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION)

k8s: docker
	kubectl set image deployment/{{.Service}} {{.Service}}=registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION)

.PHONY: gengo depend build test docker health deploy helm clean k8s
`

	t, err := template.New("makefile").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/build/"

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
ADD build/{{.Service}} build/grpc-health-probe /bin/
ENTRYPOINT [ "/bin/{{.Service}}" ]
`

	t, err = template.New("dockerfile").Parse(tpl)
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

	return nil
}
