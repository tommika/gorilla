// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package types

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestNothing(t *testing.T) {
	zero := Unit{}
	assert.Equal(t, Nothing, zero)
}
