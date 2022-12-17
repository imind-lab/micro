/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package api

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成Dockerfile
func CreateDockerfile(data *template.Data) error {
	var tpl = `FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

WORKDIR .
ADD conf /conf
COPY {{.Service}}-api grpc-health-probe /bin/
ENTRYPOINT [ "/bin/{{.Service}}-api", "server" ]
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/"
	name := "Dockerfile"

	return template.CreateFile(data, tpl, path, name)
}
