/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package dao

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Dao interface {
	DB(ctx context.Context) *gorm.DB
	ExtraDB(ctx context.Context, name string) *gorm.DB
	SetDBMock(db *gorm.DB)

	Redis() Redis
	//SetRedisMock(rdb *redis.ClusterClient)
}

type dao struct {
	Cache
	Database

	dbName string
	dbMock *gorm.DB

	//redisMock *redis.ClusterClient

	timeout time.Duration
}

func NewDao(dbName string) Dao {
	timeout := viper.GetDuration("db.timeout")
	rep := &dao{
		Cache:    NewCache(),
		Database: NewDatabase(),
		dbName:   dbName,
		timeout:  timeout,
	}
	return rep
}

func (d *dao) DB(ctx context.Context) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	ctx, _ = context.WithTimeout(ctx, d.timeout)
	return d.Database.DB(d.dbName).WithContext(ctx)
}

func (d *dao) ExtraDB(ctx context.Context, name string) *gorm.DB {
	if d.dbMock != nil {
		return d.dbMock
	}
	ctx, _ = context.WithTimeout(ctx, d.timeout)
	return d.Database.DB(name).WithContext(ctx)
}

func (d *dao) SetDBMock(db *gorm.DB) {
	d.dbMock = db
}

//func (d *dao) SetRedisMock(rdb *redis.ClusterClient) {
//	d.redisMock = rdb
//}

func (d *dao) Redis() Redis {
	//if d.redisMock != nil {
	//	return d.redisMock
	//}
	return d.Cache.Redis()
}
