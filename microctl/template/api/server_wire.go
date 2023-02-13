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

    "github.com/imind-lab/micro/v2/dao"
    "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/{{.Name}}/service"
    domain "{{.Domain}}/{{.Repo}}{{.Suffix}}/domain/{{.Name}}"
    "{{.Domain}}/{{.Repo}}{{.Suffix}}/repository/{{.Name}}/persistence"
)

func Create{{.Service}}Service() *service.{{.Service}}Service {
    panic(wire.Build(dao.NewCache, dao.NewDatabase, dao.NewDao, persistence.New{{.Service}}Repository, domain.New{{.Service}}Domain, service.New{{.Service}}Service))
}
`

	path := "./" + data.Name + "-api/server/"
	name := "wire.go"

	return template.CreateFile(data, tpl, path, name)
}
