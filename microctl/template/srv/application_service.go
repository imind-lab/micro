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

// 生成client/service.go
func CreateApplicationService(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
    "context"
    "errors"
    "fmt"
    "io"

    "github.com/go-playground/validator/v10"
    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/status"
    "github.com/imind-lab/micro/v2/tracing"
    "go.uber.org/zap"

    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
    domain "{{.Domain}}/{{.Repo}}/domain/{{.Name}}"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}/model"
)

type {{.Service}}Service struct {
    {{.Package}}.Unimplemented{{.Service}}ServiceServer

    vd *validator.Validate
    dm domain.{{.Service}}Domain
}

func New{{.Service}}Service(dm domain.{{.Service}}Domain) *{{.Service}}Service {
    svc := &{{.Service}}Service{
        dm: dm,
        vd: validator.New(),
    }

    return svc
}

// Create{{.Service}} 创建{{.Service}}
func (svc *{{.Service}}Service) Create{{.Service}}(ctx context.Context, req *{{.Package}}.Create{{.Service}}Request) (*{{.Package}}.Create{{.Service}}Response, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Create{{.Service}} request")

    rsp := &{{.Package}}.Create{{.Service}}Response{}

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
        logger.Error("Name不能为空", zap.Any("name", req.Name), zap.Error(err))
        rsp.SetCode(status.InvalidParams, "Name不能为空")
        return rsp, nil
    }
    m := model.{{.Service}}{
        Name: req.Name,
        Type: int8(req.Type),
    }
    err = svc.dm.Create{{.Service}}(ctx, m)
    if err != nil {
        msg := "创建{{.Service}}失败"
        logger.Error(msg, zap.Any("sample", m), zap.Error(err))
        var state status.Error
        if errors.As(err, &state) {
            rsp.SetCode(state.Code, state.Msg)
        } else {
            rsp.SetCode(status.DBQueryFailed, msg)
        }
    }
    rsp.SetCode(status.Success)
    return rsp, nil
}

// Get{{.Service}}ById 根据Id获取{{.Service}}
func (svc *{{.Service}}Service) Get{{.Service}}ById(ctx context.Context, req *{{.Package}}.Get{{.Service}}ByIdRequest) (*{{.Package}}.Get{{.Service}}ByIdResponse, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}ById request")

    rsp := &{{.Package}}.Get{{.Service}}ByIdResponse{}
    m, err := svc.dm.Get{{.Service}}ById(ctx, int(req.Id))
    if err != nil {
        msg := "获取{{.Service}}失败"
        logger.Error(msg, zap.Any("sample", m), zap.Error(err))
        var state status.Error
        if errors.As(err, &state) {
            rsp.SetCode(state.Code, state.Msg)
        } else {
            rsp.SetCode(status.DBQueryFailed, msg)
        }
        return rsp, nil
    }
    rsp.SetBody(status.Success, m)
    return rsp, nil
}

func (svc *{{.Service}}Service) Get{{.Service}}List0(ctx context.Context, req *{{.Package}}.Get{{.Service}}List0Request) (*{{.Package}}.Get{{.Service}}ListResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}List0 request", zap.Any("req", req))
    rsp := &{{.Package}}.Get{{.Service}}ListResponse{}
    pageNum := int(req.PageNum)
    if pageNum <= 0 {
        pageNum = 1
    }
    pageSize := int(req.PageSize)
    if pageSize <= 0 {
        pageSize = 20
    } else if pageSize > 50 {
        pageSize = 20
    }
    err := svc.vd.Var(req.Type, "gte=0,lte=3")
    if err != nil {
        msg := "请输入有效的Type"
        logger.Error(msg, zap.Int32("status", req.Type), zap.Error(err))
        rsp.SetCode(status.InvalidParams, msg)
        return rsp, nil
    }
    list, err := svc.dm.Get{{.Service}}List0(ctx, int(req.Type), pageSize, pageNum, req.IsDesc)
    if err != nil {
        var state status.Error
        if errors.As(err, &state) {
            rsp.SetCode(state.Code, state.Msg)
        } else {
            rsp.SetCode(status.DBQueryFailed, "数据查询失败")
        }
        return rsp, nil
    }
    rsp.SetBody(status.Success, list)
    return rsp, nil
}

func (svc *{{.Service}}Service) Get{{.Service}}List1(ctx context.Context, req *{{.Package}}.Get{{.Service}}List1Request) (*{{.Package}}.Get{{.Service}}ListResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Get{{.Service}}List1 request", zap.Any("req", req))
    rsp := &{{.Package}}.Get{{.Service}}ListResponse{}
    pageSize := int(req.PageSize)
    if pageSize <= 0 {
        pageSize = 20
    } else if pageSize > 50 {
        pageSize = 20
    }
    list, err := svc.dm.Get{{.Service}}List1(ctx, int(req.Type), pageSize, int(req.LastId), req.IsDesc)
    if err != nil {
        var state status.Error
        if errors.As(err, &state) {
            rsp.SetCode(state.Code, state.Msg)
        } else {
            rsp.SetCode(status.DBQueryFailed, "数据查询失败")
        }
        return rsp, nil
    }
    rsp.SetBody(status.Success, list)
    return rsp, nil
}

func (svc *{{.Service}}Service) Update{{.Service}}Type(ctx context.Context, req *{{.Package}}.Update{{.Service}}TypeRequest) (*{{.Package}}.Update{{.Service}}TypeResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Update{{.Service}}Type request")

    rsp := &{{.Package}}.Update{{.Service}}TypeResponse{}

    affected, err := svc.dm.Update{{.Service}}Type(ctx, int(req.Id), int(req.Type))
    if err != nil || affected <= 0 {
        msg := "更新{{.Service}}失败"
        logger.Error(msg, zap.Int8("affected", affected), zap.Error(err))
        rsp.SetCode(status.DBUpdateFailed, msg)
        return rsp, nil
    }
    rsp.SetCode(status.Success)
    return rsp, nil
}

func (svc *{{.Service}}Service) Delete{{.Service}}ById(ctx context.Context, req *{{.Package}}.Delete{{.Service}}ByIdRequest) (*{{.Package}}.Delete{{.Service}}ByIdResponse, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("Receive Delete{{.Service}}ById request")

    rsp := &{{.Package}}.Delete{{.Service}}ByIdResponse{}
    affected, err := svc.dm.Delete{{.Service}}ById(ctx, int(req.Id))
    if err != nil || affected <= 0 {
        msg := "更新{{.Service}}失败"
        logger.Error(msg, zap.Int8("affected", affected), zap.Error(err))
        var state status.Error
        if errors.As(err, &state) {
            rsp.SetCode(state.Code, state.Msg)
        } else {
            rsp.SetCode(status.DBUpdateFailed, msg)
        }
        return rsp, nil
    }
    rsp.SetCode(status.Success, "")
    return rsp, nil
}

func (svc *{{.Service}}Service) Get{{.Service}}ListByStream(stream {{.Package}}.{{.Service}}Service_Get{{.Service}}ListByStreamServer) error {
    logger := log.GetLogger(stream.Context())
    logger.Debug("Receive Get{{.Service}}ListByStream request")

    for {
        r, err := stream.Recv()
        logger.Debug("stream.Recv", zap.Any("r", r), zap.Error(err))
        if err == io.EOF {
            return nil
        }
        if err != nil {
            logger.Error("Recv Stream error", zap.Error(err))
            return err
        }

        if r.Id > 0 {
            m, err := svc.dm.Get{{.Service}}ById(stream.Context(), int(r.Id))
            if err != nil {
                logger.Error("Get{{.Service}}ById error", zap.Any("sample", m), zap.Error(err))
                return err
            }

            err = stream.Send(&{{.Package}}.Get{{.Service}}ListByStreamResponse{
                Index:  r.Index,
                Result: m,
            })
            if err != nil {
                logger.Error("Send Stream error", zap.Error(err))
                return err
            }
        } else {
            _ = stream.Send(&{{.Package}}.Get{{.Service}}ListByStreamResponse{
                Index:  r.Index,
                Result: nil,
            })
        }
    }
}
`

	path := "./" + data.Name + "/application/" + data.Name + "/service/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
