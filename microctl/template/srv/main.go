/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成main.go
func CreateMain(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package main

import (
	"{{.Domain}}/{{.Project}}/{{.Service}}/cmd"
)

func main() {
	cmd.Execute()
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/"
	name := "main.go"

	return template.CreateFile(data, tpl, path, name)
}
