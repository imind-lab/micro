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

// 生成server/server.go
func CreateServer(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package server

import (
	"fmt"
	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	"google.golang.org/grpc"

	"github.com/imind-lab/micro"
	grpcx "github.com/imind-lab/micro/grpc"
	"github.com/spf13/viper"
	"{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/service"
)

func Serve() error {
	svc := micro.NewService()

	// 初始化kafka代理
	//endpoint, err := broker.NewBroker(constant.MQName)
	//if err != nil {
	//	return err
	//}
	//// 设置消息队列事件处理器（可选）
	//mqHandler := subscriber.New{{.Svc}}(svc.Options().Context)
	//endpoint.Subscribe(
	//	broker.Processor{Topic: endpoint.Options().Topics["create{{.Service}}"], Handler: mqHandler.CreateHandle, Retry: 1},
	//)

	grpcCred := grpcx.NewGrpcCred()

	svc.Init(
		//micro.Broker(endpoint),
		micro.ServerCred(grpcCred.ServerCred()),
		micro.ClientCred(grpcCred.ClientCred()))

	grpcSrv := svc.GrpcServer()
	{{.Service}}.Register{{.Svc}}ServiceServer(grpcSrv, service.New{{.Svc}}Service())

	// 注册gRPC-Gateway
	endPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.grpc"))

	mux := svc.ServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(grpcCred.ClientCred())}
	err := {{.Service}}.Register{{.Svc}}ServiceHandlerFromEndpoint(svc.Options().Context, mux, endPoint, opts)
	if err != nil {
		return err
	}

	// This commentary is for scaffolding. Do not modify or delete it

	return svc.Run()
}
`

	t, err := template.New("main").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/server/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "server.go"

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
