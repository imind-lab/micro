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

// 生成domain/service.go
func CreateRepositoryInterface(data *tpl.Data) error {
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
	Get{{.Svc}}List0(ctx context.Context, status, pageSize, pageNum int, desc bool) ([]model.{{.Svc}}, int, error)
	Get{{.Svc}}List1(ctx context.Context, status, pageSize, lastId int, desc bool) ([]model.{{.Svc}}, int, error)
	Get{{.Svc}}List0Ids(ctx context.Context, status, pageSize, pageNum int, desc bool) ([]int, int, error)
	Get{{.Svc}}List1Ids(ctx context.Context, status, pageSize, lastId int, desc bool) ([]int, int, error)
	Update{{.Svc}}Status(ctx context.Context, id, status int) (int8, error)
	Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error)

	// This commentary is for scaffolding. Do not modify or delete it
}
`

	t, err := template.New("domain_service").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/repository/" + data.Service + "/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "repository.go"

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
