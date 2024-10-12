// Package gonuts provides enhanced error handling functionality
// with structured error messages, error codes, context, stack trace support, and logging integration.
package gonuts

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ErrorPlus represents an extended error type that includes an error message,
// an error code, additional context, stack trace, and the original wrapped error.
// It is designed to be immutable and thread-safe.
// Example usage of ErrorPlus with context and logging.
//
//	func example() {
//		err := errors.New("database connection failed")
//		errPlus := NewInternalError("Unable to reach the database", err).
//			WithContext("userID", 1234).
//			WithContext("retry", true)
//		errPlus.Log()
//	}
type ErrorPlus struct {
	err        error                  // Original wrapped error
	msg        string                 // Error message with context
	code       int                    // Error code, can be HTTP code or custom code
	context    map[string]interface{} // Additional contextual information
	stackTrace []string               // Stack trace at the point of error creation
	timestamp  time.Time              // Time when the error was created
}

// NewErrorPlus creates a new ErrorPlus instance by wrapping an error with a custom message and code.
// It captures the stack trace at the point of creation.
func NewErrorPlus(err error, msg string, code int) *ErrorPlus {
	return &ErrorPlus{
		err:        err,
		msg:        msg,
		code:       code,
		context:    make(map[string]interface{}),
		stackTrace: captureStackTrace(),
		timestamp:  time.Now(),
	}
}

// Error implements the error interface, returning the error message including context and original error.
func (e *ErrorPlus) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// ErrorMsg returns the error message without the original error.
func (e *ErrorPlus) Msg() string {
	return e.msg
}

// ErrorStr returns the error message as a string.
func (e *ErrorPlus) ContextStr() string {
	return GetPrettyJson(e.context)
}

// Context returns the additional context associated with the ErrorPlus.
func (e *ErrorPlus) Context() map[string]interface{} {
	return e.context
}

// CodeStr returns the error code as a string.
func (e *ErrorPlus) CodeStr() string {
	return fmt.Sprintf("%d", e.code)
}

// Code returns the error code associated with the ErrorPlus.
func (e *ErrorPlus) Code() int {
	return e.code
}

// Unwrap returns the underlying original error, allowing access to the original error for further handling.
func (e *ErrorPlus) Unwrap() error {
	return e.err
}

// Is compares the ErrorPlus with another target error, checking if the underlying error matches.
func (e *ErrorPlus) Is(target error) bool {
	return errors.Is(e.err, target)
}

// As attempts to map the ErrorPlus to a target error type, useful for type assertion.
func (e *ErrorPlus) As(target interface{}) bool {
	return errors.As(e.err, target)
}

// StackTrace returns the stack trace associated with the error.
func (e *ErrorPlus) StackTrace() []string {
	return e.stackTrace
}

// Timestamp returns the time when the error was created.
func (e *ErrorPlus) Timestamp() time.Time {
	return e.timestamp
}

// WithMsg returns a new ErrorPlus with the provided message, preserving immutability.
func (e *ErrorPlus) WithMsg(msg string) *ErrorPlus {
	return &ErrorPlus{
		err:        e.err,
		msg:        msg,
		code:       e.code,
		context:    copyContext(e.context),
		stackTrace: e.stackTrace, // Stack trace remains the same
		timestamp:  e.timestamp,
	}
}

// WithCode returns a new ErrorPlus with the provided code, preserving immutability.
func (e *ErrorPlus) WithCode(code int) *ErrorPlus {
	return &ErrorPlus{
		err:        e.err,
		msg:        e.msg,
		code:       code,
		context:    copyContext(e.context),
		stackTrace: e.stackTrace,
		timestamp:  e.timestamp,
	}
}

// WithContext returns a new ErrorPlus with the additional context, preserving immutability.
func (e *ErrorPlus) WithContext(key string, value interface{}) *ErrorPlus {
	newContext := copyContext(e.context)
	newContext[key] = value
	return &ErrorPlus{
		err:        e.err,
		msg:        e.msg,
		code:       e.code,
		context:    newContext,
		stackTrace: e.stackTrace,
		timestamp:  e.timestamp,
	}
}

// WithValues returns a new ErrorPlus with the provided message and code, preserving immutability.
func (e *ErrorPlus) WithValues(msg string, code int) *ErrorPlus {
	return &ErrorPlus{
		err:        e.err,
		msg:        msg,
		code:       code,
		context:    copyContext(e.context),
		stackTrace: e.stackTrace,
		timestamp:  e.timestamp,
	}
}

// MarshalJSON implements the json.Marshaler interface, allowing custom JSON serialization.
func (e *ErrorPlus) MarshalJSON() ([]byte, error) {
	type Alias ErrorPlus
	return json.Marshal(&struct {
		*Alias
		Error     string    `json:"error"`
		Timestamp time.Time `json:"timestamp"`
	}{
		Alias:     (*Alias)(e),
		Error:     e.err.Error(),
		Timestamp: e.timestamp,
	})
}

// Format implements the fmt.Formatter interface for custom formatting.
func (e *ErrorPlus) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "ErrorPlus:\n  Msg: %s\n  Code: %d\n  Error: %+v\n  Context: %v\n  StackTrace:\n%s", e.msg, e.code, e.err, e.context, strings.Join(e.stackTrace, "\n"))
		} else {
			fmt.Fprintf(f, "%s: %v", e.msg, e.err)
		}
	default:
		fmt.Fprintf(f, "%s: %v", e.msg, e.err)
	}
}

// captureStackTrace captures the current stack trace.
func captureStackTrace() []string {
	const maxFrames = 32
	pcs := make([]uintptr, maxFrames)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	var stackTrace []string
	for {
		frame, more := frames.Next()
		stackTrace = append(stackTrace, fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return stackTrace
}

// copyContext makes a deep copy of the context map.
func copyContext(original map[string]interface{}) map[string]interface{} {
	if original == nil {
		return nil
	}
	newContext := make(map[string]interface{})
	for k, v := range original {
		newContext[k] = v
	}
	return newContext
}

// Helper functions for common HTTP errors.

// NewNotFoundError creates a new ErrorPlus representing a 404 Not Found error.
func NewNotFoundError(msg string, err error) *ErrorPlus {
	return NewErrorPlus(err, msg, 404)
}

// NewInternalError creates a new ErrorPlus representing a 500 Internal Server Error.
func NewInternalError(msg string, err error) *ErrorPlus {
	return NewErrorPlus(err, msg, 500)
}

// NewUnauthorizedError creates a new ErrorPlus representing a 401 Unauthorized error.
func NewUnauthorizedError(msg string, err error) *ErrorPlus {
	return NewErrorPlus(err, msg, 401)
}

// NewBadRequestError creates a new ErrorPlus representing a 400 Bad Request error.
func NewBadRequestError(msg string, err error) *ErrorPlus {
	return NewErrorPlus(err, msg, 400)
}

// WithLogger allows setting a custom logger. By default, uses gonuts.L (the package's default logger).
var errorLogger = L

// SetErrorLogger allows injecting a custom logger for ErrorPlus instances.
func SetErrorLogger(logger *zap.SugaredLogger) {
	errorLogger = logger
}

// Log logs the error using the configured logger.
func (e *ErrorPlus) Log() {
	errorLogger.Errorf("%+v", e)
}

// Ensure ErrorPlus satisfies the standard library interfaces.
var _ interface {
	error
	fmt.Formatter
	json.Marshaler
	// errors.Wrapper
} = (*ErrorPlus)(nil)
