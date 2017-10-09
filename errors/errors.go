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

// errCoder discribes error with code
// error value can return error code if it implements it
type errCoder interface {
	Code() int
}

// Code returns error code if applicable or zero otherwise
func Code(err error) (code int) {
	if err == nil {
		return 0
	}
	ec, ok := err.(errCoder)
	if !ok {
		return
	}
	return ec.Code()
}

// CodeMessage returns error code and message if applicable or zero code otherwise
func CodeMessage(err error) (int, string) {
	if err == nil {
		return 0, ""
	}
	ec, ok := err.(errCoder)
	if !ok {
		return 0, err.Error()
	}
	return ec.Code(), err.Error()
}
