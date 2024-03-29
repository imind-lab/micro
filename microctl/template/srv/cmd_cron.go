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

// 生成cmd/cron.go
func CreateCmdCron(data *tpl.Data) error {
    var tpl = `package cmd

import (
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"{{.Domain}}/{{.Repo}}/cmd/cron"
)

// 计划任务方法需要幂等
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "show {{.Service}} cronjob sample",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			c := cron.New()
			vf := reflect.ValueOf(c)

			target := strings.Title(args[0])
			method := vf.MethodByName(target)

			if method.IsValid() {
				method.Call([]reflect.Value{})
			} else {
				log.Println("指定的计划任务方法不存在")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}
`

    t, err := template.New("cmd_cron").Parse(tpl)
    if err != nil {
        return err
    }

    t.Option()
    dir := "./" + data.Name + "/cmd/"

    err = os.MkdirAll(dir, os.ModePerm)
    if err != nil {
        return err
    }

    fileName := dir + "cron.go"

    f, err := os.Create(fileName)
    if err != nil {
        return err
    }
    err = t.Execute(f, data)
    if err != nil {
        return err
    }
    f.Close()

    tpl = `package cron

type Cron struct{}

func New() Cron {
	return Cron{}
}
`

    t, err = template.New("cmd_cron_cron").Parse(tpl)
    if err != nil {
        return err
    }

    t.Option()
    dir = "./" + data.Name + "/cmd/cron/"

    err = os.MkdirAll(dir, os.ModePerm)
    if err != nil {
        return err
    }

    fileName = dir + "cron.go"

    f, err = os.Create(fileName)
    if err != nil {
        return err
    }
    err = t.Execute(f, data)
    if err != nil {
        return err
    }
    f.Close()

    tpl = `package cron

import (
	"fmt"
	"time"
)

func (c Cron) EchoTime() {
	fmt.Println(time.Now())
}
`

    t, err = template.New("cmd_cron_sample").Parse(tpl)
    if err != nil {
        return err
    }

    t.Option()

    fileName = dir + data.Package + ".go"

    f, err = os.Create(fileName)
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
