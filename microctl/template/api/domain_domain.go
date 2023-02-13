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

// 生成domain/domain.go
func CreateDomainDomain(data *template.Data) error {
    var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"context"
	{{.Service}}_api "gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/application/{{.Service}}/proto"
	repository "gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/repository/{{.Service}}"
)

type {{.Svc}}Domain interface {
	Create{{.Svc}}(ctx context.Context, name string, typ int32) error
	Get{{.Svc}}ById(ctx context.Context, id int32) (*{{.Service}}_api.{{.Svc}}, error)
	Get{{.Svc}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Service}}_api.{{.Svc}}List, error)
	Get{{.Svc}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Service}}_api.{{.Svc}}List, error)
	Update{{.Svc}}Type(ctx context.Context, id, typ int32) error
	Delete{{.Svc}}ById(ctx context.Context, id int32) error
	Get{{.Svc}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Service}}_api.{{.Svc}}, error)

	// This commentary is for scaffolding. Do not modify or delete it
}

type {{.Service}}Domain struct {
	repo repository.{{.Svc}}Repository
}

func New{{.Svc}}Domain(repo repository.{{.Svc}}Repository) {{.Svc}}Domain {
	dm := {{.Service}}Domain{
		repo: repo,
	}
	return dm
}
`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/domain/" + data.Service + "/"
    name := "domain.go"

    return template.CreateFile(data, tpl, path, name)
}
