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

// 生成domain/convert.go
func CreateDomainConvert(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

func {{.Svc}}OutMap(ins []model.{{.Svc}}, fn func(model.{{.Svc}}) *{{.Service}}.{{.Svc}}) []*{{.Service}}.{{.Svc}} {
	var outs []*{{.Service}}.{{.Svc}}
	for _, in := range ins {
		outs = append(outs, fn(in))
	}
	return outs
}

func {{.Svc}}Out(in model.{{.Svc}}) *{{.Service}}.{{.Svc}} {
	if in.IsEmpty() {
		return nil
	}

	out := &{{.Service}}.{{.Svc}}{}

	out.Id = int32(in.Id)
	out.Name = in.Name
	out.ViewNum = int32(in.ViewNum)
	out.Status = int32(in.Status)
	out.CreateTime = in.CreateTime
	out.CreateDatetime = in.CreateDatetime
	out.UpdateDatetime = in.UpdateDatetime

	return out
}
`

	t, err := template.New("domain_convert").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/domain/" + data.Service + "/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "convert.go"

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
