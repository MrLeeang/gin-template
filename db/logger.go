package db

import (
	"context"
	"errors"
	"gin-template/pkg/config"
	zlogger "gin-template/pkg/logger"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type zapLogger struct {
	level                     logger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func (l zapLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l zapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level < logger.Info {
		return
	}
	// zap.L().Sugar().Debugf(msg, data...)
	zlogger.WithContext(ctx).Sugar().Debugf(msg, data...)
}

func (l zapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level < logger.Warn {
		return
	}
	// zap.L().Sugar().Warnf(msg, data...)
	zlogger.WithContext(ctx).Sugar().Warnf(msg, data...)
}

func (l zapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level < logger.Error {
		return
	}

	// zap.L().Sugar().Errorf(msg, data...)
	zlogger.WithContext(ctx).Sugar().Errorf(msg, data...)
}

func (l zapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	if l.level <= 0 {
		return
	}

	lg := zlogger.WithContext(ctx)
	sql, rows := fc()
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		// lg.Errorf("SQL: %s, Error: %s", sql, err.Error())
		// lg.Errorf(string(debug.Stack()))
		lg.Error("[GORM]", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql), zap.String("stack", string(debug.Stack())))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.level >= logger.Warn:
		// lg.Warnf("Warning: SQL: %s, Rows: %d, Duration: %s", sql, rows, time.Since(begin))
		lg.Warn("[GORM]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.level >= logger.Info:
		// debug模式，这是为了全局可以打印sql
		if config.Global.Debug {
			lg.Debug("[GORM]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
			// lg.Debugf("SQL: %s, Rows: %d, Duration: %s", sql, rows, time.Since(begin))
		} else {
			// 非debug模式下，db.Session.Debug().First(&user)，打印sql语句
			lg.Info("[GORM]", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
			// lg.Infof("SQL: %s, Rows: %d, Duration: %s", sql, rows, time.Since(begin))
		}
	}
}
