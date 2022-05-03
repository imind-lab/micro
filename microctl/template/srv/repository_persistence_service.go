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
func CreateRepositoryPersistenceService(data *tpl.Data) error {
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
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/sentinel"
	"github.com/imind-lab/micro/status"
	"github.com/imind-lab/micro/tracing"
	"github.com/imind-lab/micro/util"
	errorsx "github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
	utilx "{{.Domain}}/{{.Project}}/{{.Service}}/pkg/util"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

func (repo {{.Svc}}Repository) Create{{.Svc}}(ctx context.Context, {{.Service}} model.{{.Svc}}) (model.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	err := repo.CreateModel(ctx, &{{.Service}})
	if err != nil {
		return {{.Service}}, err
	}
	err = repo.CacheModel(ctx, {{.Service}}, constant.Cache{{.Svc}}, {{.Service}}.Id, constant.CacheMinute5)
	return {{.Service}}, err
}

// 忽略部分字段的更新
func (repo {{.Svc}}Repository) Update{{.Svc}}WithOmit(ctx context.Context, {{.Service}} model.{{.Svc}}, columns ...string) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	err := repo.UpdateWithOmit(ctx, &{{.Service}}, columns...)
	if err != nil {
		return err
	}

	return repo.DelModelCache(ctx, constant.Cache{{.Svc}}, {{.Service}}.Id)
}

// 只更新指定的部分字段
func (repo {{.Svc}}Repository) Update{{.Svc}}WithSelect(ctx context.Context, {{.Service}} model.{{.Svc}}, columns ...string) error {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	err := repo.UpdateWithSelect(ctx, &{{.Service}}, columns...)
	if err != nil {
		return err
	}
	return repo.DelModelCache(ctx, constant.Cache{{.Svc}}, {{.Service}}.Id)
}

// 根据Id获取{{.Svc}}(有缓存)
func (repo {{.Svc}}Repository) Get{{.Svc}}ById(ctx context.Context, id int) (model.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	logger := log.GetLogger(ctx)

	var m model.{{.Svc}}
	err := repo.GetModelCache(ctx, &m, constant.Cache{{.Svc}}, id)
	logger.Debug("GetModelCache", zap.Any("{{.Service}}", m), zap.Error(err))
	if err == nil {
		return m, nil
	}

	m, err = repo.Find{{.Svc}}ById(ctx, id)
	if err != nil {
		return m, errorsx.WithMessage(err, "{{.Svc}}Repository.Get{{.Svc}}ById")
	}

	if m.IsEmpty() {
		repo.CacheModelDefault(ctx, m, constant.Cache{{.Svc}}, id)
	} else {
		tool := util.NewCacheTool()
		expire := constant.CacheMinute5 + tool.RandExpire(120)
		err := repo.CacheModel(ctx, m, constant.Cache{{.Svc}}, id, expire)
		fmt.Println(err)
	}

	return m, nil
}

// 根据Id获取{{.Svc}}(无缓存)
func (repo {{.Svc}}Repository) Find{{.Svc}}ById(ctx context.Context, id int) (model.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	var m model.{{.Svc}}
	tx := repo.DB(ctx).Where("id = ?", id)
	err := tx.First(&m).Error
	if err != nil {
		logger := log.GetLogger(ctx)
		logger.Error("data select failed", zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) {
			return m, status.ErrDBDeadlineExceeded
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, nil
		}
		return m, status.ErrDBQuery
	}
	return m, nil
}

func (repo {{.Svc}}Repository) Get{{.Svc}}sCount(ctx context.Context, status int) (int64, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)

	key := utilx.CacheKey(constant.Cache{{.Svc}}Cnt, strconv.Itoa(status))
	cnt, err := repo.Redis().GetNumber(ctx, key)
	if err == nil {
		return cnt, nil
	}
	cnt, err = repo.Find{{.Svc}}sCount(ctx, status)
	if err != nil {
		return 0, errorsx.WithMessage(err, "{{.Svc}}Repository.Get{{.Svc}}sCount")
	}
	err = repo.Redis().Set(ctx, key, cnt, constant.CacheMinute5)
	if err != nil {
		logger.Error("redis.Set", zap.String("key", key), zap.Error(err))
	}
	return cnt, nil
}

func (repo {{.Svc}}Repository) Find{{.Svc}}sCount(ctx context.Context, status int) (int64, error) {
	var count int64
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Select("count(id)")
	tx = tx.Where("status=?", status)
	if err := tx.Count(&count).Error; err != nil {
		return 0, errorsx.Wrap(err, "{{.Svc}}Repository.Find{{.Svc}}sCount")
	}
	return count, nil
}

// Get{{.Svc}}List0 Get the list of {{.Service}} with cache for pageNum
func (repo {{.Svc}}Repository) Get{{.Svc}}List0(ctx context.Context, status, pageSize, pageNum int, desc bool) ([]model.{{.Svc}}, int, error) {
	logger := log.GetLogger(ctx)

	ids, cnt, err := repo.Get{{.Svc}}List0Ids(ctx, status, pageSize, pageNum, desc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Service}}sRepository.Get{{.Svc}}sList0.Get{{.Svc}}sListIds")
	}

	{{.Service}}s, err := repo.Get{{.Svc}}List4Concurrent(ctx, ids, repo.Get{{.Svc}}ById)
	logger.Debug("Get{{.Svc}}List4Concurrent", zap.Any("{{.Service}}s", {{.Service}}s), zap.Error(err))
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Svc}}Repository.Get{{.Svc}}List0.Get{{.Svc}}List4Concurrent")
	}
	return {{.Service}}s, cnt, nil
}

// Get{{.Svc}}List0Ids Get the list of {{.Service}} id with cache for pageNum
func (repo {{.Svc}}Repository) Get{{.Svc}}List0Ids(ctx context.Context, status, pageSize, pageNum int, desc bool) ([]int, int, error) {
	key := utilx.CacheKey(constant.Cache{{.Svc}}Ids, "0_", strconv.Itoa(status))
	var (
		ids []int
		cnt int
		err error
	)
	ids, cnt, err = repo.Redis().SortedSetRange(ctx, key, int64(pageSize), int64(pageNum), desc)
	if err == nil {
		return ids, cnt, nil
	}

	ids, args, err := repo.Find{{.Svc}}List0Ids(ctx, status, pageSize, pageNum, desc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Svc}}Repository.Get{{.Svc}}List0")
	}
	expire := constant.CacheMinute5 + util.RandDuration(120)
	repo.Redis().SortedSetSet(ctx, key, args, expire)
	return ids, len(args), nil
}

// Find{{.Svc}}List0Ids Get the list of {{.Service}} id without cache for pageNum
func (repo {{.Svc}}Repository) Find{{.Svc}}List0Ids(ctx context.Context, status, pageSize, pageNum int, desc bool) ([]int, []*redis.Z, error) {
	limit := 1000
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Select("id")
	tx = tx.Where("status=?", status)
	if desc {
		tx = tx.Order("create_time DESC")
	} else {
		tx = tx.Order("create_time")
	}
	tx = tx.Limit(limit)

	return repo.FetchList0ID(ctx, tx, pageSize, pageNum, limit, desc)
}

// Get{{.Svc}}List1 Get the list of {{.Service}} with cache for lastId
func (repo {{.Svc}}Repository) Get{{.Svc}}List1(ctx context.Context, status, pageSize, lastId int, desc bool) ([]model.{{.Svc}}, int, error) {
	logger := log.GetLogger(ctx)

	ids, cnt, err := repo.Get{{.Svc}}List1Ids(ctx, status, pageSize, lastId, desc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Service}}sRepository.Get{{.Svc}}sList1.Get{{.Svc}}sListIds")
	}

	{{.Service}}s, err := repo.Get{{.Svc}}List4Concurrent(ctx, ids, repo.Get{{.Svc}}ById)
	logger.Debug("Get{{.Svc}}List4Concurrent", zap.Any("{{.Service}}s", {{.Service}}s), zap.Error(err))
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Svc}}Repository.Get{{.Svc}}List1.Get{{.Svc}}List4Concurrent")
	}
	return {{.Service}}s, cnt, nil
}

// Get{{.Svc}}List1Ids Get the list of {{.Service}} id with cache for lastId
func (repo {{.Svc}}Repository) Get{{.Svc}}List1Ids(ctx context.Context, status, pageSize, lastId int, desc bool) ([]int, int, error) {
	key := utilx.CacheKey(constant.Cache{{.Svc}}Ids, "1_", strconv.Itoa(status))
	var (
		ids []int
		cnt int
		err error
	)
	ids, cnt, err = repo.Redis().SortedSetRangeByScore(ctx, key, int64(pageSize), int64(lastId), desc)
	if err == nil {
		return ids, cnt, nil
	}

	ids, args, err := repo.Find{{.Svc}}List1Ids(ctx, status, pageSize, lastId, desc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, "{{.Svc}}Repository.Get{{.Svc}}List0")
	}
	expire := constant.CacheMinute5 + util.RandDuration(120)
	repo.Redis().SortedSetSet(ctx, key, args, expire)
	return ids, len(args), nil
}

// Find{{.Svc}}List1Ids Get the list of {{.Service}} id without cache for lastId
func (repo {{.Svc}}Repository) Find{{.Svc}}List1Ids(ctx context.Context, status, pageSize, lastId int, desc bool) ([]int, []*redis.Z, error) {
	limit := 1000
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Select("id")
	tx = tx.Where("status=?", status)
	if lastId > 0 {
		if desc {
			tx = tx.Where("id<?", lastId)
		} else {
			tx = tx.Where("id>?", lastId)
		}
	}
	if desc {
		tx = tx.Order("id DESC")
	} else {
		tx = tx.Order("id")
	}
	tx = tx.Limit(limit)

	return repo.FetchList1ID(ctx, tx, pageSize)
}

func (repo {{.Svc}}Repository) Get{{.Svc}}List4Concurrent(ctx context.Context, ids []int, fn func(context.Context, int) (model.{{.Svc}}, error)) ([]model.{{.Svc}}, error) {
	logger := log.GetLogger(ctx)

	limiter := sentinel.GetHighLimiter()

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

func (repo {{.Svc}}Repository) Update{{.Svc}}Status(ctx context.Context, id, status int) (int8, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("invoke info", zap.Int("id", id), zap.Int("status", status))
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Where("id = ?", id)
	tx = tx.Update("status", status)
	if tx.Error != nil {

		return 0, errorsx.Wrap(tx.Error, "{{.Svc}}Repository.Update{{.Svc}}Status")
	}
	repo.DelModelCache(ctx, constant.Cache{{.Svc}}, id)
	return int8(tx.RowsAffected), nil
}

func (repo {{.Svc}}Repository) Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error) {
	logger := log.GetLogger(ctx)

	logger.Debug("invoke info", zap.Int("id", id))
	tx := repo.DB(ctx).Delete(&model.{{.Svc}}{}, id)
	if tx.Error != nil {
		return 0, errorsx.Wrap(tx.Error, "{{.Svc}}Repository.Delete{{.Svc}}ById")
	}
	repo.DelModelCacheAll(ctx, constant.Cache{{.Svc}}, id)
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
