/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package dao

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Dao interface {
	DB() *gorm.DB
	ExtraDB(name string) *gorm.DB
	SetDBMock(db *gorm.DB)

	Redis() *redis.ClusterClient
	SetRedisMock(rdb *redis.ClusterClient)
}

type dao struct {
	Cache
	Database

	dbName string
	dbMock *gorm.DB

	redisMock *redis.ClusterClient
}

func NewDao(dbName string) Dao {
	rep := &dao{
		Cache:    NewCache(),
		Database: NewDatabase(),
		dbName:   dbName,
	}
	return rep
}

func (d *dao) DB() *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	return d.Database.DB(d.dbName)
}

func (d *dao) ExtraDB(name string) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	return d.Database.DB(name)
}

func (d *dao) SetDBMock(db *gorm.DB) {
	d.dbMock = db
}

func (d *dao) SetRedisMock(rdb *redis.ClusterClient) {
	d.redisMock = rdb
}

func (d *dao) Redis() *redis.ClusterClient {
	if d.redisMock != nil {
		return d.redisMock
	}
	return d.Cache.Redis()
}
