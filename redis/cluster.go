package redis

import (
    "context"
    "net"
    "time"

    "github.com/redis/go-redis/v9"
)

func NewRedisCluster(conf RedisConfig) *redis.ClusterClient {
    return clusterClient(conf.Timeout, conf.Addr)
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
        DialTimeout: 8 * time.Second, //连接建立超时时间，默认5秒。
        ReadTimeout: timeout,         //读超时，默认3秒， -1表示取消读超时

        //闲置连接检查包括IdleTimeout，MaxConnAge
        ConnMaxIdleTime: 5 * time.Minute, //闲置超时，默认5分钟，-1表示取消闲置超时检查
        ConnMaxLifetime: 0 * time.Second, //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

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
