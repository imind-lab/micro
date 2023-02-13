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

// 生成client/service.go
func CreateApplicationService(data *template.Data) error {
	var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/03/03
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
    "context"
    "fmt"

    sentinel "github.com/alibaba/sentinel-golang/api"
    "github.com/alibaba/sentinel-golang/core/base"
    "github.com/go-playground/validator/v10"
    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/status"
    "github.com/imind-lab/micro/v2/tracing"
    "go.uber.org/zap"

    {{.Package}}_api "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/{{.Name}}/proto"
    domain "{{.Domain}}/{{.Repo}}{{.Suffix}}/domain/{{.Name}}"
)

type {{.Service}}Service struct {
    {{.Package}}_api.Unimplemented{{.Service}}ServiceServer

    vd *validator.Validate
    dm domain.{{.Service}}Domain
}

// New{{.Service}}Service 创建用户服务实例
func New{{.Service}}Service(dm domain.{{.Service}}Domain) *{{.Service}}Service {
    svc := &{{.Service}}Service{
        dm: dm,
        vd: validator.New(),
    }
    return svc
}

// Create{{.Service}} 创建{{.Service}}
func (svc *{{.Service}}Service) Create{{.Service}}(ctx context.Context, req *{{.Package}}_api.Create{{.Service}}Request) (*{{.Package}}_api.Create{{.Service}}Response, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Create{{.Service}} request")

    rsp := &{{.Package}}_api.Create{{.Service}}Response{}

    err := svc.vd.Struct(req)
    if err != nil {

        if _, ok := err.(*validator.InvalidValidationError); ok {
            fmt.Println(err)
        }

        for _, err := range err.(validator.ValidationErrors) {
            fmt.Println(err.Namespace())
            fmt.Println(err.Field())
            fmt.Println(err.StructNamespace())
            fmt.Println(err.StructField())
            fmt.Println(err.Tag())
            fmt.Println(err.ActualTag())
            fmt.Println(err.Kind())
            fmt.Println(err.Type())
            fmt.Println(err.Value())
            fmt.Println(err.Param())
            fmt.Println()
        }

    }

    err = svc.vd.Var(req.Name, "required,email")
    if err != nil {
        logger.Error("无效的Name", zap.Any("name", req.Name))

        rsp.SetCode(status.InvalidParams, "无效的Name")
        return rsp, nil
    }

    err = svc.dm.Create{{.Service}}(ctx, req.Name, req.Type)
    e := status.AsError(err)
    rsp.SetCode(e.Code, e.Msg)
    return rsp, nil
}

// Get{{.Service}}ById 根据Id获取{{.Service}}
func (svc *{{.Service}}Service) Get{{.Service}}ById(ctx context.Context, req *{{.Package}}_api.Get{{.Service}}ByIdRequest) (*{{.Package}}_api.Get{{.Service}}ByIdResponse, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}ById request")

    rsp := &{{.Package}}_api.Get{{.Service}}ByIdResponse{}

    data, err := svc.dm.Get{{.Service}}ById(ctx, req.Id)

    e := status.AsError(err)
    if e.Code == status.Success {
        rsp.SetBody(status.Success, data)
    } else {
        rsp.SetCode(e.Code, e.Msg)
    }
    return rsp, nil
}

func (svc *{{.Service}}Service) Get{{.Service}}List0(ctx context.Context, req *{{.Package}}_api.Get{{.Service}}List0Request) (*{{.Package}}_api.Get{{.Service}}ListResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}List0 request")

    rsp := &{{.Package}}_api.Get{{.Service}}ListResponse{}

    err := svc.vd.Struct(req)
    if err != nil {
        if _, ok := err.(*validator.InvalidValidationError); ok {
            fmt.Println(err)
        }
    }

    rateEntry, rateError := sentinel.Entry("abcd", sentinel.WithTrafficType(base.Inbound))
    if rateError != nil {
        logger.Debug("Get{{.Service}}List0限流了")

        rsp.SetCode(status.SystemError, "Get{{.Service}}List0限流了")
        return rsp, nil
    }
    defer rateEntry.Exit()

    data, err := svc.dm.Get{{.Service}}List0(ctx, req.Type, req.PageSize, req.PageNum, req.IsDesc)

    e := status.AsError(err)
    if e.Code == status.Success {
        rsp.SetBody(status.Success, data)
    } else {
        rsp.SetCode(e.Code, e.Msg)
    }
    return rsp, nil
}

func (svc *{{.Service}}Service) Get{{.Service}}List1(ctx context.Context, req *{{.Package}}_api.Get{{.Service}}List1Request) (*{{.Package}}_api.Get{{.Service}}ListResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}List1 request")

    rsp := &{{.Package}}_api.Get{{.Service}}ListResponse{}

    rateEntry, rateError := sentinel.Entry("abcd", sentinel.WithTrafficType(base.Inbound))
    if rateError != nil {
        logger.Debug("Get{{.Service}}List1限流了")

        rsp.SetCode(status.SystemError, "Get{{.Service}}List1限流了")
        return rsp, nil
    }
    defer rateEntry.Exit()

    data, err := svc.dm.Get{{.Service}}List1(ctx, req.Type, req.PageSize, req.LastId, req.IsDesc)

    e := status.AsError(err)
    if e.Code == status.Success {
        rsp.SetBody(status.Success, data)
    } else {
        rsp.SetCode(e.Code, e.Msg)
    }
    return rsp, nil
}

func (svc *{{.Service}}Service) Update{{.Service}}Type(ctx context.Context, req *{{.Package}}_api.Update{{.Service}}TypeRequest) (*{{.Package}}_api.Update{{.Service}}TypeResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Update{{.Service}}Type request")

    rsp := &{{.Package}}_api.Update{{.Service}}TypeResponse{}

    err := svc.dm.Update{{.Service}}Type(ctx, req.Id, req.Type)
    e := status.AsError(err)
    rsp.SetCode(e.Code, e.Msg)
    return rsp, nil
}

func (svc *{{.Service}}Service) Delete{{.Service}}ById(ctx context.Context, req *{{.Package}}_api.Delete{{.Service}}ByIdRequest) (*{{.Package}}_api.Delete{{.Service}}ByIdResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Delete{{.Service}}ById request")

    rsp := &{{.Package}}_api.Delete{{.Service}}ByIdResponse{}

    err := svc.dm.Delete{{.Service}}ById(ctx, req.Id)

    e := status.AsError(err)
    rsp.SetCode(e.Code, e.Msg)

    return rsp, nil
}
func (svc *{{.Service}}Service) Get{{.Service}}ListByIds(ctx context.Context, req *{{.Package}}_api.Get{{.Service}}ListByIdsRequest) (*{{.Package}}_api.Get{{.Service}}ListByIdsResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}ListByIds request")

    rsp := &{{.Package}}_api.Get{{.Service}}ListByIdsResponse{}

    data, err := svc.dm.Get{{.Service}}ListByIds(ctx, req.Ids)

    e := status.AsError(err)
    if e.Code == status.Success {
        rsp.SetBody(status.Success, data)
    } else {
        rsp.SetCode(e.Code, e.Msg)
    }
    return rsp, nil
}
`

	path := "./" + data.Name + "-api/application/" + data.Name + "/service/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
