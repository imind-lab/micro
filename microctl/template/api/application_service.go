/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package api

import (
	"github.com/imind-lab/micro/microctl/template"
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
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/status"
	"github.com/imind-lab/micro/tracing"
	"go.uber.org/zap"

	{{.Service}}_api "gitlab.imind.tech/{{.Project}}/{{.Service}}-api/application/{{.Service}}/proto"
	domain "gitlab.imind.tech/{{.Project}}/{{.Service}}-api/domain/{{.Service}}"
)

type {{.Svc}}Service struct {
	{{.Service}}_api.Unimplemented{{.Svc}}ServiceServer

	vd *validator.Validate
	dm domain.{{.Svc}}Domain
}

// New{{.Svc}}Service 创建用户服务实例
func New{{.Svc}}Service(dm domain.{{.Svc}}Domain) *{{.Svc}}Service {
	svc := &{{.Svc}}Service{
		dm: dm,
		vd: validator.New(),
	}
	return svc
}

// Create{{.Svc}} 创建{{.Svc}}
func (svc *{{.Svc}}Service) Create{{.Svc}}(ctx context.Context, req *{{.Service}}_api.Create{{.Svc}}Request) (*{{.Service}}_api.Create{{.Svc}}Response, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Create{{.Svc}} request")

	rsp := &{{.Service}}_api.Create{{.Svc}}Response{}

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

	err = svc.dm.Create{{.Svc}}(ctx, req.Name, req.Type)
	e := status.AsError(err)
	rsp.SetCode(e.Code, e.Msg)
	return rsp, nil
}

// Get{{.Svc}}ById 根据Id获取{{.Svc}}
func (svc *{{.Svc}}Service) Get{{.Svc}}ById(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}ByIdRequest) (*{{.Service}}_api.Get{{.Svc}}ByIdResponse, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}ById request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ByIdResponse{}

	data, err := svc.dm.Get{{.Svc}}ById(ctx, req.Id)

	e := status.AsError(err)
	if e.Code == status.Success {
		rsp.SetBody(status.Success, data)
	} else {
		rsp.SetCode(e.Code, e.Msg)
	}
	return rsp, nil
}

func (svc *{{.Svc}}Service) Get{{.Svc}}List0(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}List0Request) (*{{.Service}}_api.Get{{.Svc}}ListResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}List0 request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ListResponse{}

	err := svc.vd.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}
	}

	rateEntry, rateError := sentinel.Entry("abcd", sentinel.WithTrafficType(base.Inbound))
	if rateError != nil {
		logger.Debug("Get{{.Svc}}List0限流了")

		rsp.SetCode(status.SystemError, "Get{{.Svc}}List0限流了")
		return rsp, nil
	}
	defer rateEntry.Exit()

	data, err := svc.dm.Get{{.Svc}}List0(ctx, req.Type, req.PageSize, req.PageNum, req.IsDesc)

	e := status.AsError(err)
	if e.Code == status.Success {
		rsp.SetBody(status.Success, data)
	} else {
		rsp.SetCode(e.Code, e.Msg)
	}
	return rsp, nil
}

func (svc *{{.Svc}}Service) Get{{.Svc}}List1(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}List1Request) (*{{.Service}}_api.Get{{.Svc}}ListResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}List1 request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ListResponse{}

	rateEntry, rateError := sentinel.Entry("abcd", sentinel.WithTrafficType(base.Inbound))
	if rateError != nil {
		logger.Debug("Get{{.Svc}}List1限流了")

		rsp.SetCode(status.SystemError, "Get{{.Svc}}List1限流了")
		return rsp, nil
	}
	defer rateEntry.Exit()

	data, err := svc.dm.Get{{.Svc}}List1(ctx, req.Type, req.PageSize, req.LastId, req.IsDesc)

	e := status.AsError(err)
	if e.Code == status.Success {
		rsp.SetBody(status.Success, data)
	} else {
		rsp.SetCode(e.Code, e.Msg)
	}
	return rsp, nil
}

func (svc *{{.Svc}}Service) Update{{.Svc}}Type(ctx context.Context, req *{{.Service}}_api.Update{{.Svc}}TypeRequest) (*{{.Service}}_api.Update{{.Svc}}TypeResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Update{{.Svc}}Type request")

	rsp := &{{.Service}}_api.Update{{.Svc}}TypeResponse{}

	err := svc.dm.Update{{.Svc}}Type(ctx, req.Id, req.Type)
	e := status.AsError(err)
	rsp.SetCode(e.Code, e.Msg)
	return rsp, nil
}

func (svc *{{.Svc}}Service) Delete{{.Svc}}ById(ctx context.Context, req *{{.Service}}_api.Delete{{.Svc}}ByIdRequest) (*{{.Service}}_api.Delete{{.Svc}}ByIdResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Delete{{.Svc}}ById request")

	rsp := &{{.Service}}_api.Delete{{.Svc}}ByIdResponse{}

	err := svc.dm.Delete{{.Svc}}ById(ctx, req.Id)

	e := status.AsError(err)
	rsp.SetCode(e.Code, e.Msg)

	return rsp, nil
}
func (svc *{{.Svc}}Service) Get{{.Svc}}ListByIds(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}ListByIdsRequest) (*{{.Service}}_api.Get{{.Svc}}ListByIdsResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}ListByIds request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ListByIdsResponse{}

	data, err := svc.dm.Get{{.Svc}}ListByIds(ctx, req.Ids)

	e := status.AsError(err)
	if e.Code == status.Success {
		rsp.SetBody(status.Success, data)
	} else {
		rsp.SetCode(e.Code, e.Msg)
	}
	return rsp, nil
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/application/" + data.Service + "/service/"
	name := data.Service + ".go"

	return template.CreateFile(data, tpl, path, name)
}
