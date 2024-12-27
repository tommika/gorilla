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
		"discogs",
		"--artist=316130",
		"--path-to=229492",
		"--related-depth=4",
		"../../discogs/test-data/discogs-artists-xtc.xml",
	}
	main()
	assert.Equal(t, 0, exitCode)
}

func TestDebugEnabled(t *testing.T) {
	defer setup()()
	os.Args = []string{
		"discogs",
		"--debug",
		"../../discogs/test-data/discogs-artists-xtc.xml",
	}
	main()
	assert.Equal(t, 0, exitCode)
}

func TestInvalidArgs(t *testing.T) {
	defer setup()()
	os.Args = []string{"discogs", "--bogus"}
	main()
	assert.Equal(t, 1, exitCode)
}

func TestFileNotFound(t *testing.T) {
	defer setup()()
	os.Args = []string{"discogs", "--debug", "bogus.xml"}
	main()
	assert.Equal(t, 2, exitCode)
}

func TestArtistNotFound(t *testing.T) {
	defer setup()()
	os.Args = []string{
		"discogs",
		"--artist=1111111",
		"../../test-data/discogs-artists-xtc.xml",
	}
	main()
	assert.Equal(t, 2, exitCode)
}
