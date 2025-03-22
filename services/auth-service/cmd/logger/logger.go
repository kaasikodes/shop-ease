package logger

// Logger defines the interface for logging
type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}
