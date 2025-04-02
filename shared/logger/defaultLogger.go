package logger

import (
	"io"
	"log"
	"os"
)

// DefaultLogger (uses Go's built-in log package)
type DefaultLogger struct {
	logFile string
	logger  *log.Logger
}

// New creates a new instance of DefaultLogger
func New(logFileLocation string) Logger {
	file, err := os.OpenFile(logFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create a multi-writer (logs to file + terminal)
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Initialize the logger

	return &DefaultLogger{
		logFile: logFileLocation,
		logger:  log.New(multiWriter, "", log.LstdFlags),
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(v ...any) {
	l.logger.Println("INFO:", v)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(v ...any) {
	l.logger.Println("WARN:", v)
}

// Error logs an error message
func (l *DefaultLogger) Error(v ...any) {
	l.logger.Println("ERROR:", v)
}

// Fatal error
func (l *DefaultLogger) Fatal(v ...any) {
	l.logger.Fatal("FATAL ERROR:", v)
}
