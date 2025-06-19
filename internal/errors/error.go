package errors

import (
	"bytes"
	"fmt"
)

type Err struct {
	Code    ErrCode
	Message string
	Detail  interface{}
	Cause   error
}

var _ error = &Err{}

type ErrCode byte

const (
	AssertionErr ErrCode = iota
	ParseErr
	// The server was not avaliable on the requested port
	ServerNotAvaliable
	// General network error commuicating with the server
	NetworkError
	// The connection to the server was reset in the middle of an operation
	ConnectionResetError
	// The server returned an error message, but we couldn't parse it
	AdbError
	// The server returned a "device not found" error
	DeviceNotFound
	// Tried to perform an operation on a path that doesn't exist on the device
	FileNotExistError
)

func Errorf(code ErrCode, format string, args ...interface{}) error {
	return &Err{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func WrapErrf(cause error, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}

	err := cause.(*Err)
	return &Err{
		Code:    err.Code,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

func CombineErrs(msg string, code ErrCode, errs ...error) error {
	var nonNilErrs []error

	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	switch len(nonNilErrs) {
	case 0:
		return nil
	case 1:
		return nonNilErrs[0]
	default:
		return WrapErrorf(multiError(errs), code, "%s", msg)
	}
}

type multiError []error

func (errs multiError) Error() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "%d errors:[", len(errs))
	for i, err := range errs {
		buffer.WriteString(err.Error())
		if i < len(errs)-1 {
			buffer.WriteString(" U ")
		}
	}
	buffer.WriteRune(']')

	return buffer.String()
}

func WrapErrorf(cause error, code ErrCode, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}

	return &Err{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

func AssertionErrorf(format string, args ...interface{}) error {
	return &Err{
		Code:    AssertionErr,
		Message: fmt.Sprintf(format, args...),
	}
}

func (err *Err) Error() string {
	msg := fmt.Sprintf("%s: %s", err.Code, err.Message)
	if err.Detail != nil {
		msg = fmt.Sprintf("%s (%+v)", msg, err.Detail)
	}

	return msg
}

func HasErrCode(err error, code ErrCode) bool {
	switch err := err.(type) {
	case *Err:
		return err.Code == code
	default:
		return false
	}
}

/*
ErrorWithCauseChain formats err and all its cause if it's an *Err, else returns err.Error()
*/
func ErrorWithCauseChain(err error) string {
	var buffer bytes.Buffer

	for {
		if wrappedErr, ok := err.(*Err); ok && wrappedErr.Cause != nil {
			fmt.Fprintln(&buffer, wrappedErr.Error())
			fmt.Fprintln(&buffer, "cause by ")
			err = wrappedErr.Cause
		} else {
			break
		}
	}

	if err != nil {
		buffer.WriteString(err.Error())
	} else {
		buffer.WriteString("<err=nil>")
	}

	return buffer.String()
}
