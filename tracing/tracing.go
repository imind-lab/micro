/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package tracing

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func InitTracer(service string) (opentracing.Tracer, io.Closer, error) {
	tracingAgent := viper.GetString("tracing.agent")
	tracingType := viper.GetString("tracing.type")
	tracingParam := viper.GetFloat64("tracing.param")

	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  tracingType,
			Param: tracingParam,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: tracingAgent,
		},
	}

	jLogger := jaeger.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(config.Logger(jLogger), config.Metrics(jMetricsFactory))

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, err
}

func StartSpan(ctx context.Context, name string, fields ...log.Field) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	span.LogFields(fields...)

	return span, ctx
}
