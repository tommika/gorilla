// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package main

import (
	"os"
	"testing"

	"github.com/tommika/gorilla/assert"
)

func setup() func() {
	exitCode = 0
	exitSav := exit
	// fake-out exit
	exit = func(code int) { exitCode = code }
	return func() {
		exit = exitSav
		exitCode = -1
	}
}

var exitCode int = -1

func Test(t *testing.T) {
	defer setup()()
	os.Args = []string{
		"gpx",
		"../.././geo/gpx/test-data/BeaconToColdSpring.gpx",
	}
	main()
	assert.Equal(t, 0, exitCode)
}

func TestDebugEnabled(t *testing.T) {
	defer setup()()
	os.Args = []string{
		"gpx",
		"--debug",
		"../.././geo/gpx/test-data/BeaconToColdSpring.gpx",
	}
	main()
	assert.Equal(t, 0, exitCode)
}

func TestInvalidArgs(t *testing.T) {
	defer setup()()
	os.Args = []string{"gpx", "--bogus"}
	main()
	assert.Equal(t, 1, exitCode)
}

func TestFileNotFound(t *testing.T) {
	defer setup()()
	os.Args = []string{"gpx", "--debug", "bogus.xml"}
	main()
	assert.Equal(t, 2, exitCode)
}

func TestHelp(t *testing.T) {
	defer setup()()
	os.Args = []string{"gpx", "--help"}
	main()
	assert.Equal(t, 0, exitCode)
}

func TestUsage(t *testing.T) {
	defer setup()()
	os.Args = []string{"gpx"}
	main()
	assert.Equal(t, 1, exitCode)
}
