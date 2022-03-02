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

// 生成pkg/constant/cache.go
func CreatePkgConstantCache(data *tpl.Data) error {
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
	CacheCmsPromotion     = "cms_iromotion_"
	CacheCmsPromotionKeys = "cms_iromotion_keys_"
	CacheCmsPromotionCnt  = "cms_iromotion_cnt_"
	CacheCmsPromotionIds  = "cms_iromotion_ids_"
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
	CachePicture     = "third_party_map_"
	CachePictureKeys = "third_party_map_keys_"
	CachePictureCnt  = "third_party_map_cnt_"
	CachePictureIds  = "third_party_map_ids_"
)

const (
	CacheThirdPartyMap     = "third_party_map_"
	CacheThirdPartyMapKeys = "third_party_map_keys_"
	CacheThirdPartyMapCnt  = "third_party_map_cnt_"
	CacheThirdPartyMapIds  = "third_party_map_ids_"
	CacheThirdPartyMapMid  = "third_party_map_mid_"
)

const (
	CacheCategoryRelationship     = "category_relationship_"
	CacheCategoryRelationshipKeys = "category_relationship_keys_"
	CacheCategoryRelationshipCnt  = "category_relationship_cnt_"
	CacheCategoryRelationshipIds  = "category_relationship_ids_"
)

const (
	CacheD{{.Svc}} = "d_"
)
`

	t, err := template.New("cmd_server").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/pkg/constant/"

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
