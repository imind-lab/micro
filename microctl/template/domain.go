/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
	"os"
	"text/template"
)

// 生成domain
func CreateDomain(data *Data) error {
	var tpl = `/**
 *  IMindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
	"context"
	"math"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	"{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}/repository"
	"{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}/repository/model"
	"{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}/repository/persistence"
	"github.com/imind-lab/micro/dao"
)

type {{.Svc}}Domain interface {
	Create{{.Svc}}(ctx context.Context, dto *{{.Service}}.{{.Svc}}) error

	Get{{.Svc}}ById(ctx context.Context, id int32) (*{{.Service}}.{{.Svc}}, error)
	Get{{.Svc}}List(ctx context.Context, status, lastId, pageSize, page int32) (*{{.Service}}.{{.Svc}}List, error)

	Update{{.Svc}}Status(ctx context.Context, id, status int32) (int64, error)
	Update{{.Svc}}Count(ctx context.Context, id, num int32, column string) (int64, error)

	Delete{{.Svc}}ById(ctx context.Context, id int32) (int64, error)
}

type {{.Service}}Domain struct {
	dao.Cache

	repo repository.{{.Svc}}Repository
}

func New{{.Svc}}Domain() {{.Svc}}Domain {
	repo := persistence.New{{.Svc}}Repository()
	dm := {{.Service}}Domain{
		Cache: dao.NewCache(),
		repo:  repo}
	return dm
}

func (dm {{.Service}}Domain) Create{{.Svc}}(ctx context.Context, dto *{{.Service}}.{{.Svc}}) error {
	m := {{.Svc}}Dto2Model(dto)
	_, err := dm.repo.Create{{.Svc}}(ctx, m)
	return err
}

func (dm {{.Service}}Domain) Get{{.Svc}}ById(ctx context.Context, id int32) (*{{.Service}}.{{.Svc}}, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Service}}Domain"), zap.String("func", "Get{{.Svc}}ById"))

	logger.Info("{{.Service}}Domain.Get{{.Svc}}ById invoke")
	m, err := dm.repo.Get{{.Svc}}ById(ctx, id)
	return {{.Svc}}Model2Dto(m), errors.WithMessage(err, "{{.Service}}Domain.Get{{.Svc}}ById")
}

func (dm {{.Service}}Domain) Get{{.Svc}}List(ctx context.Context, status, lastId, pageSize, page int32) (*{{.Service}}.{{.Svc}}List, error) {
	list, total, err := dm.repo.Get{{.Svc}}List(ctx, status, lastId, pageSize, page)
	if err != nil {
		return nil, err
	}
	{{.Service}}s := {{.Svc}}Map(list, {{.Svc}}Model2Dto)

	var totalPage int32 = 0
	if total == 0 {
		page = 1
	} else {
		totalPage = int32(math.Ceil(float64(total) / float64(pageSize)))
	}
	{{.Service}}List := &{{.Service}}.{{.Svc}}List{}
	{{.Service}}List.Datalist = {{.Service}}s
	{{.Service}}List.Total = int32(total)
	{{.Service}}List.TotalPage = totalPage
	{{.Service}}List.CurPage = page

	return {{.Service}}List, nil
}

func (dm {{.Service}}Domain) Update{{.Svc}}Status(ctx context.Context, id, status int32) (int64, error) {
	return dm.repo.Update{{.Svc}}Status(ctx, id, status)
}

func (dm {{.Service}}Domain) Update{{.Svc}}Count(ctx context.Context, id, num int32, column string) (int64, error) {
	return dm.repo.Update{{.Svc}}Count(ctx, id, num, column)
}

func (dm {{.Service}}Domain) Delete{{.Svc}}ById(ctx context.Context, id int32) (int64, error) {
	return dm.repo.Delete{{.Svc}}ById(ctx, id)
}

func {{.Svc}}Map(pos []model.{{.Svc}}, fn func(model.{{.Svc}}) *{{.Service}}.{{.Svc}}) []*{{.Service}}.{{.Svc}} {
	var dtos []*{{.Service}}.{{.Svc}}
	for _, po := range pos {
		dtos = append(dtos, fn(po))
	}
	return dtos
}

func {{.Svc}}Model2Dto(po model.{{.Svc}}) *{{.Service}}.{{.Svc}} {
	if po.IsEmpty() {
		return nil
	}

	dto := &{{.Service}}.{{.Svc}}{}
	dto.Id = po.Id
	dto.Name = po.Name
	dto.ViewNum = po.ViewNum
	dto.Status = po.Status
	dto.CreateTime = po.CreateTime
	dto.UpdateDatetime = po.UpdateDatetime
	dto.CreateDatetime = po.CreateDatetime

	return dto
}

func {{.Svc}}Dto2Model(dto *{{.Service}}.{{.Svc}}) model.{{.Svc}} {
	if dto == nil {
		return model.{{.Svc}}{}
	}

	po := model.{{.Svc}}{}
	po.Id = dto.Id
	po.Name = dto.Name
	po.ViewNum = dto.ViewNum
	po.Status = dto.Status
	po.CreateTime = dto.CreateTime
	po.UpdateDatetime = dto.UpdateDatetime
	po.CreateDatetime = dto.CreateDatetime

	return po
}
`

	t, err := template.New("domain").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/domain/" + data.Service + "/service/"

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
