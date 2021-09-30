/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
	"os"
	"strings"
	"text/template"
)

// 生成model
func CreateModel(data *Data) error {
	var tpl = `/**
 *  {{.Svc}}Lab
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
	Id	           int32  ${backtick}gorm:"primary_key" redis:"id"${backtick}
	Name	       string ${backtick}redis:"name,omitempty"${backtick}
	ViewNum	       int32  ${backtick}redis:"view_num,omitempty"${backtick}
	Status	       int32  ${backtick}redis:"status,omitempty"${backtick}
	CreateTime     int64  ${backtick}redis:"create_time,omitempty"${backtick}
	CreateDatetime string ${backtick}redis:"create_datetime,omitempty"${backtick}
	UpdateDatetime string ${backtick}redis:"update_datetime,omitempty"${backtick}
}

func ({{.Svc}}) TableName() string {
	return "tbl_{{.Service}}"
}

func (m *{{.Svc}}) BeforeCreate(tx *gorm.DB) error {
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

	t, err := template.New("model").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/server/model/"

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
