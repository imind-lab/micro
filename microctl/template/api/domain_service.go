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

// 生成domain/service.go
func CreateDomainService(data *template.Data) error {
    var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"context"
	"github.com/pkg/errors"

	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/tracing"
	"github.com/imind-lab/micro/util"

	{{.Service}}_api "gitlab.imind.tech/{{.Project}}/{{.Service}}-api/application/{{.Service}}/proto"
)

func (dm {{.Service}}Domain) Create{{.Svc}}(ctx context.Context, name string, typ int32) error {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Create{{.Svc}} invoke")

	err := dm.repo.Create{{.Svc}}(ctx, name, typ)
	if err != nil {
		return errors.WithMessage(err, util.GetFuncName())
	}
	return nil
}

func (dm {{.Service}}Domain) Get{{.Svc}}ById(ctx context.Context, id int32) (*{{.Service}}_api.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Get{{.Svc}}ById invoke")

	m, err := dm.repo.Get{{.Svc}}ById(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, util.GetFuncName())
	}
	return {{.Svc}}Srv2Api(m), nil
}

func (dm {{.Service}}Domain) Get{{.Svc}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Service}}_api.{{.Svc}}List, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Get{{.Svc}}List0 invoke")
	list, err := dm.repo.Get{{.Svc}}List0(ctx, typ, pageSize, pageNum, isDesc)
	if err != nil {
		return nil, errors.WithMessage(err, util.GetFuncName())
	}
	return {{.Svc}}ListSrv2Api(list), nil
}

// 疑问：中间时翻上一页
func (dm {{.Service}}Domain) Get{{.Svc}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Service}}_api.{{.Svc}}List, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Get{{.Svc}}List1 invoke")
	list, err := dm.repo.Get{{.Svc}}List1(ctx, typ, pageSize, lastId, isDesc)
	if err != nil {
		return nil, errors.WithMessage(err, util.GetFuncName())
	}
	return {{.Svc}}ListSrv2Api(list), nil
}

func (dm {{.Service}}Domain) Update{{.Svc}}Type(ctx context.Context, id, typ int32) error {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Update{{.Svc}}Type invoke")

	err := dm.repo.Update{{.Svc}}Type(ctx, id, typ)
	if err != nil {
		return errors.WithMessage(err, util.GetFuncName())
	}
	return nil
}

func (dm {{.Service}}Domain) Delete{{.Svc}}ById(ctx context.Context, id int32) error {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Delete{{.Svc}}ById invoke")

	err := dm.repo.Delete{{.Svc}}ById(ctx, id)
	if err != nil {
		return errors.WithMessage(err, util.GetFuncName())
	}
	return nil
}

func (dm {{.Service}}Domain) Get{{.Svc}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Service}}_api.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Info("{{.Service}}Domain.Get{{.Svc}}ListByIds invoke")
	list, err := dm.repo.Get{{.Svc}}ListByIds(ctx, ids)
	if err != nil {
		return nil, errors.WithMessage(err, util.GetFuncName())
	}
	return {{.Svc}}Map(list, {{.Svc}}Srv2Api), nil
}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/domain/" + data.Service + "/"
    name := data.Service + ".go"

    return template.CreateFile(data, tpl, path, name)
}
