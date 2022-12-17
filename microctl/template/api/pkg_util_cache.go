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

// 生成pkg/util/cache.go
func CreatePkgUtilCache(data *tpl.Data) error {
	var tpl = `package util

import (
        "{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
        "github.com/imind-lab/micro/util"
)

func CacheKey(keys ...string) string {
        return constant.CachePrefix + util.AppendString(keys...)
}
`

	t, err := template.New("cmd_server").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/util/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "cache.go"

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
