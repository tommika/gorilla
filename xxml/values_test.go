// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xxml

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestValidDateTimeValues(t *testing.T) {
	valid := []string{
		"12:30:00",
		"12:30:00-05:00",
		"12:30:00Z",
		"2024-11-04",
		"2024-11-04-05:00",
		"2024-11-04Z",
		"2024-11-04T12:30:00",
		"2024-11-04T12:30:00-05:00",
		"2024-11-04T12:30:00Z",
	}
	for _, val := range valid {
		dt, err := parseDateTime(val)
		assert.Nil(t, err)
		t.Logf("dt=%s\n", dt)
	}
}

func TestInvalidDateTimeValues(t *testing.T) {
	invalid := []string{
		"12:30:0",
		"2024-11T12:30:00",
	}
	for _, val := range invalid {
		_, err := parseDateTime(val)
		assert.NotNil(t, err)
		t.Log(err)
	}
}
