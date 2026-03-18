package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vyolayer/vyolayer/pkg/errors"
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

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

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
		level := fmt.Sprint(entry["level"])
		message := fmt.Sprint(entry["message"])

		prefix := fmt.Sprintf("[%s] %s", level, message)
		if isColorEnabled() {
			prefix = fmt.Sprintf("%s[%s]%s %s%s%s",
				levelColor(level), level, colorReset,
				colorBlue, message, colorReset,
			)
		}

		// Pretty print in development
		jsonBytes, _ := json.MarshalIndent(entry, "", "  ")
		fmt.Printf("%s\n%s\n", prefix, string(jsonBytes))
	} else {
		// Single line JSON in production
		jsonBytes, _ := json.Marshal(entry)
		fmt.Println(string(jsonBytes))
	}
}

func isColorEnabled() bool {
	return os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"
}

func levelColor(level string) string {
	switch strings.ToUpper(level) {
	case "INFO":
		return colorGreen
	case "WARNING":
		return colorYellow
	case "ERROR", "CRITICAL":
		return colorRed
	default:
		return colorCyan
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
