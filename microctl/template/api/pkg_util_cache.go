/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
    tpl "github.com/imind-lab/micro/v2/microctl/template"
    "github.com/imind-lab/micro/v2/microctl/template/share"
)

// 生成pkg/util/cache.go
func CreatePkgUtilCache(data *tpl.Data) error {
    return share.CreatePkgUtilCache(data)
}
