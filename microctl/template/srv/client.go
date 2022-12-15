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

// 生成build/Dockerfile
func CreateClient(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on 2021/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package client

import (
	"context"
	"github.com/imind-lab/micro"

	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
)

var opts = Options{Name: "{{.Service}}", Tls: true}

type Options struct {
	Name string
	Tls  bool
}

func New(ctx context.Context, opt ...Option) ({{.Service}}.{{.Svc}}ServiceClient, func() error, error) {
	for _, o := range opt {
		o(&opts)
	}

	conn, err := micro.ClientConn(ctx, opts.Name, opts.Tls)
	if err != nil {
		return nil, nil, err
	}

	var close = func() error {
		if conn != nil {
			conn.Close()
		}
		return nil
	}

	return {{.Service}}.New{{.Svc}}ServiceClient(conn), close, nil
}

type Option func(*Options)

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Tls(tls bool) Option {
	return func(o *Options) {
		o.Tls = tls
	}
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/client/"
	name := "client.go"

	return template.CreateFile(data, tpl, path, name)
}
