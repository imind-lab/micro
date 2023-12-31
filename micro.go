package micro

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	grpcx "github.com/imind-lab/micro/v2/grpc"
)

type Service interface {
	// Name The service name
	Name() string
	// Init initialises options
	Init(context.Context, ...Option)
	// Options returns the current options
	Options() Options
	// ServeMux return grpc-gateway serveMux
	ServeMux() *runtime.ServeMux
	// GrpcServer returns the grpc server
	GrpcServer() *grpc.Server
	// Run the service
	Run(context.Context) error
	// Stop the service
	Stop(context.Context, context.CancelFunc) error
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

	once sync.Once
}

func (s *service) Name() string {
	return s.opts.Name
}

func (s *service) Init(ctx context.Context, opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}
	s.once.Do(func() {
		s.grpcServer = s.newGrpcServer(ctx)
		server, mux := s.newHttpServer(s.opts.Handlers...)
		s.httpServer = server
		s.serveMux = mux
	})
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) ServeMux() *runtime.ServeMux {
	return s.serveMux
}

func (s *service) GrpcServer() *grpc.Server {
	return s.grpcServer
}

func (s *service) Run(ctx context.Context) error {
	for _, fn := range s.opts.BeforeRun {
		fn()
	}

	gRPCEndPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.grpc"))
	httpEndPoint := fmt.Sprintf(":%d", viper.GetInt("service.port.http"))
	grpcListener, err := net.Listen("tcp", gRPCEndPoint)
	if err != nil {
		s.opts.Logger.Error().Err(err).Msg("TCP Listen error")
		return err
	}

	httpListener, err := net.Listen("tcp", httpEndPoint)
	if err != nil {
		s.opts.Logger.Error().Err(err).Msg("TCP Listen error")
	}

	rg := &run.Group{}
	rg.Add(func() error {
		s.opts.Logger.Info().Msg("GrpcServer is running" + grpcListener.Addr().String())
		err := s.grpcServer.Serve(grpcListener)
		if err != nil {
			s.opts.Logger.Error().Err(err).Msg("GrpcServer running error")
		}
		return err
	}, func(err error) {
		s.grpcServer.GracefulStop()
		s.grpcServer.Stop()
	})

	rg.Add(func() error {
		s.opts.Logger.Info().Msg("HttpServer is running" + httpListener.Addr().String())
		err := s.httpServer.Serve(httpListener)
		if err != nil {
			s.opts.Logger.Error().Err(err).Msg("HttpServer running error")
		}
		return err
	}, func(err error) {
		s.httpServer.Shutdown(ctx)
		s.httpServer.Close()
	})

	rg.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	if err := rg.Run(); err != nil {
		s.opts.Logger.Error().Err(err).Msg("run group run error")
		for _, fn := range s.opts.AfterRun {
			fn()
		}
		return err
	}
	return nil
}

func (s *service) Stop(ctx context.Context, cancel context.CancelFunc) error {
	for _, fn := range s.opts.BeforeStop {
		fn()
	}

	cancel()

	s.grpcServer.GracefulStop()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		fn()
	}

	return nil
}

func (s *service) String() string {
	return fmt.Sprintf("%s service instance", s.opts.Name)
}

func (s *service) newGrpcServer(ctx context.Context) *grpc.Server {

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	authFn := func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}
		if token != "daniel" {
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		}
		return ctx, nil
	}

	filterHealthZ := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	}

	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})

	panicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		s.opts.Logger.Error().Any("panic", p).Any("stack", debug.Stack()).Msg("recovered from panic")
		return status.Errorf(codes.Internal, "%s", p)
	}
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	logger := zerolog.Ctx(ctx)

	unaryInterceptors = append(unaryInterceptors,
		srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarContext)),
		logging.UnaryServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(logTraceID)),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(filterHealthZ)),
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(panicRecoveryHandler)),
	)
	streamInterceptors = append(streamInterceptors,
		srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarContext)),
		logging.StreamServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(logTraceID)),
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), selector.MatchFunc(filterHealthZ)),
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(panicRecoveryHandler)),
	)

	var serverOpt = []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	if s.opts.ServerCred != nil {
		serverOpt = append(serverOpt, grpc.Creds(s.opts.ServerCred))
	}
	serverOpt = append(serverOpt, grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...))

	grpcServer := grpc.NewServer(serverOpt...)

	reflection.Register(grpcServer)

	// 注册gRPC健康检测
	srv := health.NewServer()
	srv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, srv)

	return grpcServer
}

func (s *service) newHttpServer(handlers ...Handler) (*http.Server, *runtime.ServeMux) {
	mux := runtime.NewServeMux(grpcx.EnableGatewayJsonTag())

	handler := grpcx.HandlerFunc(s.grpcServer, mux)
	for _, handle := range handlers {
		handler = handle(handler)
	}

	httpServer := &http.Server{
		Handler: handler,
	}

	return httpServer, mux
}

func ClientConn(ctx context.Context, name string, tls bool) (*grpc.ClientConn, error) {
	service := viper.GetString("rpc." + name + ".service")
	port := viper.GetInt("rpc." + name + ".port")
	addr := fmt.Sprintf("%s:%d", service, port)

	retryOpts := []retry.CallOption{
		retry.WithMax(3),
		retry.WithBackoff(retry.BackoffExponential(100 * time.Millisecond)),
		retry.WithCodes(codes.NotFound, codes.Aborted),
	}

	reg := prometheus.NewRegistry()
	clMetrics := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg.MustRegister(clMetrics)

	logger := zerolog.Ctx(ctx)

	var unaryInterceptors []grpc.UnaryClientInterceptor
	var streamInterceptors []grpc.StreamClientInterceptor

	unaryInterceptors = append(unaryInterceptors,
		timeout.UnaryClientInterceptor(time.Millisecond*500),
		retry.UnaryClientInterceptor(retryOpts...),
		clMetrics.UnaryClientInterceptor(grpcprom.WithExemplarFromContext(exemplarContext)),
		logging.UnaryClientInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(logTraceID)))
	streamInterceptors = append(streamInterceptors,
		clMetrics.StreamClientInterceptor(grpcprom.WithExemplarFromContext(exemplarContext)),
		logging.StreamClientInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(logTraceID)))

	var dialOpt = []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	if tls {
		dialOpt = append(dialOpt, grpc.WithTransportCredentials(grpcx.NewGrpcCred().ClientCred()))
	} else {
		dialOpt = append(dialOpt, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	dialOpt = append(dialOpt, grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...))

	conn, err := grpc.Dial(addr, dialOpt...)

	return conn, err
}

var (
	logTraceID = func(ctx context.Context) logging.Fields {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return logging.Fields{"traceID", span.TraceID().String()}
		}
		return nil
	}

	exemplarContext = func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	// interceptorLogger adapts zerolog logger to interceptor logger.
	// This code is simple enough to be copied and not imported.
	interceptorLogger = func(l *zerolog.Logger) logging.Logger {
		return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l := l.With().Fields(fields).Logger()

			switch lvl {
			case logging.LevelDebug:
				l.Debug().Msg(msg)
			case logging.LevelInfo:
				l.Info().Msg(msg)
			case logging.LevelWarn:
				l.Warn().Msg(msg)
			case logging.LevelError:
				l.Error().Msg(msg)
			default:
				panic(fmt.Sprintf("unknown level %v", lvl))
			}
		})
	}
)
