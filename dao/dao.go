/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package dao

import (
    "context"

    "gorm.io/gorm"

    "github.com/imind-lab/micro/v2/redis"
)

type Dao interface {
    DB(ctx context.Context) *gorm.DB
    ExtraDB(ctx context.Context, name string) *gorm.DB

    Redis() redis.Redis
}

type dao struct {
    Cache
    Database
}

func NewDao(cache Cache, db Database) Dao {
    rep := &dao{
        Cache:    cache,
        Database: db,
    }
    return rep
}

func (d *dao) DB(ctx context.Context) *gorm.DB {
    ctx, _ = context.WithTimeout(ctx, d.Database.Timeout())
    return d.Database.DB("default").WithContext(ctx)
}

func (d *dao) ExtraDB(ctx context.Context, name string) *gorm.DB {
    ctx, _ = context.WithTimeout(ctx, d.Database.Timeout())
    return d.Database.DB(name).WithContext(ctx)
}

func (d *dao) Redis() redis.Redis {
    return d.Cache.Redis()
}
