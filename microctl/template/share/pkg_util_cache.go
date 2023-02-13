/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package share

import (
    "github.com/imind-lab/micro/v2/microctl/template"
)

const (
    _PkgUtilPath = "/pkg/util/"
)

// 生成pkg/util/cache.go
func CreatePkgUtilCache(data *template.Data) error {
    var tpl = `package util

import (
    "github.com/imind-lab/micro/v2/util"

    "{{.Domain}}/{{.Repo}}/pkg/constant"
)

const (
    _PkgUtilPath = "/pkg/util/"
)

func CacheKey(keys ...string) string {
    return constant.CachePrefix + util.AppendString(keys...)
}
`

    path := "./" + data.Name + data.Suffix + _PkgUtilPath
    name := "cache.go"

    return template.CreateFile(data, tpl, path, name)
}
