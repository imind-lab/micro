package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/imind-lab/micro/tracing"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
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

func ClientConn(ctx context.Context, name string, tls bool) (*grpc.ClientConn, *tracesdk.TracerProvider, error) {
	service := viper.GetString("rpc." + name + ".service")
	port := viper.GetInt("rpc." + name + ".port")
	addr := fmt.Sprintf("%s:%d", service, port)

	fmt.Println("addr", addr)

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

	provider, err := tracing.InitTracer()
	if err == nil {
		unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryClientInterceptor())
		streamInterceptors = append(streamInterceptors, otelgrpc.StreamClientInterceptor())
	}

	var dialOpt []grpc.DialOption
	if tls {
		dialOpt = append(dialOpt, grpc.WithTransportCredentials(NewGrpcCred().ClientCred()))
	} else {
		dialOpt = append(dialOpt, grpc.WithInsecure())
	}

	fmt.Println("dialOpt", dialOpt)
	dialOpt = append(dialOpt, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)))

	conn, err := grpc.Dial(addr, dialOpt...)

	return conn, provider, nil
}
