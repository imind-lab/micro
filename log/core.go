/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package log

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, format string, options ...zap.Option) *zap.Logger {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress, format)

	opts := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(0), zap.Development()}
	opts = append(opts, options...)
	return zap.New(core, opts...)
}

func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, format string) zapcore.Core {
	// 日志文件路径配置
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

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

func GetLogger(ctx context.Context, layer, fn string) *zap.Logger {
	return ctxzap.Extract(ctx).With(zap.String("layer", layer), zap.String("func", fn))
}
