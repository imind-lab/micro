/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
    "github.com/imind-lab/micro/v2/microctl/template"
)

// 生成Makefile
func CreateMakefile(data *template.Data) error {
    var tpl = `GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

ifndef LOCAL
LOCAL = true
endif

ifndef IMAGE_TAG
IMAGE_TAG = 0.0.1.1
endif

ifndef IMAGE_URL
IMAGE_URL = registry.cn-beijing.aliyuncs.com/imind/{{.Service}}-api
endif

ifndef NAMESPACE
NAMESPACE = default
endif

ifndef RPC_HOST
RPC_HOST := {{.Service}}-api.imind.tech
endif

{{.Service}}_names:={{.Service}}-api
{{.Service}}_path:=./application/{{.Service}}/proto

define process
	protoc -I. --proto_path $(1) --proto_path ./pkg/proto \
 --go_out $(1) --go_opt paths=source_relative \
 --go-grpc_out $(1) --go-grpc_opt paths=source_relative \
 --grpc-gateway_out $(1) --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=false $(2).proto
	microctl inject --path=$(1)/$(2).pb.go
	sed -i "" 's/,omitempty//g' $(1)/$(2).pb.go

endef

proto:
	$(foreach name,$({{.Service}}_names),$(call process,$({{.Service}}_path),$(name)))

wire:
	cd server && wire

depend:
	go get ./...

test:
	go test -v ./... -cover

build:
ifeq ($(LOCAL), false)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}}-api ./main.go
else
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -a -installsuffix cgo -o {{.Service}}-api ./main.go
endif

docker:
ifeq ($(LOCAL), false)
	cp $(GOPATH)/bin/grpc-health-probe .
endif
	docker build -f ./Dockerfile -t $(IMAGE_URL):$(IMAGE_TAG) .
ifeq ($(LOCAL), false)
	docker push $(IMAGE_URL):$(IMAGE_TAG)
endif

deploy:
	helm upgrade --install {{.Service}}-api ./deploy/helm/{{.Service}}-api --set image.repository=$(IMAGE_URL),image.tag=$(IMAGE_TAG),traefik.host=$(RPC_HOST) -n $(NAMESPACE)

run:
	go run main.go server

clean:
	docker rmi $(IMAGE_URL):$(IMAGE_TAG)

release:
	make proto
	make depend
	make wire
	make build
	make docker
	make deploy

all:
	make proto
	make depend
	make wire
	make test
	make build
	make docker
	make deploy

.PHONY: proto depend wire build test docker deploy run clean release all
`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/"
    name := "Makefile"

    return template.CreateFile(data, tpl, path, name)
}
