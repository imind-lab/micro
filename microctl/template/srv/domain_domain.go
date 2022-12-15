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

// 生成domain/domain.go
func CreateDomainDomain(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on 2021/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"context"
	"github.com/imind-lab/micro/dao"

	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	repository "{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

type {{.Svc}}Domain interface {
	Create{{.Svc}}(ctx context.Context, m model.{{.Svc}}) error

	Get{{.Svc}}ById(ctx context.Context, id int) (*{{.Service}}.{{.Svc}}, error)
	Get{{.Svc}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) (*{{.Service}}.{{.Svc}}List, error)
	Get{{.Svc}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) (*{{.Service}}.{{.Svc}}List, error)

	Update{{.Svc}}Type(ctx context.Context, id, typ int) (int8, error)
	Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error)

	// This commentary is for scaffolding. Do not modify or delete it
}

type {{.Service}}Domain struct {
	dao.Cache
	repo repository.{{.Svc}}Repository
}

func New{{.Svc}}Domain(repo repository.{{.Svc}}Repository) {{.Svc}}Domain {
	dm := {{.Service}}Domain{
		Cache: dao.NewCache(),
		repo:  repo,
	}
	return dm
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/domain/" + data.Service + "/"
	name := "domain.go"

	return template.CreateFile(data, tpl, path, name)
}
