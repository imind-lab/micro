/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package api

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成repository/model.go
func CreateRepositoryPersistence(data *template.Data) error {
	var tpl = `package persistence

import (
	"github.com/imind-lab/micro/dao"

	"gitlab.imind.tech/{{.Project}}/{{.Service}}-api/repository/{{.Service}}"
)

type {{.Svc}}Repository struct {
	dao.Dao
}

// New{{.Svc}}Repository create a {{.Service}} repository instance
func New{{.Svc}}Repository(dao dao.Dao) {{.Service}}.{{.Svc}}Repository {
	repo := {{.Svc}}Repository{
		Dao: dao,
	}
	return repo
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/repository/" + data.Service + "/persistence/"
	name := "persistence.go"

	return template.CreateFile(data, tpl, path, name)
}
