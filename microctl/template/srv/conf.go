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
  namespace: {{.Project}}
  name: {{.Service}}
  version: latest
  logLevel: -1
  logFormat: json
  port: #监听端口
    http: 80
    grpc: 50051


db:
  logLevel: 4
  max:
    open: 10
    idle: 5
    life: 30
  timeout: 5s
  default:
    master:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: imind
      name: imind
    replica:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: imind
      name: imind

redis:
  model: node
  timeout: 5s
  addr: '127.0.0.1:6379'
#  pass: imind
  db: 0

kafka:
  default:
    producer:
      - 'kafka.infra:9092'
    consumer:
      - 'kafka.infra:9092'
    topic:
      {{.Service}}Create: {{.Service}}_create
      {{.Service}}Update: {{.Service}}_update

tracing:
  agent:
    host: '127.0.0.1'
    port: 6831`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/conf/"
    name := "conf.yaml"

    return template.CreateFile(data, tpl, path, name)
}
