/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成build/Makefile
func CreateBuildMakefile(data *tpl.Data) error {
	var tpl = `GOPATH := $(shell go env GOPATH)
VERSION := 0.0.1.4

names:={{.Service}}-api
path:=../application/{{.Service}}/proto

define process
	protoc -I. --proto_path $(path) --proto_path ../pkg/proto \
 --go_out $(path) --go_opt paths=source_relative --go-grpc_out $(path) --go-grpc_opt paths=source_relative \
 --grpc-gateway_out $(path) --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=false $(1).proto
	microctl inject --path=$(path)/$(1).pb.go
	sed -i '' 's/,omitempty//g' $(path)/$(1).pb.go

endef

gengo:
	$(foreach name,$(names),$(call process,$(name)))

depend: gengo
	go get ../...

build: depend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}}-api ../main.go

test:
	go test -v ../... -cover

docker:
	docker build -f ./Dockerfile -t registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION) ../
	docker push registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION)

deploy: docker
	@helm upgrade --install {{.Service}}-api ../deploy/helm/{{.Service}}-api --set image.tag=$(VERSION) -n micro

clean:
	docker rmi registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION)

k8s:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}}-api ../main.go
	docker build -f ./DockerfileDev -t registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION) ../
	helm upgrade --install {{.Service}}-api ../deploy/helm/{{.Service}}-api --set image.tag=$(VERSION)
	#kubectl set image deployment/{{.Service}}-api {{.Service}}=registry.cn-beijing.aliyuncs.com/{{.Project}}/{{.Service}}-api:$(VERSION)

.PHONY: gengo depend build test docker deploy clean k8s`

	t, err := template.New("build_dockerfile").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
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

	return nil
}
