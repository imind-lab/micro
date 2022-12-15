package micro

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gopkg.in/tomb.v2"

	grpcx "github.com/imind-lab/micro/grpc"
	"github.com/imind-lab/micro/log"
)

type Service interface {
	// Name The service name
	Name() string
	// Init initialises options
	Init(...Option)
	// Options returns the current options
	Options() Options
	// ServeMux return grpc-gateway serveMux
	ServeMux() *runtime.ServeMux
	// GrpcServer returns the grpc server
	GrpcServer() *grpc.Server
	// Run the service
	Run() error
	// Stop the service
	Stop() error
	// String The service implementation
	String() string
}

func NewService(opts ...Option) Service {
	service := new(service)
	options := newOptions(opts...)
	// set opts
	service.opts = options

	return service
}

type service struct {
	opts Options

	serveMux   *runtime.ServeMux
	grpcServer *grpc.Server
	httpServer *http.Server

	ctx    context.Context
	cancel context.CancelFunc

	once sync.Once
}

func (s service) Name() string {
	return s.opts.Name
}

func (s *service) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}
	s.once.Do(func() {
		s.grpcServer = s.newGrpcServer()
		server, mux := s.newHttpServer(s.opts.Handlers...)
		s.httpServer = server
		s.serveMux = mux
	})
}

func (s service) Options() Options {
	return s.opts
}

func (s service) ServeMux() *runtime.ServeMux {
	return s.serveMux
}

func (s service) GrpcServer() *grpc.Server {
	return s.grpcServer
}

func (s service) Run() error {
	for _, fn := range s.opts.BeforeRun {
		fn()
	}

	// 初始化上下文
	s.ctx, s.cancel = context.WithCancel(s.opts.Context)

	gRPCEndPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.grpc"))
	httpEndPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.http"))
	grpcListener, err := net.Listen("tcp", gRPCEndPoint)
	if err != nil {
		s.opts.Logger.Error("TCP Listen err", zap.Error(err))
		return err
	}

	httpListener, err := net.Listen("tcp", httpEndPoint)
	if err != nil {
		s.opts.Logger.Error("TCP Listen err", zap.Error(err))
	}

	var tb1 tomb.Tomb
	tb1.Go(func() error {
		// start gRPC server
		return s.startGrpcServer(grpcListener)
	})

	var tb2 tomb.Tomb
	tb2.Go(func() error {
		// start http server
		return s.startHttpServer(httpListener)
	})

	for _, fn := range s.opts.AfterRun {
		fn()
	}

	for {
		select {
		case <-tb1.Dead():
			s.opts.Logger.Warn("tb1 Dead")
			tb1 = tomb.Tomb{}
			tb1.Go(func() error {
				return s.startGrpcServer(grpcListener)
			})
		case <-tb2.Dead():
			s.opts.Logger.Warn("tb2 Dead")
			tb2 = tomb.Tomb{}
			tb2.Go(func() error {
				return s.startHttpServer(httpListener)
			})
		}
	}

	return nil
}

func (s service) Stop() error {
	for _, fn := range s.opts.BeforeStop {
		fn()
	}

	s.cancel()

	s.grpcServer.GracefulStop()
	err := s.httpServer.Shutdown(s.opts.Context)
	if err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		fn()
	}

	return nil
}

func (s service) String() string {
	return fmt.Sprintf("%s service instance", s.opts.Name)
}

func (s service) newGrpcServer() *grpc.Server {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	customFunc := func(p interface{}) (err error) {
		if s.opts.Logger != nil {
			s.opts.Logger.Error("recovery panic", zap.String("trace", string(debug.Stack())), zap.Error(p.(error)))
		}
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	unaryInterceptors = append(unaryInterceptors, grpc_recovery.UnaryServerInterceptor(recoveryOpts...), grpc_ctxtags.UnaryServerInterceptor(), otelgrpc.UnaryServerInterceptor())
	streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor(recoveryOpts...), grpc_ctxtags.StreamServerInterceptor(), otelgrpc.StreamServerInterceptor())

	if s.opts.Logger != nil {
		opts := []grpc_zap.Option{
			grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
				// will not log gRPC calls if it was a call to healthcheck and no error was raised
				if err == nil && fullMethodName == "/grpc.health.v1.Health/Check" {
					return false
				}
				// by default everything will be logged
				return true
			}),
		}
		unaryInterceptors = append(unaryInterceptors, log.UnaryServerInterceptor(), grpc_zap.UnaryServerInterceptor(s.opts.Logger, opts...))
		streamInterceptors = append(streamInterceptors, log.StreamServerInterceptor(), grpc_zap.StreamServerInterceptor(s.opts.Logger, opts...))
	}

	var serverOpt []grpc.ServerOption
	if s.opts.ServerCred != nil {
		serverOpt = append(serverOpt, grpc.Creds(s.opts.ServerCred))
	}
	serverOpt = append(serverOpt, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))

	grpcServer := grpc.NewServer(serverOpt...)

	reflection.Register(grpcServer)

	// 注册gRPC健康检测
	srv := health.NewServer()
	srv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, srv)

	return grpcServer
}

func (s service) newHttpServer(handlers ...Handler) (*http.Server, *runtime.ServeMux) {
	mux := runtime.NewServeMux(grpcx.EnableGatewayJsonTag())

	handler := grpcx.GrpcHandlerFunc(s.grpcServer, mux)
	for _, handle := range handlers {
		handler = handle(handler)
	}

	httpServer := &http.Server{
		Handler: handler,
	}

	return httpServer, mux
}

func (s service) startGrpcServer(listener net.Listener) error {
	s.opts.Logger.Info("GrpcServer is running" + listener.Addr().String())

	if err := s.grpcServer.Serve(listener); err != nil {
		s.opts.Logger.Error("HttpServer running error", zap.Error(err))
		return err
	}
	return nil
}

func (s service) startHttpServer(listener net.Listener) error {
	s.opts.Logger.Info("HttpServer is running" + listener.Addr().String())
	if err := s.httpServer.Serve(listener); err != nil {
		s.opts.Logger.Error("HttpServer running error", zap.Error(err))
		return err
	}
	return nil
}

func ClientConn(ctx context.Context, name string, tls bool) (*grpc.ClientConn, error) {
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

	unaryInterceptors = append(unaryInterceptors, grpc_retry.UnaryClientInterceptor(retryOpts...), log.UnaryClientInterceptor(), grpc_zap.UnaryClientInterceptor(ctxzap.Extract(ctx), zapOpts...), otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(otel.GetTracerProvider())))
	streamInterceptors = append(streamInterceptors, log.StreamClientInterceptor(), grpc_zap.StreamClientInterceptor(ctxzap.Extract(ctx), zapOpts...), otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(otel.GetTracerProvider())))

	var dialOpt []grpc.DialOption
	if tls {
		dialOpt = append(dialOpt, grpc.WithTransportCredentials(grpcx.NewGrpcCred().ClientCred()))
	} else {
		dialOpt = append(dialOpt, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	dialOpt = append(dialOpt, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)))

	conn, err := grpc.Dial(addr, dialOpt...)

	return conn, err
}
