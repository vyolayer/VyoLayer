package errors

import "fmt"

// ErrorBuilder provides a fluent interface for building errors
type ErrorBuilder struct {
	error *AppError
}

// NewBuilder creates a new error builder with the given error code
func NewBuilder(code ErrorCode) *ErrorBuilder {
	return &ErrorBuilder{
		error: New(code),
	}
}

// WithMessage sets a custom message for the error
func (b *ErrorBuilder) WithMessage(message string, args ...interface{}) *ErrorBuilder {
	if len(args) > 0 {
		b.error.Message = fmt.Sprintf(message, args...)
	} else {
		b.error.Message = message
	}
	return b
}

// WithHTTPStatus sets a custom HTTP status code
func (b *ErrorBuilder) WithHTTPStatus(status int) *ErrorBuilder {
	b.error.HTTPStatus = status
	return b
}

// WithSeverity sets the severity level
func (b *ErrorBuilder) WithSeverity(severity Severity) *ErrorBuilder {
	b.error.Severity = severity
	return b
}

// WithMetadata adds a single metadata entry
func (b *ErrorBuilder) WithMetadata(key string, value interface{}) *ErrorBuilder {
	if b.error.Metadata == nil {
		b.error.Metadata = make(map[string]interface{})
	}
	b.error.Metadata[key] = value
	return b
}

// WithMetadataMap adds multiple metadata entries
func (b *ErrorBuilder) WithMetadataMap(metadata map[string]interface{}) *ErrorBuilder {
	if b.error.Metadata == nil {
		b.error.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		b.error.Metadata[k] = v
	}
	return b
}

// WithWrap wraps an existing error
func (b *ErrorBuilder) WithWrap(err error) *ErrorBuilder {
	b.error.Wrapped = err
	return b
}

// WithoutStackTrace removes the stack trace from the error
func (b *ErrorBuilder) WithoutStackTrace() *ErrorBuilder {
	b.error.StackTrace = nil
	return b
}

// Build returns the constructed AppError
func (b *ErrorBuilder) Build() *AppError {
	return b.error
}

// BuildAndReturn is a convenience method that returns the error
// This allows for inline error creation: return errors.NewBuilder(code).Build()
func (b *ErrorBuilder) BuildAndReturn() error {
	return b.error
}
