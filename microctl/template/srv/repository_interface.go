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

// 生成domain/service.go
func CreateRepositoryInterface(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"context"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

type {{.Svc}}Repository interface {
	Create{{.Svc}}(ctx context.Context, m model.{{.Svc}}) (model.{{.Svc}}, error)
	Get{{.Svc}}ById(ctx context.Context, id int) (model.{{.Svc}}, error)
	Get{{.Svc}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]model.{{.Svc}}, int, error)
	Get{{.Svc}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]model.{{.Svc}}, int, error)
	Get{{.Svc}}List0Ids(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]int, int, error)
	Get{{.Svc}}List1Ids(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]int, int, error)
	Update{{.Svc}}Type(ctx context.Context, id, typ int) (int8, error)
	Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error)

	// This commentary is for scaffolding. Do not modify or delete it
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/repository/" + data.Service + "/"
	name := "repository.go"

	return template.CreateFile(data, tpl, path, name)
}
