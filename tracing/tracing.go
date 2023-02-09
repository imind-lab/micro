/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package tracing

import (
    "context"
    "os"

    "github.com/spf13/viper"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
    "go.opentelemetry.io/otel/trace"

    "github.com/imind-lab/micro/v2/util"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
    service := viper.GetString("service.name")
    namespace := viper.GetString("service.namespace")
    version := viper.GetString("service.version")
    host := viper.GetString("tracing.agent.host")
    port := viper.GetString("tracing.agent.port")

    hostname, _ := os.Hostname()
    // Create the Jaeger exporter
    exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(host), jaeger.WithAgentPort(port)))
    if err != nil {
        return nil, err
    }
    tp := sdktrace.NewTracerProvider(
        // Always be sure to batch in production.
        sdktrace.WithBatcher(exp),
        // Record information about this application in a Resource.
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(service),
            semconv.ServiceNamespaceKey.String(namespace),
            semconv.ServiceVersionKey.String(version),
            semconv.ServiceInstanceIDKey.String(hostname),
        )),
        sdktrace.WithSampler(filterHealthSampler{}),
    )
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
    return tp, nil
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

func InitProvider(service, namespace string) (*sdktrace.TracerProvider, error) {
    host := viper.GetString("tracing.agent.host")
    port := viper.GetString("tracing.agent.port")

    hostname, _ := os.Hostname()
    // Create the Jaeger exporter
    exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(host), jaeger.WithAgentPort(port)))
    if err != nil {
        return nil, err
    }
    tp := sdktrace.NewTracerProvider(
        // Always be sure to batch in production.
        sdktrace.WithBatcher(exp),
        // Record information about this application in a Resource.
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(service),
            semconv.ServiceNamespaceKey.String(namespace),
            semconv.ServiceInstanceIDKey.String(hostname),
        )),
        sdktrace.WithSampler(filterHealthSampler{}),
    )
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
    return tp, nil
}
