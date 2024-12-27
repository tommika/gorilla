// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tommika/gorilla/discogs"
	"github.com/tommika/gorilla/xflags"
)

type Options struct {
	Debug        bool     `flag:"d|debug,Enable debug logging"`
	ArtistId     uint64   `flag:"a|artist,Artist id to search"`
	PathToId     uint64   `flag:"p|path-to,Find the path to this artist"`
	RelatedDepth int      `flag:"r|related-depth,Show related artists to this depth"`
	Files        []string `flag:"*,<artists.xml> ..."` // artists to import
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
		xflags.Usage(os.Stderr, os.Args[0], &opts)
		exit(1)
		return // needed for testing
	}
	errors := 0
	d := discogs.NewDiscogs()
	for _, fn := range opts.Files {
		if err := d.ImportArtists(fn); err != nil {
			errors += 1
		}
	}

	fmt.Printf("#artists: %d\n", d.ArtistCount())
	fmt.Printf("#nodes  : %d\n", d.ArtistNodeCount())
	fmt.Printf("#edges  : %d\n", d.ArtistEdgeCount())

	artistId := discogs.ArtistId(opts.ArtistId)
	pathToId := discogs.ArtistId(opts.PathToId)
	if opts.ArtistId != 0 {
		name := d.ArtistName(artistId)
		if len(name) == 0 {
			fmt.Fprintf(os.Stderr, "artist not found: id=%d\n", opts.ArtistId)
			errors += 1
		} else {
			fmt.Fprintf(os.Stderr, "found artist: %s [%d]\n", name, opts.ArtistId)
		}
		if opts.RelatedDepth > 0 {
			fmt.Fprintf(os.Stderr, "related artists (depth==%d):\n", opts.RelatedDepth)
			count := 0
			d.RelatedArtists(artistId, opts.RelatedDepth, func(id discogs.ArtistId) {
				count++
				fmt.Fprintf(os.Stderr, "\t[%d] %s\n", id, d.ArtistName(id))
			})
			fmt.Fprintf(os.Stderr, "found %d related artists\n", count)
		}
		if opts.PathToId != 0 {
			path := d.PathBetweenArtists(artistId, pathToId, 0)
			if len(path) == 0 {
				fmt.Fprintf(os.Stderr, "No path found to %s [%d]\n", d.ArtistName(pathToId), pathToId)
			} else {
				fmt.Fprintf(os.Stderr, "Path to [%d]%s:\n", opts.PathToId, d.ArtistName(pathToId))
				for _, id := range path {
					fmt.Fprintf(os.Stderr, "\t[%d]%s\n", id, d.ArtistName(id))
				}
			}
		}
	}
	if errors > 0 {
		exit(2)
		return // needed for testing
	}
}
