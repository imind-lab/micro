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

// 生成domain/convert.go
func CreateDomainConvert(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on 2021/06/01
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
	out.Type = int32(in.Type)
	out.CreateTime = in.CreateTime
	out.CreateDatetime = in.CreateDatetime
	out.UpdateDatetime = in.UpdateDatetime

	return out
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/domain/" + data.Service + "/"
	name := "convert.go"

	return template.CreateFile(data, tpl, path, name)
}
