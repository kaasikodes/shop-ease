package logger

import (
	"context"
)

const DefaultLogFormat = "[%s]: %s primary_identifier %s trace_id %s"
const DefaultLogMaxSizeMB = 2
const DefaultLogMaxBackups = 1
const DefaultLogMaxAgeDays = 3

var DefaultLogCompress = true

func defineLogConfig(config LogConfig) LogConfig {
	// Initialize the logger

	if len(config.Format) == 0 {
		config.Format = DefaultLogFormat
	}
	if (config.MaxAgeDays) == 0 {
		config.MaxAgeDays = DefaultLogMaxAgeDays
	}
	if (config.MaxBackups) == 0 {
		config.MaxBackups = DefaultLogMaxBackups
	}
	if (config.MaxSizeMB) == 0 {
		config.MaxSizeMB = DefaultLogMaxSizeMB
	}
	if (config.Compress) == nil {
		config.Compress = &DefaultLogCompress
	}
	return config

}

type LogConfig struct {
	LogFilePath       string
	Format            string // e.g. "[%s]: %s primary_identifier %s trace_id %s"
	PrimaryIdentifier string
	MaxSizeMB         int   // max size per log file in MB
	MaxBackups        int   // number of old files to keep
	MaxAgeDays        int   // how many days to retain logs
	Compress          *bool // whether to compress rotated logs
}

// Logger defines the interface for logging
type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
	WithContext(ctx context.Context) Logger
}
