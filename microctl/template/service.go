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

// 生成service
func CreateService(data *Data) error {
	var tpl = `/**
 *  IMindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/imind-lab/micro/broker"
	"github.com/imind-lab/micro/util"
	"go.uber.org/zap"

	"{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	domain "{{.Domain}}/{{.Project}}/{{.Service}}/domain/{{.Service}}/service"
	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
)

type {{.Svc}}Service struct {
	{{.Service}}.Unimplemented{{.Svc}}ServiceServer

	validate *validator.Validate

	dm domain.{{.Svc}}Domain
}

func New{{.Svc}}Service() *{{.Svc}}Service {
	dm := domain.New{{.Svc}}Domain()
	svc := &{{.Svc}}Service{
		dm:       dm,
		validate: validator.New(),
	}

	return svc
}

// Create{{.Svc}} 创建{{.Svc}}
func (svc *{{.Svc}}Service) Create{{.Svc}}(ctx context.Context, req *{{.Service}}.Create{{.Svc}}Request) (*{{.Service}}.Create{{.Svc}}Response, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Create{{.Svc}}"))
	logger.Debug("Receive Create{{.Svc}} request")

	rsp := &{{.Service}}.Create{{.Svc}}Response{}

	m := req.Dto
	fmt.Println("Dto", m)
	err := svc.validate.Struct(req)
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
	if m == nil {
		logger.Error("{{.Svc}}不能为空", zap.Any("params", m), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "{{.Svc}}不能为空"
		rsp.Error = err
		return rsp, nil
	}

	err = svc.validate.Var(m.Name, "required,email")
	if err != nil {
		logger.Error("Name不能为空", zap.Any("name", m.Name), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "Name不能为空"
		rsp.Error = err
		return rsp, nil
	}
	m.CreateTime = util.GetNowWithMillisecond()
	m.CreateDatetime = time.Now().Format(util.DateTimeFmt)
	m.UpdateDatetime = time.Now().Format(util.DateTimeFmt)
	err = svc.dm.Create{{.Svc}}(ctx, m)
	if err != nil {
		logger.Error("创建{{.Svc}}失败", zap.Any("{{.Service}}", m), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "创建{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = true

	endpoint, err := broker.NewBroker(constant.MQName)
	if err != nil {
		ctxzap.Error(ctx, "broker.NewBroker error", zap.Error(err))
		return rsp, err
	}
	endpoint.Publish(&broker.Message{
		Topic: endpoint.Options().Topics["create{{.Service}}"],
		Body:  []byte(fmt.Sprintf("{{.Svc}} %s Created", m.Name)),
	})

	return rsp, nil
}

// Get{{.Svc}}ById 根据Id获取{{.Svc}}
func (svc *{{.Svc}}Service) Get{{.Svc}}ById(ctx context.Context, req *{{.Service}}.Get{{.Svc}}ByIdRequest) (*{{.Service}}.Get{{.Svc}}ByIdResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Get{{.Svc}}ById"))
	logger.Debug("Receive Get{{.Svc}}ById request")

	rsp := &{{.Service}}.Get{{.Svc}}ByIdResponse{}
	m, err := svc.dm.Get{{.Svc}}ById(ctx, req.Id)
	if err != nil {
		logger.Error("获取{{.Svc}}失败", zap.Any("{{.Service}}", m), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "获取{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = true
	rsp.Dto = m
	return rsp, nil
}

func (svc *{{.Svc}}Service) Get{{.Svc}}List(ctx context.Context, req *{{.Service}}.Get{{.Svc}}ListRequest) (*{{.Service}}.Get{{.Svc}}ListResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Get{{.Svc}}List"))
	logger.Debug("Receive Get{{.Svc}}List request")
	rsp := &{{.Service}}.Get{{.Svc}}ListResponse{}

	err := svc.validate.Struct(req)
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
	err = svc.validate.Var(req.Status, "gte=0,lte=3")
	if err != nil {
		logger.Error("请输入有效的Status", zap.Int32("status", req.Status), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "请输入有效的Status"
		rsp.Error = err
		return rsp, nil
	}

	if req.Pagesize <= 0 {
		req.Pagesize = 20
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	list, err := svc.dm.Get{{.Svc}}List(ctx, req.Status, req.Lastid, req.Pagesize, req.Page)
	if err != nil {
		logger.Error("获取{{.Svc}}失败", zap.Any("list", list), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "获取{{.Svc}}List失败"
		rsp.Error = err
		return rsp, nil
	}
	rsp.Success = true
	rsp.Data = list
	return rsp, nil
}

func (svc *{{.Svc}}Service) Update{{.Svc}}Status(ctx context.Context, req *{{.Service}}.Update{{.Svc}}StatusRequest) (*{{.Service}}.Update{{.Svc}}StatusResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Update{{.Svc}}Status"))
	logger.Debug("Receive Update{{.Svc}}Status request")

	rsp := &{{.Service}}.Update{{.Svc}}StatusResponse{}
	affected, err := svc.dm.Update{{.Svc}}Status(ctx, req.Id, req.Status)
	if err != nil || affected <= 0 {
		logger.Error("更新{{.Svc}}失败", zap.Int64("affected", affected), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "更新{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}
	rsp.Success = true
	return rsp, nil
}

func (svc *{{.Svc}}Service) Update{{.Svc}}Count(ctx context.Context, req *{{.Service}}.Update{{.Svc}}CountRequest) (*{{.Service}}.Update{{.Svc}}CountResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Update{{.Svc}}Count"))
	logger.Debug("Receive Update{{.Svc}}Count request")

	rsp := &{{.Service}}.Update{{.Svc}}CountResponse{}
	affected, err := svc.dm.Update{{.Svc}}Count(ctx, req.Id, req.Num, req.Column)
	if err != nil || affected <= 0 {
		logger.Error("更新{{.Svc}}失败", zap.Int64("affected", affected), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "更新{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}
	rsp.Success = true

	endpoint, err := broker.NewBroker(constant.MQName)
	if err != nil {
		ctxzap.Error(ctx, "kafka.New error", zap.Error(err))
		return rsp, err
	}
	endpoint.Publish(&broker.Message{
		Topic: endpoint.Options().Topics["update{{.Service}}count"],
		Body:  nil,
	})

	return rsp, nil
}

func (svc *{{.Svc}}Service) Delete{{.Svc}}ById(ctx context.Context, req *{{.Service}}.Delete{{.Svc}}ByIdRequest) (*{{.Service}}.Delete{{.Svc}}ByIdResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Delete{{.Svc}}ById"))
	logger.Debug("Receive Delete{{.Svc}}ById request")

	rsp := &{{.Service}}.Delete{{.Svc}}ByIdResponse{}
	affected, err := svc.dm.Delete{{.Svc}}ById(ctx, req.Id)
	if err != nil || affected <= 0 {
		logger.Error("更新{{.Svc}}失败", zap.Int64("affected", affected), zap.Error(err))

		err := &{{.Service}}.Error{}
		err.Message = "删除{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}
	rsp.Success = true
	return rsp, nil
}

func (svc *{{.Svc}}Service) Get{{.Svc}}ListByStream(stream {{.Service}}.{{.Svc}}Service_Get{{.Svc}}ListByStreamServer) error {
	logger := ctxzap.Extract(stream.Context()).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Get{{.Svc}}ListByStream"))
	logger.Debug("Receive Get{{.Svc}}ListByStream request")

	for {
		r, err := stream.Recv()
		ctxzap.Debug(stream.Context(), "stream.Recv", zap.Any("r", r), zap.Error(err))
		if err == io.EOF {
			return nil
		}
		if err != nil {
			logger.Error("Recv Stream error", zap.Error(err))
			return err
		}

		if r.Id > 0 {
			m, err := svc.dm.Get{{.Svc}}ById(stream.Context(), r.Id)
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

	t, err := template.New("service").Parse(tpl)
	if err != nil {
		return err
	}

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
