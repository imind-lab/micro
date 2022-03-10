/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成pkg/constant/cache.go
func CreatePkgConstantCache(data *tpl.Data) error {
	var tpl = `/**
 *  IMindLab
 *
 *  Create by songli on {{.Year}}/03/03
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package constant

import "time"

const (
	CachePrefix   = "im_"
	CacheDay30    = time.Hour * 24 * 30
	CacheDay15    = time.Hour * 24 * 15
	CacheDay7     = time.Hour * 24 * 7
	CacheDay3     = time.Hour * 24 * 3
	CacheDay2     = time.Hour * 24 * 2
	CacheDay1     = time.Hour * 24
	CacheHour12   = time.Hour * 12
	CacheHour6    = time.Hour * 6
	CacheHour2    = time.Hour * 2
	CacheHour1    = time.Hour
	CacheMinute30 = time.Minute * 30
	CacheMinute10 = time.Minute * 10
	CacheMinute5  = time.Minute * 5
	CacheMinute1  = time.Minute
	CacheSecond20 = time.Second * 20
)
`

	t, err := template.New("cmd_server").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/constant/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "cache.go"

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
