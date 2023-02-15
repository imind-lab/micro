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

// 生成domain/service.go
func CreateRepositoryInterface(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Package}}

import (
	"context"

	"{{.Domain}}/{{.Repo}}/repository/{{.Name}}/model"
)

type {{.Service}}Repository interface {
	Create{{.Service}}(ctx context.Context, m model.{{.Service}}) (model.{{.Service}}, error)
	Get{{.Service}}ById(ctx context.Context, id int) (model.{{.Service}}, error)
	Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]model.{{.Service}}, int, error)
	Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]model.{{.Service}}, int, error)
	Get{{.Service}}List0Ids(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]int, int, error)
	Get{{.Service}}List1Ids(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]int, int, error)
	Update{{.Service}}Type(ctx context.Context, id, typ int) (int8, error)
	Delete{{.Service}}ById(ctx context.Context, id int) (int8, error)

	//+IMind:scaffold
}
`

	path := "./" + data.Name + "/repository/" + data.Name + "/"
	name := "repository.go"

	return template.CreateFile(data, tpl, path, name)
}
