/**
 *  MindLab
 *
 *  Create by songli on 2021/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package api

import (
	"github.com/imind-lab/micro/microctl/template"
	"github.com/imind-lab/micro/microctl/template/share"
)

func CreatePkgGoogleProtos(data *template.Data) error {
	return share.CreatePkgGoogleProtos(data, "-api")
}
