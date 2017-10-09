package errors

import "fmt"

// New returns an error with the supplied message.
func New(message string) error {
	return &fundamental{
		msg: message,
	}
}

// NewWithCode returns an error with the supplied message and code.
func NewWithCode(message string, code int) error {
	return &fundamental{
		cod: code,
		msg: message,
	}
}

// fundamental is an error that has a message and a code, but no caller.
type fundamental struct {
	cod int
	msg string
}

// Error implements std error interface
func (f *fundamental) Error() string { return f.msg }

// Code returns an integer error code
func (f *fundamental) Code() int { return f.cod }

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(code int, format string, args ...interface{}) error {
	return &fundamental{
		cod: code,
		msg: fmt.Sprintf(format, args...),
	}
}

//-----------------------------------------------------------------------------
// withCode incapsulates 'code' with ordinary error
type withCode struct {
	cod int
	err error
}

// withCode wraps code over existing error
func (w *withCode) Error() string { return fmt.Sprintf("%d : %s", w.cod, w.Error()) }

// WrapCode wraps error with code
func WrapCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &withCode{
		cod: code,
		err: err,
	}
}

// Code returns error code if applicable or zero otherwise
func Code(err error) (code int) {
	type coderer interface {
		Code() int
	}

	for err != nil {
		coderr, ok := err.(coderer)
		if !ok {
			break
		}
		code = coderr.Code()
	}
	return
}
