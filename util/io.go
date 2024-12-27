// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// ReadWords
func ReadFileAsWords(path string) (words []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}
	err = scanner.Err()
	return
}

// ReadLineBuffered reads a line of input from the given reader.
// If an EOL is encountered, ReadLineBuffered returns ("", io.EOL)
func ReadLineBuffered(in io.Reader, buff []byte) (string, error) {
	n := 0
	for n < len(buff) {
		m, err := in.Read(buff[n : n+1])
		if err != nil || buff[n] == '\n' {
			if errors.Is(err, io.EOF) && n > 0 {
				// delay returning EOF
				err = nil
			}
			return string(buff[:n+m]), err
		}
		n += m
	}
	return string(buff[:n]), nil
}
