// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xxml

// Examples used in README.md

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

var artistsXml = `
<?xml version="1.0"?>
<artists>
    <artist id="2112">
        <name>Rush</name>
        <description xml:space="preserve">
           Rush was a progressive rock band formed in Toronto Canada in 1968.
		   
		   Rush were known for their virtuosic musicianship, complex
		   compositions and eclectic lyrical motifs, which drew primarily on
		   science fiction, fantasy and philosophy 
		</description>
        <member id="2113">John Rutsey</member> <!-- original drummer -->
        <member id="2114">Neal Peart</member>
        <member id="2115">Alex Lifeson</member>
        <member id="2116">Geddy Lee</member>
    </artist>
</artists>
`

type ArtistName struct {
	Id   string `x:",attr"` // attribute
	Name string `x:",cdata"`
}

type Artist struct {
	Id          string `x:",attr"` // attribute
	Name        string
	Description string
	Members     []ArtistName `x:"member"`
}

type ArtistsDoc struct {
	Artists []Artist `x:"artist"` // sequence of elements
}

func TestArtistsExample(t *testing.T) {
	artists := ArtistsDoc{}
	err := ReadXmlWithRootName(stringReader(artistsXml), &artists, XmlName("", "artists"))
	assert.Nil(t, err)
	logJson(t, artists)
}

// Use a callback rather than having to keep all of the artists in memory
type ArtistsCallback struct {
	Artist func(a *Artist) // artist sequence callback
}

func TestArtistsCallbackExample(t *testing.T) {
	artists := ArtistsCallback{}
	artists.Artist = func(a *Artist) {
		logJson(t, a)
	}
	err := ReadXmlWithRootName(stringReader(artistsXml), &artists, XmlName("", "artists"))
	assert.Nil(t, err)
}

var namespaceXml = `
<?xml version="1.0"?>
<doc 
		xmlns="example.com/ns1" 
		xmlns:ns2="example.com/ns2" 
		xmlns:ns3="example.com/ns3"
		ns3:a1="wow">
	<e1>Hello</e1>
	<ns2:e2>Goodbye</ns2:e2>
</doc>
`

type DocWithNamespace struct {
	_  struct{} `x:"example.com/ns1"`       // specify default namespace
	A1 string   `x:",attr,example.com/ns3"` // uses ns3
	E1 string   // uses default namespace
	E2 string   `x:",,example.com/ns2"` // uses ns2
}

func TestNamespaceExample(t *testing.T) {
	doc := DocWithNamespace{}
	err := ReadXmlWithRootName(stringReader(namespaceXml), &doc, XmlName("example.com/ns1", "doc"))
	assert.Nil(t, err)
	logJson(t, doc)
}
