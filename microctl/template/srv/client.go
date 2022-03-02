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

// 生成build/Dockerfile
func CreateClient(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
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

	t, err := template.New("client_client").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
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

	return nil
}
