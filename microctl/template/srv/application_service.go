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

// 生成client/service.go
func CreateApplicationService(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
	"context"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/status"
	"github.com/imind-lab/micro/tracing"
	"go.uber.org/zap"

	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	service "{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}"
)

type {{.Svc}}Service struct {
	{{.Service}}.Unimplemented{{.Svc}}ServiceServer

	vd *validator.Validate
	dm service.{{.Svc}}Domain
}

func New{{.Svc}}Service() *{{.Svc}}Service {
	dm := service.New{{.Svc}}Domain()
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

	//info, articles, infos, configs, promotions, categories, emails, menus, pictures, maps := Export{{.Svc}}Info(req)

	rsp.SetCode(status.Success, "")
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
		logger.Error("获取{{.Svc}}失败", zap.Any("{{.Service}}", m), zap.Error(err))
		rsp.SetCode(status.DBQueryFailed, "获取{{.Svc}}失败")
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
	list, err := svc.dm.Get{{.Svc}}List0(ctx, int(req.Status), pageSize, pageNum, req.Order)
	if err != nil {
		rsp.SetCode(status.InternalError, "服务器内部错误")
		return rsp, nil
	}
	rsp.SetBody(status.Success, list)
	return rsp, nil
}

func (svc *{{.Svc}}Service) Get{{.Svc}}List1(ctx context.Context, req *{{.Service}}.Get{{.Svc}}List1Request) (*{{.Service}}.Get{{.Svc}}ListResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Get{{.Svc}}List0 request", zap.Any("req", req))
	rsp := &{{.Service}}.Get{{.Svc}}ListResponse{}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 50 {
		pageSize = 20
	}
	list, err := svc.dm.Get{{.Svc}}List1(ctx, int(req.Status), pageSize, int(req.LastId), req.Order)
	if err != nil {
		rsp.SetCode(status.InternalError, "服务器内部错误")
		return rsp, nil
	}
	rsp.SetBody(status.Success, list)
	return rsp, nil
}

func (svc *{{.Svc}}Service) Update{{.Svc}}Status(ctx context.Context, req *{{.Service}}.Update{{.Svc}}StatusRequest) (*{{.Service}}.Update{{.Svc}}StatusResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Update{{.Svc}}Status request")

	rsp := &{{.Service}}.Update{{.Svc}}StatusResponse{}

	return rsp, nil
}

func (svc *{{.Svc}}Service) Delete{{.Svc}}ById(ctx context.Context, req *{{.Service}}.Delete{{.Svc}}ByIdRequest) (*{{.Service}}.Delete{{.Svc}}ByIdResponse, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("Receive Delete{{.Svc}}ById request")

	rsp := &{{.Service}}.Delete{{.Svc}}ByIdResponse{}
	affected, err := svc.dm.Delete{{.Svc}}ById(ctx, int(req.Id))
	if err != nil || affected <= 0 {
		logger.Error("更新{{.Svc}}失败", zap.Int8("affected", affected), zap.Error(err))
		rsp.SetCode(status.DBUpdateFailed, "更新{{.Svc}}失败")
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

	t, err := template.New("application_service").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/application/" + data.Service + "/service/"

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
