// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
// The graph package includes data structures and algorithms for analyzing graphs of objects.
package graph

import (
	"fmt"
	"io"

	"github.com/tommika/gorilla/algorithms/matrix"
	"github.com/tommika/gorilla/util"
)

// GraphType is used to indicate if a graph is directed or undirected
type GraphType = bool

const (
	// Undirected is an undirected graph
	Undirected GraphType = false
	// Directed is a directed graph
	Directed GraphType = true
)

type Weight interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type Graph[N comparable, W Weight] struct {
	directed GraphType
	numEdges int
	nodes    map[N]*node[N, W]
}

// NewGraph creates an empty graph of the given type (Directed or Undirected)
func NewGraph[N comparable, W Weight](directed GraphType) *Graph[N, W] {
	g := &Graph[N, W]{
		numEdges: 0,
		nodes:    make(map[N]*node[N, W], 0),
		directed: directed,
	}
	return g
}

// GraphFromAdjacencyMatrix constructs a graph from the given adjacency matrix. If the
// graph is undirected, the only cells above the diagonal are considered.
func GraphFromAdjacencyMatrix[W Weight](m matrix.Matrix[W], directed GraphType) *Graph[int, W] {
	g := NewGraph[int, W](directed)
	for r := range m.NumRows() {
		for c := range m.NumCols() {
			w := m.Get(r, c)
			if r != c && w != util.Zero[W]() {
				g.AddEdge(r, c, w)
			}
		}
	}
	return g
}

// Fprintf outputs the graph to the given writer
func (g *Graph[N, W]) Fprint(out io.Writer) {
	var dir rune
	if g.directed {
		dir = '>'
	}
	fmt.Printf("#nodes=%d, #edges=%d, directed=%t\n", g.NodeCount(), g.EdgeCount(), g.directed)
	for l, v := range g.nodes {
		for _, e := range v.outgoing {
			fmt.Fprintf(out, "[%v] --(%v)--%c [%v]\n", l, e.weight, dir, e.node)
		}
	}
}
func (g *Graph[N, W]) Type() GraphType {
	return g.directed
}

func (g *Graph[N, W]) NodeCount() int {
	return len(g.nodes)
}

func (g *Graph[N, W]) EdgeCount() int {
	return g.numEdges
}

func (g *Graph[N, W]) AddEdge(from, to N, weight W) {
	if from == to {
		// REVIEW: may want to return some indication of this edge being ignored
		return
	}
	nFrom, found := g.nodes[from]
	if !found {
		nFrom = g.newNode(from)
		g.nodes[from] = nFrom
	}
	nTo, found := g.nodes[to]
	if !found {
		nTo = g.newNode(to)
		g.nodes[to] = nTo
	}

	// When adding an edge, we do not distinguish between directed and undirected.
	// Rather, the edge is added as if it were directed, by adding it to the
	// source nodes outgoing list, and the target nodes incoming list.
	// If the graph is undirected, then we consider both lists when looking for a
	// edge between two nodes; i.e., there is no difference between incoming and
	// outgoing, so we need to look in both lists. This is handled transparently
	// by the method below,
	nFrom.outgoing = append(nFrom.outgoing, g.newEdge(to, weight))
	nTo.incoming = append(nTo.incoming, g.newEdge(from, weight))

	g.numEdges++
}

// HasNode determines if the graph includes the given node.
func (g *Graph[N, W]) HasNode(n N) bool {
	_, found := g.nodes[n]
	return found
}

// HasEdge determines if there is an edge from one node to another.
func (g *Graph[N, W]) HasEdge(from, to N) bool {
	return g.findEdge(from, to) != nil
}

// GetEdgeWeight determines the weight of the edge from one node to another.
// If the edge does not exist, found will be false and the weight is the
// zero-value.
func (g *Graph[N, W]) GetEdgeWeight(from, to N) (weight W, found bool) {
	e := g.findEdge(from, to)
	if e == nil {
		return
	}
	return e.weight, true
}

// outgoing returns the outgoing edges from the given node. If the graph
// is undirected, then incoming edges are considered to be outgoing edges
// as well.
func (g *Graph[N, W]) outgoing(from N) []edge[N, W] {
	edges := g.nodes[from].outgoing
	if !g.directed {
		edges = append(edges, g.nodes[from].incoming...)
	}
	return edges
}

// incoming returns the incoming edges from the given node. If the graph
// is undirected, then outgoing edges are considered to be incoming edges
// as well.
func (g *Graph[N, W]) incoming(from N) []edge[N, W] {
	edges := g.nodes[from].incoming
	if !g.directed {
		edges = append(edges, g.nodes[from].outgoing...)
	}
	return edges
}

// findEdge determines if there is an edge between two nodes.
// If there is, return a pointer to the edge, otherwise return nil.
// In the case of an undirected graph, we need to check both
// incoming and outgoing edges.
func (g *Graph[N, W]) findEdge(from, to N) (theEdge *edge[N, W]) {
	if nFrom, found := g.nodes[from]; found {
		theEdge = g.findEdgeInSlice(nFrom.outgoing, to)
		if theEdge == nil && !g.directed {
			// undirected graph; also check incoming edges
			theEdge = g.findEdgeInSlice(nFrom.incoming, to)
		}
	}
	return
}

func (g *Graph[N, W]) findEdgeInSlice(edges []edge[N, W], to N) *edge[N, W] {
	for i := range edges {
		if edges[i].node == to {
			return &(edges[i])
		}
	}
	return nil
}

type edge[N, W comparable] struct {
	node   N
	weight W
}

type node[N, W comparable] struct {
	name     N
	outgoing []edge[N, W]
	incoming []edge[N, W]
}

func (g *Graph[N, W]) newNode(name N) *node[N, W] {
	return &node[N, W]{
		name: name,
	}
}

func (g *Graph[N, W]) newEdge(n N, weight W) edge[N, W] {
	return edge[N, W]{
		node:   n,
		weight: weight,
	}
}
