// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package graph

import (
	"os"
	"strconv"
	"testing"

	"github.com/tommika/gorilla/algorithms/matrix"
	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/util"
)

func TestSelfEdge(t *testing.T) {
	g := NewGraph[string, int](Directed)
	g.AddEdge("A", "A", 0)
	g.Fprint(os.Stderr)
	assert.Equal(t, 0, 0)
	assert.Equal(t, 0, 0)
}

func TestGraphFromAdjacencyMatrix(t *testing.T) {
	const ms = `0 1 0 0
	            0 0 1 0
	            1 0 0 0
	            0 0 0 0`
	m := matrix.ParseMatrix(ms, strconv.Atoi)
	m.Fprintf(os.Stdout)

	g1 := GraphFromAdjacencyMatrix[int](m, Directed)
	g1.Fprint(os.Stderr)
	assert.Equal(t, 3, g1.NodeCount())
	assert.Equal(t, 3, g1.EdgeCount())
	assert.Equal(t, 1, len(g1.incoming(0)))
	assert.Equal(t, 1, len(g1.outgoing(0)))

	g2 := GraphFromAdjacencyMatrix[int](m, Undirected)
	g2.Fprint(os.Stderr)
	assert.Equal(t, 3, g2.NodeCount())
	assert.Equal(t, 3, g2.EdgeCount())
	assert.Equal(t, 2, len(g2.incoming(0)))
	assert.Equal(t, 2, len(g2.outgoing(0)))
}

func TestDirected(t *testing.T) {
	const ams = `0 1 1 0 
	             0 0 1 0
	             1 0 0 0
	             0 0 0 0`
	am := matrix.ParseMatrix(ams, strconv.Atoi)
	am.Fprintf(os.Stdout)

	g := GraphFromAdjacencyMatrix[int](am, Directed)
	g.Fprint(os.Stderr)

	assert.Equal(t, 3, g.NodeCount())
	assert.Equal(t, 4, g.EdgeCount())

	assert.True(t, g.HasEdge(0, 1))
	assert.True(t, util.IsOk(g.GetEdgeWeight(0, 1)))
	assert.True(t, g.HasEdge(0, 2))
	assert.True(t, g.HasEdge(1, 2))
	assert.True(t, g.HasEdge(2, 0))

	assert.False(t, g.HasEdge(1, 0))
	assert.False(t, util.IsOk(g.GetEdgeWeight(1, 0)))
	assert.False(t, g.HasEdge(2, 1))

	assert.False(t, g.HasEdge(1, 3))
	assert.False(t, g.HasEdge(3, 1))
}

func TestUndirected(t *testing.T) {
	const ams = `0 1 0 0 
	             0 0 1 0
	             0 0 0 1
	             0 0 0 0`
	am := matrix.ParseMatrix(ams, strconv.Atoi)
	am.Fprintf(os.Stdout)
	g := GraphFromAdjacencyMatrix[int](am, Undirected)
	g.Fprint(os.Stderr)

	assert.Equal(t, 4, g.NodeCount())
	assert.Equal(t, 3, g.EdgeCount())

	assert.True(t, g.HasEdge(0, 1))
	assert.True(t, g.HasEdge(1, 0))
	assert.False(t, g.HasEdge(0, 2))
	assert.False(t, g.HasEdge(2, 0))
	assert.True(t, g.HasEdge(1, 2))
	assert.True(t, g.HasEdge(2, 1))
	assert.True(t, g.HasEdge(2, 3))
	assert.True(t, g.HasEdge(3, 2))
	assert.False(t, g.HasEdge(0, 3))
	assert.False(t, g.HasEdge(3, 0))
	assert.False(t, g.HasEdge(1, 3))
	assert.False(t, g.HasEdge(3, 1))

}
