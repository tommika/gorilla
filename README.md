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

### Applications and Command Line Tools

[gpx](./cmd/gpx) - Summarize a GPX file (leverages [geo/gpx](./geo/gpx) package)

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
