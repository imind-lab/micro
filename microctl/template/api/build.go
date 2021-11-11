/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package api

import (
	"os"
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

depend: gengo
	go get ../...

build: depend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}}-api ../main.go

test:
	go test -v ../... -cover

docker: build health
	docker build -f ./Dockerfile -t 348681422678.dkr.ecr.ap-southeast-1.amazonaws.com/{{.Project}}/{{.Service}}-api:$(VERSION) ../
	docker push 348681422678.dkr.ecr.ap-southeast-1.amazonaws.com/{{.Project}}/{{.Service}}-api:$(VERSION)
	rm -rf {{.Service}}-api grpc-health-probe

health:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o grpc-health-probe ../pkg/grpc-health-probe/main.go

deploy: docker
	helm upgrade {{.Service}}-api ../deploy/helm/{{.Service}}-api --set image.tag=$(VERSION)

clean:
	docker rmi 348681422678.dkr.ecr.ap-southeast-1.amazonaws.com/{{.Project}}/{{.Service}}-api:$(VERSION)

k8s:
	kubectl set image deployment/{{.Service}}-api {{.Service}}=348681422678.dkr.ecr.ap-southeast-1.amazonaws.com/{{.Project}}/{{.Service}}-api:$(VERSION)

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
COPY build/{{.Service}}-api build/grpc-health-probe /bin/
ENTRYPOINT [ "/bin/{{.Service}}-api", "server" ]
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

	return nil
}
