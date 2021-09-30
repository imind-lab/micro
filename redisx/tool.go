/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package redisx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func SetSortedSet(ctx context.Context, cli *redis.Client, key string, args []*redis.Z, expire time.Duration) error {
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

func GetNumber(ctx context.Context, cli *redis.Client, key string) (int64, error) {
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

func HGetAll(ctx context.Context, cli *redis.Client, key string, value interface{}) error {
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

func HGet(ctx context.Context, cli *redis.Client, key string, value interface{}, fields ...string) error {
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

func ZRevRange(ctx context.Context, cli *redis.Client, key string, start, stop int64) ([]int32, error) {
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			data := cli.ZRevRange(ctx, key, start, stop)
			if data.Err() == nil {
				var ids []int32
				err := data.ScanSlice(&ids)
				if err == nil {
					return ids, nil
				}
			}
		case "none":
		default:
			return []int32{}, nil
		}
		return nil, errors.New("Data does not exist")
	}
	return nil, err
}

func ZRevRangeWithCard(ctx context.Context, cli *redis.Client, key string, lastid, pagesize, page int32) ([]int32, int, error) {
	rtype, err := cli.Type(ctx, key).Result()
	if err == nil {
		switch rtype {
		case "zset":
			reply := cli.ZCard(ctx, key)
			if reply.Err() == nil {
				if lastid > 0 {
					data := cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: fmt.Sprintf("%d", lastid-1), Min: "-inf", Offset: 0, Count: int64(pagesize)})
					if data.Err() == nil {
						var ids []int32
						err := data.ScanSlice(&ids)
						if err == nil {
							return ids, int(reply.Val()), nil
						}
					}
				} else {
					start := (page - 1) * pagesize
					data := cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: "+inf", Min: "-inf", Offset: int64(start), Count: int64(pagesize)})
					if data.Err() == nil {
						var ids []int32
						err := data.ScanSlice(&ids)
						if err == nil {
							return ids, int(reply.Val()), nil
						}
					}
				}
			}
		case "none":
		default:
			return []int32{}, 0, nil
		}
		return nil, 0, errors.New("Data does not exist")
	}
	return nil, 0, err
}

func SetSet(ctx context.Context, cli *redis.Client, key string, args []interface{}, expire time.Duration) error {
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

func SetHashTable(ctx context.Context, cli *redis.Client, key string, m interface{}, expire time.Duration) error {
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
