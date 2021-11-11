package micro

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/imind-lab/micro/broker"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/credentials"
	"io"
)

type Options struct {
	Name   string
	Broker broker.Broker

	Context context.Context
	Tracer  opentracing.Tracer
	Closer  io.Closer

	Logger *zap.Logger

	ClientCred credentials.TransportCredentials
	ServerCred credentials.TransportCredentials

	BeforeRun  []func() error
	BeforeStop []func() error
	AfterRun   []func() error
	AfterStop  []func() error

	Signal bool
}

func newOptions(opts ...Option) Options {
	name := viper.GetString("service.name")
	project := viper.GetString("service.project")

	logPath := viper.GetString("log.path")
	logLevel := viper.GetInt("log.level")
	logAge := viper.GetInt("log.age")
	logSize := viper.GetInt("log.size")
	logBackup := viper.GetInt("log.backup")
	logCompress := viper.GetBool("log.compress")
	logFormat := viper.GetString("log.format")

	logger := log.NewLogger(logPath, zapcore.Level(logLevel), logSize, logBackup, logAge, logCompress, logFormat, zap.Fields(zap.String("project", project), zap.String("service", name)))
	ctx := ctxzap.ToContext(context.Background(), logger)

	traceName := viper.GetString("tracing.name.server")
	// 初始化调用链追踪
	tracer, closer, _ := tracing.InitTracer(traceName)

	opt := Options{
		Name:    name,
		Context: ctx,
		Logger:  logger,
		Tracer:  tracer,
		Closer:  closer,
		Signal:  true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Broker to be used for service
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

// ServerCred to be used for service
func ServerCred(cred credentials.TransportCredentials) Option {
	return func(o *Options) {
		o.ServerCred = cred
	}
}

// ClientCred to be used for service
func ClientCred(cred credentials.TransportCredentials) Option {
	return func(o *Options) {
		o.ClientCred = cred
	}
}

type Option func(*Options)

// Before and Afters

// BeforeRun run funcs before service starts
func BeforeRun(fn func() error) Option {
	return func(o *Options) {
		o.BeforeRun = append(o.BeforeRun, fn)
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterRun run funcs after service starts
func AfterRun(fn func() error) Option {
	return func(o *Options) {
		o.AfterRun = append(o.AfterRun, fn)
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
