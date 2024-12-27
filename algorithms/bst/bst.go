// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License

// The bst package implements a binary search tree, with the following
// associative-array (map) operations:
// * Get
// * Put
// * Delete
// * Size
//
// and the following ordered-collection operations:
// * VisitInOrder
// * Min
// * Max
//
// Trees can be either unbalanced or balanced. Balanced trees are implemented
// using a Red-Black tree algorithm.
//
// This implementation does not maintain parent pointers within the nodes of the
// tree.  As such, many operations (in particular re-balancing after insertion
// and deletion) must remember the path from the root to an impacted node for
// the duration of the operation.

package bst

import (
	"cmp"
	"fmt"

	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/must"
)

type BST[K cmp.Ordered, V any] struct {
	root    *node[K, V]
	size    int
	balance bool
	// sentinels
	leftLeaf  node[K, V]
	rightLeaf node[K, V]
}

// NewBST creates a new unbalanced or balanced tree.
func NewBST[K cmp.Ordered, V any](balance bool) *BST[K, V] {
	t := &BST[K, V]{
		balance: balance,
	}
	if balance {
		t.leftLeaf.col = black
		t.rightLeaf.col = black
	}
	return t
}

// Get returns the value associated with the given key.
// If not found, returns the zero value for the value type
// and ok=false.
func (t *BST[K, V]) Get(key K) (val V, found bool) {
	n, _ := t.root.findNode(key, false)
	if n != nil {
		val = n.val
		found = true
	}
	return
}

func (t *BST[K, V]) MustGet(key K) (val V) {
	return must.BeOk(t.Get(key))
}

// Put adds the given key value pair to the tree.
// If the key already exists in the tree, the given
// value replaces the existing value.
func (t *BST[K, V]) Put(key K, val V) {
	t.insert(key, val)
}

// Delete removes the given key, and associated value,
// from the tree, and returns true if the key existed
// and false if it did not.
func (t *BST[K, V]) Delete(key K) (ok bool) {
	return t.delete(key)
}

// Size returns the number of key/value pairs currently
// in the tree.
func (t *BST[K, V]) Size() int {
	return t.size
}

// Min returns the smallest key in the tree.
func (t *BST[K, V]) Min() (minKey K, ok bool) {
	return t.root.min()
}

// Max returns the largest key in the tree.
func (t *BST[K, V]) Max() (maxKey K, ok bool) {
	return t.root.max()
}

// VisitInOrder performs an in-order traversal of all key/value pairs in the
// tree.
func (t *BST[K, V]) VisitInOrder(v types.Visitor[K, V]) {
	t.root.visitInOrder(v)
}

// insert will insert the given key,val pair into the tree rooted at the given
// node and return a pointer to the (potentially new) root of the tree.  If the
// key already exists in the tree, the current val associated with the key is
// replaced with the given val.
func (t *BST[K, V]) insert(key K, val V) {
	if t.root == nil {
		t.root = newNode(key, val)
		t.size += 1
		if t.balance {
			t.root.col = black
		}
		return
	}
	// nodePath holds the path of nodes from the root down to the node that was
	// inserted. This list is maintained only if we're balancing the tree.
	var nodePath nodeList[K, V]
	if t.balance {
		nodePath = append(nodePath, nil)
	}
	// comp will hold the result of comparison where the value is to be inserted
	var comp int = 0
	// Use iteration to find where the key should exist in the tree.
	x := t.root
	// p holds the immediate parent of x
	var p *node[K, V] = nil
	for x != nil {
		if t.balance {
			nodePath = append(nodePath, x)
		}
		p = x
		comp = cmp.Compare(key, x.key)
		if comp == 0 {
			// Key already exists in tree; replace value and return
			x.val = val
			return
		}
		if comp < 1 {
			x = x.left
		} else {
			x = x.right
		}
	}
	// at this point, p is the leaf where the new node is to be inserted
	x = newNode(key, val)
	t.size += 1
	// determine which side it belongs
	if comp < 1 {
		p.left = x
	} else {
		p.right = x
	}
	if t.balance {
		x.col = black
		nodePath = append(nodePath, x)
		// balance the tree
		t.balanceInsert(nodePath)
	}
}

func (t *BST[K, V]) delete(key K) (deleted bool) {
	if t.root == nil {
		return
	}
	// Since we do not maintain parent pointers in the node itself, we need to
	// remember the path from the root down to the node that we're splicing-out.
	// For a simple binary tree, we really only need the immediate parent. But for
	// a balanced tree, we need the entire path. So we'll just save the entire
	// path as path. The first element of the path is always nil, representing
	// the parent of the root.
	n, path := t.root.findNode(key, true)
	if n == nil {
		// Not found; nothing to do
		return
	}
	must.BeTrue(len(path) > 1)
	must.BeTrue(path[len(path)-1] == n)

	// At this point, we know we're deleting a node
	deleted = true
	t.size -= 1

	// Three cases to consider:
	// 1. n is a leaf node    - remove it from its parent
	// 2. n has one child (c) - replace n with the subtree at c
	// 3. n has two children  - replace n with the its successor, s, and remove s
	if n.left == nil && n.right == nil {
		// case (1) - leaf node
		if len(path) == 2 {
			// removing the root; tree is now empty
			t.root = nil
		} else {
			nParent := path[len(path)-2]
			if nParent.left == n {
				nParent.left = nil
				// use sentinel for leaf in path
				path[len(path)-1] = &t.leftLeaf
			} else {
				nParent.right = nil
				// use sentinel for leaf in path
				path[len(path)-1] = &t.rightLeaf
			}
			if t.balance && n.col == black {
				// The node we removed is black; need to re-balance
				t.balanceDelete(path)
			}
		}
		return
	}

	if n.left == nil || n.right == nil {
		// case (2) - one child
		var c *node[K, V]
		if n.left != nil {
			c = n.left
		} else {
			c = n.right
		}
		if len(path) == 2 {
			// removing root; c is the new root
			t.root = c
		} else {
			path[len(path)-1] = c
			nParent := path[len(path)-2]
			if nParent.left == n {
				nParent.left = c
			} else {
				nParent.right = c
			}
			if t.balance && n.col == black {
				// The node we removed is black; need to re-balance
				// On the path, replace the node with its child
				t.balanceDelete(path)
			}
		}
		return
	}
	// Cases (1) & (2) above can be handled together, but I'm leaving them
	// separate for now for clarity

	// case (3) - two children. replace content of n with its successor, s, and remove s

	// Find the successor, s
	s := n.right
	for s.left != nil {
		path = append(path, s)
		s = s.left
	}
	path = append(path, s)
	// Unlink s from its parent, replace with s's right child.
	// s will have at most one child, and that child (if present) will be on the
	// right.
	sParent := path[len(path)-2]
	path[len(path)-1] = s.right
	if sParent.left == s {
		// left side
		sParent.left = s.right
		if s.right == nil {
			// Use sentinel in path for leaf
			path[len(path)-1] = &t.leftLeaf
		}
	} else {
		// right side
		sParent.right = s.right
		if s.right == nil {
			// Use sentinel in path for leaf
			path[len(path)-1] = &t.rightLeaf
		}
	}
	// Copy fields from s to n (except for color)
	n.copy(s)
	if t.balance && s.col == black {
		// The node we removed is black; need to re-balance
		t.balanceDelete(path)
	}
	return
}

// The following are potentially expensive operations that traverse
// the entire tree, and are intended to be used internally for testing and
// diagnostics.

func (t *BST[K, V]) count() int {
	return t.root.count()
}

func (t *BST[K, V]) height() int {
	return t.root.height()
}

func (t *BST[K, V]) validate() error {
	if !t.balance {
		// not a balanced tree; nothing to validate
		return nil
	}
	if t.root.color() != black {
		return fmt.Errorf("root must be black")
	}
	_, err := t.root.blackHeight()
	return err
}
