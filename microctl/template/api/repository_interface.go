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

// 生成domain/service.go
func CreateRepositoryInterface(data *template.Data) error {
	var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/06/01
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Package}}

import (
    "context"

    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
)

type {{.Service}}Repository interface {
    Create{{.Service}}(ctx context.Context, name string, typ int32) error
    Get{{.Service}}ById(ctx context.Context, id int32) (*{{.Package}}.{{.Service}}, error)
    Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int32, isDesc bool) (*{{.Package}}.{{.Service}}List, error)
    Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int32, isDesc bool) (*{{.Package}}.{{.Service}}List, error)
    Update{{.Service}}Type(ctx context.Context, id, typ int32) error
    Delete{{.Service}}ById(ctx context.Context, id int32) error
    Get{{.Service}}ListByIds(ctx context.Context, ids []int32) ([]*{{.Package}}.{{.Service}}, error)

    // This commentary is for scaffolding. Do not modify or delete it
}
`

	path := "./" + data.Name + "-api/repository/" + data.Name + "/"
	name := "repository.go"

	return template.CreateFile(data, tpl, path, name)
}
