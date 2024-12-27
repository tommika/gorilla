// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package resources

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestReadResourceAsString(t *testing.T) {
	_, err := ReadResourceAsString("README.md")
	assert.Nil(t, err)
}

func TestReadResourceNotFound(t *testing.T) {
	_, err := ReadResourceAsString("bogus")
	assert.NotNil(t, err)
}

func TestOpenResourceNotFound(t *testing.T) {
	_, err := OpenResource("bogus")
	assert.NotNil(t, err)
}
