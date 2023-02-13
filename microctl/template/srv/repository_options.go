/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
    "os"
    "text/template"

    tpl "github.com/imind-lab/micro/v2/microctl/template"
)

// 生成domain/service.go
func CreateRepositoryOptions(data *tpl.Data) error {
    var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"time"
)

type ObjectByIdOptions struct {
	RandExpire time.Duration
}

func NewObjectByIdOptions(randExpire time.Duration) *ObjectByIdOptions {
	return &ObjectByIdOptions{RandExpire: randExpire}
}

type ObjectByIdOption func(*ObjectByIdOptions)

func ObjectByIdRandExpire(expire time.Duration) ObjectByIdOption {
	return func(o *ObjectByIdOptions) {
		o.RandExpire = expire
	}
}
`

    t, err := template.New("domain_service").Parse(tpl)
    if err != nil {
        return err
    }

    t.Option()
    dir := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "/repository/" + data.Service + "/"

    err = os.MkdirAll(dir, os.ModePerm)
    if err != nil {
        return err
    }

    fileName := dir + "options.go"

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
