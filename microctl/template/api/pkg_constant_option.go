/**
 *  MindLab
 *
 *  Create by songli on 2021/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成pkg/constant/option.go
func CreatePkgConstantOption(data *tpl.Data) error {
	var tpl = `/**
 *  IMindLab
 *
 *  Create by songli on {{.Year}}/03/03
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package constant

import (
	"time"
)

//CRequestTimeout 并发请求超时时间
const CRequestTimeout = time.Second * 10

const DBName = "imind"
const Realtime = false

const MQName = "business"
const GreetQueueLen = 32
`

	t, err := template.New("cmd_server").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/constant/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "option.go"

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
