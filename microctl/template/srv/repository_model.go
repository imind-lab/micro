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

// 生成repository/model.go
func CreateRepositoryModel(data *template.Data) error {
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

const _DTFormat = "2006-01-02 15:04:05"

type {{.Svc}} struct {
	Id             int    ${backtick}gorm:"primary_key" redis:"id"${backtick}
	Name           string ${backtick}redis:"name,omitempty"${backtick}
	ViewNum        int16  ${backtick}redis:"view_num,omitempty"${backtick}
	Type           int8   ${backtick}redis:"type,omitempty"${backtick}
	CreateTime     uint32 ${backtick}redis:"create_time,omitempty"${backtick}
	CreateDatetime string ${backtick}redis:"create_datetime,omitempty"${backtick}
	UpdateDatetime string ${backtick}redis:"update_datetime,omitempty"${backtick}
}

func ({{.Svc}}) TableName() string {
	return "tbl_{{.Service}}"
}

func (m *{{.Svc}}) BeforeCreate(tx *gorm.DB) error {
	m.CreateTime = uint32(time.Now().Unix())
	m.CreateDatetime = time.Now().Format(_DTFormat)
	m.UpdateDatetime = time.Now().Format(_DTFormat)
	return nil
}

func (m *{{.Svc}}) BeforeUpdate(tx *gorm.DB) error {
	m.UpdateDatetime = time.Now().Format(_DTFormat)
	return nil
}

func (m {{.Svc}}) IsEmpty() bool {
	return reflect.DeepEqual(m, {{.Svc}}{})
}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/repository/" + data.Service + "/model/"

    name := data.Service + ".go"

    return template.CreateFile(data, tpl, path, name)
}
