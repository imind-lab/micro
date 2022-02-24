/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func SetSortedSet(ctx context.Context, cli *redis.ClusterClient, key string, args []*redis.Z, expire time.Duration) error {
	if len(args) == 0 {
		err := cli.Set(ctx, key, "", expire).Err()
		return err
	}
	err := cli.ZAdd(ctx, key, args...).Err()
	if err != nil {
		return err
	}
	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetNumber(ctx context.Context, cli *redis.ClusterClient, key string) (int64, error) {
	reply := cli.Get(ctx, key)
	if reply.Err() != nil {
		return 0, reply.Err()
	}
	cnt, err := reply.Int64()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func HGetAll(ctx context.Context, cli *redis.ClusterClient, key string, value interface{}) error {
	reply := cli.HGetAll(ctx, key)
	v, err := reply.Result()
	if err != nil {
		return err
	}
	if len(v) > 0 {
		err := reply.Scan(value)
		if err == nil {
			return nil
		}
	}
	return errors.New("Data does not exist")
}

func HGet(ctx context.Context, cli *redis.ClusterClient, key string, value interface{}, fields ...string) error {
	if len(fields) > 0 {
		reply := cli.HMGet(ctx, key, fields...)
		v, err := reply.Result()
		if err != nil {
			return err
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
			return err
		}
		if len(v) > 0 {
			err := reply.Scan(value)
			if err == nil {
				return nil
			}
		}
	}
	return errors.New("Data does not exist")
}

func ZRevRange(ctx context.Context, cli *redis.ClusterClient, key string, start, stop int64) ([]int, error) {
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			data := cli.ZRevRange(ctx, key, start, stop)
			if data.Err() == nil {
				var ids []int
				err := data.ScanSlice(&ids)
				if err == nil {
					return ids, nil
				}
			}
		case "none":
		default:
			return []int{}, nil
		}
		return nil, errors.New("Data does not exist")
	}
	return nil, err
}

func ZRevRangeWithCard(ctx context.Context, cli *redis.ClusterClient, key string, lastId, pageSize, pageNum int64) ([]int, int, error) {
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			reply := cli.ZCard(ctx, key)
			if reply.Err() == nil {
				if lastId > 0 {
					data := cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatInt(lastId-1, 10), Min: "-inf", Offset: 0, Count: pageSize})
					if data.Err() == nil {
						var ids []int
						err := data.ScanSlice(&ids)
						if err == nil {
							return ids, int(reply.Val()), nil
						}
					}
				} else {
					start := (pageNum - 1) * pageSize
					data := cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: "+inf", Min: "-inf", Offset: start, Count: pageSize})
					if data.Err() == nil {
						var ids []int
						err := data.ScanSlice(&ids)
						if err == nil {
							return ids, int(reply.Val()), nil
						}
					}
				}
			}
		case "none":
		default:
			return []int{}, 0, nil
		}
		return nil, 0, errors.New("Data does not exist")
	}
	return nil, 0, err
}

func SetSet(ctx context.Context, cli *redis.ClusterClient, key string, args []interface{}, expire time.Duration) error {
	if len(args) == 0 {
		err := cli.Set(ctx, key, "", expire).Err()
		return err

	}
	err := cli.SAdd(ctx, key, args...).Err()
	if err != nil {
		return err
	}
	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetHashTable(ctx context.Context, cli *redis.ClusterClient, key string, m interface{}, expire time.Duration) error {
	err := cli.HMSet(ctx, key, FlatStruct(m)).Err()
	if err != nil {
		return err
	}

	err = cli.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}

	return nil
}

func DelKeys(ctx context.Context, cli *redis.ClusterClient, key string) error {
	keys, err := cli.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	err = cli.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}
