package micro

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

type Handler func(next http.Handler) http.Handler

type Options struct {
	Name string

	Context context.Context

	Logger *zap.Logger

	ClientCred credentials.TransportCredentials
	ServerCred credentials.TransportCredentials

	BeforeRun  []func() error
	BeforeStop []func() error
	AfterRun   []func() error
	AfterStop  []func() error

	Handlers []Handler

	Signal bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
		Signal:  true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Name to be used for service
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Context to be used for service
func Context(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

// Logger to be used for service
func Logger(logger *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
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

// HttpHandler add handler for Http
func HttpHandler(handlers ...Handler) Option {
	return func(o *Options) {
		o.Handlers = append(o.Handlers, handlers...)
	}
}
