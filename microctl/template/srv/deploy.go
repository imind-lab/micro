/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package srv

import (
	tpl "github.com/imind-lab/micro/microctl/template"
	"github.com/imind-lab/micro/microctl/template/share"
)

func CreateDeploy(data *tpl.Data) error {
	return share.CreateDeploy(data)
}
