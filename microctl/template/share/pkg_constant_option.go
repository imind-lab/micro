/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package share

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成pkg/constant/option.go
func CreatePkgConstantOption(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package constant

import (
	"time"
)

// CRequestTimeout 并发请求超时时间
const CRequestTimeout = time.Second * 10

const DBName = "imind"
const Realtime = false

const MQName = "business"
const GreetQueueLen = 32
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/pkg/constant/"
	name := "option.go"

	return template.CreateFile(data, tpl, path, name)
}