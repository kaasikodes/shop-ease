package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/otel/trace"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DefaultLogger (uses Go's built-in log package)
type DefaultLogger struct {
	logFile           string
	logger            *log.Logger
	format            string
	primaryIdentifier string
	traceID           string
}

// New creates a new instance of DefaultLogger
func New(config LogConfig) Logger {

	// Initialize the logger

	config = defineLogConfig(config)
	lumberjackLogger := &lumberjack.Logger{ //for auto deletion of log files
		Filename:   config.LogFilePath,
		MaxSize:    config.MaxSizeMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   *config.Compress,
	}
	// Create a multi-writer (logs to file + terminal)
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)

	return &DefaultLogger{
		logFile:           config.LogFilePath,
		logger:            log.New(multiWriter, "", log.LstdFlags),
		format:            config.Format,
		primaryIdentifier: config.PrimaryIdentifier,
	}
}

// with context for the primary purpose of appending trace id to the log, might be refactored later to be precise in its functional description
func (l *DefaultLogger) WithContext(ctx context.Context) Logger {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	traceID := sc.TraceID().String()

	return &DefaultLogger{
		logger:            l.logger,
		format:            l.format,
		logFile:           l.logFile,
		primaryIdentifier: l.primaryIdentifier,
		traceID:           traceID,
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(v ...any) {
	msg := fmt.Sprintf(l.format, "INFO", fmt.Sprint(v...), l.primaryIdentifier, l.traceID)

	l.logger.Println(msg)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(v ...any) {
	msg := fmt.Sprintf(l.format, "WARN", fmt.Sprint(v...), l.primaryIdentifier, l.traceID)

	l.logger.Println(msg)
}

// Error logs an error message
func (l *DefaultLogger) Error(v ...any) {
	msg := fmt.Sprintf(l.format, "ERROR", fmt.Sprint(v...), l.primaryIdentifier, l.traceID)

	l.logger.Println(msg)
}

// Fatal error
func (l *DefaultLogger) Fatal(v ...any) {
	msg := fmt.Sprintf(l.format, "FATAL ERROR", fmt.Sprint(v...), l.primaryIdentifier, l.traceID)

	l.logger.Fatal(msg)

}
