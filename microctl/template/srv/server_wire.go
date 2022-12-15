/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成server/server.go
func CreateServerWire(data *template.Data) error {
	var tpl = `//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/service"
	domain "{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/persistence"{{if .MQ}}
	"github.com/imind-lab/micro/broker"                            {{end}}
	"github.com/imind-lab/micro/dao"
)

func Create{{.Svc}}Service({{if .MQ}}bk broker.Broker{{end}}) *service.{{.Svc}}Service {
	panic(wire.Build(dao.NewCache, dao.NewDatabase, dao.NewDao, persistence.New{{.Svc}}Repository, domain.New{{.Svc}}Domain, service.New{{.Svc}}Service))
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/server/"
	name := "wire.go"

	return template.CreateFile(data, tpl, path, name)
}
