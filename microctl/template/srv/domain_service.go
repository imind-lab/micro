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
func CreateDomainService(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"context"
	"math"

	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/tracing"
	"github.com/pkg/errors"

	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

func (dm {{.Service}}Domain) Create{{.Svc}}(ctx context.Context, m model.{{.Svc}}) error {
	_, err := dm.repo.Create{{.Svc}}(ctx, m)
	return err
}

func (dm {{.Service}}Domain) Get{{.Svc}}ById(ctx context.Context, id int) (*{{.Service}}.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Get{{.Svc}}ById invoke")

	m, err := dm.repo.Get{{.Svc}}ById(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "{{.Service}}sDomain.Get{{.Svc}}sById")
	}
	return {{.Svc}}Out(m), nil
}

func (dm {{.Service}}Domain) Get{{.Svc}}List0(ctx context.Context, status, pageSize, pageNum int, desc bool) (*{{.Service}}.{{.Svc}}List, error) {
	list, total, err := dm.repo.Get{{.Svc}}List0(ctx, status, pageSize, pageNum, desc)
	if err != nil {
		return nil, err
	}
	{{.Service}}s := {{.Svc}}OutMap(list, {{.Svc}}Out)

	var totalPage int32 = 0
	if total == 0 {
		pageNum = 1
	} else {
		totalPage = int32(math.Ceil(float64(total) / float64(pageSize)))
	}
	{{.Service}}List := &{{.Service}}.{{.Svc}}List{}
	{{.Service}}List.Datalist = {{.Service}}s
	{{.Service}}List.Total = int32(total)
	{{.Service}}List.TotalPage = totalPage
	{{.Service}}List.CurPage = int32(pageNum)

	return {{.Service}}List, nil
}

// 疑问：中间时翻上一页
func (dm {{.Service}}Domain) Get{{.Svc}}List1(ctx context.Context, status, pageSize, lastId int, desc bool) (*{{.Service}}.{{.Svc}}List, error) {
	list, total, err := dm.repo.Get{{.Svc}}List1(ctx, status, pageSize, lastId, desc)
	if err != nil {
		return nil, err
	}
	{{.Service}}s := {{.Svc}}OutMap(list, {{.Svc}}Out)

	var totalPage int32 = 0
	if total > 0 {
		totalPage = int32(math.Ceil(float64(total) / float64(pageSize)))
	}
	{{.Service}}List := &{{.Service}}.{{.Svc}}List{}
	{{.Service}}List.Datalist = {{.Service}}s
	{{.Service}}List.Total = int32(total)
	{{.Service}}List.TotalPage = totalPage
	{{.Service}}List.CurPage = 1

	return {{.Service}}List, nil
}

func (dm {{.Service}}Domain) Update{{.Svc}}Status(ctx context.Context, id, status int) (int8, error) {
	return dm.repo.Update{{.Svc}}Status(ctx, id, status)
}

func (dm {{.Service}}Domain) Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error) {
	return dm.repo.Delete{{.Svc}}ById(ctx, id)
}
`

	t, err := template.New("domain_service").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/domain/" + data.Service + "/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + data.Service + ".go"

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
