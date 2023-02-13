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

// 生成server/server.go
func CreateServerWire(data *template.Data) error {
    var tpl = `//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/imind-lab/micro/dao"
	"gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/application/{{.Service}}/service"
	domain "gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/domain/{{.Service}}"
	"gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/repository/{{.Service}}/persistence"
)

func Create{{.Svc}}Service() *service.{{.Svc}}Service {
	panic(wire.Build(dao.NewCache, dao.NewDatabase, dao.NewDao, persistence.New{{.Svc}}Repository, domain.New{{.Svc}}Domain, service.New{{.Svc}}Service))
}
`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/server/"
    name := "wire.go"

    return template.CreateFile(data, tpl, path, name)
}
