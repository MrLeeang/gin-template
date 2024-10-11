package logger

import (
	"gin-template/pkg/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

var (
	Debug func(msg string, fields ...zap.Field)
	Info  func(msg string, fields ...zap.Field)
	Warn  func(msg string, fields ...zap.Field)
	Error func(msg string, fields ...zap.Field)
)

var (
	Debugf func(template string, args ...interface{})
	Infof  func(template string, args ...interface{})
	Warnf  func(template string, args ...interface{})
	Errorf func(template string, args ...interface{})
)

func InitializeLogger() *zap.Logger {
	// 自定义时间编码器
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// 自定义 Encoder 配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大写编码器
		EncodeTime:     customTimeEncoder,                // ISO8601 时间格式
		EncodeDuration: zapcore.StringDurationEncoder,    // 字符串时间编码器
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 完整文件路径编码器
	}

	var level zapcore.Level
	if config.Global.Debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	// 创建 Core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 控制台编码器
		zapcore.AddSync(os.Stdout),
		level, // 日志级别
	)

	Logger = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(Logger)

	zap.RedirectStdLog(Logger)

	Debug = Logger.Debug
	Info = Logger.Info
	Warn = Logger.Warn
	Error = Logger.Error

	Debugf = Logger.Sugar().Debugf
	Infof = Logger.Sugar().Infof
	Warnf = Logger.Sugar().Warnf
	Errorf = Logger.Sugar().Errorf

	return Logger
}
