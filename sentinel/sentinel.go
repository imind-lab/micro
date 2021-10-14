package sentinel

import (
	"go.uber.org/zap"

	"github.com/alibaba/sentinel-golang/pkg/datasource/k8s"
)

var opts = Options{0, "micro"}

type Sentinel struct {
	ds     *k8s.DataSource
	logger *zap.Logger
}

func NewSentinel(logger *zap.Logger, opt ...Option) (*Sentinel, error) {
	for _, o := range opt {
		o(&opts)
	}
	ds, err := k8s.NewDataSource(opts.Namespace)
	if err != nil {
		logger.Warn("k8s.NewDataSource error", zap.Error(err))
		return nil, err
	}

	switch opts.Type {
	case 1:
		err = ds.RegisterController(k8s.CircuitBreakerRulesCRD, "sentinel-circuitbreaker-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController CircuitBreakerRulesCRD error", zap.Error(err))
			return nil, err
		}
	case 2:
		err = ds.RegisterController(k8s.FlowRulesCRD, "sentinel-flow-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController FlowRulesCRD error", zap.Error(err))
			return nil, err
		}
	case 3:
		err = ds.RegisterController(k8s.CircuitBreakerRulesCRD, "sentinel-circuitbreaker-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController CircuitBreakerRulesCRD error", zap.Error(err))
			return nil, err
		}
		err = ds.RegisterController(k8s.FlowRulesCRD, "sentinel-flow-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController FlowRulesCRD error", zap.Error(err))
			return nil, err
		}
		err = ds.RegisterController(k8s.IsolationRulesCRD, "sentinel-circuitbreaker-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController CircuitBreakerRulesCRD error", zap.Error(err))
			return nil, err
		}
		err = ds.RegisterController(k8s.HotspotRulesCRD, "sentinel-flow-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController FlowRulesCRD error", zap.Error(err))
			return nil, err
		}
		err = ds.RegisterController(k8s.SystemRulesCRD, "sentinel-flow-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController FlowRulesCRD error", zap.Error(err))
			return nil, err
		}
	default:
		err = ds.RegisterController(k8s.CircuitBreakerRulesCRD, "sentinel-circuitbreaker-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController CircuitBreakerRulesCRD error", zap.Error(err))
			return nil, err
		}
		err = ds.RegisterController(k8s.FlowRulesCRD, "sentinel-flow-rules")
		if err != nil {
			logger.Warn("k8s.RegisterController FlowRulesCRD error", zap.Error(err))
			return nil, err
		}
	}

	ds.Run()

	return &Sentinel{
		ds:     ds,
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
