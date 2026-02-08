package errors

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

// Severity represents the severity level of an error
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityCritical Severity = "CRITICAL"
)

// StackFrame represents a single frame in the stack trace
type StackFrame struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// AppError is the core error type for the application
type AppError struct {
	// Code is the unique error code
	Code ErrorCode `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// HTTPStatus is the HTTP status code to return
	HTTPStatus int `json:"httpStatus"`

	// Severity indicates the severity of the error
	Severity Severity `json:"severity"`

	// StackTrace contains the call stack where the error occurred
	StackTrace []StackFrame `json:"stackTrace,omitempty"`

	// Metadata contains additional context about the error
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Wrapped is the underlying error that was wrapped
	Wrapped error `json:"-"`

	// Timestamp is when the error occurred
	Timestamp time.Time `json:"timestamp"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Wrapped
}

// Is checks if the error matches the target error code
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(key string, value interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithMetadataMap adds multiple metadata entries
func (e *AppError) WithMetadataMap(metadata map[string]interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		e.Metadata[k] = v
	}
	return e
}

// GetMetadata retrieves a metadata value by key
func (e *AppError) GetMetadata(key string) (interface{}, bool) {
	if e.Metadata == nil {
		return nil, false
	}
	val, ok := e.Metadata[key]
	return val, ok
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) []StackFrame {
	const maxDepth = 32
	var frames []StackFrame

	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pcs)

	if n == 0 {
		return frames
	}

	pcs = pcs[:n]
	callersFrames := runtime.CallersFrames(pcs)

	for {
		frame, more := callersFrames.Next()
		frames = append(frames, StackFrame{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		})

		if !more {
			break
		}
	}

	return frames
}

// New creates a new AppError with the given error code
func New(code ErrorCode) *AppError {
	metadata := GetMetadata(code)

	return &AppError{
		Code:       code,
		Message:    metadata.DefaultMessage,
		HTTPStatus: metadata.HTTPStatus,
		Severity:   metadata.Severity,
		StackTrace: captureStackTrace(3),
		Metadata:   make(map[string]interface{}),
		Timestamp:  time.Now(),
	}
}

// NewWithMessage creates a new AppError with a custom message
func NewWithMessage(code ErrorCode, message string, args ...interface{}) *AppError {
	metadata := GetMetadata(code)

	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: metadata.HTTPStatus,
		Severity:   metadata.Severity,
		StackTrace: captureStackTrace(3),
		Metadata:   make(map[string]interface{}),
		Timestamp:  time.Now(),
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(err error, code ErrorCode, message string, args ...interface{}) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, return it as is
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	metadata := GetMetadata(code)

	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: metadata.HTTPStatus,
		Severity:   metadata.Severity,
		StackTrace: captureStackTrace(3),
		Metadata:   make(map[string]interface{}),
		Wrapped:    err,
		Timestamp:  time.Now(),
	}
}

// WrapWithContext wraps an error with additional context metadata
func WrapWithContext(err error, code ErrorCode, message string, metadata map[string]interface{}) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, add metadata and return
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.WithMetadataMap(metadata)
	}

	meta := GetMetadata(code)

	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: meta.HTTPStatus,
		Severity:   meta.Severity,
		StackTrace: captureStackTrace(3),
		Metadata:   metadata,
		Wrapped:    err,
		Timestamp:  time.Now(),
	}
}

// Is checks if an error has a specific error code
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// As attempts to convert an error to AppError
func As(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetCode returns the error code from an error
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ErrInternalUnexpected
}

// GetHTTPStatus returns the HTTP status code from an error
func GetHTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus
	}
	return 500
}

// GetSeverity returns the severity level from an error
func GetSeverity(err error) Severity {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Severity
	}
	return SeverityError
}
