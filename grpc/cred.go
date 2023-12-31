/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package grpc

import (
	"crypto/tls"
	"os"

	"golang.org/x/net/http2"
	"google.golang.org/grpc/credentials"
)

type GrpcCred interface {
	ServerCred() credentials.TransportCredentials
	ClientCred() credentials.TransportCredentials
	GetTLSConfig() *tls.Config
	Options() Options
}

type cred struct {
	stc  credentials.TransportCredentials
	ctc  credentials.TransportCredentials
	cnf  *tls.Config
	opts Options
}

func NewGrpcCred(opt ...Option) GrpcCred {
	opts := NewOptions()
	for _, o := range opt {
		o(&opts)
	}

	return &cred{
		opts: opts,
	}
}

func (c *cred) ServerCred() credentials.TransportCredentials {
	if c.stc == nil {
		stc, err := credentials.NewServerTLSFromFile(c.opts.ServerCert, c.opts.ServerKey)
		if err != nil {
			c.opts.Logger.Error().Err(err).Msg("NewServerTLSFromFile error")
		}
		c.stc = stc
	}
	return c.stc
}

func (c *cred) ClientCred() credentials.TransportCredentials {
	if c.ctc == nil {
		ctc, err := credentials.NewClientTLSFromFile(c.opts.ServerCert, c.opts.Domain)
		if err != nil {
			c.opts.Logger.Error().Err(err).Msg("NewClientTLSFromFile error")
		}
		c.ctc = ctc
	}
	return c.ctc
}

func (c *cred) GetTLSConfig() *tls.Config {
	if c.cnf == nil {
		var certKeyPair *tls.Certificate
		cert, err := os.ReadFile(c.opts.ServerCert)
		if err != nil {
			c.opts.Logger.Error().Err(err).Msg("ReadFile error")
		}
		key, err := os.ReadFile(c.opts.ServerKey)
		if err != nil {
			c.opts.Logger.Error().Err(err).Msg("ReadFile error")
		}

		pair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			c.opts.Logger.Error().Err(err).Msg("X509KeyPair error")
		}

		certKeyPair = &pair

		c.cnf = &tls.Config{
			Certificates: []tls.Certificate{*certKeyPair},
			NextProtos:   []string{http2.NextProtoTLS},
		}
	}
	return c.cnf
}

func (c *cred) Options() Options {
	return c.opts
}
