/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tp "github.com/imind-lab/micro/microctl/template"
)

// 生成service
func CreateService(data *tp.Data) error {
	var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	{{.Service}}_api "{{.Domain}}/{{.Project}}/{{.Service}}-api/server/proto/{{.Service}}-api"
	{{.Service}}Client "{{.Domain}}/{{.Project}}/{{.Service}}/client"
	"{{.Domain}}/{{.Project}}/{{.Service}}/server/proto/{{.Service}}"
	sentinelx "github.com/imind-lab/micro/sentinel"
)

type {{.Svc}}Service struct {
	{{.Service}}_api.Unimplemented{{.Svc}}ServiceServer

	validate *validator.Validate

	ds *sentinelx.Sentinel
}

//New{{.Svc}}Service 创建用户服务实例
func New{{.Svc}}Service(logger *zap.Logger) *{{.Svc}}Service {
	ds, _ := sentinelx.NewSentinel(logger)
	svc := &{{.Svc}}Service{
		ds:       ds,
		validate: validator.New(),
	}
	return svc
}

// Create{{.Svc}} 创建{{.Svc}}
func (svc *{{.Svc}}Service) Create{{.Svc}}(ctx context.Context, req *{{.Service}}_api.Create{{.Svc}}Request) (*{{.Service}}_api.Create{{.Svc}}Response, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Create{{.Svc}}"))
	logger.Debug("Receive Create{{.Svc}} request")

	rsp := &{{.Service}}_api.Create{{.Svc}}Response{}

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
	err = svc.validate.Var(m, "required")
	fmt.Println("validate", err)
	if m == nil {
		logger.Error("{{.Svc}}不能为空", zap.Any("params", m))

		err := &{{.Service}}_api.Error{}
		err.Message = "{{.Svc}}不能为空"
		rsp.Error = err
		return rsp, nil
	}

	err = svc.validate.Var(m.Name, "required,email")
	fmt.Println("validate", m.Name, err)
	if len(m.Name) == 0 {
		logger.Error("Name不能为空", zap.Any("name", m.Name))

		err := &{{.Service}}_api.Error{}
		err.Message = "Name不能为空"
		rsp.Error = err
		return rsp, nil
	}

	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Create{{.Svc}} Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)

	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("服务器繁忙，请稍候再试", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}
	{{.Service}}Srv := {{.Svc}}Gw2Srv(req.Dto)
	resp, err := {{.Service}}Cli.Create{{.Svc}}(ctx, &{{.Service}}.Create{{.Svc}}Request{
		Dto: {{.Service}}Srv,
	})
	if err != nil {
		logger.Error("{{.Service}}Cli.Create{{.Svc}} error", zap.Any("{{.Service}}", m), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "创建{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = resp.Success
	rsp.Error = ErrorSrv2Gw(resp.Error)
	return rsp, nil
}

// Get{{.Svc}}ById 根据Id获取{{.Svc}}
func (svc *{{.Svc}}Service) Get{{.Svc}}ById(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}ByIdRequest) (*{{.Service}}_api.Get{{.Svc}}ByIdResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Get{{.Svc}}ById"))
	logger.Debug("Receive Get{{.Svc}}ById request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ByIdResponse{}

	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Get{{.Svc}}ById Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)
	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("{{.Service}}Client.New error", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	sentinelEntry, blockError := sentinel.Entry("test1")
	if blockError != nil {
		logger.Error("触发熔断降级", zap.Any("TriggeredRule", blockError.TriggeredRule()), zap.Any("TriggeredValue", blockError.TriggeredValue()))

		err := &{{.Service}}_api.Error{}
		err.Message = "获取{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}
	defer sentinelEntry.Exit()

	resp, err := {{.Service}}Cli.Get{{.Svc}}ById(ctx, &{{.Service}}.Get{{.Svc}}ByIdRequest{
		Id: req.Id,
	})
	ctxzap.Debug(ctx, "{{.Service}}Cli.Get{{.Svc}}ById", zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		logger.Error("{{.Service}}Cli.Get{{.Svc}}ById error", zap.Any("resp", resp), zap.Error(err))

		sentinelEntry.SetError(err)

		err := &{{.Service}}_api.Error{}
		err.Message = "获取{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = resp.Success
	rsp.Dto = {{.Svc}}Srv2Gw(resp.Dto)
	rsp.Error = ErrorSrv2Gw(resp.Error)

	return rsp, nil
}

func (svc *{{.Svc}}Service) Get{{.Svc}}List(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}ListRequest) (*{{.Service}}_api.Get{{.Svc}}ListResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Get{{.Svc}}List"))
	logger.Debug("Receive Get{{.Svc}}List request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ListResponse{}

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

	rateEntry, rateError := sentinel.Entry("abcd", sentinel.WithTrafficType(base.Inbound))
	if rateError != nil {
		ctxzap.Debug(ctx, "Get{{.Svc}}List限流了")
		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}
	defer rateEntry.Exit()
	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Get{{.Svc}}List Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)

	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("{{.Service}}Client.New error", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	resp, err := {{.Service}}Cli.Get{{.Svc}}List(ctx, &{{.Service}}.Get{{.Svc}}ListRequest{
		Status:   req.Status,
		Lastid:   req.Lastid,
		Pagesize: req.Pagesize,
		Page:     req.Page,
	})
	if err != nil {
		logger.Error("{{.Service}}Cli.Get{{.Svc}}List error", zap.Any("resp", resp), zap.Error(err))
		err := &{{.Service}}_api.Error{}
		err.Message = "获取{{.Svc}}List失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = resp.Success
	rsp.Data = {{.Svc}}ListSrv2Gw(resp.Data)
	rsp.Error = ErrorSrv2Gw(resp.Error)
	return rsp, nil
}

func (svc *{{.Svc}}Service) Update{{.Svc}}Status(ctx context.Context, req *{{.Service}}_api.Update{{.Svc}}StatusRequest) (*{{.Service}}_api.Update{{.Svc}}StatusResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Update{{.Svc}}Status"))
	logger.Debug("Receive Update{{.Svc}}Status request")

	rsp := &{{.Service}}_api.Update{{.Svc}}StatusResponse{}

	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Update{{.Svc}}Status Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)

	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("{{.Service}}Client.New error", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	resp, err := {{.Service}}Cli.Update{{.Svc}}Status(ctx, &{{.Service}}.Update{{.Svc}}StatusRequest{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		logger.Error("{{.Service}}Cli.Update{{.Svc}}Status error", zap.Any("resp", resp), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "更新{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = resp.Success
	rsp.Error = ErrorSrv2Gw(resp.Error)

	return rsp, nil
}

func (svc *{{.Svc}}Service) Update{{.Svc}}Count(ctx context.Context, req *{{.Service}}_api.Update{{.Svc}}CountRequest) (*{{.Service}}_api.Update{{.Svc}}CountResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Update{{.Svc}}Count"))
	logger.Debug("Receive Update{{.Svc}}Count request")

	rsp := &{{.Service}}_api.Update{{.Svc}}CountResponse{}

	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Update{{.Svc}}Count Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)

	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("{{.Service}}Client.New error", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	resp, err := {{.Service}}Cli.Update{{.Svc}}Count(ctx, &{{.Service}}.Update{{.Svc}}CountRequest{
		Id:     req.Id,
		Num:    req.Num,
		Column: req.Column,
	})
	if err != nil {
		logger.Error("{{.Service}}Cli.Update{{.Svc}}Count error", zap.Any("resp", resp), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "更新{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = resp.Success
	rsp.Error = ErrorSrv2Gw(resp.Error)

	return rsp, nil
}

func (svc *{{.Svc}}Service) Delete{{.Svc}}ById(ctx context.Context, req *{{.Service}}_api.Delete{{.Svc}}ByIdRequest) (*{{.Service}}_api.Delete{{.Svc}}ByIdResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Delete{{.Svc}}ById"))
	logger.Debug("Receive Delete{{.Svc}}ById request")

	rsp := &{{.Service}}_api.Delete{{.Svc}}ByIdResponse{}

	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Delete{{.Svc}}ById Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)

	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("{{.Service}}Client.New error", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	resp, err := {{.Service}}Cli.Delete{{.Svc}}ById(ctx, &{{.Service}}.Delete{{.Svc}}ByIdRequest{
		Id: req.Id,
	})
	if err != nil {
		logger.Error("{{.Service}}Cli.Delete{{.Svc}}ById error", zap.Any("resp", resp), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "删除{{.Svc}}失败"
		rsp.Error = err
		return rsp, nil
	}

	rsp.Success = resp.Success
	rsp.Error = ErrorSrv2Gw(resp.Error)

	return rsp, nil
}
func (svc *{{.Svc}}Service) Get{{.Svc}}ListByIds(ctx context.Context, req *{{.Service}}_api.Get{{.Svc}}ListByIdsRequest) (*{{.Service}}_api.Get{{.Svc}}ListByIdsResponse, error) {
	logger := ctxzap.Extract(ctx).With(zap.String("layer", "{{.Svc}}Service"), zap.String("func", "Get{{.Svc}}ListByIds"))
	logger.Debug("Receive Get{{.Svc}}ListByIds request")

	rsp := &{{.Service}}_api.Get{{.Svc}}ListByIdsResponse{}

	uid := 0
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uids := meta.Get("uid")
		if len(uids) > 0 {
			uid, _ = strconv.Atoi(uids[0])
		}
	}
	ctxzap.Debug(ctx, "Get{{.Svc}}ListByIds Metadata", zap.Any("meta", meta), zap.Int("uid", uid), zap.Bool("ok", ok))

	ctx = metadata.NewOutgoingContext(ctx, meta)

	{{.Service}}Cli, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error("{{.Service}}Client.New error", zap.Any("{{.Service}}Cli", {{.Service}}Cli), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	data := make([]*{{.Service}}_api.{{.Svc}}, len(req.Ids))

	streamClient, err := {{.Service}}Cli.Get{{.Svc}}ListByStream(ctx)
	if err != nil {
		logger.Error("{{.Service}}Cli.Get{{.Svc}}ListByStream error", zap.Any("streamClient", streamClient), zap.Error(err))

		err := &{{.Service}}_api.Error{}
		err.Message = "服务器繁忙，请稍候再试"
		rsp.Error = err
		return rsp, nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			resp, err := streamClient.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				ctxzap.Error(ctx, "Get{{.Svc}}ListByStream Recv error", zap.Error(err))
				return
			}
			fmt.Println("Recv", resp.Index, resp.Result)
			data[resp.Index] = {{.Svc}}Srv2Gw(resp.Result)
		}
	}()

	for key, val := range req.Ids {
		_ = streamClient.Send(&{{.Service}}.Get{{.Svc}}ListByStreamRequest{
			Index: int32(key),
			Id:    val,
		})
	}
	streamClient.CloseSend()
	wg.Wait()

	for _, m := range data {
		if m != nil {
			rsp.Data = append(rsp.Data, m)
		}
	}

	rsp.Success = true

	return rsp, nil
}

func (svc *{{.Svc}}Service) Close() {
	if svc.ds != nil {
		svc.ds.Close()
	}
}
`

	t, err := template.New("service").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/server/service/"

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

	tpl = `/**
 *  {{.Project}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package service

import (
        {{.Service}}_api "{{.Domain}}/{{.Project}}/{{.Service}}-api/server/proto/{{.Service}}-api"
        "{{.Domain}}/{{.Project}}/{{.Service}}/server/proto/{{.Service}}"
)

func {{.Svc}}Map(pos []*{{.Service}}.{{.Svc}}, fn func(*{{.Service}}.{{.Svc}}) *{{.Service}}_api.{{.Svc}}) []*{{.Service}}_api.{{.Svc}} {
        var dtos []*{{.Service}}_api.{{.Svc}}
        for _, po := range pos {
                dtos = append(dtos, fn(po))
        }
        return dtos
}

func {{.Svc}}Gw2Srv(po *{{.Service}}_api.{{.Svc}}) *{{.Service}}.{{.Svc}} {
        if po == nil {
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

func {{.Svc}}Srv2Gw(dto *{{.Service}}.{{.Svc}}) *{{.Service}}_api.{{.Svc}} {
        if dto == nil {
                return nil
        }

        po := &{{.Service}}_api.{{.Svc}}{}
        po.Id = dto.Id
        po.Name = dto.Name
        po.ViewNum = dto.ViewNum
        po.Status = dto.Status
        po.CreateTime = dto.CreateTime
        po.UpdateDatetime = dto.UpdateDatetime
        po.CreateDatetime = dto.CreateDatetime

        return po
}

func {{.Svc}}ListSrv2Gw(dto *{{.Service}}.{{.Svc}}List) *{{.Service}}_api.{{.Svc}}List {
        if dto == nil {
                return nil
        }

        po := &{{.Service}}_api.{{.Svc}}List{}
        po.Total = dto.Total
        po.TotalPage = dto.TotalPage
        po.CurPage = dto.CurPage
        po.Datalist = {{.Svc}}Map(dto.Datalist, {{.Svc}}Srv2Gw)

        return po
}

func ErrorSrv2Gw(err *{{.Service}}.Error) *{{.Service}}_api.Error {
        if err == nil {
                return nil
        }

        po := &{{.Service}}_api.Error{}
        po.Message = err.Message
        po.Code = err.Code

        return po
}
`

	t, err = template.New("gw-convert").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "convert.go"

	f, err = os.Create(fileName)
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
