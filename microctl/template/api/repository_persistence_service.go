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
    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/status"
    "github.com/imind-lab/micro/v2/tracing"
    "go.uber.org/zap"

    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
    {{.Svc}}Client "{{.Domain}}/{{.Repo}}/client"
)

const _NewClientError = "{{.Svc}}Client.New error"

func (repo {{.Service}}Repository) Create{{.Service}}(ctx context.Context, name string, typ int32) error {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))
        return status.ErrGRPCInternal
    }
    defer close()

    resp, err := {{.Svc}}Cli.Create{{.Service}}(ctx, &{{.Package}}.Create{{.Service}}Request{
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

// 根据Id获取{{.Service}}(有缓存)
func (repo {{.Service}}Repository) Get{{.Service}}ById(ctx context.Context, id int32) (*{{.Package}}.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    sentinelEntry, blockError := sentinel.Entry("test1")
    if blockError != nil {
        logger.Error("Get{{.Service}}ById触发熔断")
        return nil, status.ErrCircuitBreak
    }
    defer sentinelEntry.Exit()

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))
        return nil, status.ErrGRPCInternal
    }
    defer close()

    resp, err := {{.Svc}}Cli.Get{{.Service}}ById(ctx, &{{.Package}}.Get{{.Service}}ByIdRequest{
        Id: id,
    })
    logger.Debug("{{.Svc}}Cli.Get{{.Service}}ById", zap.Any("resp", resp), zap.Error(err))
    if err != nil {
        logger.Error("{{.Svc}}Cli.Get{{.Service}}ById error", zap.Any("resp", resp), zap.Error(err))
        sentinelEntry.SetError(err)

        return nil, err
    }

    return resp.Data, nil
}

func (repo {{.Service}}Repository) Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Package}}.{{.Service}}List, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    sentinelEntry, blockError := sentinel.Entry("test1")
    if blockError != nil {
        logger.Error("Get{{.Service}}List0触发熔断")
        return nil, status.ErrCircuitBreak
    }
    defer sentinelEntry.Exit()

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))

        return nil, status.ErrGRPCInternal
    }
    defer close()

    resp, err := {{.Svc}}Cli.Get{{.Service}}List0(ctx, &{{.Package}}.Get{{.Service}}List0Request{
        Type:     typ,
        PageSize: pageSize,
        PageNum:  pageNum,
        IsDesc:   isDesc,
    })
    logger.Debug("{{.Svc}}Cli.Get{{.Service}}List0", zap.Any("resp", resp), zap.Error(err))
    if err != nil {
        logger.Error("{{.Svc}}Cli.Get{{.Service}}List0 error", zap.Any("resp", resp), zap.Error(err))
        sentinelEntry.SetError(err)

        return nil, err
    }

    return resp.Data, nil
}

func (repo {{.Service}}Repository) Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Package}}.{{.Service}}List, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    sentinelEntry, blockError := sentinel.Entry("test1")
    if blockError != nil {
        logger.Error("Get{{.Service}}List1触发熔断")
        return nil, status.ErrCircuitBreak
    }
    defer sentinelEntry.Exit()

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))

        return nil, status.ErrGRPCInternal
    }
    defer close()

    resp, err := {{.Svc}}Cli.Get{{.Service}}List1(ctx, &{{.Package}}.Get{{.Service}}List1Request{
        Type:     typ,
        PageSize: pageSize,
        LastId:   lastId,
        IsDesc:   isDesc,
    })
    logger.Debug("{{.Svc}}Cli.Get{{.Service}}List1", zap.Any("resp", resp), zap.Error(err))
    if err != nil {
        logger.Error("{{.Svc}}Cli.Get{{.Service}}List1 error", zap.Any("resp", resp), zap.Error(err))

        sentinelEntry.SetError(err)

        return nil, err
    }

    return resp.Data, nil
}

func (repo {{.Service}}Repository) Update{{.Service}}Type(ctx context.Context, id, typ int32) error {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))
        return status.ErrGRPCInternal
    }
    defer close()

    resp, err := {{.Svc}}Cli.Update{{.Service}}Type(ctx, &{{.Package}}.Update{{.Service}}TypeRequest{
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

func (repo {{.Service}}Repository) Delete{{.Service}}ById(ctx context.Context, id int32) error {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))
        return status.ErrGRPCInternal
    }
    defer close()

    resp, err := {{.Svc}}Cli.Delete{{.Service}}ById(ctx, &{{.Package}}.Delete{{.Service}}ByIdRequest{
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

func (repo {{.Service}}Repository) Get{{.Service}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Package}}.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    {{.Svc}}Cli, close, err := {{.Svc}}Client.New(ctx)
    if err != nil {
        logger.Error(_NewClientError, zap.Error(err))
        return nil, status.ErrGRPCInternal
    }
    defer close()

    stream, err := {{.Svc}}Cli.Get{{.Service}}ListByStream(ctx)
    if err != nil {
        logger.Error("{{.Svc}}Cli.Get{{.Service}}ListByStream error", zap.Any("stream", stream), zap.Error(err))
        return nil, status.ErrGRPCInternal
    }

    data := make([]*{{.Package}}.{{.Service}}, len(ids))

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
                logger.Error("Get{{.Service}}ListByStream Recv error", zap.Error(err))
                return
            }
            fmt.Println("Recv", resp.Index, resp.Result)
            data[resp.Index] = resp.Result
        }
    }()

    for key, val := range ids {
        _ = stream.Send(&{{.Package}}.Get{{.Service}}ListByStreamRequest{
            Index: int32(key),
            Id:    val,
        })
    }
    stream.CloseSend()
    wg.Wait()
    return data, nil
}
`

	path := "./" + data.Name + "-api/repository/" + data.Name + "/persistence/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
