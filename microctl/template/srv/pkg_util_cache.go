/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
    "os"
    "text/template"

    tpl "github.com/imind-lab/micro/v2/microctl/template"
)

// 生成pkg/util/cache.go
func CreatePkgUtilCache(data *tpl.Data) error {
    var tpl = `package util

import (
	"github.com/imind-lab/micro/v2/util"

	"{{.Domain}}/{{.Repo}}/pkg/constant"
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
    dir := "./" + data.Name + "/pkg/util/"

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
