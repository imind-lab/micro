/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
	"os"
	"text/template"
)

// 生成Client
func CreateClient(data *Data) error {
	// 生成client.go
	var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package client

import (
	"context"
	"strconv"
)

var {{.Service}}s map[string]*{{.Service}}Client

var opts = Options{Name: "{{.Service}}", Tls: true}

func init() {
	{{.Service}}s = make(map[string]*{{.Service}}Client)
}

type Options struct {
	Name string
	Tls  bool
}

func New(ctx context.Context, opt ...Option) (*{{.Service}}Client, error) {
	for _, o := range opt {
		o(&opts)
	}
	key := opts.Name + strconv.FormatBool(opts.Tls)
	{{.Service}}Client, ok := {{.Service}}s[key]
	if !ok {
		{{.Service}}Client, err := New{{.Svc}}Client(ctx, opts.Name, opts.Tls)
		if err == nil {
			{{.Service}}s[key] = {{.Service}}Client
		}
		return {{.Service}}Client, err
	}
	return {{.Service}}Client, nil
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

func Close() {
	for _, client := range {{.Service}}s {
		client.Close()
	}
}
`

	t, err := template.New("client").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/client/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "client.go"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成template.go
	tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package client

import (
	"context"
	"{{.Domain}}/{{.Project}}/{{.Service}}/server/proto/{{.Service}}"
	"github.com/imind-lab/micro/grpc"
	"io"
)

type {{.Service}}Client struct {
	{{.Service}}.{{.Svc}}ServiceClient
	closer io.Closer
}

func New{{.Svc}}Client(ctx context.Context, name string, tls bool) (*{{.Service}}Client, error) {
	conn, closer, err := grpc.ClientConn(ctx, name, tls)
	if err != nil {
		return nil, err
	}
	return &{{.Service}}Client{
		{{.Svc}}ServiceClient: {{.Service}}.New{{.Svc}}ServiceClient(conn),
		closer:	     closer,
	}, nil
}

func (tc *{{.Service}}Client) Close() error {
	return tc.closer.Close()
}
`

	t, err = template.New("clientmain").Parse(tpl)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + data.Service + ".go"
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
