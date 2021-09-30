package api

import (
	"os"
	"text/template"

	tp "github.com/imind-lab/micro/microctl/template"
)

// 生成server
func CreateServer(data *tp.Data) error {
	var tpl = `package server

import (
	"fmt"

	"{{.Domain}}/spf13/viper"
	"google.golang.org/grpc"

	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
	"{{.Domain}}/{{.Project}}/{{.Service}}/server/proto/{{.Service}}"
	"{{.Domain}}/{{.Project}}/{{.Service}}/server/service"
	"{{.Domain}}/{{.Project}}/{{.Service}}/server/subscriber"
	"{{.Domain}}/{{.Project}}/micro"
	"{{.Domain}}/{{.Project}}/micro/broker"
	"{{.Domain}}/{{.Project}}/micro/grpcx"
)

func Serve() error {
	svc := micro.NewService()

	// 初始化kafka代理
	endpoint, err := broker.NewBroker(constant.MQName)
	if err != nil {
		return err
	}
	// 设置消息队列事件处理器（可选）
	mqHandler := subscriber.New{{.Svc}}(svc.Options().Context)
	endpoint.Subscribe(
		broker.Processor{Topic: endpoint.Options().Topics["createuser"], Handler: mqHandler.CreateHandle, Retry: 1},
		broker.Processor{Topic: endpoint.Options().Topics["updateusercount"], Handler: mqHandler.UpdateCountHandle, Retry: 0},
	)

	grpcCred := grpcx.NewGrpcCred()

	svc.Init(
		micro.Broker(endpoint),
		micro.ServerCred(grpcCred.ServerCred()),
		micro.ClientCred(grpcCred.ClientCred()))

	grpcSrv := svc.GrpcServer()
	{{.Service}}.Register{{.Svc}}ServiceServer(grpcSrv, service.New{{.Svc}}Service())

	// 注册gRPC-Gateway
	endPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.grpc"))
	fmt.Println(endPoint)

	mux := svc.ServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(grpcCred.ClientCred())}
	err = {{.Service}}.Register{{.Svc}}ServiceHandlerFromEndpoint(svc.Options().Context, mux, endPoint, opts)
	if err != nil {
		return err
	}
	return svc.Run()
}
`

	t, err := template.New("server").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/server/"

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
