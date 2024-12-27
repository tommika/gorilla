<!--
Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License 
-->

![alt gorilla](./doc/gorilla-icon.png "Gorilla")

Gorilla
=======

Computer Science workouts, in Go!

This is intended to be a repository of Go (Golang) reference source code.  It's
where I go to figure-out how things are done, in Go.

All of the code in this repository is of my own creation, with minimal
dependency on third party packages.

I'm planning to cover these general categories:

* Algorithms and data structures
* Data Processing
* Computer Graphics 
* IT and Security
* General Utilities

Pre-reqs
--------

* Go v1.23+


Quick Start
-----------

Build and test everything
```bash
make
```

Take a look at test code coverage:
```bash
open tmp/coverage/coverage.html
```

Run performance benchmarks:
```bash
make bench
```

What you'll find inside
-----------------------

### Elementary Algorithms and Data Structures
[bst](./algorithms/bst) - Generic binary search tree data structure and
algorithms (balanced and unbalanced.)

[heap](./algorithms/heap) - Generic heap data structure and algorithms,
including heap sort and priority queue implementations.

[graph](./algorithms/graph) - Generic graph data structure and algorithms.

[queue](./algorithms/queue) - Generic queue/FIFO abstraction, with multiple
implementations: circular array, slice, and linked-list.

### Utility Packages
[xflags](./xflags) - A powerful declarative model for Go's built-in command-line
processing package (flags.)

### Data Processing Packages
[xxml](./xxml/README.md) - Decoding of XML-formatted data. Well-suited for parsing very
large files with minimal memory footprint. Includes a novel declarative event
processing mechanism. Lots of good stuff in here related to Go reflection
and custom struct tags.

[geo](./geo) - Geospatial tools including handling of [gpx](./geo/gpx) (GPS Exchange Format) data.

[discogs](./discogs/README.md) - Tools for working with the [Discogs](https://www.discogs.com/) database of music discographies.

### Misc Packages

[security](./security) - Security related code, including handling of JWT and JWKS data.

### Applications

[gpx](./cmd/gpx) - Summarize a GPX file (leverages [geo/gpx](./geo/gpx) package)
```
$  ./bin/gpx ./geo/gpx/test-data/BeaconToColdSpring.gpx
<-- output -->
{
  "NumPoints": 954,
  "StartPt": {
    "Lat": 41.493350844830275,
    "Lon": -73.96022438071668,
    "Ele": 40.993408203125,
    "Time": "2015-12-31T10:52:58Z"
  },
  "EndPt": {
    "Lat": 41.42669535242021,
    "Lon": -73.96546960808337,
    "Ele": 15.5184326171875,
    "Time": "2015-12-31T16:11:09Z"
  },
  "Extent": {
    "MinLat": 41.42669535242021,
    "MaxLat": 41.4934463147074,
    "MinLon": -73.97041861899197,
    "MaxLon": -73.94215350039303
  },
  "Duration": 19091,
  "MinElePt": {
    "Lat": 41.42682728357613,
    "Lon": -73.96564370021224,
    "Ele": 10.231201171875,
    "Time": "2015-12-31T16:09:03Z"
  },
  "MaxElePt": {
    "Lat": 41.48140596225858,
    "Lon": -73.94462255761027,
    "Ele": 494.2550048828125,
    "Time": "2015-12-31T12:26:20Z"
  },
  "EleGain": 973.334228515625,
  "Distance": 12347.400339673362,
  "GridDistance": 12369.67370356709,
  "Area": 0,
  "Loopiness": 0,
  "IsClosed": false
}
```

Notes On Go Language Features
-----------------------------

## Core
### Arrays and Slices
Nice article explaining Go's array and slice types:
https://go.dev/blog/slices-intro

## Generics
### Generic alias types
Not generally available in Go until v1.24, this will allow the following:
```
type Matrix[T any] = [][]T
```
See: https://go.dev/blog/alias-names



Notes
-----

\* The Gorilla icon image is AI-generated using [Google Gemini](https://gemini.google.com/)
