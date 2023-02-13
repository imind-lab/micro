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

// 生成main.go
func CreateMain(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package main

import (
    "{{.Domain}}/{{.Repo}}/cmd"
)

func main() {
    cmd.Execute()
}
`

	path := "./" + data.Name + "/"
	name := "main.go"

	return template.CreateFile(data, tpl, path, name)
}
