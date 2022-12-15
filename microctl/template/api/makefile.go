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

// 生成Makefile
func CreateMakefile(data *tpl.Data) error {
	var tpl = `GOARCH := $(shell go env GOARCH)

ifdef CI_COMMIT_SHORT_SHA
TAG := $(CI_COMMIT_SHORT_SHA)
else
TAG := 0.0.1.0
endif

ifdef IMAGE_URL
URL := $(IMAGE_URL)
else
URL := 348681422678.dkr.ecr.ap-southeast-1.amazonaws.com/backend/{{.Service}}-api
endif

ifdef IMAGE_SERVER
SERVER := $(IMAGE_SERVER)
else
SERVER := 348681422678.dkr.ecr.ap-southeast-1.amazonaws.com
endif

ifdef PROFILE
AWS_PROFILE := $(PROFILE)
else
AWS_PROFILE := uat-sg-profile
endif

ifdef REGION
AWS_REGION := $(REGION)
else
AWS_REGION := ap-southeast-1
endif

ifdef EKS
AWS_EKS := $(EKS)
else
AWS_EKS := uat-sg-eks-cluster
endif

ifdef NAMESPACE
NS := $(NAMESPACE)
else
NS := rainbow
endif

names:={{.Service}}-api
path:=./application/{{.Service}}/proto

define process
	protoc -I. --proto_path $(path) --proto_path ./pkg/proto \
 --go_out $(path) --go_opt paths=source_relative \
 --go-grpc_out $(path) --go-grpc_opt paths=source_relative \
 --grpc-gateway_out $(path) --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=false $(1).proto
	microctl inject --path=$(path)/$(1).pb.go
	sed -i '' 's/,omitempty//g' $(path)/$(1).pb.go

endef

gengo:
	$(foreach name,$(names),$(call process,$(name)))

depend:
	go get ./...

test:
	go test -v ./... -cover

build: depend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}}-api ./main.go

docker:
	docker build -f ./Dockerfile -t $(URL):$(TAG) .
	aws --profile $(AWS_PROFILE) --region $(AWS_REGION) ecr get-login-password | docker login --password-stdin  --username AWS $(SERVER)
	#docker push $(URL):$(TAG)

deploy: docker
	helm upgrade --install {{.Service}}-api ./deploy/helm/{{.Service}}-api --set image.tag=$(TAG),image.repository=$(URL) -n $(NS)

run:
	go run main.go server

clean:
	docker rmi $(URL):$(TAG)

k8s:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -a -installsuffix cgo -o {{.Service}}-api ./main.go
	docker build -f ./DockerfileDev -t $(URL):$(TAG) ./
	kubectl set image deployment/{{.Service}}-api {{.Service}}=$(URL):$(TAG)
	rm {{.Service}}-api

.PHONY: gengo depend build test docker deploy run clean k8s
`

	t, err := template.New("dockerfile").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/"

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
