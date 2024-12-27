// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License

package bst

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/tommika/gorilla/algorithms/types"
)

type color uint8

const (
	none color = iota
	black
	red
)

// node is a node in a binary search tree.
// This implementation does not maintain parent pointers.
// Where parent relationships are needed, in particular during
// tree balancing operations, a path of nodes from root to
// node is temporarily maintained.
type node[K cmp.Ordered, V any] struct {
	key   K
	val   V
	left  *node[K, V]
	right *node[K, V]
	col   color
}

func (node *node[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("(key=%v,col=%d)", node.key, node.col))
	return sb.String()
}

type nodeList[K cmp.Ordered, V any] []*node[K, V]

func (l nodeList[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, n := range l {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", n))
	}
	sb.WriteString("]")
	return sb.String()
}

func (l nodeList[K, V]) swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func newNode[K cmp.Ordered, V any](key K, val V) *node[K, V] {
	return &node[K, V]{
		key:   key,
		val:   val,
		left:  nil,
		right: nil,
	}
}

// Copy wil copy the fields from one node to another.
// The node color is NOT copied
func (node *node[K, V]) copy(from *node[K, V]) {
	node.key = from.key
	node.val = from.val
	// DO NOT copy the color
}

// count will determine the number of nodes in the tree
// rooted at the given node.
func (root *node[K, V]) count() (count int) {
	root.visitInOrder(func(_ K, _ V) {
		count++
	})
	return
}

// visitInOrder will visit all nodes in the tree, in order,
// and for each node, call the given visitor function.
func (n *node[K, V]) visitInOrder(visitor types.Visitor[K, V]) {
	if n == nil {
		return
	}
	n.left.visitInOrder(visitor)
	visitor(n.key, n.val)
	n.right.visitInOrder(visitor)
}

// height will determine the height of the tree rooted at the given node.
func (node *node[K, V]) height() (height int) {
	if node == nil {
		return
	}
	hLeft := node.left.height()
	hRight := node.right.height()
	height = 1 + max(hLeft, hRight)
	return
}

func (root *node[K, V]) max() (maxKey K, ok bool) {
	if root == nil {
		return
	}
	n := root
	for n.right != nil {
		n = n.right
	}
	return n.key, true
}

func (root *node[K, V]) min() (minKey K, ok bool) {
	if root == nil {
		return
	}
	n := root
	for n.left != nil {
		n = n.left
	}
	return n.key, true
}

// findNode will search the tree (rooted a the given node)
// for the given key. If found, returns a pointer to the node containing
// the key. If savePath is true, then the path from the root to the node
// is returned as a nodeList
func (root *node[K, V]) findNode(key K, savePath bool) (found *node[K, V], path nodeList[K, V]) {
	if savePath {
		path = append(path, nil)
	}
	n := root
	for n != nil && found == nil {
		if savePath {
			path = append(path, n)
		}
		comp := cmp.Compare(key, n.key)
		if comp == 0 {
			found = n
		} else {
			if comp < 1 {
				n = n.left
			} else {
				n = n.right
			}
		}
	}
	return
}

// color returns the color of the node.
// External leaf nodes (nil) are black.
func (n *node[K, V]) color() (c color) {
	if n == nil {
		c = black
	} else {
		c = n.col
	}
	return
}

// func (n *node[K, V]) leftChild() (left *node[K, V]) {
// 	if n != nil {
// 		left = n.left
// 	}
// 	return
// }

// func (n *node[K, V]) rightChild() (right *node[K, V]) {
// 	if n != nil {
// 		right = n.right
// 	}
// 	return
// }

func (node *node[K, V]) blackHeight() (blackHeight int, err error) {
	if err = node.validateRBNode(); err != nil {
		return 0, err
	}
	if node == nil {
		return 1, nil
	}
	var hT int
	if blackHeight, err = node.left.blackHeight(); err == nil {
		if hT, err = node.right.blackHeight(); err == nil {
			if blackHeight != hT {
				// property #4: black height must be the same for all paths
				err = fmt.Errorf("property #4 violated: inconsistent black height: hLeft=%d, hRight=%d", blackHeight, hT)
			} else if node.col == black {
				blackHeight += 1
			}
		}
	}
	return
}

func (node *node[K, V]) validateRBNode() error {
	// property #2: leaf nodes (nil) are black
	if node == nil {
		// leaf node; ok
		return nil
	}
	// property #1: every node is either red or black
	if !(node.col == red || node.col == black) {
		// invalid
		return fmt.Errorf("property #1 violated: col=%d", node.col)
	}
	// property #3: if a node is red, then both children are black
	if node.col == red && !(node.left.color() == black && node.right.color() == black) {
		// invalid
		return fmt.Errorf("property #3 violated: left.col=%d, right.col=%d", node.left.color(), node.right.color())
	}
	return nil
}
