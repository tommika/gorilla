// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package gpx

import (
	"os"

	"github.com/tommika/gorilla/xxml"
)

// ReadGpxDocument reads a GPX document from the given named file.
func ReadGpxDocument(path string) (*Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadGpxDocumentFromStream(file)
}

// ReadGpxDocumentFromStream reads a GPX document from the given file descriptor.
func ReadGpxDocumentFromStream(in *os.File) (*Document, error) {
	gpx := Document{}
	err := xxml.ReadXmlWithRootName(in, &gpx, GPX_NAME)
	return &gpx, err
}
