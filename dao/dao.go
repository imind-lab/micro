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
	WriteDB(ctx context.Context) *gorm.DB
	ReadDB(ctx context.Context) *gorm.DB

	ExtraWriteDB(ctx context.Context, name string) *gorm.DB
	ExtraReadDB(ctx context.Context, name string) *gorm.DB

	SetRealTime(realTime bool)
	SetDBMock(db *gorm.DB)

	Redis() *redis.Client

	SetRedisMock(rdb *redis.Client)
}

type dao struct {
	MySQL
	Cache

	dbName    string
	dbMock    *gorm.DB
	realTime  bool
	redisMock *redis.Client
}

func NewRepository(dbName string, realTime bool) Dao {
	rep := &dao{
		MySQL:    NewMySQL(),
		Cache:    NewCache(),
		dbName:   dbName,
		realTime: realTime,
	}
	return rep
}

func (d *dao) WriteDB(ctx context.Context) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	return d.MySQL.WriteDB(d.dbName).WithContext(ctx)
}

func (d *dao) ReadDB(ctx context.Context) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	if d.realTime {
		return d.MySQL.WriteDB(d.dbName).WithContext(ctx)
	}
	return d.MySQL.ReadDB(d.dbName).WithContext(ctx)
}

func (d *dao) ExtraWriteDB(ctx context.Context, name string) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	return d.MySQL.WriteDB(name).WithContext(ctx)
}

func (d *dao) ExtraReadDB(ctx context.Context, name string) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	if d.realTime {
		return d.MySQL.WriteDB(name).WithContext(ctx)
	}
	return d.MySQL.ReadDB(name).WithContext(ctx)
}

func (d *dao) SetDBMock(db *gorm.DB) {
	d.dbMock = db
}

func (d *dao) SetRealTime(realTime bool) {
	d.realTime = realTime
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
