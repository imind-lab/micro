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

package {{.Service}}

import (
	{{.Service}}_api "gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/application/{{.Service}}/proto"
	{{.Service}} "gitlab.imind.tech/{{.Repo}}/{{.Service}}/application/{{.Service}}/proto"
)

func {{.Svc}}Map(pos []*{{.Service}}.{{.Svc}}, fn func(*{{.Service}}.{{.Svc}}) *{{.Service}}_api.{{.Svc}}) []*{{.Service}}_api.{{.Svc}} {
	var dtos []*{{.Service}}_api.{{.Svc}}
	for _, po := range pos {
		dtos = append(dtos, fn(po))
	}
	return dtos
}

func {{.Svc}}Api2Srv(po *{{.Service}}_api.{{.Svc}}) *{{.Service}}.{{.Svc}} {
	if po == nil {
		return nil
	}

	dto := &{{.Service}}.{{.Svc}}{}
	dto.Id = po.Id
	dto.Name = po.Name
	dto.ViewNum = po.ViewNum
	dto.Type = po.Type
	dto.CreateTime = po.CreateTime
	dto.UpdateDatetime = po.UpdateDatetime
	dto.CreateDatetime = po.CreateDatetime

	return dto
}

func {{.Svc}}Srv2Api(dto *{{.Service}}.{{.Svc}}) *{{.Service}}_api.{{.Svc}} {
	if dto == nil {
		return nil
	}

	po := &{{.Service}}_api.{{.Svc}}{}
	po.Id = dto.Id
	po.Name = dto.Name
	po.ViewNum = dto.ViewNum
	po.Type = dto.Type
	po.CreateTime = dto.CreateTime
	po.UpdateDatetime = dto.UpdateDatetime
	po.CreateDatetime = dto.CreateDatetime

	return po
}

func {{.Svc}}ListSrv2Api(dto *{{.Service}}.{{.Svc}}List) *{{.Service}}_api.{{.Svc}}List {
	if dto == nil {
		return nil
	}

	po := &{{.Service}}_api.{{.Svc}}List{}
	po.Total = dto.Total
	po.TotalPage = dto.TotalPage
	po.CurPage = dto.CurPage
	po.Datalist = {{.Svc}}Map(dto.Datalist, {{.Svc}}Srv2Api)

	return po
}
`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/domain/" + data.Service + "/"
    name := "convert.go"

    return template.CreateFile(data, tpl, path, name)
}
