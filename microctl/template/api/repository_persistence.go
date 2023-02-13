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

// 生成repository/model.go
func CreateRepositoryPersistence(data *template.Data) error {
	var tpl = `package persistence

import (
    "github.com/imind-lab/micro/v2/dao"

    "{{.Domain}}/{{.Repo}}{{.Suffix}}/repository/{{.Name}}"
)

type {{.Service}}Repository struct {
    dao.Dao
}

// New{{.Service}}Repository create a {{.Svc}} repository instance
func New{{.Service}}Repository(dao dao.Dao) {{.Package}}.{{.Service}}Repository {
    repo := {{.Service}}Repository{
        Dao: dao,
    }
    return repo
}
`

	path := "./" + data.Name + "-api/repository/" + data.Name + "/persistence/"
	name := "persistence.go"

	return template.CreateFile(data, tpl, path, name)
}
