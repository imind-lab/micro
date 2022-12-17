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

// 生成main.go
func CreateMain(data *template.Data) error {
	var tpl = `/**
 *  ImindLab
 *
 *  Create by songli on {{.Year}}/03/03
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package main

import (
	"gitlab.imind.tech/{{.Project}}/{{.Service}}-api/cmd"
)

func main() {
	cmd.Execute()
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/"
	name := "main.go"

	return template.CreateFile(data, tpl, path, name)
}
