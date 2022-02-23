/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package srv

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成docker
func CreateBuild(data *tpl.Data) error {
	// 生成Makefile
	var tpl = `GOPATH := $(shell go env GOPATH)
VERSION := 0.0.1.0

gengo:
	protoc -I. --proto_path ../application/{{.Service}}/proto --proto_path ../pkg/proto \
 --go_out ../application/{{.Service}}/proto --go_opt paths=source_relative --go-grpc_out ../application/{{.Service}}/proto --go-grpc_opt paths=source_relative \
 --grpc-gateway_out ../application/{{.Service}}/proto --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=false {{.Service}}.proto
	@microctl inject --path=../application/{{.Service}}/proto/{{.Service}}.pb.go
	@sed -i '' 's/,omitempty//g' ../application/{{.Service}}/proto/{{.Service}}.pb.go

depend: gengo
	go get ../...

build: depend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}} ../main.go

test:
	go test -v ../... -cover

docker: gengo
	docker build -f ./Dockerfile -t registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION) ../
	#docker push registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION)

deploy: docker
	@helm upgrade --install {{.Service}} ../deploy/helm/{{.Service}} --set image.tag=$(VERSION)

clean:
	docker rmi registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION)

k8s: docker
	kubectl set image deployment/{{.Service}} {{.Service}}=registry.cn-beijing.aliyuncs.com/imind/{{.Service}}:$(VERSION)

.PHONY: gengo depend build test docker deploy clean k8s
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
	tpl = `FROM golang:1.17.7-alpine3.15 as builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add build-base gcc git openssh binutils-gold
WORKDIR /go/src/{{.Domain}}/{{.Project}}/{{.Service}}/
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}} main.go
RUN go get {{.Domain}}/grpc-ecosystem/grpc-health-probe

FROM alpine:3.15
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

WORKDIR .
ADD conf /conf
COPY --from=builder /go/src/{{.Domain}}/{{.Project}}/{{.Service}}/{{.Service}} /go/bin/grpc-health-probe /bin/
ENTRYPOINT [ "/bin/{{.Service}}", "server" ]
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
