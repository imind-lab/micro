/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package tracing

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/imind-lab/micro/v2/util"
)

func InitTracer(ctx context.Context) (func(context.Context) error, error) {
	service := viper.GetString("service.name")
	namespace := viper.GetString("service.namespace")
	version := viper.GetString("service.version")
	endpoint := viper.GetString("tracing.endpoint")

	hostname, _ := os.Hostname()

	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exporter),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceNamespaceKey.String(namespace),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceInstanceIDKey.String(hostname),
		)),
		sdktrace.WithSampler(filterHealthSampler{}),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp.Shutdown, nil
}

func StartSpan(ctx context.Context, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	layer, name := util.GetPtrFuncName()
	return otel.Tracer(layer).Start(ctx, util.AppendString(layer, ".", name), opts...)
}

type filterHealthSampler struct{}

func (as filterHealthSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if p.Name == "grpc.health.v1.Health/Check" {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	}
	return sdktrace.SamplingResult{
		Decision:   sdktrace.RecordAndSample,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (as filterHealthSampler) Description() string {
	return "filterHealthSampler"
}
