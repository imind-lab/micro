/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package broker

import (
    "context"
    "crypto/tls"

    "github.com/spf13/viper"
)

type Options struct {
    Secure    bool
    TLSConfig *tls.Config

    ProducerAddr []string
    ConsumerAddr []string
    Topics       map[string]string
    GroupId      string

    Context context.Context
}

type Option func(*Options)

func ProducerAddr(addrs ...string) Option {
    return func(o *Options) {
        o.ProducerAddr = addrs
    }
}

func ConsumerAddr(addrs ...string) Option {
    return func(o *Options) {
        o.ConsumerAddr = addrs
    }
}

// Secure communication with the broker
func Secure(b bool) Option {
    return func(o *Options) {
        o.Secure = b
    }
}

// TLSConfig Specify TLS Config
func TLSConfig(t *tls.Config) Option {
    return func(o *Options) {
        o.TLSConfig = t
    }
}

func Context(ctx context.Context) Option {
    return func(o *Options) {
        o.Context = ctx
    }
}

func GroupId(group string) Option {
    return func(o *Options) {
        o.GroupId = group
    }
}

func NewOptions(name string) Options {
    name = "kafka." + name
    producerAddr := viper.GetStringSlice(name + ".producer")
    consumerAddr := viper.GetStringSlice(name + ".consumer")
    topics := viper.GetStringMapString(name + ".topic")
    opts := Options{
        ProducerAddr: producerAddr,
        ConsumerAddr: consumerAddr,
        Topics:       topics,
        Context:      context.Background(),
        GroupId:      viper.GetString(name + ".groupid"),
    }
    return opts
}
