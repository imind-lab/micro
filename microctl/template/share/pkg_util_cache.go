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

// 生成pkg/util/cache.go
func CreatePkgUtilCache(data *template.Data) error {
	var tpl = `package util

import (
	"github.com/imind-lab/micro/util"

	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
)

func CacheKey(keys ...string) string {
	return constant.CachePrefix + util.AppendString(keys...)
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/pkg/util/"
	name := "cache.go"

	return template.CreateFile(data, tpl, path, name)
}
