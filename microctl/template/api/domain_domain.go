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

package {{.Package}}

import (
    "context"
    {{.Package}}_api "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/{{.Name}}/proto"
    repository "{{.Domain}}/{{.Repo}}{{.Suffix}}/repository/{{.Name}}"
)

type {{.Service}}Domain interface {
    Create{{.Service}}(ctx context.Context, name string, typ int32) error
    Get{{.Service}}ById(ctx context.Context, id int32) (*{{.Package}}_api.{{.Service}}, error)
    Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Package}}_api.{{.Service}}List, error)
    Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Package}}_api.{{.Service}}List, error)
    Update{{.Service}}Type(ctx context.Context, id, typ int32) error
    Delete{{.Service}}ById(ctx context.Context, id int32) error
    Get{{.Service}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Package}}_api.{{.Service}}, error)

    // This commentary is for scaffolding. Do not modify or delete it
}

type {{.Svc}}Domain struct {
    repo repository.{{.Service}}Repository
}

func New{{.Service}}Domain(repo repository.{{.Service}}Repository) {{.Service}}Domain {
    dm := {{.Svc}}Domain{
        repo: repo,
    }
    return dm
}
`

	path := "./" + data.Name + "-api/domain/" + data.Name + "/"
	name := "domain.go"

	return template.CreateFile(data, tpl, path, name)
}
