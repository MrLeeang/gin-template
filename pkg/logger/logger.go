package logger

import (
	"fmt"
	"gin-template/pkg/config"
	"io"
	"os"
	"path"
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

	// 创建一个文件输出
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	writer := io.MultiWriter(file, os.Stdout)

	var level zapcore.Level
	if config.Global.Server.Debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
		writer = io.MultiWriter(file)

	}

	writeSyncer := zapcore.AddSync(writer)

	// 创建 Core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 控制台编码器
		writeSyncer,                              // 输出到文件
		level,                                    // 日志级别
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
