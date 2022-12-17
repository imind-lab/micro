/**
 *  MindLab
 *
 *  Create by songli on 2021/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package srv

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成client/service.go
func CreateClientService(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package client

import (
	"context"

	"github.com/imind-lab/rainbow"
	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type {{.Service}}Client struct {
	{{.Service}}.{{.Svc}}ServiceClient
	context  context.Context
	provider *tracesdk.TracerProvider
}

func New{{.Svc}}Client(ctx context.Context, name string, tls bool) (*{{.Service}}Client, error) {
	conn, provider, err := rainbow.ClientConn(ctx, name, tls)
	if err != nil {
		return nil, err
	}
	return &{{.Service}}Client{
		{{.Svc}}ServiceClient: {{.Service}}.New{{.Svc}}ServiceClient(conn),
		context:             ctx,
		provider:            provider,
	}, nil
}

func (cli *{{.Service}}Client) Close() error {
	return cli.provider.Shutdown(cli.context)
}
`

	t, err := template.New("client_service").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/client/"

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
