/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
    "github.com/imind-lab/micro/v2/microctl/template"
    "github.com/imind-lab/micro/v2/microctl/template/share"
)

func CreatePkgGoogleProtos(data *template.Data) error {
    return share.CreatePkgGoogleProtos(data, "-api")
}
