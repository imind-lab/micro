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

// 生成server/server.go
func CreateServer(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/imind-lab/micro"
	{{if .MQ}}
	"github.com/imind-lab/micro/broker"{{end}}
	grpcx "github.com/imind-lab/micro/grpc"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/tracing"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	{{if .MQ}}
	"{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/event/subscriber"{{end}}
	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
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

	runtime.SetBlockProfileRate(1)

	// initialize tracing
	provider, err := tracing.InitTracer()
	if err != nil {
		return err
	}
	defer provider.Shutdown(ctx)
	{{if .MQ}}
	// initialize kafka broker
	endpoint, err := broker.NewBroker()
	if err != nil {
		return err
	}
	defer endpoint.Close()

	// set up the handler for MQ
	mqHandler := subscriber.New{{.Svc}}(ctx)
	endpoint.Subscribe(
		broker.Processor{Topic: endpoint.Options().Topics["{{.Service}}create"], Handler: mqHandler.CreateHandle, Retry: 1},
	)
	{{end}}
	grpcCred := grpcx.NewGrpcCred()

	svc := micro.NewService()
	svc.Init(
		micro.Context(ctx),
		micro.Logger(logger),
		micro.Name(conf.Name),
		micro.ServerCred(grpcCred.ServerCred()),
		micro.ClientCred(grpcCred.ClientCred()))
	//micro.HttpHandler(AuthHandler))

	// 注册gRPC-Gateway
	endPoint := fmt.Sprintf(":%d", conf.Port.Grpc)

	mux := svc.ServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(grpcCred.ClientCred())}
	err = {{.Service}}.Register{{.Svc}}ServiceHandlerFromEndpoint(ctx, mux, endPoint, opts)
	if err != nil {
		return err
	}

	grpcSrv := svc.GrpcServer()

	{{.Service}}Svc := Create{{.Svc}}Service({{if .MQ}}endpoint{{end}})
	{{.Service}}.Register{{.Svc}}ServiceServer(grpcSrv, {{.Service}}Svc)

	// This commentary is for scaffolding. Do not modify or delete it

	return svc.Run()
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("auth") != "auth" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/server/"
	name := "server.go"

	return template.CreateFile(data, tpl, path, name)
}
