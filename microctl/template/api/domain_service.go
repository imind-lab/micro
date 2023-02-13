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

package {{.Package}}

import (
    "context"
    "github.com/pkg/errors"

    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/tracing"
    "github.com/imind-lab/micro/v2/util"

    {{.Package}}_api "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/{{.Name}}/proto"
)

func (dm {{.Svc}}Domain) Create{{.Service}}(ctx context.Context, name string, typ int32) error {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Create{{.Service}} invoke")

    err := dm.repo.Create{{.Service}}(ctx, name, typ)
    if err != nil {
        return errors.WithMessage(err, util.GetFuncName())
    }
    return nil
}

func (dm {{.Svc}}Domain) Get{{.Service}}ById(ctx context.Context, id int32) (*{{.Package}}_api.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Get{{.Service}}ById invoke")

    m, err := dm.repo.Get{{.Service}}ById(ctx, id)
    if err != nil {
        return nil, errors.WithMessage(err, util.GetFuncName())
    }
    return {{.Service}}Srv2Api(m), nil
}

func (dm {{.Svc}}Domain) Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Package}}_api.{{.Service}}List, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Get{{.Service}}List0 invoke")
    list, err := dm.repo.Get{{.Service}}List0(ctx, typ, pageSize, pageNum, isDesc)
    if err != nil {
        return nil, errors.WithMessage(err, util.GetFuncName())
    }
    return {{.Service}}ListSrv2Api(list), nil
}

// 疑问：中间时翻上一页
func (dm {{.Svc}}Domain) Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Package}}_api.{{.Service}}List, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Get{{.Service}}List1 invoke")
    list, err := dm.repo.Get{{.Service}}List1(ctx, typ, pageSize, lastId, isDesc)
    if err != nil {
        return nil, errors.WithMessage(err, util.GetFuncName())
    }
    return {{.Service}}ListSrv2Api(list), nil
}

func (dm {{.Svc}}Domain) Update{{.Service}}Type(ctx context.Context, id, typ int32) error {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Update{{.Service}}Type invoke")

    err := dm.repo.Update{{.Service}}Type(ctx, id, typ)
    if err != nil {
        return errors.WithMessage(err, util.GetFuncName())
    }
    return nil
}

func (dm {{.Svc}}Domain) Delete{{.Service}}ById(ctx context.Context, id int32) error {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Delete{{.Service}}ById invoke")

    err := dm.repo.Delete{{.Service}}ById(ctx, id)
    if err != nil {
        return errors.WithMessage(err, util.GetFuncName())
    }
    return nil
}

func (dm {{.Svc}}Domain) Get{{.Service}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Package}}_api.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("{{.Svc}}Domain.Get{{.Service}}ListByIds invoke")
    list, err := dm.repo.Get{{.Service}}ListByIds(ctx, ids)
    if err != nil {
        return nil, errors.WithMessage(err, util.GetFuncName())
    }
    return {{.Service}}Map(list, {{.Service}}Srv2Api), nil
}
`

	path := "./" + data.Name + "-api/domain/" + data.Name + "/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
