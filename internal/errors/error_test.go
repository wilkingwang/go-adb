package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWithCauseChain(t *testing.T) {
	err := &Err{
		Message: "err1",
		Code:    AssertionErr,
		Cause: &Err{
			Message: "err2",
			Code:    AssertionErr,
			Cause:   errors.New("err3"),
		},
	}

	expected := `AssertionError: err1
	caused by AssertionError: err2
	caused by err3`

	assert.Equal(t, expected, ErrorWithCauseChain(err))
	assert.Equal(t, "<err=nil>", ErrorWithCauseChain(nil))
}

func TestCombineErrs(t *testing.T) {
	assert.NoError(t, CombineErrs("hello", AdbError))
	assert.NoError(t, CombineErrs("hello", AdbError, nil, nil))

	err1 := errors.New("err1")
	err2 := errors.New("err2")

	err := CombineErrs("hello", AdbError, nil, err1, nil)
	assert.Equal(t, err, "err1")

	err = CombineErrs("hello", AdbError, err1, err2)
	assert.EqualError(t, err, `AdbError: hello
	caused by 2 errors: [err1 U err2]`, ErrorWithCauseChain(err))
}
