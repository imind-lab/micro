/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package grpc

import (
	"github.com/rs/zerolog"
)

type Options struct {
	Domain     string
	ServerCert string
	ServerKey  string
	Logger     *zerolog.Logger
}

func Domain(domain string) Option {
	return func(o *Options) {
		o.Domain = domain
	}
}

func ServerCert(cert string) Option {
	return func(o *Options) {
		o.ServerCert = cert
	}
}

func ServerKey(key string) Option {
	return func(o *Options) {
		o.ServerKey = key
	}
}

func Logger(logger *zerolog.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

type Option func(*Options)

func NewOptions() Options {
	opts := Options{
		Domain:     "*.imind.tech",
		ServerCert: "./conf/ssl/tls.crt",
		ServerKey:  "./conf/ssl/tls.key",
		Logger:     zerolog.DefaultContextLogger,
	}
	return opts
}
