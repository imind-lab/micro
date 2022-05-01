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
	var tpl = `package persistence

import (
	"context"
	"errors"
	"{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
	utilx "{{.Domain}}/{{.Project}}/{{.Service}}/pkg/util"
	"{{.Domain}}/{{.Project}}/{{.Service}}/repository/{{.Service}}"
	"github.com/go-redis/redis/v8"
	"github.com/imind-lab/micro/dao"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/status"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type {{.Svc}}Repository struct {
	dao.Dao
}

// New{{.Svc}}Repository create a {{.Service}} repository instance
func New{{.Svc}}Repository() {{.Service}}.{{.Svc}}Repository {
	rep := dao.NewDao(constant.DBName)
	repo := {{.Svc}}Repository{
		Dao: rep,
	}
	return repo
}

// CreateModel store the model object in the database
func (repo {{.Svc}}Repository) CreateModel(ctx context.Context, m interface{}) error {
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
func (repo {{.Svc}}Repository) UpdateWithOmit(ctx context.Context, m interface{}, columns ...string) error {
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
func (repo {{.Svc}}Repository) UpdateWithSelect(ctx context.Context, m interface{}, columns ...string) error {
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
func (repo {{.Svc}}Repository) CacheModel(ctx context.Context, m interface{}, preKey string, id int, expire time.Duration) error {
	key := utilx.CacheKey(preKey, strconv.Itoa(id))
	err := repo.Redis().HashTableSet(ctx, key, m, expire)
	if err != nil {
		return err
	}
	keys := utilx.CacheKey(preKey, "keys_", strconv.Itoa(id))
	return repo.Redis().SAdd(ctx, keys, key)
}

func (repo {{.Svc}}Repository) CacheModelDefault(ctx context.Context, m interface{}, preKey string, id int) error {
	key := utilx.CacheKey(preKey, strconv.Itoa(id))
	return repo.Redis().HashTableSet(ctx, key, m, constant.CacheMinute1)
}

func (repo {{.Svc}}Repository) GetModelCache(ctx context.Context, m interface{}, preKey string, id int) error {
	key := utilx.CacheKey(preKey, strconv.Itoa(id))
	return repo.Redis().HashTableGet(ctx, key, m)
}

func (repo {{.Svc}}Repository) DelModelCache(ctx context.Context, preKey string, id int) error {
	key := utilx.CacheKey(preKey, strconv.Itoa(id))
	err := repo.Redis().Del(ctx, key)
	if err != nil {
		logger := log.GetLogger(ctx)
		logger.Warn("Del.error", zap.String("key", key), zap.Error(err))
	}
	return nil
}

func (repo {{.Svc}}Repository) DelModelCacheAll(ctx context.Context, preKey string, id int) error {
	key := utilx.CacheKey(preKey, "keys_", strconv.Itoa(id))
	err := repo.Redis().SetDelKeys(ctx, key)
	if err != nil {
		logger := log.GetLogger(ctx)
		logger.Warn("SetDelKeys.error", zap.String("key", key), zap.Error(err))
	}
	return nil
}

func (repo {{.Svc}}Repository) FetchList0ID(ctx context.Context, tx *gorm.DB, pageSize, pageNum, limit int, desc bool) ([]int, []*redis.Z, error) {
	logger := log.GetLogger(ctx)

	rows, err := tx.Rows()
	if err != nil {
		logger.Error("Data select failed", zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, status.ErrDBDeadlineExceeded
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return []int{}, []*redis.Z{}, nil
		}
		return nil, nil, status.ErrDBQuery
	}
	defer rows.Close()

	var ids []int
	var args []*redis.Z
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
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, nil, status.ErrDBDeadlineExceeded
			}
			return nil, nil, status.ErrDBUpdate
		}

		check := false
		if index >= start {
			check = true
		}
		if check {
			if len(ids) < pageSize {
				ids = append(ids, id)
			}
		}
		idx = index
		if desc {
			idx = max - index
		}
		args = append(args, &redis.Z{Score: idx, Member: id})
		index++
	}
	if err = rows.Err(); err != nil {
		logger.Error("Data err failed", zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, status.ErrDBDeadlineExceeded
		}
		return nil, nil, status.ErrDBUpdate
	}
	return ids, args, nil
}

func (repo {{.Svc}}Repository) FetchList1ID(ctx context.Context, tx *gorm.DB, pageSize int) ([]int, []*redis.Z, error) {
	logger := log.GetLogger(ctx)

	rows, err := tx.Rows()
	if err != nil {
		logger.Error("Data select failed", zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, status.ErrDBDeadlineExceeded
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return []int{}, []*redis.Z{}, nil
		}
		return nil, nil, status.ErrDBQuery
	}
	defer rows.Close()

	var ids []int
	var args []*redis.Z
	for rows.Next() {
		var (
			id int
		)
		err = rows.Scan(&id)
		if err != nil {
			logger.Error("Data scan failed", zap.Error(err))
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, nil, status.ErrDBDeadlineExceeded
			}
			return nil, nil, status.ErrDBUpdate
		}

		if len(ids) < pageSize {
			ids = append(ids, id)
		}

		args = append(args, &redis.Z{Score: float64(id), Member: id})
	}
	if err = rows.Err(); err != nil {
		logger.Error("Data err failed", zap.Error(err))
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, status.ErrDBDeadlineExceeded
		}
		return nil, nil, status.ErrDBUpdate
	}
	return ids, args, nil
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

	fileName := dir + "persistence.go"

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
