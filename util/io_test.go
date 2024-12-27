// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestReadFileAsWords(t *testing.T) {
	words, err := ReadFileAsWords("../algorithms/test-data/words.txt")
	assert.Nil(t, err)
	assert.Equal(t, 93398, len(words))

	words, err = ReadFileAsWords("bogus.txt")
	assert.NotNil(t, err)
	assert.Nil(t, words)
}

func stringReader(s string) io.Reader {
	return bytes.NewReader([]byte(s))
}

const input1 = `
Hello
There.
This is a long line`

func TestReadLine(t *testing.T) {
	buff := make([]byte, 8)
	in := stringReader(input1)
	var err error
	var line string
	for err == nil {
		line, err = ReadLineBuffered(in, buff)
		if err == nil {
			t.Log(strings.TrimSpace(line))
		}
	}
	assert.Equal(t, 0, len(line))
	assert.True(t, errors.Is(err, io.EOF))
	line, err = ReadLineBuffered(in, buff)
	assert.True(t, errors.Is(err, io.EOF))
	assert.Equal(t, 0, len(line))
}
