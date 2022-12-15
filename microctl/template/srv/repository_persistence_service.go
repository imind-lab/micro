/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成repository/model.go
func CreateRepositoryPersistenceService(data *template.Data) error {
	var tpl = `/**
 *  {{.Svc}}
 *
 *  Create by songli on 2021/06/01
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
	errorsx "github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	{{if .MQ}}
	"github.com/imind-lab/micro/broker"{{end}}
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/sentinel"
	"github.com/imind-lab/micro/status"
	"github.com/imind-lab/micro/tracing"
	"github.com/imind-lab/micro/util"

	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
	utilx "{{.Domain}}/{{.Project}}/{{.Service}}/pkg/util"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}/model"
)

const _CDType = "type=?"

func (repo {{.Svc}}Repository) Create{{.Svc}}(ctx context.Context, {{.Service}} model.{{.Svc}}) (model.{{.Svc}}, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	err := repo.CreateModel(ctx, &{{.Service}})
	if err != nil {
		return {{.Service}}, err
	}
	logger := log.GetLogger(ctx)
	err = repo.CacheModel(ctx, &{{.Service}}, constant.Cache{{.Svc}}, {{.Service}}.Id, constant.CacheMinute5)
	if err != nil {
		logger.Warn("CacheModel error", zap.Error(err))
	}
	{{if .MQ}}
	err = repo.broker.Publish(&broker.Message{
		Topic: repo.broker.Options().Topics["{{.Service}}create"],
		Body:  []byte(fmt.Sprintf("{{.Svc}} %s Created", {{.Service}}.Name)),
	})
	if err != nil {
		logger.Error("Kafka publish error", zap.Error(err))
	}{{end}}
	return {{.Service}}, nil
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

	var {{.Service}} model.{{.Svc}}
	err := repo.GetModelCache(ctx, &{{.Service}}, constant.Cache{{.Svc}}, id)
	logger.Debug("GetModelCache", zap.Any("{{.Service}}", {{.Service}}), zap.Error(err))
	if err == nil {
		return {{.Service}}, nil
	}

	{{.Service}}, err = repo.Find{{.Svc}}ById(ctx, id)
	if err != nil {
		return {{.Service}}, errorsx.WithMessage(err, util.GetFuncName())
	}

	if {{.Service}}.IsEmpty() {
		repo.CacheModelDefault(ctx, &{{.Service}}, constant.Cache{{.Svc}}, id)
	} else {
		tool := util.NewCacheTool()
		expire := constant.CacheMinute5 + tool.RandExpire(120)
		err := repo.CacheModel(ctx, &{{.Service}}, constant.Cache{{.Svc}}, id, expire)
		if err != nil {
			logger.Warn("CacheModel", zap.Error(err))
		}
	}

	return {{.Service}}, nil
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
			return m, status.ErrRecordNotFound
		}
		return m, status.ErrDBQuery
	}
	return m, nil
}

func (repo {{.Svc}}Repository) Get{{.Svc}}sCount(ctx context.Context, typ int) (int64, error) {
	ctx, span := tracing.StartSpan(ctx)
	span.End()

	logger := log.GetLogger(ctx)

	key := utilx.CacheKey(constant.Cache{{.Svc}}Cnt, strconv.Itoa(typ))
	cnt, err := repo.Redis().GetNumber(ctx, key)
	if err == nil {
		return cnt, nil
	}
	cnt, err = repo.Find{{.Svc}}sCount(ctx, typ)
	if err != nil {
		return 0, errorsx.WithMessage(err, util.GetFuncName())
	}
	err = repo.Redis().Set(ctx, key, cnt, constant.CacheMinute5)
	if err != nil {
		logger.Error("redis.Set", zap.String("key", key), zap.Error(err))
	}
	return cnt, nil
}

func (repo {{.Svc}}Repository) Find{{.Svc}}sCount(ctx context.Context, typ int) (int64, error) {
	var count int64
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Select("count(id)")
	tx = tx.Where(_CDType, typ)
	if err := tx.Count(&count).Error; err != nil {
		return 0, errorsx.Wrap(err, util.GetFuncName())
	}
	return count, nil
}

// Get{{.Svc}}List0 Get the list of {{.Service}} with cache for pageNum
func (repo {{.Svc}}Repository) Get{{.Svc}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]model.{{.Svc}}, int, error) {
	logger := log.GetLogger(ctx)

	ids, cnt, err := repo.Get{{.Svc}}List0Ids(ctx, typ, pageSize, pageNum, isDesc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
	}

	{{.Service}}s, err := repo.Get{{.Svc}}List4Concurrent(ctx, ids, repo.Get{{.Svc}}ById)
	logger.Debug("Get{{.Svc}}List4Concurrent", zap.Any("{{.Service}}s", {{.Service}}s), zap.Error(err))
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
	}
	return {{.Service}}s, cnt, nil
}

// Get{{.Svc}}List0Ids Get the list of {{.Service}} id with cache for pageNum
func (repo {{.Svc}}Repository) Get{{.Svc}}List0Ids(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]int, int, error) {
	key := utilx.CacheKey(constant.Cache{{.Svc}}Ids, "0_", strconv.Itoa(typ))
	var (
		ids []int
		cnt int
		err error
	)
	cnt, err = repo.Redis().SortedSetRange(ctx, key, int64(pageSize), int64(pageNum), isDesc, &ids)
	fmt.Println("SortedSetRange", cnt, err, ids)
	if err == nil {
		return ids, cnt, nil
	}

	ids, args, err := repo.Find{{.Svc}}List0Ids(ctx, typ, pageSize, pageNum, isDesc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
	}
	expire := constant.CacheMinute5 + util.RandDuration(120)
	repo.Redis().SortedSetSet(ctx, key, args, expire)
	return ids, len(args), nil
}

// Find{{.Svc}}List0Ids Get the list of {{.Service}} id without cache for pageNum
func (repo {{.Svc}}Repository) Find{{.Svc}}List0Ids(ctx context.Context, typ, pageSize, pageNum int, desc bool) ([]int, []*redis.Z, error) {
	limit := 1000
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Select("id")
	tx = tx.Where(_CDType, typ)
	if desc {
		tx = tx.Order("create_time DESC")
	} else {
		tx = tx.Order("create_time")
	}
	tx = tx.Limit(limit)

	return repo.FetchList0ID(ctx, tx, pageSize, pageNum, limit, desc)
}

// Get{{.Svc}}List1 Get the list of {{.Service}} with cache for lastId
func (repo {{.Svc}}Repository) Get{{.Svc}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]model.{{.Svc}}, int, error) {
	logger := log.GetLogger(ctx)

	ids, cnt, err := repo.Get{{.Svc}}List1Ids(ctx, typ, pageSize, lastId, isDesc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
	}

	{{.Service}}s, err := repo.Get{{.Svc}}List4Concurrent(ctx, ids, repo.Get{{.Svc}}ById)
	logger.Debug("Get{{.Svc}}List4Concurrent", zap.Any("{{.Service}}s", {{.Service}}s), zap.Error(err))
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
	}
	return {{.Service}}s, cnt, nil
}

// Get{{.Svc}}List1Ids Get the list of {{.Service}} id with cache for lastId
func (repo {{.Svc}}Repository) Get{{.Svc}}List1Ids(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]int, int, error) {
	key := utilx.CacheKey(constant.Cache{{.Svc}}Ids, "1_", strconv.Itoa(typ))
	var (
		ids []int
		cnt int
		err error
	)
	cnt, err = repo.Redis().SortedSetRangeByScore(ctx, key, int64(pageSize), int64(lastId), isDesc, &ids)
	if err == nil {
		return ids, cnt, nil
	}

	ids, args, err := repo.Find{{.Svc}}List1Ids(ctx, typ, pageSize, lastId, isDesc)
	if err != nil {
		return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
	}
	expire := constant.CacheMinute5 + util.RandDuration(120)
	repo.Redis().SortedSetSet(ctx, key, args, expire)
	return ids, len(args), nil
}

// Find{{.Svc}}List1Ids Get the list of {{.Service}} id without cache for lastId
func (repo {{.Svc}}Repository) Find{{.Svc}}List1Ids(ctx context.Context, typ, pageSize, lastId int, desc bool) ([]int, []*redis.Z, error) {
	limit := 1000
	tx := repo.DB(ctx).Model(model.{{.Svc}}{}).Select("id")
	tx = tx.Where(_CDType, typ)
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

func (repo {{.Svc}}Repository) Update{{.Svc}}Type(ctx context.Context, id, typ int) (int8, error) {
	logger := log.GetLogger(ctx)
	logger.Debug("invoke info", zap.Int("id", id), zap.Int("type", typ))
	tx := repo.DB(ctx).Model(&model.{{.Svc}}{}).Where("id = ?", id)
	tx = tx.Update("type", typ)
	if tx.Error != nil {
		return 0, errorsx.Wrap(tx.Error, util.GetFuncName())
	}
	repo.DelModelCache(ctx, constant.Cache{{.Svc}}, id)
	return int8(tx.RowsAffected), nil
}

func (repo {{.Svc}}Repository) Delete{{.Svc}}ById(ctx context.Context, id int) (int8, error) {
	logger := log.GetLogger(ctx)

	logger.Debug("invoke info", zap.Int("id", id))
	tx := repo.DB(ctx).Delete(&model.{{.Svc}}{}, id)
	if tx.Error != nil {
		return 0, errorsx.Wrap(tx.Error, util.GetFuncName())
	}
	repo.DelModelCacheAll(ctx, constant.Cache{{.Svc}}, id)
	return int8(tx.RowsAffected), nil
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/repository/" + data.Service + "/persistence/"
	name := data.Service + ".go"

	return template.CreateFile(data, tpl, path, name)
}
