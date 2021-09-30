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
	"math/rand"
	"sync"
	"time"

	"gorm.io/gorm/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbOnce   sync.Once
	dbClient MySQL
)

const dsnFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&multiStatements=true&interpolateParams=true&parseTime=True&loc=Local"

// 定义MySQL接口方法
type MySQL interface {
	WriteDB(name string) *gorm.DB
	ReadDB(name string) *gorm.DB
}

type RWDB struct {
	WOnce sync.Once
	Write *gorm.DB
	ROnce sync.Once
	Read  []*gorm.DB
}

// 创建MySQL接口实例
func NewMySQL() MySQL {
	dbOnce.Do(func() {
		dbClient = &database{
			dbs: make(map[string]*RWDB),
		}
	})
	return dbClient
}

type database struct {
	dbs map[string]*RWDB
}

// 获取指定数据库的写连接
// @name 数据库别名
func (d *database) WriteDB(name string) *gorm.DB {
	db, ok := d.dbs[name]
	if ok {
		db.WOnce.Do(func() {
			var err error
			db.Write, err = initDB(name)
			if err != nil {
				log.Fatalf("initDB error: %s, %v", name, err)
			}
		})
		return db.Write
	}

	db = &RWDB{}
	db.WOnce.Do(func() {
		var err error
		db.Write, err = initDB(name)
		if err != nil {
			log.Fatalf("initDB error: %s, %v", name, err)
		}
	})
	d.dbs[name] = db

	return db.Write
}

// 获取指定数据库的读连接（多读实例之间随机分配）
// @name 数据库别名
func (d *database) ReadDB(name string) *gorm.DB {
	db, ok := d.dbs[name]
	if ok {
		db.ROnce.Do(func() {
			var err error
			db.Read, err = initDBs(name)
			if err != nil {
				log.Fatalf("initDB error: %s, %v", name, err)
			}
		})
		return randDB(db.Read)
	}

	db = &RWDB{}
	db.ROnce.Do(func() {
		var err error
		db.Read, err = initDBs(name)
		if err != nil {
			log.Fatalf("initDB error: %s, %v", name, err)
		}
	})
	d.dbs[name] = db

	return randDB(db.Read)
}

//
func initDB(name string) (*gorm.DB, error) {
	host, port, user, pass, dbname, err := readConfig(name)
	if err != nil {
		return nil, err
	}
	return openDB(host, user, pass, dbname, port)
}

// 根据参数打开数据库连接
func openDB(host, user, pass, dbname string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf(dsnFormat, user, pass, host, port, dbname)

	logMode := viper.GetInt("db.logMode")
	if logMode < 1 || logMode > 4 {
		logMode = 1
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logMode)),
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("can't open database: %s", dsn))
	} else {
		dbConfig(db)
	}
	return db, err
}

func readConfig(name string) (string, int, string, string, string, error) {
	hostKey := "db." + name + ".write.host"
	if !viper.IsSet(hostKey) {
		return "", 0, "", "", "", fmt.Errorf("%s write configuration not exist", name)
	}
	host := viper.GetString(hostKey)
	portKey := "db." + name + ".write.port"
	port := viper.GetInt(portKey)
	userKey := "db." + name + ".write.user"
	user := viper.GetString(userKey)
	passKey := "db." + name + ".write.pass"
	pass := viper.GetString(passKey)
	dbKey := "db." + name + ".write.name"
	db := viper.GetString(dbKey)
	return host, port, user, pass, db, nil
}

type dbInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

func initDBs(name string) ([]*gorm.DB, error) {
	infos, err := readConfigs(name)
	if err != nil {
		return nil, err
	}
	dbs := make([]*gorm.DB, 0, len(infos))
	for _, info := range infos {
		db, err := openDB(info.Host, info.User, info.Pass, info.Name, info.Port)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, err
}

func readConfigs(name string) ([]dbInfo, error) {
	var infos []dbInfo
	dbsKey := "db." + name + ".read"
	if !viper.IsSet(dbsKey) {
		return nil, fmt.Errorf("%s read configuration not exist", name)
	}
	err := viper.UnmarshalKey(dbsKey, &infos)
	return infos, err
}

// 配置数据库连接池参数
func dbConfig(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(8)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.Ping()
	}
}

// 从一组连接中随机返回一个连接
func randDB(dbs []*gorm.DB) *gorm.DB {
	cnt := len(dbs)
	if cnt == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	idx := rand.Intn(cnt)
	return dbs[idx]
}
