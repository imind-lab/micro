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

// 生成go.mod
func CreateMod(data *tpl.Data) error {
	var tpl = `module {{.Domain}}/{{.Project}}/{{.Service}}-api

go 1.16

replace {{.Domain}}/{{.Project}}/{{.Service}} => ../../../{{.Domain}}/{{.Project}}/{{.Service}}

replace github.com/imind-lab/micro => ../../../github.com/imind-lab/micro

require (
	github.com/alibaba/sentinel-golang v1.0.4
	github.com/go-playground/validator/v10 v10.10.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.8.0
	github.com/imind-lab/micro v0.0.0-{{.Year}}0502151153-ee574a04d410
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.10.1
	{{.Domain}}/{{.Project}}/{{.Service}} v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-{{.Year}}0127200216-cd36cc0744dd
	google.golang.org/genproto v0.0.0-{{.Year}}0314164441-57ef72a4c106
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
)
`

	t, err := template.New("main").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "go.mod"

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
