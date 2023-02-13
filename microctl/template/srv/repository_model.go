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
 *  {{.Service}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package model

import (
    "reflect"
    "time"

    "gorm.io/gorm"
)

const _DTFormat = "2006-01-02 15:04:05"

type {{.Service}} struct {
    Id             int    ${backtick}gorm:"primary_key" redis:"id"${backtick}
    Name           string ${backtick}redis:"name,omitempty"${backtick}
    ViewNum        int16  ${backtick}redis:"view_num,omitempty"${backtick}
    Type           int8   ${backtick}redis:"type,omitempty"${backtick}
    CreateTime     uint32 ${backtick}redis:"create_time,omitempty"${backtick}
    CreateDatetime string ${backtick}redis:"create_datetime,omitempty"${backtick}
    UpdateDatetime string ${backtick}redis:"update_datetime,omitempty"${backtick}
}

func ({{.Service}}) TableName() string {
    return "tbl_sample"
}

func (m *{{.Service}}) BeforeCreate(tx *gorm.DB) error {
    m.CreateTime = uint32(time.Now().Unix())
    m.CreateDatetime = time.Now().Format(_DTFormat)
    m.UpdateDatetime = time.Now().Format(_DTFormat)
    return nil
}

func (m *{{.Service}}) BeforeUpdate(tx *gorm.DB) error {
    m.UpdateDatetime = time.Now().Format(_DTFormat)
    return nil
}

func (m {{.Service}}) IsEmpty() bool {
    return reflect.DeepEqual(m, {{.Service}}{})
}
`

	path := "./" + data.Name + "/repository/" + data.Name + "/model/"

	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
