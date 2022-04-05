/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package dao

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var (
	dbOnce   sync.Once
	dbClient Database
)

const dsnFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&multiStatements=true&interpolateParams=true&parseTime=True&loc=Local"

// Database 定义Database接口方法
type Database interface {
	DB(name string) *gorm.DB
}

type MySQL struct {
	Once sync.Once
	DB   *gorm.DB
}

// NewDatabase 创建MySQL接口实例
func NewDatabase() Database {
	dbOnce.Do(func() {
		dbClient = &database{
			dbs: make(map[string]*MySQL),
		}
	})
	return dbClient
}

type database struct {
	dbs map[string]*MySQL
}

// 获取指定数据库的写连接
// @name 数据库别名
func (d *database) DB(name string) *gorm.DB {
	db, ok := d.dbs[name]
	if ok {
		db.Once.Do(func() {
			var err error
			db.DB, err = openDB(name)
			if err != nil {
				log.Fatalf("initDB error: %s, %v", name, err)
			}
		})
		return db.DB
	}

	db = &MySQL{}
	db.Once.Do(func() {
		var err error
		db.DB, err = openDB(name)
		if err != nil {
			log.Fatalf("initDB error: %s, %v", name, err)
		}
	})
	d.dbs[name] = db

	return db.DB
}

// 根据参数打开数据库连接
func openDB(name string) (*gorm.DB, error) {
	host, port, user, pass, dbname, err := readConfig(name, true)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf(dsnFormat, user, pass, host, port, dbname)

	logMode := viper.GetInt("db.logLevel")
	if logMode < 1 || logMode > 4 {
		logMode = 1
	}
	tablePrefix := viper.GetString("db." + name + ".tablePrefix")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logMode)),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("can't open database: %s", dsn))
	} else {
		maxOpen := viper.GetInt("db.max.open")
		maxIdle := viper.GetInt("db.max.idle")
		maxLife := viper.GetInt("db.max.life")

		host, port, user, pass, dbname, err := readConfig(name, false)
		if err == nil {
			dsn = fmt.Sprintf(dsnFormat, user, pass, host, port, dbname)

			db.Use(
				dbresolver.Register(dbresolver.Config{
					Replicas: []gorm.Dialector{mysql.Open(dsn)},
					Policy:   dbresolver.RandomPolicy{},
				}).
					SetMaxOpenConns(maxOpen).
					SetMaxIdleConns(maxIdle).
					SetConnMaxLifetime(time.Duration(maxLife) * time.Minute),
			)
		}

		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxOpenConns(maxOpen)
			sqlDB.SetMaxIdleConns(maxIdle)
			sqlDB.SetConnMaxLifetime(time.Duration(maxLife) * time.Minute)
			sqlDB.Ping()
		}
	}
	return db, err
}

func readConfig(name string, master bool) (string, int, string, string, string, error) {
	typ := "replica"
	if master {
		typ = "master"
	}
	hostKey := "db." + name + "." + typ + ".host"
	if !viper.IsSet(hostKey) {
		return "", 0, "", "", "", fmt.Errorf("%s write configuration not exist", name)
	}
	host := viper.GetString(hostKey)
	portKey := "db." + name + "." + typ + ".port"
	port := viper.GetInt(portKey)
	userKey := "db." + name + "." + typ + ".user"
	user := viper.GetString(userKey)
	passKey := "db." + name + "." + typ + ".pass"
	pass := viper.GetString(passKey)
	dbKey := "db." + name + "." + typ + ".name"
	db := viper.GetString(dbKey)
	return host, port, user, pass, db, nil
}
