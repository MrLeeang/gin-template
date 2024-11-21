package logger

import (
	"context"
	"gin-template/pkg/config"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Debug func(msg string, fields ...zap.Field)
	Info  func(msg string, fields ...zap.Field)
	Warn  func(msg string, fields ...zap.Field)
	Error func(msg string, fields ...zap.Field)
	Panic func(msg string, fields ...zap.Field)
)

var (
	Debugf func(template string, args ...interface{})
	Infof  func(template string, args ...interface{})
	Warnf  func(template string, args ...interface{})
	Errorf func(template string, args ...interface{})
	Panicf func(template string, args ...interface{})
)

type Logger struct {
	log *zap.Logger
}

func (logger *Logger) WithContext(ctx context.Context) *zap.Logger {
	return logger.log.With(zap.Any("trace_id", ctx.Value("trace_id")))
}
func (logger *Logger) Sync() error {
	return myLogger.log.Sync()
}

func WithContext(ctx context.Context) *zap.Logger {
	return myLogger.WithContext(ctx)
}

var myLogger *Logger

func InitializeZapLogger() *Logger {

	// 自定义时间编码器
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	zConfig := zap.NewProductionConfig()
	zConfig.EncoderConfig.EncodeTime = customTimeEncoder
	zConfig.EncoderConfig.TimeKey = "time"

	var level zapcore.Level

	if config.Global.Debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	zConfig.Level = zap.NewAtomicLevelAt(level)

	// 创建 logger
	logger, err := zConfig.Build()

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	Debug = logger.Debug
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
	Panic = logger.Panic

	Debugf = logger.Sugar().Debugf
	Infof = logger.Sugar().Infof
	Warnf = logger.Sugar().Warnf
	Errorf = logger.Sugar().Errorf
	Panicf = logger.Sugar().Panicf

	myLogger = &Logger{log: logger}

	return myLogger
}
