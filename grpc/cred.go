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

    "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
    "go.uber.org/zap"
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
            ctxzap.Error(c.opts.Context, "NewServerTLSFromFile error", zap.Error(err))
        }
        c.stc = stc
    }
    return c.stc
}

func (c *cred) ClientCred() credentials.TransportCredentials {
    if c.ctc == nil {
        ctc, err := credentials.NewClientTLSFromFile(c.opts.ServerCert, c.opts.Domain)
        if err != nil {
            ctxzap.Error(c.opts.Context, "NewClientTLSFromFile error", zap.Error(err))
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
            ctxzap.Error(c.opts.Context, "ReadFile error", zap.Error(err))
        }
        key, err := os.ReadFile(c.opts.ServerKey)
        if err != nil {
            ctxzap.Error(c.opts.Context, "ReadFile error", zap.Error(err))
        }

        pair, err := tls.X509KeyPair(cert, key)
        if err != nil {
            ctxzap.Error(c.opts.Context, "X509KeyPair error", zap.Error(err))
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
