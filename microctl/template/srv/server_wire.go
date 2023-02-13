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

// 生成server/server.go
func CreateServerWire(data *template.Data) error {
	var tpl = `//go:build wireinject
// +build wireinject

package server

import (
    "github.com/google/wire"

   {{if .MQ}}"github.com/imind-lab/micro/v2/broker"{{end}}
    "github.com/imind-lab/micro/v2/dao"
    "{{.Domain}}/{{.Repo}}/application/{{.Name}}/service"
    domain "{{.Domain}}/{{.Repo}}/domain/{{.Name}}"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}/persistence"
)

func Create{{.Service}}Service({{if .MQ}}bk broker.Broker{{end}}) *service.{{.Service}}Service {
    panic(wire.Build(dao.NewCache, dao.NewDatabase, dao.NewDao, persistence.New{{.Service}}Repository, domain.New{{.Service}}Domain, service.New{{.Service}}Service))
}
`

	path := "./" + data.Name + "/server/"
	name := "wire.go"

	return template.CreateFile(data, tpl, path, name)
}
