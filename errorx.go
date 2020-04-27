// Package errorx provides frequently used error constructs in API development, with context
package errorx

import (
	"fmt"
	"io"

	stderrors "errors"

	exterrors "emperror.dev/errors"
	"github.com/jinzhu/copier"
)

// E defines common error information for inspecting
// and displaying to various format
type E struct {
	e       error
	Code    string `json:"code"`
	Message string `json:"message"`
}

// New returns stub error E
func New(msg string) *E {
	return &E{
		e: stderrors.New(msg),
	}
}

func (e E) Error() string {
	if e.Message != "" {
		return e.Message + ": " + e.e.Error()
	}
	if e.Code != "" {
		return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
	}
	return e.e.Error()
}

type withFormat interface {
	Format(fmt.State, rune)
}

// Format calls wrapped error with Format() of its own
// if wrapped error is not nil
func (e E) Format(s fmt.State, verb rune) {
	withFormatError, ok := e.e.(withFormat)
	if ok {
		withFormatError.Format(s, verb)
	} else {
		io.WriteString(s, e.e.Error())
	}
}

func wrap(err error, e *E) *E {
	switch typedE := err.(type) {
	case E:
		copier.Copy(e, &typedE)
	case *E:
		copier.Copy(e, typedE)
	default:
		return e
	}
	return e
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) *E {
	return Wrapf(err, msg)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) *E {
	if err == nil {
		return nil
	}
	e := wrap(err, &E{})
	e.e = exterrors.Wrapf(err, format, args...)
	return e
}

// Unwrap returns underlying error: 1st level of nested errors
func (e E) Unwrap() error {
	return stderrors.Unwrap(stderrors.Unwrap(e.e))
}
