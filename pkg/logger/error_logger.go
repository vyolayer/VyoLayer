package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"vyolayer/pkg/errors"
)

// Logger interface for error logging
type Logger interface {
	LogError(err *errors.AppError, requestID string)
	LogInfo(message string, fields map[string]interface{})
	LogWarning(message string, fields map[string]interface{})
}

// DefaultLogger is a simple console logger
type DefaultLogger struct {
	isDevelopment bool
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(isDevelopment bool) *DefaultLogger {
	return &DefaultLogger{
		isDevelopment: isDevelopment,
	}
}

// LogError logs an error with appropriate formatting
func (l *DefaultLogger) LogError(err *errors.AppError, requestID string) {
	if err == nil {
		return
	}

	logEntry := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"level":      string(err.Severity),
		"code":       err.Code,
		"message":    err.Message,
		"httpStatus": err.HTTPStatus,
		"requestID":  requestID,
	}

	// Add metadata
	if len(err.Metadata) > 0 {
		logEntry["metadata"] = err.Metadata
	}

	// Add wrapped error if present
	if err.Wrapped != nil {
		logEntry["wrappedError"] = err.Wrapped.Error()
	}

	// Include stack trace for critical errors or in development
	if l.isDevelopment || err.Severity == errors.SeverityCritical || err.Severity == errors.SeverityError {
		if len(err.StackTrace) > 0 {
			logEntry["stackTrace"] = err.StackTrace
		}
	}

	l.logJSON(logEntry)
}

// LogInfo logs an informational message
func (l *DefaultLogger) LogInfo(message string, fields map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "INFO",
		"message":   message,
	}

	for k, v := range fields {
		logEntry[k] = v
	}

	l.logJSON(logEntry)
}

// LogWarning logs a warning message
func (l *DefaultLogger) LogWarning(message string, fields map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "WARNING",
		"message":   message,
	}

	for k, v := range fields {
		logEntry[k] = v
	}

	l.logJSON(logEntry)
}

// logJSON outputs a log entry as JSON
func (l *DefaultLogger) logJSON(entry map[string]interface{}) {
	if l.isDevelopment {
		// Pretty print in development
		jsonBytes, _ := json.MarshalIndent(entry, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		// Single line JSON in production
		jsonBytes, _ := json.Marshal(entry)
		fmt.Println(string(jsonBytes))
	}
}

// Global logger instance
var globalLogger Logger

// InitLogger initializes the global logger
func InitLogger(isDevelopment bool) {
	globalLogger = NewDefaultLogger(isDevelopment)
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Initialize with default settings if not already initialized
		isDev := os.Getenv("ENV") != "production"
		globalLogger = NewDefaultLogger(isDev)
	}
	return globalLogger
}

// SetLogger sets a custom logger
func SetLogger(logger Logger) {
	globalLogger = logger
}

// Convenience functions that use the global logger

// LogError logs an error using the global logger
func LogError(err *errors.AppError, requestID string) {
	GetLogger().LogError(err, requestID)
}

// LogInfo logs an info message using the global logger
func LogInfo(message string, fields map[string]interface{}) {
	GetLogger().LogInfo(message, fields)
}

// LogWarning logs a warning using the global logger
func LogWarning(message string, fields map[string]interface{}) {
	GetLogger().LogWarning(message, fields)
}

// LogStandardError logs a standard Go error
func LogStandardError(err error, requestID string) {
	if err == nil {
		return
	}

	// Convert to AppError if possible
	if appErr, ok := errors.As(err); ok {
		LogError(appErr, requestID)
		return
	}

	// Otherwise log as a standard error
	log.Printf("[ERROR] [RequestID: %s] %v", requestID, err)
}
