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

// 生成domain/domain.go
func CreateDomainDomain(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Package}}

import (
    "context"
    "github.com/imind-lab/micro/v2/dao"

    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
    repository "{{.Domain}}/{{.Repo}}/repository/{{.Name}}"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}/model"
)

type {{.Service}}Domain interface {
    Create{{.Service}}(ctx context.Context, m model.{{.Service}}) error

    Get{{.Service}}ById(ctx context.Context, id int) (*{{.Package}}.{{.Service}}, error)
    Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) (*{{.Package}}.{{.Service}}List, error)
    Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) (*{{.Package}}.{{.Service}}List, error)

    Update{{.Service}}Type(ctx context.Context, id, typ int) (int8, error)
    Delete{{.Service}}ById(ctx context.Context, id int) (int8, error)

    //+IMindScaffold! Do not modify or delete it
}

type sampleDomain struct {
    dao.Cache
    repo repository.{{.Service}}Repository
}

func New{{.Service}}Domain(repo repository.{{.Service}}Repository) {{.Service}}Domain {
    dm := sampleDomain{
        Cache: dao.NewCache(),
        repo:  repo,
    }
    return dm
}
`

	path := "./" + data.Name + "/domain/" + data.Name + "/"
	name := "domain.go"

	return template.CreateFile(data, tpl, path, name)
}
