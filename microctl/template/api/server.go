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

// 生成server/server.go
func CreateServer(data *template.Data) error {
	var tpl = `package server

import (
    "context"
    "fmt"
    "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
    "github.com/imind-lab/micro/v2"
    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/tracing"
    "github.com/spf13/viper"
    {{.Package}}_api "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/{{.Name}}/proto"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "google.golang.org/grpc"

    grpcx "github.com/imind-lab/micro/v2/grpc"
)

type Port struct {
    Http int ${backtick}yaml:"http"${backtick}
    Grpc int ${backtick}yaml:"grpc"${backtick}
}

type Config struct {
    Name      string ${backtick}yaml:"name"${backtick}
    Namespace string ${backtick}yaml:"namespace"${backtick}
    LogLevel  int    ${backtick}yaml:"logLevel"${backtick}
    LogFormat string ${backtick}yaml:"logFormat"${backtick}
    Port      Port   ${backtick}yaml:"port"${backtick}
}

func Serve() error {
    var conf Config
    if err := viper.UnmarshalKey("service", &conf); err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    logger := log.NewLogger(zapcore.Level(conf.LogLevel), conf.LogFormat, zap.Fields(zap.String("namespace", conf.Namespace), zap.String("service", conf.Name)))
    ctx := ctxzap.ToContext(context.Background(), logger)

    // 初始化调用链追踪
    provider, err := tracing.InitTracer()
    if err != nil {
        return err
    }
    defer provider.Shutdown(ctx)

    grpcCred := grpcx.NewGrpcCred()

    svc := micro.NewService()

    svc.Init(
        micro.Context(ctx),
        micro.Logger(logger),
        micro.Name(conf.Name),
        micro.ServerCred(grpcCred.ServerCred()),
        micro.ClientCred(grpcCred.ClientCred()))

    // 注册gRPC-Gateway
    endPoint := fmt.Sprintf(":%d", conf.Port.Grpc)

    mux := svc.ServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(grpcCred.ClientCred())}
    err = {{.Package}}_api.Register{{.Service}}ServiceHandlerFromEndpoint(svc.Options().Context, mux, endPoint, opts)
    if err != nil {
        return err
    }

    grpcSrv := svc.GrpcServer()

    {{.Svc}}Svc := Create{{.Service}}Service()
    {{.Package}}_api.Register{{.Service}}ServiceServer(grpcSrv, {{.Svc}}Svc)

    // This commentary is for scaffolding. Do not modify or delete it

    return svc.Run()
}
`

	path := "./" + data.Name + "-api/server/"
	name := "server.go"

	return template.CreateFile(data, tpl, path, name)
}
