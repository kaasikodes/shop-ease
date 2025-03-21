package logger

import (
	"go.uber.org/zap"
)

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger initializes a Zap logger
func NewZapLogger() Logger {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flush buffer
	return &ZapLogger{
		logger: zapLogger.Sugar(),
	}
}

func (l *ZapLogger) Info(v ...any) {
	l.logger.Infow("INFO", "msg", v)
}

func (l *ZapLogger) Warn(v ...any) {
	l.logger.Warnw("WARN", "msg", v)
}

func (l *ZapLogger) Error(v ...any) {
	l.logger.Errorw("ERROR", "msg", v)
}
