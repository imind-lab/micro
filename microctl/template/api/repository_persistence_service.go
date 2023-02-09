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

// 生成repository/model.go
func CreateRepositoryPersistenceService(data *template.Data) error {
    var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package persistence

import (
	"context"
	"fmt"
	"io"
	"sync"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"go.uber.org/zap"

	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/status"
	"github.com/imind-lab/micro/tracing"

	{{.Service}} "gitlab.imind.tech/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	{{.Service}}Client "gitlab.imind.tech/{{.Project}}/{{.Service}}/client"
)

const _NewClientError = "{{.Service}}Client.New error"

func (repo {{.Svc}}Repository) Create{{.Svc}}(ctx context.Context, name string, typ int32) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))
		return status.ErrGRPCInternal
	}
	defer close()

	resp, err := {{.Service}}Cli.Create{{.Svc}}(ctx, &{{.Service}}.Create{{.Svc}}Request{
		Name: name,
		Type: typ,
	})
	if err != nil {
		return err
	}
	if code := status.Code(resp.Code); code != status.Success {
		return status.New(code)
	}
	return nil
}

// 根据Id获取{{.Svc}}(有缓存)
func (repo {{.Svc}}Repository) Get{{.Svc}}ById(ctx context.Context, id int32) (*{{.Service}}.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	sentinelEntry, blockError := sentinel.Entry("test1")
	if blockError != nil {
		logger.Error("Get{{.Svc}}ById触发熔断")
		return nil, status.ErrCircuitBreak
	}
	defer sentinelEntry.Exit()

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))
		return nil, status.ErrGRPCInternal
	}
	defer close()

	resp, err := {{.Service}}Cli.Get{{.Svc}}ById(ctx, &{{.Service}}.Get{{.Svc}}ByIdRequest{
		Id: id,
	})
	logger.Debug("{{.Service}}Cli.Get{{.Svc}}ById", zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		logger.Error("{{.Service}}Cli.Get{{.Svc}}ById error", zap.Any("resp", resp), zap.Error(err))
		sentinelEntry.SetError(err)

		return nil, err
	}

	return resp.Data, nil
}

func (repo {{.Svc}}Repository) Get{{.Svc}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Service}}.{{.Svc}}List, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	sentinelEntry, blockError := sentinel.Entry("test1")
	if blockError != nil {
		logger.Error("Get{{.Svc}}List0触发熔断")
		return nil, status.ErrCircuitBreak
	}
	defer sentinelEntry.Exit()

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))

		return nil, status.ErrGRPCInternal
	}
	defer close()

	resp, err := {{.Service}}Cli.Get{{.Svc}}List0(ctx, &{{.Service}}.Get{{.Svc}}List0Request{
		Type:     typ,
		PageSize: pageSize,
		PageNum:  pageNum,
		IsDesc:   isDesc,
	})
	logger.Debug("{{.Service}}Cli.Get{{.Svc}}List0", zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		logger.Error("{{.Service}}Cli.Get{{.Svc}}List0 error", zap.Any("resp", resp), zap.Error(err))
		sentinelEntry.SetError(err)

		return nil, err
	}

	return resp.Data, nil
}

func (repo {{.Svc}}Repository) Get{{.Svc}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Service}}.{{.Svc}}List, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	sentinelEntry, blockError := sentinel.Entry("test1")
	if blockError != nil {
		logger.Error("Get{{.Svc}}List1触发熔断")
		return nil, status.ErrCircuitBreak
	}
	defer sentinelEntry.Exit()

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))

		return nil, status.ErrGRPCInternal
	}
	defer close()

	resp, err := {{.Service}}Cli.Get{{.Svc}}List1(ctx, &{{.Service}}.Get{{.Svc}}List1Request{
		Type:     typ,
		PageSize: pageSize,
		LastId:   lastId,
		IsDesc:   isDesc,
	})
	logger.Debug("{{.Service}}Cli.Get{{.Svc}}List1", zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		logger.Error("{{.Service}}Cli.Get{{.Svc}}List1 error", zap.Any("resp", resp), zap.Error(err))

		sentinelEntry.SetError(err)

		return nil, err
	}

	return resp.Data, nil
}

func (repo {{.Svc}}Repository) Update{{.Svc}}Type(ctx context.Context, id, typ int32) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))
		return status.ErrGRPCInternal
	}
	defer close()

	resp, err := {{.Service}}Cli.Update{{.Svc}}Type(ctx, &{{.Service}}.Update{{.Svc}}TypeRequest{
		Id:   id,
		Type: typ,
	})
	if err != nil {
		return err
	}
	if code := status.Code(resp.Code); code != status.Success {
		return status.New(code)
	}
	return nil
}

func (repo {{.Svc}}Repository) Delete{{.Svc}}ById(ctx context.Context, id int32) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))
		return status.ErrGRPCInternal
	}
	defer close()

	resp, err := {{.Service}}Cli.Delete{{.Svc}}ById(ctx, &{{.Service}}.Delete{{.Svc}}ByIdRequest{
		Id: id,
	})
	if err != nil {
		return err
	}
	if code := status.Code(resp.Code); code != status.Success {
		return status.New(code)
	}
	return nil
}

func (repo {{.Svc}}Repository) Get{{.Svc}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Service}}.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	{{.Service}}Cli, close, err := {{.Service}}Client.New(ctx)
	if err != nil {
		logger.Error(_NewClientError, zap.Error(err))
		return nil, status.ErrGRPCInternal
	}
	defer close()

	stream, err := {{.Service}}Cli.Get{{.Svc}}ListByStream(ctx)
	if err != nil {
		logger.Error("{{.Service}}C.Get{{.Svc}}ListByStream error", zap.Any("stream", stream), zap.Error(err))
		return nil, status.ErrGRPCInternal
	}

	data := make([]*{{.Service}}.{{.Svc}}, len(ids))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				logger.Error("Get{{.Svc}}ListByStream Recv error", zap.Error(err))
				return
			}
			fmt.Println("Recv", resp.Index, resp.Result)
			data[resp.Index] = resp.Result
		}
	}()

	for key, val := range ids {
		_ = stream.Send(&{{.Service}}.Get{{.Svc}}ListByStreamRequest{
			Index: int32(key),
			Id:    val,
		})
	}
	stream.CloseSend()
	wg.Wait()
	return data, nil
}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/repository/" + data.Service + "/persistence/"
    name := data.Service + ".go"

    return template.CreateFile(data, tpl, path, name)
}
