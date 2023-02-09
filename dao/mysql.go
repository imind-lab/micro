/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
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

type Max struct {
    Open int `yaml:"open"`
    Idle int `yaml:"idle"`
    Life int `yaml:"life"`
}

type Connection struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
    User string `yaml:"user"`
    Pass string `yaml:"pass"`
    Name string `yaml:"name"`
}

type Instance struct {
    Master  Connection `yaml:"master"`
    Replica Connection `yaml:"replica"`
}

const dsnFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&multiStatements=true&interpolateParams=true&parseTime=True&loc=Local"

// Database 定义Database接口方法
type Database interface {
    DB(name string) *gorm.DB
    Timeout() time.Duration
}

type MySQL struct {
    Once sync.Once
    DB   *gorm.DB
}

// NewDatabase 创建MySQL接口实例
func NewDatabase() Database {
    timeout := viper.GetDuration("db.timeout")
    dbClient := &database{
        timeout: timeout,
        dbs:     make(map[string]*MySQL),
    }
    return dbClient
}

type database struct {
    timeout time.Duration
    dbs     map[string]*MySQL
}

// 获取指定数据库的写连接
// @name 数据库别名
func (d *database) DB(name string) *gorm.DB {
    db, ok := d.dbs[name]
    if ok {
        db.Once.Do(func() {
            var err error
            db.DB, err = d.openDB(name)
            if err != nil {
                log.Fatalf("initDB error: %s, %v", name, err)
            }
        })
        return db.DB
    }

    db = &MySQL{}
    db.Once.Do(func() {
        var err error
        db.DB, err = d.openDB(name)
        if err != nil {
            log.Fatalf("initDB error: %s, %v", name, err)
        }
    })
    d.dbs[name] = db

    return db.DB
}

// 根据参数打开数据库连接
func (d *database) openDB(name string) (*gorm.DB, error) {
    var instance Instance
    if err := viper.UnmarshalKey("db."+name, &instance); err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }
    dsn := fmt.Sprintf(dsnFormat, instance.Master.User, instance.Master.Pass, instance.Master.Host, instance.Master.Port, instance.Master.Name)

    logLevel := viper.GetInt("db.logLevel")
    if logLevel < 1 || logLevel > 4 {
        logLevel = 1
    }
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
    })
    if err != nil {
        log.Fatal(fmt.Sprintf("can't open database: %s", dsn))
    } else {
        dsn = fmt.Sprintf(dsnFormat, instance.Replica.User, instance.Replica.Pass, instance.Replica.Host, instance.Replica.Port, instance.Replica.Name)

        var max Max
        if err := viper.UnmarshalKey("db.max", &max); err != nil {
            panic(fmt.Errorf("Fatal error config file: %s \n", err))
        }
        db.Use(
            dbresolver.Register(dbresolver.Config{
                Replicas: []gorm.Dialector{mysql.Open(dsn)},
                Policy:   dbresolver.RandomPolicy{},
            }).
                SetMaxOpenConns(max.Open).
                SetMaxIdleConns(max.Idle).
                SetConnMaxLifetime(time.Duration(max.Life) * time.Minute),
        )

        sqlDB, err := db.DB()
        if err == nil {
            sqlDB.SetMaxOpenConns(max.Open)
            sqlDB.SetMaxIdleConns(max.Idle)
            sqlDB.SetConnMaxLifetime(time.Duration(max.Life) * time.Minute)
            sqlDB.Ping()
        }
    }
    return db, err
}

func (d *database) Timeout() time.Duration {
    return d.timeout
}
