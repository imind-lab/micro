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
 *  {{.Service}}
 *
 *  Create by songli on 2023/02/10
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package subscriber

import (
    "context"

    "github.com/imind-lab/micro/v2/broker"
    "github.com/imind-lab/micro/v2/log"
    "go.uber.org/zap"
)

type {{.Service}} struct {
    ctx context.Context
}

func New{{.Service}}(ctx context.Context) *{{.Service}} {
    svc := &{{.Service}}{ctx}
    return svc
}

func (svc *{{.Service}}) CreateHandle(msg *broker.Message) error {
    logger := log.GetLogger(svc.ctx)
    logger.Debug("sample_create", zap.String("key", msg.Key), zap.String("body", string(msg.Body)))
    return nil
}

func (svc *{{.Service}}) UpdateCountHandle(msg *broker.Message) error {
    logger := log.GetLogger(svc.ctx)
    logger.Debug("sample_update_count", zap.String("key", msg.Key), zap.String("body", string(msg.Body)))
    return nil
}
`

	path := "./" + data.Name + "/application/" + data.Name + "/event/subscriber/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
