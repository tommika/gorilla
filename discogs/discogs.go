// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
// This package provides utilities for working with the
// [Discogs](https://www.discogs.com/) database.
package discogs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tommika/gorilla/algorithms/graph"
	"github.com/tommika/gorilla/algorithms/heap"
	"github.com/tommika/gorilla/xxml"
)

type ArtistId = uint32

type Artist struct {
	Id      ArtistId
	Name    string
	Members MemberList
	Groups  GroupList
}

type MemberList struct {
	Ids   []ArtistId `x:"id"`
	Names []Name     `x:"name"`
}

type GroupList struct {
	Names []Name `x:"name"`
}

type Name struct {
	Id   ArtistId `x:",attr"`
	Name string   `x:",cdata"`
}

type Discogs struct {
	aG     *graph.Graph[ArtistId, uint8]
	aMap   map[ArtistId]string
	aNames []Name
}

func NewDiscogs() *Discogs {
	return &Discogs{
		aG:     graph.NewGraph[ArtistId, uint8](false),
		aMap:   map[ArtistId]string{},
		aNames: []Name{},
	}
}

func (d *Discogs) ArtistCount() int {
	return len(d.aMap)
}

func (d *Discogs) ArtistName(id ArtistId) string {
	return d.aMap[id]
}

func (d *Discogs) ArtistEdgeCount() int {
	return d.aG.EdgeCount()
}

func (d *Discogs) ArtistNodeCount() int {
	return d.aG.NodeCount()
}

func (d *Discogs) RelatedArtists(from ArtistId, maxDepth int, cb func(id ArtistId)) {
	bft := d.aG.BFS(from, maxDepth)
	bft.VisitNodes(func(id ArtistId) {
		cb(id)
	})
}

func (d *Discogs) PathBetweenArtists(from, to ArtistId, maxDepth int) []ArtistId {
	bft := d.aG.BFS(from, maxDepth)
	path, _ := bft.FindPath(to)
	return path
}

func (d *Discogs) ImportArtists(fileName string) error {
	errors := 0
	if in, err := os.Open(fileName); err != nil {
		log.Printf("error: %s\n", err)
		errors += 1
	} else {
		defer in.Close()
		i := 0
		maxLen := 0
		var maxId ArtistId
		ReadArtists(in, func(a *Artist) {
			d.aMap[a.Id] = a.Name
			for _, gn := range a.Groups.Names {
				d.aG.AddEdge(a.Id, gn.Id, 1)
			}
			if i%10000 == 0 {
				s := fmt.Sprintf("%d: %s", i, a.Name)
				if len(s) > maxLen {
					maxLen = len(s)
				}
				fmt.Printf("\r%-*s", maxLen, s)
			}
			i += 1
			if a.Id > maxId {
				maxId = a.Id
			}
		})
		fmt.Printf("\r%*s%c\n", maxLen, "", 0)
		fmt.Printf("maxId=%d\n", maxId)
	}
	var err error
	if errors > 0 {
		err = fmt.Errorf("%d errors encountered during import", errors)
	} else {
		d.aNames = []Name{}
		for id, name := range d.aMap {
			d.aNames = append(d.aNames, Name{id, name})
		}
		fmt.Printf("sorting artist names...\n")
		heap.HeapSort(d.aNames, func(a, b Name) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	return err
}

// ReadArtists reads a Discogs artists dump (XML) from the given input stream,
// and for each artist in the stream, calls the given func.
func ReadArtists(in io.Reader, artistFunc func(a *Artist)) error {
	type artistsDoc struct {
		Artist func(*Artist)
	}
	doc := artistsDoc{}
	doc.Artist = artistFunc
	return xxml.ReadXmlWithRootName(in, &doc, xxml.XmlName("", "artists"))
}
