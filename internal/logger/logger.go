package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger with the given environment
func Init(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = config.Build()
	if err != nil {
		os.Exit(1)
	}
}

// Info logs a message at info level
func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

// Error logs a message at error level
func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

// Debug logs a message at debug level
func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

// Warn logs a message at warn level
func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

// Fatal logs a message at fatal level
func Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func With(fields ...zapcore.Field) *zap.Logger {
	return log.With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}
