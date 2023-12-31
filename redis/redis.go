package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/imind-lab/micro/v2/status"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisConfig struct {
	Model   string        `yaml:"model"`
	Timeout time.Duration `yaml:"timeout"`
	Addr    string        `yaml:"addr"`
	Pass    string        `yaml:"pass"`
	DB      int           `yaml:"db"`
}

type Redis interface {
	GetNumber(ctx context.Context, key string) (int64, error)

	HashTableGet(ctx context.Context, key string, value interface{}, fields ...string) error
	HashTableGetAll(ctx context.Context, key string, value interface{}) error
	HashTableSet(ctx context.Context, key string, p interface{}, expire time.Duration) error

	SetSet(ctx context.Context, key string, args []interface{}, expire time.Duration) error
	SetDelKeys(ctx context.Context, key string) error

	SortedSetRange(ctx context.Context, key string, pageSize, pageNum int64, isDesc bool, out interface{}) (int, error)
	SortedSetRangeByScore(ctx context.Context, key string, pageSize, lastId int64, isDesc bool, out interface{}) (int, error)
	SortedSetSet(ctx context.Context, key string, args []redis.Z, expire time.Duration) error

	Del(ctx context.Context, keys ...string) error
	Get(ctx context.Context, key string) *redis.StringCmd
	SAdd(ctx context.Context, key string, members ...interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	ZRem(ctx context.Context, key string, members ...interface{}) error
}

type cmdable struct {
	redis.Cmdable
	timeout time.Duration
}

func NewRedis(cmd redis.Cmdable, timeout time.Duration) Redis {
	return cmdable{Cmdable: cmd, timeout: timeout}
}

func (cmd cmdable) GetNumber(ctx context.Context, key string) (int64, error) {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	reply := cmd.Get(ctx, key)
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

func (cmd cmdable) HashTableGet(ctx context.Context, key string, value interface{}, fields ...string) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	if len(fields) > 0 {
		return cmd.hashTableMGet(ctx, key, value, fields...)
	}
	return cmd.hashTableGetAll(ctx, key, value)
}

func (cmd cmdable) hashTableMGet(ctx context.Context, key string, value interface{}, fields ...string) error {
	reply := cmd.HMGet(ctx, key, fields...)
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

func (cmd cmdable) hashTableGetAll(ctx context.Context, key string, value interface{}) error {
	reply := cmd.HGetAll(ctx, key)
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

func (cmd cmdable) HashTableGetAll(ctx context.Context, key string, value interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	reply := cmd.HGetAll(ctx, key)
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

func (cmd cmdable) HashTableSet(ctx context.Context, key string, p interface{}, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	var m map[string]interface{}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		TagName:  "redis",
		Result:   &m,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(p); err != nil {
		return err
	}
	err = cmd.HSet(ctx, key, m).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}

	err = cmd.Expire(ctx, key, expire).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}

	return nil
}

func (cmd cmdable) SetSet(ctx context.Context, key string, args []interface{}, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	if len(args) == 0 {
		return cmd.Set(ctx, key, "", expire)
	}
	err := cmd.SAdd(ctx, key, args...)
	if err != nil {
		return err
	}
	err = cmd.Expire(ctx, key, expire).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	return nil
}
func (cmd cmdable) SetDelKeys(ctx context.Context, key string) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	keys, err := cmd.SMembers(ctx, key).Result()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	return cmd.Del(ctx, keys...)
}

func (cmd cmdable) SortedSetRange(ctx context.Context, key string, pageSize, pageNum int64, isDesc bool, out interface{}) (int, error) {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	rtype, err := cmd.Type(ctx, key).Result()
	if err != nil {
		return SortedSetError(ctx, err)
	}
	switch rtype {
	case "zset":
		var cnt int64
		cnt, err = cmd.ZCard(ctx, key).Result()
		if err != nil {
			return SortedSetError(ctx, err)
		}
		start := (pageNum - 1) * pageSize
		// 是offset的索引而不是获取的长度，所以应该减一
		stop := start + pageSize - 1
		var data *redis.StringSliceCmd
		if isDesc {
			data = cmd.ZRevRange(ctx, key, start, stop)
		} else {
			data = cmd.ZRange(ctx, key, start, stop)
		}
		err = data.Err()
		if err != nil {
			return SortedSetError(ctx, err)
		}
		err = data.ScanSlice(out)
		if err != nil {
			return 0, err
		}
		return int(cnt), nil
	case "none":
	default:
		return 0, nil
	}
	return 0, status.ErrRedisDataNotExist
}

func (cmd cmdable) SortedSetRangeByScore(ctx context.Context, key string, pageSize, lastId int64, isDesc bool, out interface{}) (int, error) {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	rtype, err := cmd.Type(ctx, key).Result()
	if err != nil {
		return SortedSetError(ctx, err)
	}
	switch rtype {
	case "zset":
		var cnt int64
		cnt, err = cmd.ZCard(ctx, key).Result()
		if err != nil {
			return SortedSetError(ctx, err)
		}
		var data *redis.StringSliceCmd
		if isDesc {
			data = cmd.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatInt(lastId-1, 10), Min: "-inf", Offset: 0, Count: pageSize})
		} else {
			data = cmd.ZRangeByScore(ctx, key, &redis.ZRangeBy{Max: "+inf", Min: strconv.FormatInt(lastId+1, 10), Offset: 0, Count: pageSize})
		}
		err = data.Err()
		if err != nil {
			return SortedSetError(ctx, err)
		}
		err = data.ScanSlice(out)
		if err != nil {
			return 0, err
		}
		return int(cnt), nil
	case "none":
	default:
		return 0, nil
	}
	return 0, status.ErrRedisDataNotExist

}

func (cmd cmdable) SortedSetSet(ctx context.Context, key string, args []redis.Z, expire time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	if len(args) == 0 {
		return cmd.Set(ctx, key, "", expire)
	}
	err := cmd.ZAdd(ctx, key, args...).Err()
	if err != nil {
		return CheckDeadline(ctx, err)
	}
	err = cmd.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cmd cmdable) Del(ctx context.Context, keys ...string) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	err := cmd.Cmdable.Del(ctx, keys...).Err()
	return CheckDeadline(ctx, err)
}

func (cmd cmdable) Get(ctx context.Context, key string) *redis.StringCmd {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	return cmd.Cmdable.Get(ctx, key)
}

func (cmd cmdable) SAdd(ctx context.Context, key string, members ...interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	err := cmd.Cmdable.SAdd(ctx, key, members...).Err()
	return CheckDeadline(ctx, err)
}

func (cmd cmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	err := cmd.Cmdable.Set(ctx, key, value, expiration).Err()
	return CheckDeadline(ctx, err)
}

func (cmd cmdable) ZRem(ctx context.Context, key string, members ...interface{}) error {
	ctx, _ = context.WithTimeout(ctx, cmd.timeout)

	err := cmd.Cmdable.ZRem(ctx, key, members...).Err()
	return CheckDeadline(ctx, err)
}

func CheckDeadline(ctx context.Context, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger := log.Ctx(ctx)
		logger.Error().Err(err).Msg("RedisDeadlineExceeded")
		return status.ErrRedisDeadlineExceeded
	}
	return err
}

func SortedSetError(ctx context.Context, err error) (int, error) {
	if errors.Is(err, context.DeadlineExceeded) {
		logger := log.Ctx(ctx)
		logger.Error().Err(err).Msg("RedisDeadlineExceeded")
		return 0, status.ErrRedisDeadlineExceeded
	}
	return 0, err
}
