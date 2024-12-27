<!--
Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
-->
xxml
====

The xxml package reads XML-formatted data using Go structures and custom field
tags.  


struct tags for XML processing
------------------------------

Tags are of the form:
```
   x:"<name>,<kind>,<namespace>"
```
where `kind` is one of:
* `elem` (default)
* `attr`
* `cdata`

All of these are optional, and if not specified, defaults to the following
```
   x:"<field-name>,elem,<def-namespace>"
```

For example, say we want to decode an XML document that describes a music collection:

```xml
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
        <member id="2113">John Rutsey</member>
        <member id="2114">Neal Peart</member>
        <member id="2115">Alex Lifeson</member>
        <member id="2116">Geddy Lee</member>
    </artist>
</artists>

```

We can use the following model in Go :
```go
type ArtistName struct {
	Id   string `x:",attr"` 
	Name string `x:",cdata"`
}

type Artist struct {
	Id          string       `x:",attr"`
	Name        string
	Members     []ArtistName `x:"member"`
	Description string       `x:",cdata"`
}

type ArtistsDoc struct {
	Artists []Artist `x:"artist"` // sequence of elements
}
```

We can then read this as follows:
```go
artists := ArtistsDoc{}
err := ReadXmlWithRootName(inputFile, &artists, XmlName("", "artists"))
```

Stream processing callbacks
---------------------------

Let's say we're processing a document that has a large number of elements.
Using the example above, rather than having to read the entire sequence of
artists into memory before processing, we can use a callback and process
each artist as it appears in the stream:

```go
type ArtistsCallback struct {
	Artist func(a *Artist) // artist sequence callback
}

... 
artists := ArtistsCallback{}
artists.Artist = func(a *Artist) {
	// do something with this artist
	...
}
err := ReadXmlWithRootName(inputFile, &artists, XmlName("", "artists"))

```


Namespaces
----------
The default XML namespace can be specified using the field `_`.
If no namespace is specified (at the field level or using the default)
then the XML document's current default namespace is used.


For example, given this XML,

```xml
<?xml version="1.0"?>
<doc 
		xmlns="example.com/ns1" 
		xmlns:ns2="example.com/ns2" 
		xmlns:ns3="example.com/ns3"
		ns3:a1="wow">
	<e1>Hello</e1>
	<ns2:e2>Goodbye</ns2:e2>
</doc>
```

we could have this model in Go:

```go
type DocWithNamespace struct {
	_  struct{} `x:"example.com/ns1"`       // specify default namespace
	A1 string   `x:",attr,example.com/ns3"` // uses ns3
	E1 string   // uses default namespace
	E2 string   `x:",,example.com/ns2"` // uses ns2
}
```


TODO
----

- Embedded structs

Do we really need another XML deserializer?
-------------------------------------------
It is true that Go has built-in support for deserializing XML to a struct.
However, that implementation has several limitations:
* requires reading the entire XML document into memory before decoding
* requires decoding the entire XML document before the application can start processing
* rigid namespace handling 
* ... 

My implementation does not have these limitations.  Internally, it uses Go's
pull-based XML event parser (encoding.xml.Decoder) to read data incrementally
from an XML stream, not requiring entire XML document to be in memory.
Further, it supports a novel declarative callback mechanism, allowing
applications to process the XML documents as a stream of decoded, application
specific data structures

It is also very instructive to implement such a thing, to understand the
mechanics of building a declarative model. This implementation makes heavy use
of Go reflection and custom struct tags.

References
----------
* https://www.w3.org/TR/2006/REC-xml11-20060816/
* https://www.xml.com/axml/testaxml.htm
