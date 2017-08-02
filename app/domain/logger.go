package domain

import "github.com/gobuffalo/buffalo"

// Logger interface is used throughout Marvin
// to log all the things.
type Logger interface {
	// WithField will return a logger that will log the given field with value for each log message
	WithField(string, interface{}) Logger
	// WithFields is equivalent to calling WithField for each key-value pair in the map.
	WithFields(map[string]interface{}) Logger
	// Log a debug message with a format string
	Debugf(string, ...interface{})
	//Log an info message with a format string
	Infof(string, ...interface{})
	//Log a warning message with a format string
	Warnf(string, ...interface{})
	//Log an error message with a format string
	Errorf(string, ...interface{})
	//Log a fatal message with a format string
	Fatalf(string, ...interface{})
	// Log a debug message
	Debug(...interface{})
	// Log an info mesage
	Info(...interface{})
	// Log a warning message
	Warn(...interface{})
	// Log an error message
	Error(...interface{})
	// Log a fatal message
	Fatal(...interface{})
	// 😱 Panic with this message
	Panic(...interface{})
}

type loggerWrapper struct {
	buffalo.Logger
}

// LoggerFromBuffaloLogger creates a marvin-specific Logger from a buffalo.Logger
func LoggerFromBuffaloLogger(l buffalo.Logger) Logger {
	lw := loggerWrapper{l}
	return &lw
}

// WithFields eturns a logger that logs the given fields with each message
func (l *loggerWrapper) WithFields(fields map[string]interface{}) Logger {
	return LoggerFromBuffaloLogger(l.Logger.WithFields(fields))
}

// WithField returns a logger that logs the field with every message
func (l *loggerWrapper) WithField(name string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{name: value})
}
