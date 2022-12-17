/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package share

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成pkg/constant/cache.go
func CreatePkgConstantCache(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package constant

import "time"

const (
	CachePrefix   = "{{.Service}}_"
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

const (
	CacheArticle     = "article_"
	CacheArticleKeys = "article_keys_"
	CacheArticleCnt  = "article_cnt_"
	CacheArticleIds  = "article_ids_"
)

const (
	CacheCategory     = "category_"
	CacheCategoryKeys = "category_keys_"
	CacheCategoryCnt  = "category_cnt_"
	CacheCategoryIds  = "category_ids_"
)

const (
	CacheCmsInfo     = "cms_info_"
	CacheCmsInfoKeys = "cms_info_keys_"
	CacheCmsInfoCnt  = "cms_info_cnt_"
	CacheCmsInfoIds  = "cms_info_ids_"
)

const (
	CacheConfig     = "config_"
	CacheConfigKeys = "config_keys_"
	CacheConfigCnt  = "config_cnt_"
	CacheConfigIds  = "config_ids_"
)

const (
	CacheEmail     = "email_"
	CacheEmailKeys = "email_keys_"
	CacheEmailCnt  = "email_cnt_"
	CacheEmailIds  = "email_ids_"
)

const (
	CacheMenu     = "menu_"
	CacheMenuKeys = "menu_keys_"
	CacheMenuCnt  = "menu_cnt_"
	CacheMenuIds  = "menu_ids_"
)

const (
	Cache{{.Svc}}     = "m_"
	Cache{{.Svc}}Keys = "m_keys_"
	Cache{{.Svc}}Cnt  = "m_cnt_"
	Cache{{.Svc}}Ids  = "m_ids_"
)

const (
	CacheD{{.Svc}} = "d_"
)
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/pkg/constant/"
	name := "cache.go"

	return template.CreateFile(data, tpl, path, name)
}
