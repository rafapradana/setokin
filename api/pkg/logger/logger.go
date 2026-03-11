// Package logger provides structured logging using Zap.
package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init initializes the global logger based on APP_ENV.
func Init() {
	once.Do(func() {
		var err error
		if os.Getenv("APP_ENV") == "production" {
			log, err = zap.NewProduction()
		} else {
			log, err = zap.NewDevelopment()
		}
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})
}

// Get returns the global logger instance.
func Get() *zap.Logger {
	if log == nil {
		Init()
	}
	return log
}

// Info logs an info message with optional fields.
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs an error message with optional fields.
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Warn logs a warning message with optional fields.
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug logs a debug message with optional fields.
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync flushes any buffered log entries.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
