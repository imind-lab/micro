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

// 生成domain/service.go
func CreateDomainService(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package {{.Package}}

import (
    "context"
    "math"

    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/tracing"
    "github.com/imind-lab/micro/v2/util"
    "github.com/pkg/errors"

    {{.Package}} "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}/model"
)

func (dm sampleDomain) Create{{.Service}}(ctx context.Context, m model.{{.Service}}) error {
    _, err := dm.repo.Create{{.Service}}(ctx, m)
    return err
}

func (dm sampleDomain) Get{{.Service}}ById(ctx context.Context, id int) (*{{.Package}}.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)
    logger.Info("sampleDomain.Get{{.Service}}ById invoke")

    m, err := dm.repo.Get{{.Service}}ById(ctx, id)
    if err != nil {
        return nil, errors.WithMessage(err, util.GetFuncName())
    }
    return {{.Service}}Out(m), nil
}

func (dm sampleDomain) Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) (*{{.Package}}.{{.Service}}List, error) {
    list, total, err := dm.repo.Get{{.Service}}List0(ctx, typ, pageSize, pageNum, isDesc)
    if err != nil {
        return nil, err
    }
    samples := {{.Service}}OutMap(list, {{.Service}}Out)

    var totalPage int32 = 0
    if total == 0 {
        pageNum = 1
    } else {
        totalPage = int32(math.Ceil(float64(total) / float64(pageSize)))
    }
    sampleList := &{{.Package}}.{{.Service}}List{}
    sampleList.Datalist = samples
    sampleList.Total = int32(total)
    sampleList.TotalPage = totalPage
    sampleList.CurPage = int32(pageNum)

    return sampleList, nil
}

// 疑问：中间时翻上一页
func (dm sampleDomain) Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) (*{{.Package}}.{{.Service}}List, error) {
    list, total, err := dm.repo.Get{{.Service}}List1(ctx, typ, pageSize, lastId, isDesc)
    if err != nil {
        return nil, err
    }
    samples := {{.Service}}OutMap(list, {{.Service}}Out)

    var totalPage int32 = 0
    if total > 0 {
        totalPage = int32(math.Ceil(float64(total) / float64(pageSize)))
    }
    sampleList := &{{.Package}}.{{.Service}}List{}
    sampleList.Datalist = samples
    sampleList.Total = int32(total)
    sampleList.TotalPage = totalPage
    sampleList.CurPage = 1

    return sampleList, nil
}

func (dm sampleDomain) Update{{.Service}}Type(ctx context.Context, id, typ int) (int8, error) {
    return dm.repo.Update{{.Service}}Type(ctx, id, typ)
}

func (dm sampleDomain) Delete{{.Service}}ById(ctx context.Context, id int) (int8, error) {
    return dm.repo.Delete{{.Service}}ById(ctx, id)
}
`

	path := "./" + data.Name + "/domain/" + data.Name + "/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
