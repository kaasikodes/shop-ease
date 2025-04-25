package logger

import "context"

const DefaultLogFormat = "[%s]: %s primary_identifier %s trace_id %s"

type LogConfig struct {
	LogFilePath       string
	Format            string // e.g. "[%s]: %s primary_identifier %s trace_id %s"
	PrimaryIdentifier string
}

// Logger defines the interface for logging
type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
	WithContext(ctx context.Context) Logger
}
