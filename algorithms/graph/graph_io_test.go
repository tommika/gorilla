// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package graph

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/tommika/gorilla/algorithms/matrix"
	"github.com/tommika/gorilla/assert"
)

func writeEdge(out io.Writer, from, to int, weight int) (int, error) {
	return fmt.Fprintf(out, "%d|%d|%d\n", from, to, weight)
}

func readEdge(in io.Reader) (from, to int, weight int, err error) {
	_, err = fmt.Fscanf(in, "%d|%d|%d\n", &from, &to, &weight)
	return
}

func stringReader(s string) io.Reader {
	return bytes.NewReader([]byte(s))
}

func testGraphIO(t *testing.T, g *Graph[int, int]) {
	buffer := bytes.Buffer{}
	out := bufio.NewWriter(&buffer)
	g.WriteEdges(out, writeEdge)
	out.Flush()
	edges := buffer.String()
	t.Logf("edges:\n%s\n", edges)
	gT := NewGraph[int, int](g.Type())
	gT.ReadEdges(stringReader(edges), readEdge)
	gT.Fprint(os.Stderr)

	assert.Equal(t, g.EdgeCount(), gT.EdgeCount())
	assert.Equal(t, g.NodeCount(), gT.NodeCount())

}

func TestGraphIO(t *testing.T) {
	const ams = `0 1 0 0 
	             0 0 1 0
	             0 0 0 1
	             0 0 0 0`
	am := matrix.ParseMatrix(ams, strconv.Atoi)
	am.Fprintf(os.Stdout)
	g1 := GraphFromAdjacencyMatrix[int](am, Undirected)
	g1.Fprint(os.Stderr)
	testGraphIO(t, g1)

	g2 := GraphFromAdjacencyMatrix[int](am, Directed)
	g2.Fprint(os.Stderr)
	testGraphIO(t, g2)

}
