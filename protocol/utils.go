package protocol

import (
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/wilkingwang/go-adb/internal/errors"
)

// ErrorResponseDetails is an error message returned by the adb server
type ErrorResponseDetails struct {
	Request   string
	ServerMsg string
}

// deviceNotFoundMessagePattern match all device not found error message returned by adb server
var deviceNotFoundMessagePattern = regexp.MustCompile(`device( '.*')? not found`)

func adbServerError(request string, serverMsg string) error {
	var msg string
	if request == "" {
		msg = fmt.Sprintf("server error: %s", serverMsg)
	} else {
		msg = fmt.Sprintf("server error for %s request: %s", request, serverMsg)
	}

	errCode := errors.AdbError
	if deviceNotFoundMessagePattern.MatchString(serverMsg) {
		errCode = errors.DeviceNotFound
	}

	return &errors.Err{
		Code:    errCode,
		Message: msg,
		Detail: ErrorResponseDetails{
			Request:   request,
			ServerMsg: serverMsg,
		},
	}
}

// IsAdbServerErrorMatching return true if err is an *Err with code AdbError and for which predict returns true when passed Details.ServerMsg
func IsAdbServerErrorMatching(err error, predicate func(string) bool) bool {
	if err, ok := err.(*errors.Err); ok && err.Code == errors.AdbError {
		return predicate(err.Detail.(ErrorResponseDetails).ServerMsg)
	}

	return false
}

func errInCompleteMessage(desc string, actual int, expected int) error {
	return &errors.Err{
		Code:    errors.ConnectionResetError,
		Message: fmt.Sprintf("incomplete %s: read %d bytes, expecting %d", desc, actual, expected),
		Detail: struct {
			ActualReadBytes int
			ExpectedBytes   int
		}{
			ActualReadBytes: actual,
			ExpectedBytes:   expected,
		},
	}
}

// writeAll writes all of data to w
func writeAll(w io.Writer, data []byte) error {
	offset := 0

	for offset < len(data) {
		n, err := w.Write(data[offset:])
		if err != nil {
			return errors.WrapErrorf(err, errors.NetworkError, "error writing %d bytes at offset %d", len(data), offset)
		}

		offset += n
	}

	return nil
}

// MultiCloseable wraps c in a ReadWriteCloser can by safely closed muliple times
func MultiCloseable(c io.ReadWriteCloser) io.ReadWriteCloser {
	return &multiCloseable{ReadWriteCloser: c}
}

type multiCloseable struct {
	io.ReadWriteCloser
	closeOnce sync.Once
	err       error
}

func (c *multiCloseable) Close() error {
	c.closeOnce.Do(func() {
		c.err = c.ReadWriteCloser.Close()
	})

	return c.err
}
