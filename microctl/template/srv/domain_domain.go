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

// 生成domain/domain.go
func CreateDomainDomain(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Service}}

import (
	"context"

	"github.com/imind-lab/micro/dao"
	
	{{.Service}} "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto"
	repository "{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/persistence"
)

type {{.Svc}}Domain interface {
	Create{{.Svc}}(ctx context.Context, m model.{{.Svc}}) error

	Get{{.Svc}}ById(ctx context.Context, id int) (*{{.Service}}.{{.Svc}}, error)
	Get{{.Svc}}List0(ctx context.Context, status, pageSize, pageNum int, desc bool) (*{{.Service}}.{{.Svc}}List, error)
	Get{{.Svc}}List1(ctx context.Context, status, pageSize, lastId int, desc bool) (*{{.Service}}.{{.Svc}}List, error)

	Update{{.Svc}}Status(ctx context.Context, id, status int) (int8, error)
	Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error)

	// This commentary is for scaffolding. Do not modify or delete it
}

type {{.Service}}Domain struct {
	dao.Cache
	repo repository.{{.Svc}}Repository
}

func New{{.Svc}}Domain() {{.Svc}}Domain {
	repo := persistence.New{{.Svc}}Repository()
	dm := {{.Service}}Domain{
		Cache: dao.NewCache(),
		repo:  repo,
	}
	return dm
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

	fileName := dir + "domain.go"

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
