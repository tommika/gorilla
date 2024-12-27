// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package resources

import (
	"embed"
	"io/fs"
	"path/filepath"
)

//go:embed static
var staticFS embed.FS

func ReadResource(rscName string) ([]byte, error) {
	path := filepath.Join("static", rscName)
	return staticFS.ReadFile(path)
}

func ReadResourceAsString(rscName string) (str string, err error) {
	b, err := ReadResource(rscName)
	if err != nil {
		return
	}
	return string(b), nil
}

func OpenResource(path string) (fs.File, error) {
	return staticFS.Open(filepath.Join("static", path))
}
