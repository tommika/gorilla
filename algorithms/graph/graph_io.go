// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package graph

import (
	"errors"
	"io"
)

type EdgeWriter[N, W any] func(out io.Writer, from, to N, weight W) (int, error)
type EdgeReader[N, W any] func(in io.Reader) (from, to N, weight W, err error)

func (g *Graph[N, W]) WriteEdges(out io.Writer, writeEdge EdgeWriter[N, W]) error {
	for v, node := range g.nodes {
		for _, e := range node.outgoing {
			if _, err := writeEdge(out, v, e.node, e.weight); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Graph[N, W]) ReadEdges(in io.Reader, readEdge EdgeReader[N, W]) error {
	for {
		if from, to, weight, err := readEdge(in); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return err
		} else {
			g.AddEdge(from, to, weight)
		}
	}
}
