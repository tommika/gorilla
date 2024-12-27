// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package graph

import (
	"os"
	"strconv"
	"testing"

	"github.com/tommika/gorilla/algorithms/matrix"
	"github.com/tommika/gorilla/assert"
)

func TestBFSiDirected(t *testing.T) {
	const ms = `0 1 0 0 
	            0 0 1 0
	            1 0 0 0
	            0 0 0 0`
	m := matrix.ParseMatrix(ms, strconv.Atoi)
	m.Fprintf(os.Stdout)
	g1 := GraphFromAdjacencyMatrix[int](m, Directed)
	g1.Fprint(os.Stderr)

	path, weight := g1.FindPath(0, 0)
	assert.NotNil(t, path)
	assert.Equal(t, weight, 0)

	bft := g1.BFS(0, 0)
	assert.Equal(t, 3, bft.NodeCount())
	count := 0
	bft.VisitNodes(func(n int) {
		count++
	})
	assert.Equal(t, 3, count)

	bft = g1.BFS(0, 0)
	assert.Equal(t, 3, bft.NodeCount())

	bft = g1.BFS(0, 1)
	assert.Equal(t, 2, bft.NodeCount())

	path, weight = g1.FindPath(0, 2)
	t.Logf("path: %v, weight: %v\n", path, weight)
	assert.NotNil(t, path)
	assert.Equal(t, 3, len(path))
	assert.Equal(t, 2, weight)

	path, weight = g1.FindPath(0, 3)
	assert.Nil(t, path)
	assert.Equal(t, 0, weight)

	path, weight = g1.FindPath(4, 5)
	assert.Nil(t, path)
	assert.Equal(t, 0, weight)

}

func TestBFSUndirected(t *testing.T) {
	const ms = `0 1 0 0 
	            0 0 1 0
	            1 0 0 0
	            0 0 0 0`
	m := matrix.ParseMatrix(ms, strconv.Atoi)
	m.Fprintf(os.Stdout)
	g1 := GraphFromAdjacencyMatrix[int](m, Undirected)
	g1.Fprint(os.Stderr)

	path, weight := g1.FindPath(0, 0)
	assert.NotNil(t, path)
	assert.Equal(t, weight, 0)

	bft := g1.BFS(0, 0)
	assert.Equal(t, 3, bft.NodeCount())
	count := 0
	bft.VisitNodes(func(n int) {
		count++
	})
	assert.Equal(t, 3, count)

	bft = g1.BFS(0, 100)
	assert.Equal(t, 3, bft.NodeCount())

	bft = g1.BFS(0, 1)
	assert.Equal(t, 3, bft.NodeCount())

	path, weight = g1.FindPath(0, 2)
	t.Logf("path: %v, weight: %v\n", path, weight)
	assert.NotNil(t, path)
	assert.Equal(t, 2, len(path))
	assert.Equal(t, 1, weight)

	path, weight = g1.FindPath(2, 0)
	t.Logf("path: %v, weight: %v\n", path, weight)
	assert.NotNil(t, path)
	assert.Equal(t, 2, len(path))
	assert.Equal(t, 1, weight)

	path, weight = g1.FindPath(0, 3)
	assert.Nil(t, path)
	assert.Equal(t, 0, weight)

	path, weight = g1.FindPath(4, 5)
	assert.Nil(t, path)
	assert.Equal(t, 0, weight)

}
