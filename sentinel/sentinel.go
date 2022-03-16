package sentinel

import (
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"go.uber.org/zap"
)

var opts = Options{0, "micro"}

type Sentinel struct {
	ds     datasource.DataSource
	logger *zap.Logger
}

func NewSentinel(logger *zap.Logger, opt ...Option) (*Sentinel, error) {
	for _, o := range opt {
		o(&opts)
	}

	return &Sentinel{
		logger: logger,
	}, nil
}

func (sent *Sentinel) Close() {
	if sent.ds != nil {
		sent.ds.Close()
	} else {
		sent.logger.Warn("Sentinel DataSource is nil while close")
	}
}

type Options struct {
	Type      int
	Namespace string
}

type Option func(*Options)

func Type(typ int) Option {
	return func(o *Options) {
		o.Type = typ
	}
}

func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}
