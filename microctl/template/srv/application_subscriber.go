/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package srv

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成client/service.go
func CreateApplicationSubscriber(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package subscriber

import (
	"context"
	"github.com/imind-lab/micro/log"

	"github.com/imind-lab/micro/broker"
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

	t, err := template.New("application_service").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/application/" + data.Service + "/event/subscriber/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + data.Service + ".go"

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
