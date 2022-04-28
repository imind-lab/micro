/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/imind-lab/micro/log"
	redisx "github.com/imind-lab/micro/redis"
	"github.com/imind-lab/micro/status"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	cacheOnce   sync.Once
	cacheClient *cache
)

type Cache interface {
	Redis() Redis
}

type cache struct {
	redisClient Redis
}

func NewCache() Cache {
	cacheOnce.Do(func() {
		var client Redis
		model := viper.GetString("redis.model")
		if model == "node" {
			client = NewRedisNode()
		} else {
			client = NewRedisCluster()
		}

		cacheClient = &cache{}
		cacheClient.redisClient = client

	})
	return cacheClient
}

func (c *cache) Redis() Redis {
	return c.redisClient
}

type Redis interface {
	GetNumber(ctx context.Context, key string) (int64, error)

	HashTableGet(ctx context.Context, key string, value interface{}, fields ...string) error
	HashTableGetAll(ctx context.Context, key string, value interface{}) error
	HashTableSet(ctx context.Context, key string, m interface{}, expire time.Duration) error

	SetSet(ctx context.Context, key string, args []interface{}, expire time.Duration) error
	SetDelKeys(ctx context.Context, key string) error

	SortedSetRange(ctx context.Context, key string, pageSize, pageNum int64, desc bool) ([]int, int, error)
	SortedSetRangeByScore(ctx context.Context, key string, pageSize, lastId int64, desc bool) ([]int, int, error)
	SortedSetSet(ctx context.Context, key string, args []*redis.Z, expire time.Duration) error

	Del(ctx context.Context, keys ...string) error
	//Get(ctx context.Context, key string) *redis.StringCmd
	SAdd(ctx context.Context, key string, members ...interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	ZRem(ctx context.Context, key string, members ...interface{}) error
}

type redisNode struct {
	*redis.Client
	timeout time.Duration
}

func NewRedisNode() redisNode {
	addr := viper.GetString("redis.addr")
	pass := viper.GetString("redis.pass")
	db := viper.GetInt("redis.db")
	timeout := viper.GetDuration("redis.timeout")
	rdb := redisClient(addr, pass, db, timeout)
	return redisNode{Client: rdb, timeout: timeout * time.Second}
}

func (cli redisNode) GetNumber(ctx context.Context, key string) (int64, error) {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	reply := cli.Get(ctx, key)
	if err := reply.Err(); err != nil {
		err = CheckDeadline(ctx, err)
		return 0, err
	}
	cnt, err := reply.Int64()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (cli redisNode) HashTableGet(ctx context.Context, key string, value interface{}, fields ...string) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	if len(fields) > 0 {
		reply := cli.HMGet(ctx, key, fields...)
		v, err := reply.Result()
		if err != nil {
			return CheckDeadline(ctx, err)
		}
		if len(v) > 0 {
			err := reply.Scan(value)
			if err == nil {
				return nil
			}
		}
	} else {
		reply := cli.HGetAll(ctx, key)
		v, err := reply.Result()
		if err != nil {
			return CheckDeadline(ctx, err)
		}
		if len(v) > 0 {
			err := reply.Scan(value)
			if err == nil {
				return nil
			}
		}
	}
	return status.ErrRedisDataNotExist
}

func (cli redisNode) HashTableGetAll(ctx context.Context, key string, value interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	reply := cli.HGetAll(ctx, key)
	v, err := reply.Result()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	if len(v) > 0 {
		err := reply.Scan(value)
		if err == nil {
			return nil
		}
	}
	return status.ErrRedisDataNotExist
}

func (cli redisNode) HashTableSet(ctx context.Context, key string, m interface{}, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.HMSet(ctx, key, redisx.FlatStruct(m)).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}

	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}

	return nil
}

func (cli redisNode) SetSet(ctx context.Context, key string, args []interface{}, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	if len(args) == 0 {
		return cli.Set(ctx, key, "", expire)
	}
	err := cli.SAdd(ctx, key, args...)
	if err != nil {
		return err
	}
	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	return nil
}
func (cli redisNode) SetDelKeys(ctx context.Context, key string) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	keys, err := cli.SMembers(ctx, key).Result()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	return cli.Del(ctx, keys...)
}

func (cli redisNode) SortedSetRange(ctx context.Context, key string, pageSize, pageNum int64, desc bool) ([]int, int, error) {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	logger := log.GetLogger(ctx)
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			reply := cli.ZCard(ctx, key)
			if reply.Err() == nil {
				start := (pageNum - 1) * pageSize
				stop := start + pageSize
				var data *redis.StringSliceCmd
				if desc {
					data = cli.ZRevRange(ctx, key, start, stop)
				} else {
					data = cli.ZRange(ctx, key, start, stop)
				}
				if data.Err() == nil {
					var ids []int
					err := data.ScanSlice(&ids)
					if err == nil {
						return ids, int(reply.Val()), nil
					}
				}
			} else if errors.Is(reply.Err(), context.DeadlineExceeded) {
				logger.Error("RedisDeadlineExceeded", zap.Error(err))
				return nil, 0, status.ErrRedisDeadlineExceeded
			}
		case "none":
		default:
			return []int{}, 0, nil
		}
		return nil, 0, status.ErrRedisDataNotExist
	} else if errors.Is(err, context.DeadlineExceeded) {
		logger.Error("RedisDeadlineExceeded", zap.Error(err))
		return nil, 0, status.ErrRedisDeadlineExceeded
	}
	return nil, 0, err
}

func (cli redisNode) SortedSetRangeByScore(ctx context.Context, key string, pageSize, lastId int64, desc bool) ([]int, int, error) {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	logger := log.GetLogger(ctx)
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			reply := cli.ZCard(ctx, key)
			if reply.Err() == nil {
				var data *redis.StringSliceCmd
				if desc {
					data = cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatInt(lastId-1, 10), Min: "-inf", Offset: 0, Count: pageSize})
				} else {
					data = cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{Max: "+inf", Min: strconv.FormatInt(lastId+1, 10), Offset: 0, Count: pageSize})
				}
				if data.Err() == nil {
					var ids []int
					err := data.ScanSlice(&ids)
					if err == nil {
						return ids, int(reply.Val()), nil
					}
				}
			} else if errors.Is(reply.Err(), context.DeadlineExceeded) {
				logger.Error("RedisDeadlineExceeded", zap.Error(err))
				return nil, 0, status.ErrRedisDeadlineExceeded
			}
		case "none":
		default:
			return []int{}, 0, nil
		}
		return nil, 0, status.ErrRedisDataNotExist
	} else if errors.Is(err, context.DeadlineExceeded) {
		logger.Error("RedisDeadlineExceeded", zap.Error(err))
		return nil, 0, status.ErrRedisDeadlineExceeded
	}
	return nil, 0, err
}

func (cli redisNode) SortedSetSet(ctx context.Context, key string, args []*redis.Z, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	if len(args) == 0 {
		return cli.Set(ctx, key, "", expire)
	}
	err := cli.ZAdd(ctx, key, args...).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cli redisNode) Del(ctx context.Context, keys ...string) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.Client.Del(ctx, keys...).Err()
	return CheckDeadline(ctx, err)
}

//func (cli redisNode) Get(ctx context.Context, key string) *redis.StringCmd {
//	ctx, _ = context.WithTimeout(ctx, cli.timeout)
//
//	return cli.Client.Get(ctx, key)
//}

func (cli redisNode) SAdd(ctx context.Context, key string, members ...interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.Client.SAdd(ctx, key, members...).Err()
	return CheckDeadline(ctx, err)
}

func (cli redisNode) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.Client.Set(ctx, key, value, expiration).Err()
	return CheckDeadline(ctx, err)
}

func (cli redisNode) ZRem(ctx context.Context, key string, members ...interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.Client.ZRem(ctx, key, members...).Err()
	return CheckDeadline(ctx, err)
}

func redisClient(addr, pass string, db int, timeout time.Duration) *redis.Client {
	fmt.Println(addr, pass, db)
	return redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp", //网络类型，tcp or unix，默认tcp
		Addr:     addr,  //主机名+冒号+端口，默认localhost:6379
		Password: pass,  //密码
		DB:       db,    // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout: 5 * time.Second,       //连接建立超时时间，默认5秒。
		ReadTimeout: timeout * time.Second, //读超时，默认3秒， -1表示取消读超时

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		Dialer: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial(network, addr)
		},
		//钩子函数
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Printf("conn=%v\n", cn)
			return nil
		},
	})
}

type redisCluster struct {
	*redis.ClusterClient
	timeout time.Duration
}

func NewRedisCluster() redisCluster {
	addr := viper.GetString("redis.addr")
	timeout := viper.GetDuration("redis.timeout")
	rdb := clusterClient(timeout, addr)
	return redisCluster{rdb, timeout * time.Second}
}

func (cli redisCluster) GetNumber(ctx context.Context, key string) (int64, error) {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	reply := cli.Get(ctx, key)
	if err := reply.Err(); err != nil {
		err = CheckDeadline(ctx, err)
		return 0, err
	}
	cnt, err := reply.Int64()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (cli redisCluster) HashTableGet(ctx context.Context, key string, value interface{}, fields ...string) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	if len(fields) > 0 {
		reply := cli.HMGet(ctx, key, fields...)
		v, err := reply.Result()
		if err != nil {
			return CheckDeadline(ctx, err)
		}
		if len(v) > 0 {
			err := reply.Scan(value)
			if err == nil {
				return nil
			}
		}
	} else {
		reply := cli.HGetAll(ctx, key)
		v, err := reply.Result()
		if err != nil {
			return CheckDeadline(ctx, err)
		}
		if len(v) > 0 {
			err := reply.Scan(value)
			if err == nil {
				return nil
			}
		}
	}
	return status.ErrRedisDataNotExist
}

func (cli redisCluster) HashTableGetAll(ctx context.Context, key string, value interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	reply := cli.HGetAll(ctx, key)
	v, err := reply.Result()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	if len(v) > 0 {
		err := reply.Scan(value)
		if err == nil {
			return nil
		}
	}
	return status.ErrRedisDataNotExist
}

func (cli redisCluster) HashTableSet(ctx context.Context, key string, m interface{}, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.HMSet(ctx, key, redisx.FlatStruct(m)).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}

	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}

	return nil
}

func (cli redisCluster) SetSet(ctx context.Context, key string, args []interface{}, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	if len(args) == 0 {
		return cli.Set(ctx, key, "", expire)
	}
	err := cli.SAdd(ctx, key, args...)
	if err != nil {
		return err
	}
	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	return nil
}
func (cli redisCluster) SetDelKeys(ctx context.Context, key string) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	keys, err := cli.SMembers(ctx, key).Result()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	return cli.Del(ctx, keys...)
}

func (cli redisCluster) SortedSetRange(ctx context.Context, key string, pageSize, pageNum int64, desc bool) ([]int, int, error) {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	logger := log.GetLogger(ctx)
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			reply := cli.ZCard(ctx, key)
			if reply.Err() == nil {
				start := (pageNum - 1) * pageSize
				stop := start + pageSize
				var data *redis.StringSliceCmd
				if desc {
					data = cli.ZRevRange(ctx, key, start, stop)
				} else {
					data = cli.ZRange(ctx, key, start, stop)
				}
				if data.Err() == nil {
					var ids []int
					err := data.ScanSlice(&ids)
					if err == nil {
						return ids, int(reply.Val()), nil
					}
				}
			} else if errors.Is(reply.Err(), context.DeadlineExceeded) {
				logger.Error("RedisDeadlineExceeded", zap.Error(err))
				return nil, 0, status.ErrRedisDeadlineExceeded
			}
		case "none":
		default:
			return []int{}, 0, nil
		}
		return nil, 0, status.ErrRedisDataNotExist
	} else if errors.Is(err, context.DeadlineExceeded) {
		logger.Error("RedisDeadlineExceeded", zap.Error(err))
		return nil, 0, status.ErrRedisDeadlineExceeded
	}
	return nil, 0, err
}

func (cli redisCluster) SortedSetRangeByScore(ctx context.Context, key string, pageSize, lastId int64, desc bool) ([]int, int, error) {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	logger := log.GetLogger(ctx)
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			reply := cli.ZCard(ctx, key)
			if reply.Err() == nil {
				var data *redis.StringSliceCmd
				if desc {
					data = cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatInt(lastId-1, 10), Min: "-inf", Offset: 0, Count: pageSize})
				} else {
					data = cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{Max: "+inf", Min: strconv.FormatInt(lastId+1, 10), Offset: 0, Count: pageSize})
				}
				if data.Err() == nil {
					var ids []int
					err := data.ScanSlice(&ids)
					if err == nil {
						return ids, int(reply.Val()), nil
					}
				}
			} else if errors.Is(reply.Err(), context.DeadlineExceeded) {
				logger.Error("RedisDeadlineExceeded", zap.Error(err))
				return nil, 0, status.ErrRedisDeadlineExceeded
			}
		case "none":
		default:
			return []int{}, 0, nil
		}
		return nil, 0, status.ErrRedisDataNotExist
	} else if errors.Is(err, context.DeadlineExceeded) {
		logger.Error("RedisDeadlineExceeded", zap.Error(err))
		return nil, 0, status.ErrRedisDeadlineExceeded
	}
	return nil, 0, err
}

func (cli redisCluster) SortedSetSet(ctx context.Context, key string, args []*redis.Z, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	if len(args) == 0 {
		return cli.Set(ctx, key, "", expire)
	}
	err := cli.ZAdd(ctx, key, args...).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cli redisCluster) Del(ctx context.Context, keys ...string) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.ClusterClient.Del(ctx, keys...).Err()
	return CheckDeadline(ctx, err)
}

//func (cli redisNode) Get(ctx context.Context, key string) *redis.StringCmd {
//	ctx, _ = context.WithTimeout(ctx, cli.timeout)
//
//	return cli.ClusterClient.Get(ctx, key)
//}

func (cli redisCluster) SAdd(ctx context.Context, key string, members ...interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.ClusterClient.SAdd(ctx, key, members...).Err()
	return CheckDeadline(ctx, err)
}

func (cli redisCluster) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.ClusterClient.Set(ctx, key, value, expiration).Err()
	return CheckDeadline(ctx, err)
}

func (cli redisCluster) ZRem(ctx context.Context, key string, members ...interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cli.timeout)

	err := cli.ClusterClient.ZRem(ctx, key, members...).Err()
	return CheckDeadline(ctx, err)
}

func clusterClient(timeout time.Duration, addrs ...string) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          addrs,
		MaxRedirects:   0,
		ReadOnly:       false,
		RouteByLatency: false,
		RouteRandomly:  false,
		PoolFIFO:       false,

		//连接池容量及闲置连接数量
		PoolSize:     100, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,  //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout: 8 * time.Second,       //连接建立超时时间，默认5秒。
		ReadTimeout: timeout * time.Second, //读超时，默认3秒， -1表示取消读超时

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		Dialer: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			netDialer := &net.Dialer{
				Timeout:   8 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial(network, addr)
		},
		//钩子函数
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			//fmt.Printf("conn=%v\n", cn)
			return nil
		},
	})
}

func CheckDeadline(ctx context.Context, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger := log.GetLogger(ctx)
		logger.Error("RedisDeadlineExceeded", zap.Error(err))
		return status.ErrRedisDeadlineExceeded
	}
	return err
}
