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
func CreateRepositoryPersistenceService(data *template.Data) error {
	var tpl = `/**
 *  {{.Service}}
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

   {{if .MQ}}
    "github.com/imind-lab/micro/v2/broker"{{end}}
    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/sentinel"
    "github.com/imind-lab/micro/v2/status"
    "github.com/imind-lab/micro/v2/tracing"
    "github.com/imind-lab/micro/v2/util"
    errorsx "github.com/pkg/errors"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "{{.Domain}}/{{.Repo}}/pkg/constant"
    utilx "{{.Domain}}/{{.Repo}}/pkg/util"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}/model"
)

const _CDType = "type=?"

func (repo {{.Service}}Repository) Create{{.Service}}(ctx context.Context, {{.Svc}} model.{{.Service}}) (model.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    err := repo.CreateModel(ctx, &{{.Svc}})
    if err != nil {
        return {{.Svc}}, err
    }
    logger := log.GetLogger(ctx)
    err = repo.CacheModel(ctx, &{{.Svc}}, constant.Cache{{.Service}}, {{.Svc}}.Id, constant.CacheMinute5)
    if err != nil {
        logger.Warn("CacheModel error", zap.Error(err))
    }
   {{if .MQ}}
    err = repo.broker.Publish(&broker.Message{
        Topic: repo.broker.Options().Topics["{{.Package}}_create"],
        Body:  []byte(fmt.Sprintf("{{.Service}} %s Created", {{.Svc}}.Name)),
    })
    if err != nil {
        logger.Error("Kafka publish error", zap.Error(err))
    }{{end}}
    return {{.Svc}}, nil
}

// 忽略部分字段的更新
func (repo {{.Service}}Repository) Update{{.Service}}WithOmit(ctx context.Context, {{.Svc}} model.{{.Service}}, columns ...string) error {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    err := repo.UpdateWithOmit(ctx, &{{.Svc}}, columns...)
    if err != nil {
        return err
    }

    return repo.DelModelCache(ctx, constant.Cache{{.Service}}, {{.Svc}}.Id)
}

// 只更新指定的部分字段
func (repo {{.Service}}Repository) Update{{.Service}}WithSelect(ctx context.Context, {{.Svc}} model.{{.Service}}, columns ...string) error {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    err := repo.UpdateWithSelect(ctx, &{{.Svc}}, columns...)
    if err != nil {
        return err
    }
    return repo.DelModelCache(ctx, constant.Cache{{.Service}}, {{.Svc}}.Id)
}

// 根据Id获取{{.Service}}(有缓存)
func (repo {{.Service}}Repository) Get{{.Service}}ById(ctx context.Context, id int) (model.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    logger := log.GetLogger(ctx)

    var {{.Svc}} model.{{.Service}}
    err := repo.GetModelCache(ctx, &{{.Svc}}, constant.Cache{{.Service}}, id)
    logger.Debug("GetModelCache", zap.Any("{{.Package}}", {{.Svc}}), zap.Error(err))
    if err == nil {
        return {{.Svc}}, nil
    }

    {{.Svc}}, err = repo.Find{{.Service}}ById(ctx, id)
    if err != nil {
        return {{.Svc}}, errorsx.WithMessage(err, util.GetFuncName())
    }

    if {{.Svc}}.IsEmpty() {
        repo.CacheModelDefault(ctx, &{{.Svc}}, constant.Cache{{.Service}}, id)
    } else {
        tool := util.NewCacheTool()
        expire := constant.CacheMinute5 + tool.RandExpire(120)
        err := repo.CacheModel(ctx, &{{.Svc}}, constant.Cache{{.Service}}, id, expire)
        if err != nil {
            logger.Warn("CacheModel", zap.Error(err))
        }
    }

    return {{.Svc}}, nil
}

// 根据Id获取{{.Service}}(无缓存)
func (repo {{.Service}}Repository) Find{{.Service}}ById(ctx context.Context, id int) (model.{{.Service}}, error) {
    ctx, span := tracing.StartSpan(ctx)
    defer span.End()

    var m model.{{.Service}}
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

func (repo {{.Service}}Repository) Get{{.Service}}sCount(ctx context.Context, typ int) (int64, error) {
    ctx, span := tracing.StartSpan(ctx)
    span.End()

    logger := log.GetLogger(ctx)

    key := utilx.CacheKey(constant.Cache{{.Service}}Cnt, strconv.Itoa(typ))
    cnt, err := repo.Redis().GetNumber(ctx, key)
    if err == nil {
        return cnt, nil
    }
    cnt, err = repo.Find{{.Service}}sCount(ctx, typ)
    if err != nil {
        return 0, errorsx.WithMessage(err, util.GetFuncName())
    }
    err = repo.Redis().Set(ctx, key, cnt, constant.CacheMinute5)
    if err != nil {
        logger.Error("redis.Set", zap.String("key", key), zap.Error(err))
    }
    return cnt, nil
}

func (repo {{.Service}}Repository) Find{{.Service}}sCount(ctx context.Context, typ int) (int64, error) {
    var count int64
    tx := repo.DB(ctx).Model(model.{{.Service}}{}).Select("count(id)")
    tx = tx.Where(_CDType, typ)
    if err := tx.Count(&count).Error; err != nil {
        return 0, errorsx.Wrap(err, util.GetFuncName())
    }
    return count, nil
}

// Get{{.Service}}List0 Get the list of {{.Service}} with cache for pageNum
func (repo {{.Service}}Repository) Get{{.Service}}List0(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]model.{{.Service}}, int, error) {
    logger := log.GetLogger(ctx)

    ids, cnt, err := repo.Get{{.Service}}List0Ids(ctx, typ, pageSize, pageNum, isDesc)
    if err != nil {
        return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
    }

    {{.Svc}}s, err := repo.Get{{.Service}}List4Concurrent(ctx, ids, repo.Get{{.Service}}ById)
    logger.Debug("Get{{.Service}}List4Concurrent", zap.Any("{{.Package}}_list", {{.Svc}}s), zap.Error(err))
    if err != nil {
        return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
    }
    return {{.Svc}}s, cnt, nil
}

// Get{{.Service}}List0Ids Get the list of {{.Service}} id with cache for pageNum
func (repo {{.Service}}Repository) Get{{.Service}}List0Ids(ctx context.Context, typ, pageSize, pageNum int, isDesc bool) ([]int, int, error) {
    key := utilx.CacheKey(constant.Cache{{.Service}}Ids, "0_", strconv.Itoa(typ))
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

    ids, args, err := repo.Find{{.Service}}List0Ids(ctx, typ, pageSize, pageNum, isDesc)
    if err != nil {
        return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
    }
    expire := constant.CacheMinute5 + util.RandDuration(120)
    repo.Redis().SortedSetSet(ctx, key, args, expire)
    return ids, len(args), nil
}

// Find{{.Service}}List0Ids Get the list of {{.Service}} id without cache for pageNum
func (repo {{.Service}}Repository) Find{{.Service}}List0Ids(ctx context.Context, typ, pageSize, pageNum int, desc bool) ([]int, []redis.Z, error) {
    limit := 1000
    tx := repo.DB(ctx).Model(model.{{.Service}}{}).Select("id")
    tx = tx.Where(_CDType, typ)
    if desc {
        tx = tx.Order("create_time DESC")
    } else {
        tx = tx.Order("create_time")
    }
    tx = tx.Limit(limit)

    return repo.FetchList0ID(ctx, tx, pageSize, pageNum, limit, desc)
}

// Get{{.Service}}List1 Get the list of {{.Service}} with cache for lastId
func (repo {{.Service}}Repository) Get{{.Service}}List1(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]model.{{.Service}}, int, error) {
    logger := log.GetLogger(ctx)

    ids, cnt, err := repo.Get{{.Service}}List1Ids(ctx, typ, pageSize, lastId, isDesc)
    if err != nil {
        return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
    }

    {{.Svc}}s, err := repo.Get{{.Service}}List4Concurrent(ctx, ids, repo.Get{{.Service}}ById)
    logger.Debug("Get{{.Service}}List4Concurrent", zap.Any("{{.Package}}_list", {{.Svc}}s), zap.Error(err))
    if err != nil {
        return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
    }
    return {{.Svc}}s, cnt, nil
}

// Get{{.Service}}List1Ids Get the list of {{.Service}} id with cache for lastId
func (repo {{.Service}}Repository) Get{{.Service}}List1Ids(ctx context.Context, typ, pageSize, lastId int, isDesc bool) ([]int, int, error) {
    key := utilx.CacheKey(constant.Cache{{.Service}}Ids, "1_", strconv.Itoa(typ))
    var (
        ids []int
        cnt int
        err error
    )
    cnt, err = repo.Redis().SortedSetRangeByScore(ctx, key, int64(pageSize), int64(lastId), isDesc, &ids)
    if err == nil {
        return ids, cnt, nil
    }

    ids, args, err := repo.Find{{.Service}}List1Ids(ctx, typ, pageSize, lastId, isDesc)
    if err != nil {
        return nil, 0, errorsx.WithMessage(err, util.GetFuncName())
    }
    expire := constant.CacheMinute5 + util.RandDuration(120)
    repo.Redis().SortedSetSet(ctx, key, args, expire)
    return ids, len(args), nil
}

// Find{{.Service}}List1Ids Get the list of {{.Service}} id without cache for lastId
func (repo {{.Service}}Repository) Find{{.Service}}List1Ids(ctx context.Context, typ, pageSize, lastId int, desc bool) ([]int, []redis.Z, error) {
    limit := 1000
    tx := repo.DB(ctx).Model(model.{{.Service}}{}).Select("id")
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

func (repo {{.Service}}Repository) Get{{.Service}}List4Concurrent(ctx context.Context, ids []int, fn func(context.Context, int) (model.{{.Service}}, error)) ([]model.{{.Service}}, error) {
    logger := log.GetLogger(ctx)

    limiter := sentinel.GetHighLimiter()

    var wg sync.WaitGroup

    count := len(ids)
    outputs := make([]*concurrent{{.Service}}Output, count)
    wg.Add(count)

    for idx, id := range ids {
        err := limiter.Wait(context.Background())
        if err != nil {
            logger.Warn("limiter wait error", zap.Error(err))
        }
        go func(idx int, id int, wg *sync.WaitGroup) {
            defer wg.Done()
            {{.Svc}}, err := fn(ctx, id)
            outputs[idx] = &concurrent{{.Service}}Output{
                object: {{.Svc}},
                err:    err,
            }
        }(idx, id, &wg)
    }
    wg.Wait()

    {{.Svc}}s := make([]model.{{.Service}}, 0, count)
    for _, output := range outputs {
        if output.err == nil {
            {{.Svc}}s = append({{.Svc}}s, output.object)
        }
    }
    return {{.Svc}}s, nil
}

type concurrent{{.Service}}Output struct {
    object model.{{.Service}}
    err    error
}

func (repo {{.Service}}Repository) Update{{.Service}}Type(ctx context.Context, id, typ int) (int8, error) {
    logger := log.GetLogger(ctx)
    logger.Debug("invoke info", zap.Int("id", id), zap.Int("type", typ))
    tx := repo.DB(ctx).Model(&model.{{.Service}}{}).Where("id = ?", id)
    tx = tx.Update("type", typ)
    if tx.Error != nil {
        return 0, errorsx.Wrap(tx.Error, util.GetFuncName())
    }
    repo.DelModelCache(ctx, constant.Cache{{.Service}}, id)
    return int8(tx.RowsAffected), nil
}

func (repo {{.Service}}Repository) Delete{{.Service}}ById(ctx context.Context, id int) (int8, error) {
    logger := log.GetLogger(ctx)

    logger.Debug("invoke info", zap.Int("id", id))
    tx := repo.DB(ctx).Delete(&model.{{.Service}}{}, id)
    if tx.Error != nil {
        return 0, errorsx.Wrap(tx.Error, util.GetFuncName())
    }
    repo.DelModelCacheAll(ctx, constant.Cache{{.Service}}, id)
    return int8(tx.RowsAffected), nil
}
`

	path := "./" + data.Name + "/repository/" + data.Name + "/persistence/"
	name := data.Package + ".go"

	return template.CreateFile(data, tpl, path, name)
}
