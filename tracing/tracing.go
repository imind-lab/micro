/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package tracing

import (
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

func InitTracer() (*tracesdk.TracerProvider, error) {
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
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceNamespaceKey.String(namespace),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceInstanceIDKey.String(hostname),
		)),
	)
	return tp, nil
}

func GetTrace(name string) trace.Tracer {
	return otel.Tracer(name)
}
