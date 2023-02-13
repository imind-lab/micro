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

// 生成conf/conf.yaml
func CreateConf(data *template.Data) error {
    var tpl = `global:
  rate:
    high:
      limit: 10
      capacity: 10
    low:
      limit: 10
      capacity: 10
  profile:
    rate: 1

service:
  namespace: {{.Repo}}
  name: {{.Service}}-api
  version: latest
  logLevel: -1
  logFormat: json
  port: #监听端口
    http: 8080
    grpc: 50052

tracing:
  agent:
    host: '127.0.0.1'
    port: 6831

rpc:
  {{.Service}}:
    service: 127.0.0.1
    port: 50051

`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/conf/"
    name := "conf.yaml"

    return template.CreateFile(data, tpl, path, name)
}
