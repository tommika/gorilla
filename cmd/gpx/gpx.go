// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tommika/gorilla/geo/gpx"
	"github.com/tommika/gorilla/must"
	"github.com/tommika/gorilla/xflags"
)

type Options struct {
	Debug bool     `flag:"d|debug,Enable debug logging"`
	Files []string `flag:"*,<file.gpx> ..."`
}

var exit = os.Exit

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	opts := Options{}
	if err := xflags.ParseArgs(os.Args[0], &opts, os.Args[1:]); err != nil {
		ec := 1
		if err == flag.ErrHelp {
			ec = 0
		}
		exit(ec)
		return // needed for testing
	}
	if len(opts.Files) == 0 {
		xflags.Usage(os.Stdout, os.Args[0], &opts)
		exit(1)
		return // needed for testing
	}

	if opts.Debug {
		fmt.Fprintf(os.Stderr, "options: %+v\n", opts)
	}
	errors := 0

	for _, file := range opts.Files {
		if doc, err := gpx.ReadGpxDocument(file); err != nil {
			errors += 1
			fmt.Fprintf(os.Stderr, "error %s\n", err)
		} else {
			summary := must.NotBeAnError(json.MarshalIndent(doc.Summarize(), "", "  "))
			fmt.Print(string(summary))
		}
	}
	if errors > 0 {
		exit(2)
		return // needed for testing
	}
}
