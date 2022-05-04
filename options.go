package micro

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/credentials"

	"github.com/imind-lab/micro/broker"
	"github.com/imind-lab/micro/log"
	"github.com/imind-lab/micro/tracing"
)

type Options struct {
	Name   string
	Broker broker.Broker

	Context        context.Context
	TracerProvider *tracesdk.TracerProvider

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
	namespace := viper.GetString("service.namespace")

	logPath := viper.GetString("log.path")
	logLevel := viper.GetInt("log.level")
	logAge := viper.GetInt("log.age")
	logSize := viper.GetInt("log.size")
	logBackup := viper.GetInt("log.backup")
	logCompress := viper.GetBool("log.compress")
	logFormat := viper.GetString("log.format")

	logger := log.NewLogger(logPath, zapcore.Level(logLevel), logSize, logBackup, logAge, logCompress, logFormat, zap.Fields(zap.String("namespace", namespace), zap.String("service", name)))
	ctx := ctxzap.ToContext(context.Background(), logger)

	// 初始化调用链追踪
	provider, _ := tracing.InitTracer()

	opt := Options{
		Name:           name,
		Context:        ctx,
		Logger:         logger,
		TracerProvider: provider,
		Signal:         true,
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
