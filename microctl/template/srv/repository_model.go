/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package srv

import (
	"os"
	"strings"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成repository/model.go
func CreateRepositoryModel(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package model

import (
	"gorm.io/gorm"
	"reflect"
	"time"
)

type {{.Svc}} struct {
	Id             int    ${backtick}gorm:"primary_key" redis:"id"${backtick}
	Name           string ${backtick}redis:"name,omitempty"${backtick}
	ViewNum        int16  ${backtick}redis:"view_num,omitempty"${backtick}
	Status         int8   ${backtick}redis:"status,omitempty"${backtick}
	CreateTime     uint32 ${backtick}redis:"create_time,omitempty"${backtick}
	CreateDatetime string ${backtick}redis:"create_datetime,omitempty"${backtick}
	UpdateDatetime string ${backtick}redis:"update_datetime,omitempty"${backtick}
}

func ({{.Svc}}) TableName() string {
	return "tbl_{{.Service}}"
}

func (m *{{.Svc}}) BeforeCreate(tx *gorm.DB) error {
	m.CreateTime = uint32(time.Now().Unix())
	m.CreateDatetime = time.Now().Format("2006-01-02 15:04:05")
	m.UpdateDatetime = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (m *{{.Svc}}) BeforeUpdate(tx *gorm.DB) error {
	m.UpdateDatetime = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (m {{.Svc}}) IsEmpty() bool {
	return reflect.DeepEqual(m, {{.Svc}}{})
}
`
	tpl = strings.Replace(tpl, "${backtick}", "`", -1)
	t, err := template.New("repository_model").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/repository/" + data.Service + "/model/"

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
