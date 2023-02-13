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
func CreateRepositoryPersistence(data *template.Data) error {
	var tpl = `package persistence

import (
    "context"
    "errors"
    "strconv"
    "time"

   {{if .MQ}}"github.com/imind-lab/micro/v2/broker"{{end}}
    "github.com/imind-lab/micro/v2/dao"
    "github.com/imind-lab/micro/v2/log"
    "github.com/imind-lab/micro/v2/status"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "{{.Domain}}/{{.Repo}}/pkg/constant"
    utilx "{{.Domain}}/{{.Repo}}/pkg/util"
    "{{.Domain}}/{{.Repo}}/repository/{{.Name}}"
)

type {{.Service}}Repository struct {
    dao.Dao

   {{if .MQ}}broker broker.Broker{{end}}
}

// New{{.Service}}Repository create a sample repository instance
func New{{.Service}}Repository(dao dao.Dao{{if .MQ}}, broker broker.Broker{{end}}) {{.Package}}.{{.Service}}Repository {
    repo := {{.Service}}Repository{
        Dao:                         dao,
       {{if .MQ}}broker: broker,{{end}}
    }
    return repo
}

// CreateModel store the model object in the database
func (repo {{.Service}}Repository) CreateModel(ctx context.Context, m interface{}) error {
    if err := repo.DB(ctx).Create(m).Error; err != nil {
        logger := log.GetLogger(ctx)
        logger.Error("Data insert failed", zap.Error(err))
        if errors.Is(err, context.DeadlineExceeded) {
            return status.ErrDBDeadlineExceeded
        }
        return status.ErrDBCreate
    }
    return nil
}

// UpdateWithOmit updates the other fields of the model to the database in addition to the specified fields
func (repo {{.Service}}Repository) UpdateWithOmit(ctx context.Context, m interface{}, columns ...string) error {
    if err := repo.DB(ctx).Omit(columns...).Save(m).Error; err != nil {
        logger := log.GetLogger(ctx)
        logger.Error("Data update failed", zap.Error(err))
        if errors.Is(err, context.DeadlineExceeded) {
            return status.ErrDBDeadlineExceeded
        }
        return status.ErrDBUpdate
    }
    return nil
}

// UpdateWithSelect update the fields specified by the model to the database
func (repo {{.Service}}Repository) UpdateWithSelect(ctx context.Context, m interface{}, columns ...string) error {
    if err := repo.DB(ctx).Select(columns).Save(m).Error; err != nil {
        logger := log.GetLogger(ctx)
        logger.Error("Data update failed", zap.Error(err))
        if errors.Is(err, context.DeadlineExceeded) {
            return status.ErrDBDeadlineExceeded
        }
        return status.ErrDBUpdate
    }
    return nil
}

// CacheModel
func (repo {{.Service}}Repository) CacheModel(ctx context.Context, m interface{}, preKey string, id int, expire time.Duration) error {
    key := utilx.CacheKey(preKey, strconv.Itoa(id))
    err := repo.Redis().HashTableSet(ctx, key, m, expire)
    if err != nil {
        return err
    }
    keys := utilx.CacheKey(preKey, "keys_", strconv.Itoa(id))
    return repo.Redis().SAdd(ctx, keys, key)
}

func (repo {{.Service}}Repository) CacheModelDefault(ctx context.Context, m interface{}, preKey string, id int) error {
    key := utilx.CacheKey(preKey, strconv.Itoa(id))
    return repo.Redis().HashTableSet(ctx, key, m, constant.CacheMinute1)
}

func (repo {{.Service}}Repository) GetModelCache(ctx context.Context, m interface{}, preKey string, id int) error {
    key := utilx.CacheKey(preKey, strconv.Itoa(id))
    return repo.Redis().HashTableGet(ctx, key, m)
}

func (repo {{.Service}}Repository) DelModelCache(ctx context.Context, preKey string, id int) error {
    key := utilx.CacheKey(preKey, strconv.Itoa(id))
    err := repo.Redis().Del(ctx, key)
    if err != nil {
        logger := log.GetLogger(ctx)
        logger.Warn("Del.error", zap.String("key", key), zap.Error(err))
    }
    return nil
}

func (repo {{.Service}}Repository) DelModelCacheAll(ctx context.Context, preKey string, id int) error {
    key := utilx.CacheKey(preKey, "keys_", strconv.Itoa(id))
    err := repo.Redis().SetDelKeys(ctx, key)
    if err != nil {
        logger := log.GetLogger(ctx)
        logger.Warn("SetDelKeys.error", zap.String("key", key), zap.Error(err))
    }
    return nil
}

func (repo {{.Service}}Repository) FetchList0ID(ctx context.Context, tx *gorm.DB, pageSize, pageNum, limit int, desc bool) ([]int, []redis.Z, error) {
    logger := log.GetLogger(ctx)

    rows, err := tx.Rows()
    if err != nil {
        logger.Error("Data select failed", zap.Error(err))
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, nil, status.ErrDBDeadlineExceeded
        } else if errors.Is(err, gorm.ErrRecordNotFound) {
            return []int{}, []redis.Z{}, nil
        }
        return nil, nil, status.ErrDBQuery
    }
    defer rows.Close()

    var ids []int
    var args []redis.Z
    start := float64((pageNum - 1) * pageSize)
    var index float64 = 0
    var max = float64(limit)
    idx := index
    for rows.Next() {
        var (
            id int
        )
        err = rows.Scan(&id)
        if err != nil {
            logger.Error("Data scan failed", zap.Error(err))
            return nil, nil, status.ErrDBUpdate
        }
        if index >= start && len(ids) < pageSize {
            ids = append(ids, id)
        }
        idx = index
        if desc {
            idx = max - index
        }
        args = append(args, redis.Z{Score: idx, Member: id})
        index++
    }

    if err = rows.Err(); err != nil {
        logger.Error("Data err failed", zap.Error(err))
        return nil, nil, status.ErrDBUpdate
    }
    return ids, args, nil
}

func (repo {{.Service}}Repository) FetchList1ID(ctx context.Context, tx *gorm.DB, pageSize int) ([]int, []redis.Z, error) {
    logger := log.GetLogger(ctx)

    rows, err := tx.Rows()
    if err != nil {
        logger.Error("Data select failed", zap.Error(err))
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, nil, status.ErrDBDeadlineExceeded
        } else if errors.Is(err, gorm.ErrRecordNotFound) {
            return []int{}, []redis.Z{}, nil
        }
        return nil, nil, status.ErrDBQuery
    }
    defer rows.Close()

    var ids []int
    var args []redis.Z
    for rows.Next() {
        var (
            id int
        )
        err = rows.Scan(&id)
        if err != nil {
            logger.Error("Data scan failed", zap.Error(err))
            return nil, nil, status.ErrDBUpdate
        }

        if len(ids) < pageSize {
            ids = append(ids, id)
        }

        args = append(args, redis.Z{Score: float64(id), Member: id})
    }

    if err = rows.Err(); err != nil {
        logger.Error("Data err failed", zap.Error(err))
        return nil, nil, status.ErrDBUpdate
    }
    return ids, args, nil
}
`

	path := "./" + data.Name + "/repository/" + data.Name + "/persistence/"
	name := "persistence.go"

	return template.CreateFile(data, tpl, path, name)
}
