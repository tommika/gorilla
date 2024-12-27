// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestCurrentFunc(t *testing.T) {
	p, n := CallerFuncName(0)
	assert.Equal(t, "github.com/tommika/gorilla/util", p)
	assert.Equal(t, "TestCurrentFunc", n)
}

func TestParseFuncName(t *testing.T) {
	data := [][]string{
		{"example.com/my/pack.MyFunc", "example.com/my/pack", "MyFunc"},
		{"example.com/MyFunc", "example.com", "MyFunc"},
		{"pack.MyFunc", "pack", "MyFunc"},
		{"main/MyFunc", "main", "MyFunc"},
		{"example.com/my/pack.init.0", "example.com/my/pack", "init.0"},
	}
	for _, d := range data {
		pkg, name := ParseFuncName(d[0])
		assert.Equal(t, d[1], pkg)
		assert.Equal(t, d[2], name)
	}
}
