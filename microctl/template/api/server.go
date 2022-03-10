/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成server/server.go
func CreateServer(data *tpl.Data) error {
	var tpl = `package server

import (
	"fmt"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"{{.Domain}}/{{.Project}}/{{.Service}}-api/application/{{.Service}}/proto"
	"{{.Domain}}/{{.Project}}/{{.Service}}-api/application/{{.Service}}/service"
	"github.com/imind-lab/micro"
	grpcx "github.com/imind-lab/micro/grpc"
)

func Serve() error {
	svc := micro.NewService()

	grpcCred := grpcx.NewGrpcCred()

	svc.Init(
		micro.ServerCred(grpcCred.ServerCred()),
		micro.ClientCred(grpcCred.ClientCred()))

	grpcSrv := svc.GrpcServer()
	{{.Service}}_api.Register{{.Svc}}ServiceServer(grpcSrv, service.New{{.Svc}}Service(svc.Options().Logger))

	// 注册gRPC-Gateway
	endPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.grpc"))
	fmt.Println(endPoint)

	mux := svc.ServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(grpcCred.ClientCred())}
	err := {{.Service}}_api.Register{{.Svc}}ServiceHandlerFromEndpoint(svc.Options().Context, mux, endPoint, opts)
	if err != nil {
		return err
	}
	return svc.Run()
}
`

	t, err := template.New("main").Parse(tpl)
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
