package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger initializes a Zap logger
func NewZapLogger(config LogConfig) Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", config.LogFilePath}
	zapLogger, err := cfg.Build()
	if err != nil {
		panic("failed to build zap logger: " + err.Error())

	}
	defer zapLogger.Sync() // flush buffer
	return &ZapLogger{
		logger: zapLogger.Sugar().With("primary_identifier", config.PrimaryIdentifier),
	}
}
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	traceID := sc.TraceID().String()

	// Attach trace ID as a field
	return &ZapLogger{
		logger: l.logger.With("trace_id", traceID),
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
func (l *ZapLogger) Fatal(v ...any) {
	l.logger.Fatal("ERROR", "msg", v)
}
