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

// 生成conf/conf.yaml
func CreateConf(data *tpl.Data) error {
	var tpl = `service:
  namespace: {{.Project}}
  name: {{.Service}}
  version: latest
  logLevel: -2
  port: #监听端口
    http: 80
    grpc: 50051
  rate:
    high:
      limit: 10
      capacity: 10
    low:
      limit: 10
      capacity: 10
  profile:
    rate: 1

db:
  logLevel: 4
  max:
    open: 10
    idle: 5
    life: 30
  timeout: 5s
  imind:
    master:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: imind123
      name: {{.Service}}
    replica:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: imind123
      name: {{.Service}}

redis:
  model: node
  timeout: 5s
  addr: '127.0.0.1:56627'
  pass: imind456
  db: 0

kafka:
  business:
    producer:
      - '127.0.0.1:9092'
    consumer:
      - '127.0.0.1:9092'
    topic:
      {{.Service}}Create: {{.Service}}_create
      {{.Service}}Update: {{.Service}}_update

tracing:
  agent:
    host: '127.0.0.1'
    port: 6831

log:
  path: './logs/ms.log'
  level: -1
  age: 7
  size: 128
  backup: 30
  compress: true
  format: json
`

	t, err := template.New("conf_conf").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/conf/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "conf.yaml"

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
