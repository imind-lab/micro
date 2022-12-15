/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package dao

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/imind-lab/micro/redis"
)

type Cache interface {
	Redis() redis.Redis
}

type cache struct {
	redisClient redis.Redis
}

func NewCache() Cache {
	var conf redis.RedisConfig
	if err := viper.UnmarshalKey("redis", &conf); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	var client redis.Redis
	if conf.Model == "node" {
		node := redis.NewRedisNode(conf)
		client = redis.NewRedis(node, conf.Timeout)
	} else {
		cluster := redis.NewRedisCluster(conf)
		client = redis.NewRedis(cluster, conf.Timeout)
	}

	cacheClient := &cache{}
	cacheClient.redisClient = client

	return cacheClient
}

func (c *cache) Redis() redis.Redis {
	return c.redisClient
}
