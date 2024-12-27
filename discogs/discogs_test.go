// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package discogs

import (
	"fmt"
	"testing"

	"github.com/tommika/gorilla/assert"
)

func Test(t *testing.T) {
	d := NewDiscogs()
	err := d.ImportArtists("./test-data/discogs-artists-xtc.xml")
	assert.Nil(t, err)
	assert.Equal(t, 3, d.ArtistCount())
	assert.Equal(t, "XTC", d.ArtistName(15118))
	fmt.Printf("#artists: %d\n", d.ArtistCount())
	fmt.Printf("#nodes  : %d\n", d.ArtistNodeCount())
	fmt.Printf("#edges  : %d\n", d.ArtistEdgeCount())
	fmt.Printf("related artists\n")
	d.RelatedArtists(316130, 3, func(id ArtistId) {
		fmt.Printf("\t%s[%d]\n", d.ArtistName(id), id)
	})
	path := d.PathBetweenArtists(316130, 229492, 0)
	assert.Equal(t, 3, len(path))
	fmt.Printf("path\n")
	for _, id := range path {
		fmt.Printf("\t%s\n", d.ArtistName(id))
	}

}

func TestFileNotFound(t *testing.T) {
	d := NewDiscogs()
	err := d.ImportArtists("../bogus")
	assert.NotNil(t, err)
}
