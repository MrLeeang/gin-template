package logger

import (
	"fmt"
	"gin-template/pkg/config"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MyFormatter MyFormatter
type MyFormatter struct{}

// Format Format
func (s *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	// 日志格式
	msg := fmt.Sprintf("%s [%s] %s:%v	%s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Caller.File, entry.Caller.Line, entry.Message)
	return []byte(msg), nil
}

// Logger 定义日志
func logger() *logrus.Logger {
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

	//实例化
	logger := logrus.New()

	//设置日志级别
	if config.Global.Server.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	//设置输出控制和文件
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	// 日志输出
	logger.SetOutput(fileAndStdoutWriter)

	// 显示行号
	logger.SetReportCaller(true)

	// 设置日志格式
	logger.SetFormatter(new(MyFormatter))

	return logger
}

// Logger 打印日志调用
// var Logger = logger()

func newLogger() *zap.Logger {

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

	lg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(lg)
	return lg
}

// 创建一个新的 zap logger
var Logger = newLogger()

// var Logger, _ = zap.NewProduction()

var Info = Logger.Info
var Infof = Logger.Sugar().Infof
var Error = Logger.Error
var Errorf = Logger.Sugar().Errorf
var Debug = Logger.Debug
var Debugf = Logger.Sugar().Debugf
var Warn = Logger.Warn
var Warnf = Logger.Sugar().Warnf
var Panic = Logger.Panic
var Panicf = Logger.Sugar().Panicf
