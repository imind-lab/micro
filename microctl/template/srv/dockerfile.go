/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
    "github.com/imind-lab/micro/v2/microctl/template"
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
COPY {{.Service}} grpc-health-probe /bin/
ENTRYPOINT [ "/bin/{{.Service}}", "server" ]
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/"
    name := "Dockerfile"

    return template.CreateFile(data, tpl, path, name)
}
