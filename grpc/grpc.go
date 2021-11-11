package grpc

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/imind-lab/micro/tracing"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func EnableGatewayJsonTag() runtime.ServeMuxOption {
	return runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	})
}

func ClientConn(ctx context.Context, name string, tls bool) (*grpc.ClientConn, io.Closer, error) {
	service := viper.GetString("rpc." + name + ".service")
	port := viper.GetInt("rpc." + name + ".port")
	addr := fmt.Sprintf("%s:%d", service, port)

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(3 * time.Second),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
	}

	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(grpc_zap.DurationToDurationField),
	}

	var unaryInterceptors []grpc.UnaryClientInterceptor
	var streamInterceptors []grpc.StreamClientInterceptor

	unaryInterceptors = append(unaryInterceptors, grpc_retry.UnaryClientInterceptor(retryOpts...), grpc_zap.UnaryClientInterceptor(ctxzap.Extract(ctx), zapOpts...))
	streamInterceptors = append(streamInterceptors, grpc_zap.StreamClientInterceptor(ctxzap.Extract(ctx), zapOpts...))

	traceName := viper.GetString("tracing.name.client")
	tracer, closer, err := tracing.InitTracer(traceName)
	if err == nil {
		filterFunc := grpc_opentracing.WithFilterFunc(func(ctx context.Context, fullMethodName string) bool {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if fullMethodName == "/grpc.health.v1.Health/Check" {
				return false
			}
			// by default everything will be logged
			return true
		})

		unaryInterceptors = append(unaryInterceptors, grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer), filterFunc))
		streamInterceptors = append(streamInterceptors, grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer), filterFunc))
	}

	var dialOpt []grpc.DialOption
	if tls {
		dialOpt = append(dialOpt, grpc.WithTransportCredentials(NewGrpcCred().ClientCred()))
	} else {
		dialOpt = append(dialOpt, grpc.WithInsecure())
	}

	dialOpt = append(dialOpt, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)))

	conn, err := grpc.Dial(addr, dialOpt...)

	return conn, closer, nil
}
