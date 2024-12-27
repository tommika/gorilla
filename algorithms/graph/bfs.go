// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package graph

import (
	"github.com/tommika/gorilla/algorithms/queue"
	"github.com/tommika/gorilla/must"
	"github.com/tommika/gorilla/util"
)

// BFTree is the tree constructed in a breadth-first search of the graph.
type BFTree[N comparable, W Weight] struct {
	root    N // source node
	nodeMap map[N]bftNode[N, W]
	nodes   []N // in BF order
}

// bftNode is a node in the tree that is constructed in a breadth-first
// search of the graph.
type bftNode[N comparable, W Weight] struct {
	pred   N
	weight W
}

// FindPath determines if there is a path between two nodes in the graph.  If
// so, the path and the accumulated weight along that path are returned.
func (g *Graph[N, W]) FindPath(s, v N) (path []N, weight W) {
	bft := g.BFS(s, 0)
	return bft.FindPath(v)
}

// BFS performs a breadth-first search of the graph, starting
// from the given node, and returns the resulting breadth-first tree.
func (g *Graph[N, W]) BFS(s N, maxDepth int) (bft BFTree[N, W]) {
	if !g.HasNode(s) {
		return
	}
	type bfqItem[N comparable] struct {
		u     N
		depth int
	}
	bft.init(s)
	q := queue.DynamicCircularArrayQueue[bfqItem[N]]{}
	q.Enqueue(bfqItem[N]{u: s, depth: 0})
	for q.Size() != 0 {
		head := q.MustDequeue()
		bft.nodes = append(bft.nodes, head.u)
		if maxDepth == 0 || head.depth < maxDepth {
			for _, e := range g.outgoing(head.u) {
				if !bft.hasNode(e.node) {
					// first time visiting this node
					bft.addEdge(e, head.u)
					q.Enqueue(bfqItem[N]{u: e.node, depth: head.depth + 1})
				}
			}
		}
	}
	return
}

// FindPath determines if the tree contains a path to the given node. If so,
// the path and the accumulated weight along that path are returned.
func (bft *BFTree[N, W]) FindPath(v N) (path []N, weight W) {
	if len(bft.nodeMap) == 0 {
		return
	}
	node, ok := bft.nodeMap[v]
	if !ok {
		// no path
		return
	}
	pathT := []N{}
	var weightT W
	for v != bft.root {
		pathT = append(pathT, v)
		weightT += node.weight
		v = node.pred
		node, ok = bft.nodeMap[v]
		must.BeTrue(ok)
	}
	pathT = append(pathT, bft.root)
	// we built the path from bottom-up, so need to reverse it
	// before returning.
	path = util.ReverseSlice(pathT)
	weight = weightT
	return
}

func (bft *BFTree[N, W]) VisitNodes(visit func(v N)) {
	for _, v := range bft.nodes {
		visit(v)
	}
}

func (bft *BFTree[N, W]) NodeCount() int {
	must.BeEqual(len(bft.nodes), len(bft.nodeMap))
	return len(bft.nodes)
}

func (bft *BFTree[N, W]) init(root N) {
	bft.root = root
	bft.nodeMap = map[N]bftNode[N, W]{}
	bft.nodeMap[root] = bftNode[N, W]{}
}

func (bft *BFTree[N, W]) hasNode(n N) bool {
	_, found := bft.nodeMap[n]
	return found
}
func (bft *BFTree[N, W]) addEdge(e edge[N, W], pred N) {
	bft.nodeMap[e.node] = bftNode[N, W]{
		weight: e.weight, // weight of edge leading to this node
		pred:   pred,
	}
}
