/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package log

import (
	"context"
	"os"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/imind-lab/micro/util"
)

var debugEnabled bool

func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, format string, options ...zap.Option) *zap.Logger {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress, format)

	opts := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(0), zap.Development()}
	opts = append(opts, options...)
	return zap.New(core, opts...)
}

func newCore(filePath string, initLevel zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, format string) zapcore.Core {
	// 日志文件路径配置
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}

	// 设置日志级别
	//atomicLevel := zap.NewAtomicLevel()
	//atomicLevel.SetLevel(level)
	atomicLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= initLevel || debugEnabled
	})

	// 公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		StacktraceKey:  "stacktrace",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		//EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeCaller: zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:   zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewCore(
		encoder, // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}

func GetLogger(ctx context.Context) *zap.Logger {
	layer, name := util.GetPtrFuncName()
	return ctxzap.Extract(ctx).With(zap.String("layer", layer), zap.String("func", name))
}

func EnableDebug() {
	debugEnabled = true
}

func DisableDebug() {
	debugEnabled = false
}

// ServerInterceptor returns a new unary server interceptors that adds trace-id to the context
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = newContextForCall(ctx)
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server that adds trace-id to the context
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := newContextForCall(stream.Context())
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

// UnaryClientInterceptor returns a new unary client interceptor that adds trace-id to the context of external gRPC calls.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = newContextForCall(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that adds trace-id to the context of external gRPC calls.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = newContextForCall(ctx)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func newContextForCall(ctx context.Context) context.Context {
	traceId := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if tid := md.Get("trace-id"); len(tid) > 0 {
			traceId = tid[0]
		}
	} else {
		md = metadata.MD{}
	}

	if traceId == "" {
		traceId = uuid.New().String()
	}

	md.Set("trace-id", traceId)

	tags := grpc_ctxtags.NewTags()
	tags = tags.Set("traceId", traceId)
	ctx = grpc_ctxtags.SetInContext(ctx, tags)

	return metadata.NewIncomingContext(ctx, md)
}
