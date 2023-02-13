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
    "{{.Domain}}/{{.Repo}}{{.Suffix}}/cmd"
)

func main() {
    cmd.Execute()
}
`

	path := "./" + data.Name + "-api/"
	name := "main.go"

	return template.CreateFile(data, tpl, path, name)
}
