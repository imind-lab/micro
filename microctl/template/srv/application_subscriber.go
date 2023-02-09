/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
    "github.com/imind-lab/micro/v2/microctl/template"
)

// 生成client/service.go
func CreateApplicationSubscriber(data *template.Data) error {
    var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package subscriber

import (
	"context"

	"github.com/imind-lab/micro/broker"
	"github.com/imind-lab/micro/log"
	"go.uber.org/zap"
)

type {{.Svc}} struct {
	ctx context.Context
}

func New{{.Svc}}(ctx context.Context) *{{.Svc}} {
	svc := &{{.Svc}}{ctx}
	return svc
}

func (svc *{{.Svc}}) CreateHandle(msg *broker.Message) error {
	logger := log.GetLogger(svc.ctx)
	logger.Debug("{{.Service}}_create", zap.String("key", msg.Key), zap.String("body", string(msg.Body)))
	return nil
}

func (svc *{{.Svc}}) UpdateCountHandle(msg *broker.Message) error {
	logger := log.GetLogger(svc.ctx)
	logger.Debug("{{.Service}}_update_count", zap.String("key", msg.Key), zap.String("body", string(msg.Body)))
	return nil
}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/application/" + data.Service + "/event/subscriber/"
    name := data.Service + ".go"

    return template.CreateFile(data, tpl, path, name)
}
