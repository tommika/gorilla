// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"errors"
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestSafeErrStr(t *testing.T) {
	expected := "Wow"
	assert.Equal(t, expected, SafeErrStr(errors.New(expected)))
	assert.Equal(t, "", SafeErrStr(nil))
}

func TestOptionalString(t *testing.T) {
	val := "Good"
	def := "Golly"
	assert.Equal(t, val, OptionalString(val, def))
	assert.Equal(t, def, OptionalString("", def))
}

func TestNilOrEmpty(t *testing.T) {
	notEmpty := "not empty"
	empty := ""
	assert.False(t, NilOrEmptyString(&notEmpty))
	assert.True(t, NilOrEmptyString(&empty))
	assert.True(t, NilOrEmptyString(nil))
}

func TestSplitAndTrim(t *testing.T) {
	input := " this ,is ,  a,,list"
	split := SplitAndTrim(input, ",")
	assert.Equal(t, 5, len(split))
	assert.Equal(t, "this", split[0])
	assert.Equal(t, "is", split[1])
	assert.Equal(t, "a", split[2])
	assert.Equal(t, "", split[3])
	assert.Equal(t, "list", split[4])
}

func TestBoolStrings(t *testing.T) {
	for _, v := range trueVals {
		assert.True(t, IsBool(v))
		assert.True(t, IsTrue(v))
		assert.False(t, IsFalse(v))
	}
	for _, v := range falseVals {
		assert.True(t, IsBool(v))
		assert.False(t, IsTrue(v))
		assert.True(t, IsFalse(v))
	}
}
