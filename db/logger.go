package db

import (
	"context"
	"errors"
	"gin-template/pkg/config"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type ContextFn func(ctx context.Context) []zapcore.Field

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	Context                   ContextFn
}

func NewZapLogger() *Logger {
	return &Logger{
		ZapLogger:                 zap.L(),
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             1 * time.Second,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		Context:                   nil,
	}
}

func (l *Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.LogLevel = level
	return l
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger(ctx).Sugar().Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger(ctx).Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger(ctx).Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		logger.Error("[GORM]", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		// logger.Sugar().Errorf("%s\telapsed:%s\trows:%d\tsql:%s\n%s", err, elapsed, rows, sql, string(debug.Stack()))
		logger.Sugar().Errorf(string(debug.Stack()))

	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		logger.Warn("[GORM]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		// logger.Sugar().Warnf("elapsed:%s\trows:%d\tsql:%s", elapsed, rows, sql)
	case l.LogLevel >= gormlogger.Info:
		// debug模式，这是为了全局可以打印sql
		sql, rows := fc()
		if config.Global.Debug {
			logger.Debug("[GORM]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		} else {
			// 非debug模式下，db.Session.Debug().First(&user)，打印sql语句
			logger.Info("[GORM]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
		}
		// logger.Sugar().Debugf("elapsed:%s\trows:%d\tsql:%s", elapsed, rows, sql)
	}
}

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)

func (l Logger) logger(ctx context.Context) *zap.Logger {
	logger := l.ZapLogger
	if l.Context != nil {
		fields := l.Context(ctx)
		logger = logger.With(fields...)
	}

	if l.SkipCallerLookup {
		return logger
	}

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return logger
}
