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

// 生成go.mod
func CreateMod(data *tpl.Data) error {
	var tpl = `module {{.Domain}}/{{.Project}}/{{.Service}}

go 1.16

replace (
	github.com/imind-lab/micro => ../../../github.com/imind-lab/micro
)

require (
	github.com/go-playground/validator/v10 v10.10.0
	github.com/go-redis/redis/v8 v8.11.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/imind-lab/micro v0.0.0-{{.Year}}0224071238-db0be8430c34
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	go.uber.org/zap v1.19.1
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.15
)
`

	t, err := template.New("main").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/"

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
