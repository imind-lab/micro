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

// 生成domain/convert.go
func CreateDomainConvert(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Package}}

import (
    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}/model"
)

func {{.Service}}OutMap(ins []model.{{.Service}}, fn func(model.{{.Service}}) *{{.Package}}.{{.Service}}) []*{{.Package}}.{{.Service}} {
    var outs []*{{.Package}}.{{.Service}}
    for _, in := range ins {
        outs = append(outs, fn(in))
    }
    return outs
}

func {{.Service}}Out(in model.{{.Service}}) *{{.Package}}.{{.Service}} {
    if in.IsEmpty() {
        return nil
    }

    out := &{{.Package}}.{{.Service}}{}

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

	path := "./" + data.Name + "/domain/" + data.Name + "/"
	name := "convert.go"

	return template.CreateFile(data, tpl, path, name)
}
