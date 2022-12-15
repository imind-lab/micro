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

// 生成client/service.go
func CreateApplicationService(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on 2021/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/status"
	"github.com/imind-lab/micro/tracing"
	"go.uber.org/zap"

	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	domain "{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

type {{.Svc}}Service struct {
	{{.Service}}.Unimplemented{{.Svc}}ServiceServer

	vd *validator.Validate
	dm domain.{{.Svc}}Domain
}

func New{{.Svc}}Service(dm domain.{{.Svc}}Domain) *{{.Svc}}Service {
	svc := &{{.Svc}}Service{
		dm: dm,
		vd: validator.New(),
	}

	return svc
}

// Create{{.Svc}} 创建{{.Svc}}
func (svc *{{.Svc}}Service) Create{{.Svc}}(ctx context.Context, req *{{.Service}}.Create{{.Svc}}Request) (*{{.Service}}.Create{{.Svc}}Response, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Create{{.Svc}} request")

	rsp := &{{.Service}}.Create{{.Svc}}Response{}

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
	m := model.{{.Svc}}{
		Name: req.Name,
		Type: int8(req.Type),
	}
	err = svc.dm.Create{{.Svc}}(ctx, m)
	if err != nil {
		msg := "创建{{.Svc}}失败"
		logger.Error(msg, zap.Any("{{.Service}}", m), zap.Error(err))
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

// Get{{.Svc}}ById 根据Id获取{{.Svc}}
func (svc *{{.Svc}}Service) Get{{.Svc}}ById(ctx context.Context, req *{{.Service}}.Get{{.Svc}}ByIdRequest) (*{{.Service}}.Get{{.Svc}}ByIdResponse, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}ById request")

	rsp := &{{.Service}}.Get{{.Svc}}ByIdResponse{}
	m, err := svc.dm.Get{{.Svc}}ById(ctx, int(req.Id))
	if err != nil {
		msg := "获取{{.Svc}}失败"
		logger.Error(msg, zap.Any("{{.Service}}", m), zap.Error(err))
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

func (svc *{{.Svc}}Service) Get{{.Svc}}List0(ctx context.Context, req *{{.Service}}.Get{{.Svc}}List0Request) (*{{.Service}}.Get{{.Svc}}ListResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}List0 request", zap.Any("req", req))
	rsp := &{{.Service}}.Get{{.Svc}}ListResponse{}
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
	list, err := svc.dm.Get{{.Svc}}List0(ctx, int(req.Type), pageSize, pageNum, req.IsDesc)
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

func (svc *{{.Svc}}Service) Get{{.Svc}}List1(ctx context.Context, req *{{.Service}}.Get{{.Svc}}List1Request) (*{{.Service}}.Get{{.Svc}}ListResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}List1 request", zap.Any("req", req))
	rsp := &{{.Service}}.Get{{.Svc}}ListResponse{}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 50 {
		pageSize = 20
	}
	list, err := svc.dm.Get{{.Svc}}List1(ctx, int(req.Type), pageSize, int(req.LastId), req.IsDesc)
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

func (svc *{{.Svc}}Service) Update{{.Svc}}Type(ctx context.Context, req *{{.Service}}.Update{{.Svc}}TypeRequest) (*{{.Service}}.Update{{.Svc}}TypeResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Update{{.Svc}}Type request")

	rsp := &{{.Service}}.Update{{.Svc}}TypeResponse{}

	affected, err := svc.dm.Update{{.Svc}}Type(ctx, int(req.Id), int(req.Type))
	if err != nil || affected <= 0 {
		msg := "更新{{.Svc}}失败"
		logger.Error(msg, zap.Int8("affected", affected), zap.Error(err))
		rsp.SetCode(status.DBUpdateFailed, msg)
		return rsp, nil
	}
	rsp.SetCode(status.Success)
	return rsp, nil
}

func (svc *{{.Svc}}Service) Delete{{.Svc}}ById(ctx context.Context, req *{{.Service}}.Delete{{.Svc}}ByIdRequest) (*{{.Service}}.Delete{{.Svc}}ByIdResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Delete{{.Svc}}ById request")

	rsp := &{{.Service}}.Delete{{.Svc}}ByIdResponse{}
	affected, err := svc.dm.Delete{{.Svc}}ById(ctx, int(req.Id))
	if err != nil || affected <= 0 {
		msg := "更新{{.Svc}}失败"
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

func (svc *{{.Svc}}Service) Get{{.Svc}}ListByStream(stream {{.Service}}.{{.Svc}}Service_Get{{.Svc}}ListByStreamServer) error {
	logger := log.GetLogger(stream.Context())
	logger.Debug("Receive Get{{.Svc}}ListByStream request")

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
			m, err := svc.dm.Get{{.Svc}}ById(stream.Context(), int(r.Id))
			if err != nil {
				logger.Error("Get{{.Svc}}ById error", zap.Any("{{.Service}}", m), zap.Error(err))
				return err
			}

			err = stream.Send(&{{.Service}}.Get{{.Svc}}ListByStreamResponse{
				Index:  r.Index,
				Result: m,
			})
			if err != nil {
				logger.Error("Send Stream error", zap.Error(err))
				return err
			}
		} else {
			_ = stream.Send(&{{.Service}}.Get{{.Svc}}ListByStreamResponse{
				Index:  r.Index,
				Result: nil,
			})
		}
	}
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/application/" + data.Service + "/service/"
	name := data.Service + ".go"

	return template.CreateFile(data, tpl, path, name)
}
