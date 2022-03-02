/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package srv

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成build/Dockerfile
func CreateBuildDockerfile(data *tpl.Data) error {
	var tpl = `FROM golang:1.17.7-alpine3.15 as builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add build-base gcc git openssh binutils-gold
WORKDIR /go/src/{{.Domain}}/{{.Project}}/{{.Service}}/
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o {{.Service}} main.go
RUN go get github.com/grpc-ecosystem/grpc-health-probe

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

	t, err := template.New("build_dockerfile").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/build/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "Dockerfile"

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
