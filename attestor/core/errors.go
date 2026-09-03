package core

import (
	"errors"
	"fmt"
)

// ErrorKind is the stable domain classification consumed by the shared gRPC
// adapter. Plugins may retain native errors for logs and metrics by wrapping
// them in Error.
type ErrorKind uint8

const (
	ErrorInvalidArgument ErrorKind = iota + 1
	ErrorNotFound
	ErrorFailedPrecondition
	ErrorUnavailable
	ErrorInternal
)

// Error carries a stable class without erasing the underlying error.
type Error struct {
	Kind ErrorKind
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err == nil {
		return e.Msg
	}
	if e.Msg == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

// NewError creates a classified error. A nil cause is valid for pure input or
// route errors.
func NewError(kind ErrorKind, message string, cause error) error {
	return &Error{Kind: kind, Msg: message, Err: cause}
}

// ErrorKindOf returns the stable class of err, defaulting to Internal so a
// new plugin cannot accidentally expose an unclassified upstream error as a
// successful or retryable response.
func ErrorKindOf(err error) ErrorKind {
	var classified *Error
	if errors.As(err, &classified) && classified != nil {
		return classified.Kind
	}
	return ErrorInternal
}
