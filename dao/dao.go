/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Dao interface {
	DB(ctx context.Context) *gorm.DB
	ExtraDB(ctx context.Context, name string) *gorm.DB
	SetDBMock(db *gorm.DB)

	Redis() *redis.Client
	SetRedisMock(rdb *redis.Client)
}

type dao struct {
	Cache
	Database

	dbName string
	dbMock *gorm.DB

	redisMock *redis.Client
}

func NewDao(dbName string) Dao {
	rep := &dao{
		Cache:    NewCache(),
		Database: NewDatabase(),
		dbName:   dbName,
	}
	return rep
}

func (d *dao) DB(ctx context.Context) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	return d.Database.DB(d.dbName).WithContext(ctx)
}

func (d *dao) ExtraDB(ctx context.Context, name string) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	return d.Database.DB(name).WithContext(ctx)
}

func (d *dao) SetDBMock(db *gorm.DB) {
	d.dbMock = db
}

func (d *dao) SetRedisMock(rdb *redis.Client) {
	d.redisMock = rdb
}

func (d *dao) Redis() *redis.Client {
	if d.redisMock != nil {
		return d.redisMock
	}
	return d.Cache.Redis()
}
