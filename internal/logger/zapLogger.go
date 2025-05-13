package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger initializes a Zap logger
func NewZapLogger(config LogConfig) Logger {
	// Initialize the logger

	config = defineLogConfig(config)

	lumberjackLogger := &lumberjack.Logger{ //for auto deletion of log files
		Filename:   config.LogFilePath,
		MaxSize:    config.MaxSizeMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   *config.Compress,
	}
	// Define log outputs
	writer := zapcore.AddSync(lumberjackLogger)
	console := zapcore.AddSync(os.Stdout)

	// Set up encoder (JSON format with detailed timestamp and caller)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.LevelKey = "level"
	encoderCfg.NameKey = "logger"
	encoderCfg.CallerKey = "caller"
	encoderCfg.MessageKey = "msg"
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	// Create core with both stdout and file
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.NewMultiWriteSyncer(console, writer), zap.InfoLevel)

	// Add options: caller info, stack traces, and sugar for structured logging
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar().With("primary_identifier", config.PrimaryIdentifier)

	defer zapLogger.Sync() // flush buffer
	return &ZapLogger{
		logger: zapLogger,
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
