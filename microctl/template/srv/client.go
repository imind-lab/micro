/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/v2/microctl/template"
)

// 生成build/Dockerfile
func CreateClient(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package client

import (
    "context"
    "github.com/imind-lab/micro/v2"
    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
)

var opts = Options{Name: "sample", Tls: true}

type Options struct {
    Name string
    Tls  bool
}

func New(ctx context.Context, opt ...Option) ({{.Package}}.{{.Service}}ServiceClient, func() error, error) {
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

    return {{.Package}}.New{{.Service}}ServiceClient(conn), close, nil
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

	path := "./" + data.Name + "/client/"
	name := "client.go"

	return template.CreateFile(data, tpl, path, name)
}
