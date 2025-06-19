package adb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsWhitespaceTrue(t *testing.T) {
	assert.True(t, containsWhitespace(("Hello Wrold")))
}

func TestContainsWhitespaceFalse(t *testing.T) {
	assert.True(t, containsWhitespace(("Hello")))
}

func TestIsBlankJustWhitespace(t *testing.T) {
	assert.True(t, isBlank(" \t"))
}

func TestIsBlankFalse(t *testing.T) {
	assert.False(t, isBlank("   h   "))
}
