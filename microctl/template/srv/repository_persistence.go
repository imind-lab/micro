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

// 生成repository/model.go
func CreateRepositoryPersistence(data *tpl.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package persistence

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/imind-lab/micro/dao"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/sentinel"
	"github.com/imind-lab/micro/tracing"
	"github.com/imind-lab/micro/util"
	errorsx "github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
	utilx "{{.Domain}}/{{.Project}}/{{.Service}}/pkg/util"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

type {{.Service}}Repository struct {
	dao.Dao
}

//New{{.Svc}}Repository 创建用户仓库实例
func New{{.Svc}}Repository() {{.Service}}.{{.Svc}}Repository {
	rep := dao.NewDao(constant.DBName)
	repo := {{.Service}}Repository{
		Dao: rep,
	}
	return repo
}

func (repo {{.Service}}Repository) Create{{.Svc}}(ctx context.Context, m model.{{.Svc}}) (model.{{.Svc}}, error) {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Create{{.Svc}}")
	defer span.Finish()

	if err := repo.DB().Create(&m).Error; err != nil {
		return m, errorsx.Wrap(err, "{{.Service}}Repository.Create{{.Svc}}")
	}
	repo.Cache{{.Svc}}(ctx, m)
	return m, nil
}

// 忽略部分字段的更新
func (repo {{.Service}}Repository) Update{{.Svc}}WithOmit(ctx context.Context, m model.{{.Svc}}, columns ...string) error {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Update{{.Svc}}WithOmit")
	defer span.Finish()

	if err := repo.DB().Omit(columns...).Save(&m).Error; err != nil {
		return errorsx.Wrap(err, "{{.Service}}Repository.Update{{.Svc}}WithOmit")
	}
	repo.Cache{{.Svc}}(ctx, m)
	return nil
}

// 只更新指定的部分字段
func (repo {{.Service}}Repository) Update{{.Svc}}WithSelect(ctx context.Context, m model.{{.Svc}}, columns ...string) error {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Update{{.Svc}}WithSelect")
	defer span.Finish()

	if err := repo.DB().Select(columns).Save(&m).Error; err != nil {
		return errorsx.Wrap(err, "{{.Service}}Repository.Update{{.Svc}}WithSelect")
	}
	repo.DelCache{{.Svc}}(ctx, int(m.Id))
	return nil
}

func (repo {{.Service}}Repository) Cache{{.Svc}}(ctx context.Context, m model.{{.Svc}}) error {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Cache{{.Svc}}")
	defer span.Finish()

	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Cache{{.Svc}}")
	key := utilx.CacheKey(constant.Cache{{.Svc}}, strconv.Itoa(int(m.Id)))
	expire := constant.CacheMinute5
	err := repo.Redis().HashTableSet(ctx, key, m, expire)
	if err != nil {
		logger.Warn("SetHashTable.error", zap.Error(err))
	}
	keys := utilx.CacheKey(constant.Cache{{.Svc}}Keys, strconv.Itoa(int(m.Id)))
	err = repo.Redis().SAdd(ctx, keys, key).Err()
	if err != nil {
		logger.Warn("SAdd.error", zap.Error(err))
	}
	return nil
}

func (repo {{.Service}}Repository) DelCache{{.Svc}}(ctx context.Context, id int) error {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Cache{{.Svc}}")
	defer span.Finish()

	key := utilx.CacheKey(constant.Cache{{.Svc}}Keys, strconv.Itoa(id))
	err := repo.Redis().SetDelKeys(ctx, key)
	if err != nil {
		logger := log.GetLogger(ctx, "{{.Service}}Repository", "DelCache{{.Svc}}")
		logger.Warn("SetDelKeys.error", zap.Error(err))
	}
	return nil
}

// 根据Id获取{{.Svc}}(有缓存)
func (repo {{.Service}}Repository) Get{{.Svc}}ById(ctx context.Context, id int, opt ...{{.Service}}.ObjectByIdOption) (model.{{.Svc}}, error) {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Get{{.Svc}}ById")
	defer span.Finish()

	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Get{{.Svc}}ById")
	opts := {{.Service}}.NewObjectByIdOptions(util.RandDuration(120))
	for _, o := range opt {
		o(opts)
	}

	var m model.{{.Svc}}
	key := utilx.CacheKey(constant.Cache{{.Svc}}, strconv.Itoa(id))
	err := repo.Redis().HashTableGet(ctx, key, &m)
	logger.Debug("redis.HGetAll", zap.Any("{{.Service}}", m), zap.String("key", key), zap.Error(err))
	if err == nil {
		return m, nil
	}

	m, err = repo.Find{{.Svc}}ById(ctx, id)
	if err != nil {
		return m, errorsx.WithMessage(err, "{{.Service}}Repository.Get{{.Svc}}ById")
	}

	expire := constant.CacheMinute5 + opts.RandExpire
	if m.IsEmpty() {
		expire = constant.CacheMinute1
	}
	repo.Redis().HashTableSet(ctx, key, m, expire)
	return m, nil
}

// 根据Id获取{{.Svc}}(无缓存)
func (repo {{.Service}}Repository) Find{{.Svc}}ById(ctx context.Context, id int) (model.{{.Svc}}, error) {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Find{{.Svc}}ById")
	defer span.Finish()

	var m model.{{.Svc}}
	err := repo.DB().Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, nil
		}
		return m, errorsx.Wrap(err, "{{.Service}}Repository.Find{{.Svc}}ById")
	}
	return m, nil
}

func (repo {{.Service}}Repository) Get{{.Svc}}sCount(ctx context.Context, {{.Service}}Id int) (int64, error) {
	span, ctx := tracing.StartSpan(ctx, "{{.Service}}Repository.Get{{.Svc}}sCount")
	defer span.Finish()

	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Get{{.Svc}}sCount")

	key := utilx.CacheKey(constant.Cache{{.Svc}}Cnt, strconv.Itoa({{.Service}}Id))
	cnt, err := repo.Redis().GetNumber(ctx, key)
	if err == nil {
		return cnt, nil
	}
	cnt, err = repo.Find{{.Svc}}sCount(ctx, {{.Service}}Id)
	if err != nil {
		return 0, errorsx.WithMessage(err, "{{.Service}}Repository.Get{{.Svc}}sCount")
	}
	err = repo.Redis().Set(ctx, key, cnt, constant.CacheMinute5).Err()
	if err != nil {
		logger.Error("redis.Set", zap.String("key", key), zap.Error(err))
	}
	return cnt, nil
}

func (repo {{.Service}}Repository) Find{{.Svc}}sCount(ctx context.Context, {{.Service}}Id int) (int64, error) {
	var count int64
	tx := repo.DB().Model(model.{{.Svc}}{}).Select("count(id)")
	tx = tx.Where("{{.Service}}_id=?", {{.Service}}Id)
	if err := tx.Count(&count).Error; err != nil {
		return 0, errorsx.Wrap(err, "{{.Service}}Repository.Find{{.Svc}}sCount")
	}
	return count, nil
}

// 获取{{.Svc}}列表(有缓存)
func (repo {{.Service}}Repository) Get{{.Svc}}List(ctx context.Context, status, lastId, pageSize, pageNum int, desc bool) ([]model.{{.Svc}}, int, error) {
	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Get{{.Svc}}List")

	ids, cnt, err := repo.Get{{.Svc}}ListIds(ctx, status, lastId, pageSize, pageNum, desc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Service}}sRepository.Get{{.Svc}}sList.Get{{.Svc}}sListIds")
	}

	ctx1, cancel := context.WithTimeout(ctx, constant.CRequestTimeout)
	defer cancel()

	{{.Service}}s, err := repo.Get{{.Svc}}List4Concurrent(ctx1, ids, repo.Get{{.Svc}}ById)
	logger.Debug("Get{{.Svc}}List4Concurrent", zap.Any("{{.Service}}s", {{.Service}}s), zap.Error(err))
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Service}}Repository.Get{{.Svc}}List.Get{{.Svc}}List4Concurrent")
	}
	return {{.Service}}s, cnt, nil
}

// 获取{{.Svc}}Id列表(有缓存)
func (repo {{.Service}}Repository) Get{{.Svc}}ListIds(ctx context.Context, status, lastId, pageSize, pageNum int, desc bool) ([]int, int, error) {
	key := utilx.CacheKey(constant.Cache{{.Svc}}Ids, strconv.Itoa(status))
	var (
		ids []int
		cnt int
		err error
	)
	if lastId > 0 {
		ids, cnt, err = repo.Redis().SortedSetRangeByScore(ctx, key, int64(lastId), int64(pageSize), desc)
	} else {
		ids, cnt, err = repo.Redis().SortedSetRange(ctx, key, int64(pageSize), int64(pageNum), desc)
	}
	if err == nil {
		return ids, cnt, nil
	}

	ids, args, err := repo.Find{{.Svc}}ListIds(ctx, status, lastId, pageSize, pageNum, desc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Service}}Repository.Get{{.Svc}}List")
	}
	expire := constant.CacheMinute5 + util.RandDuration(120)
	repo.Redis().SortedSetSet(ctx, key, args, expire)
	return ids, len(args), nil
}

// 获取{{.Svc}}Id列表(无缓存)
func (repo {{.Service}}Repository) Find{{.Svc}}ListIds(ctx context.Context, status, lastId, pageSize, pageNum int, desc bool) ([]int, []*redis.Z, error) {
	fmt.Println(status, lastId, pageSize, pageNum, desc)
	tx := repo.DB().Model(model.{{.Svc}}{}).Select("id")
	tx = tx.Where("status=?", status)
	if desc {
		tx = tx.Order("create_time DESC")
	} else {
		tx = tx.Order("create_time")
	}
	tx = tx.Limit(1000)
	rows, err := tx.Rows()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []int{}, []*redis.Z{}, nil
		}
		return nil, nil, errorsx.Wrap(err, "{{.Service}}Repository.Find{{.Svc}}ListIds.Rows")
	}
	defer rows.Close()

	var ids []int
	var args []*redis.Z
	start := float64((pageNum - 1) * pageSize)
	var index float64 = 0
	idx := index
	for rows.Next() {
		var (
			id int
		)
		err = rows.Scan(&id)
		if err != nil {
			return nil, nil, errorsx.Wrap(err, "{{.Service}}Repository.Find{{.Svc}}ListIds.Scan")
		}

		check := false
		if lastId == 0 {
			if index >= start {
				check = true
			}
		} else if lastId > id {
			check = true
		}
		if check {
			if len(ids) < pageSize {
				ids = append(ids, id)
			}
		}
		idx = index
		if desc {
			idx = 10000 - index
		}
		args = append(args, &redis.Z{Score: idx, Member: id})
		index++
	}
	if err = rows.Err(); err != nil {
		return nil, nil, errorsx.Wrap(err, "{{.Service}}Repository.Find{{.Svc}}ListIds.Err")
	}
	return ids, args, nil
}

func (repo {{.Service}}Repository) Get{{.Svc}}List4Concurrent(ctx context.Context, ids []int, fn func(context.Context, int, ...{{.Service}}.ObjectByIdOption) (model.{{.Svc}}, error)) ([]model.{{.Svc}}, error) {
	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Get{{.Svc}}List4Concurrent")

	limiter := sentinel.GetLimiter()

	var wg sync.WaitGroup

	count := len(ids)
	outputs := make([]*concurrent{{.Svc}}Output, count)
	wg.Add(count)

	for idx, id := range ids {
		err := limiter.Wait(context.Background())
		if err != nil {
			logger.Warn("limiter wait error", zap.Error(err))
		}
		go func(idx int, id int, wg *sync.WaitGroup) {
			defer wg.Done()
			{{.Service}}, err := fn(ctx, id)
			outputs[idx] = &concurrent{{.Svc}}Output{
				object: {{.Service}},
				err:    err,
			}
		}(idx, id, &wg)
	}
	wg.Wait()

	{{.Service}}s := make([]model.{{.Svc}}, 0, count)
	for _, output := range outputs {
		if output.err == nil {
			{{.Service}}s = append({{.Service}}s, output.object)
		}
	}
	return {{.Service}}s, nil
}

type concurrent{{.Svc}}Output struct {
	object model.{{.Svc}}
	err    error
}

func (repo {{.Service}}Repository) Update{{.Svc}}Status(ctx context.Context, id, status int) (int8, error) {
	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Update{{.Svc}}Status")

	logger.Debug("invoke info", zap.Int("id", id), zap.Int("status", status))
	tx := repo.DB().Model(model.{{.Svc}}{}).Where("id = ?", id)
	tx = tx.Update("status", status)
	if tx.Error != nil {
		return 0, errorsx.Wrap(tx.Error, "{{.Service}}Repository.Update{{.Svc}}Status")
	}
	key := utilx.CacheKey(constant.Cache{{.Svc}}, strconv.Itoa(int(id)))
	reply, err := repo.Redis().Del(ctx, key).Result()
	if err != nil {
		logger.Warn("Del Cache", zap.String("key", key), zap.Int64("reply", reply), zap.Error(err))
	}
	return int8(tx.RowsAffected), nil
}

func (repo {{.Service}}Repository) Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error) {
	logger := log.GetLogger(ctx, "{{.Service}}Repository", "Delete{{.Svc}}ById")

	logger.Debug("invoke info", zap.Int("id", id))
	tx := repo.DB().Delete(&model.{{.Svc}}{}, id)
	if tx.Error != nil {
		return 0, errorsx.Wrap(tx.Error, "{{.Service}}Repository.Delete{{.Svc}}ById")
	}
	key := utilx.CacheKey(constant.Cache{{.Svc}}, strconv.Itoa(id))
	reply, err := repo.Redis().Del(ctx, key).Result()
	logger.Debug("Del Cache", zap.String("key", key), zap.Int64("reply", reply), zap.Error(err))

	return int8(tx.RowsAffected), nil
}
`

	t, err := template.New("repository_model").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/repository/" + data.Service + "/persistence/"

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
