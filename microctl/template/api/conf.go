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

// 生成conf/conf.yaml
func CreateConf(data *tpl.Data) error {
	var tpl = `service:
  namespace: {{.Project}}
  name: {{.Service}}-api
  version: latest
  logLevel: -2
  port: #监听端口
    http: 88
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

tracing:
  agent: '172.16.50.50:6831'
  type: const
  param: 1
  name:
    client: {{.Project}}-{{.Service}}-api-cli
    server: {{.Project}}-{{.Service}}-api-srv

log:
  path: './logs/ms.log'
  level: -1
  age: 7
  size: 128
  backup: 30
  compress: true
  format: json

rpc:
  {{.Service}}:
    service: {{.Service}}
    port: 50051

`

	t, err := template.New("conf_conf").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/conf/"

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
