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
	"gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/cmd"
)

func main() {
	cmd.Execute()
}
`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/"
    name := "main.go"

    return template.CreateFile(data, tpl, path, name)
}
