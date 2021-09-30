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

func (rep *dao) WriteDB(ctx context.Context) *gorm.DB {
	if rep.dbMock != nil {
		return rep.dbMock
	}
	return rep.MySQL.WriteDB(rep.dbName).WithContext(ctx)
}

func (rep *dao) ReadDB(ctx context.Context) *gorm.DB {
	if rep.dbMock != nil {
		return rep.dbMock
	}
	if rep.realTime {
		return rep.MySQL.WriteDB(rep.dbName).WithContext(ctx)
	}
	return rep.MySQL.ReadDB(rep.dbName).WithContext(ctx)
}

func (rep *dao) ExtraWriteDB(ctx context.Context, name string) *gorm.DB {
	if rep.dbMock != nil {
		return rep.dbMock
	}
	return rep.MySQL.WriteDB(name).WithContext(ctx)
}

func (rep *dao) ExtraReadDB(ctx context.Context, name string) *gorm.DB {
	if rep.dbMock != nil {
		return rep.dbMock
	}
	if rep.realTime {
		return rep.MySQL.WriteDB(name).WithContext(ctx)
	}
	return rep.MySQL.ReadDB(name).WithContext(ctx)
}

func (rep *dao) SetDBMock(db *gorm.DB) {
	rep.dbMock = db
}

func (repo *dao) SetRealTime(realTime bool) {
	repo.realTime = realTime
}

func (rep *dao) SetRedisMock(rdb *redis.Client) {
	rep.redisMock = rdb
}

func (rep *dao) Redis() *redis.Client {
	if rep.redisMock != nil {
		return rep.redisMock
	}
	return rep.Cache.Redis()
}
