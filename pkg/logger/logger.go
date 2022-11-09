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
	//设置日志级别
	if config.Config.Debug {
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
var Logger = logger()

var Info = Logger.Info
var Error = Logger.Error
var Debug = Logger.Debug
var Warn = Logger.Warn
var Panic = Logger.Panic
