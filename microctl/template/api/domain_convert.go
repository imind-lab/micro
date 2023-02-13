/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
	"github.com/imind-lab/micro/v2/microctl/template"
)

// 生成domain/convert.go
func CreateDomainConvert(data *template.Data) error {
	var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/03/03
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Package}}

import (
    {{.Package}}_api "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/{{.Name}}/proto"
    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
)

func {{.Service}}Map(pos []*{{.Package}}.{{.Service}}, fn func(*{{.Package}}.{{.Service}}) *{{.Package}}_api.{{.Service}}) []*{{.Package}}_api.{{.Service}} {
    var dtos []*{{.Package}}_api.{{.Service}}
    for _, po := range pos {
        dtos = append(dtos, fn(po))
    }
    return dtos
}

func {{.Service}}Api2Srv(po *{{.Package}}_api.{{.Service}}) *{{.Package}}.{{.Service}} {
    if po == nil {
        return nil
    }

    dto := &{{.Package}}.{{.Service}}{}
    dto.Id = po.Id
    dto.Name = po.Name
    dto.ViewNum = po.ViewNum
    dto.Type = po.Type
    dto.CreateTime = po.CreateTime
    dto.UpdateDatetime = po.UpdateDatetime
    dto.CreateDatetime = po.CreateDatetime

    return dto
}

func {{.Service}}Srv2Api(dto *{{.Package}}.{{.Service}}) *{{.Package}}_api.{{.Service}} {
    if dto == nil {
        return nil
    }

    po := &{{.Package}}_api.{{.Service}}{}
    po.Id = dto.Id
    po.Name = dto.Name
    po.ViewNum = dto.ViewNum
    po.Type = dto.Type
    po.CreateTime = dto.CreateTime
    po.UpdateDatetime = dto.UpdateDatetime
    po.CreateDatetime = dto.CreateDatetime

    return po
}

func {{.Service}}ListSrv2Api(dto *{{.Package}}.{{.Service}}List) *{{.Package}}_api.{{.Service}}List {
    if dto == nil {
        return nil
    }

    po := &{{.Package}}_api.{{.Service}}List{}
    po.Total = dto.Total
    po.TotalPage = dto.TotalPage
    po.CurPage = dto.CurPage
    po.Datalist = {{.Service}}Map(dto.Datalist, {{.Service}}Srv2Api)

    return po
}
`

	path := "./" + data.Name + "-api/domain/" + data.Name + "/"
	name := "convert.go"

	return template.CreateFile(data, tpl, path, name)
}
