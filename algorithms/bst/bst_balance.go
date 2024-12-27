// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package bst

import (
	"github.com/tommika/gorilla/must"
)

// balanceInsert re-balances the tree after insertion of a new node.
// `path` is the path from the root to the node that was inserted into the tree.
// path[0] is always nil, and represents the parent of the root.
func (t *BST[K, V]) balanceInsert(path nodeList[K, V]) {
	must.BeTrue(len(path) > 1)
	must.BeNil(path[0])
	i := len(path) - 1
	path[i].col = red
	for i > 1 && path[i-1].col == red {
		x := path[i]
		// violation of RB-property -- node and its parent are both red
		// need to adjust
		if path[i-2].left == path[i-1] {
			// parent is a left-child of grand-parent
			// uncle y is on the right
			y := path[i-2].right
			if y.color() == red {
				// parent and uncle are both red
				// re-color and continue at grand-parent
				y.col = black
				path[i-1].col = black
				path[i-2].col = red
				// continue loop at grandparent
				i -= 2
			} else {
				// uncle is black
				if path[i-1].right == x {
					t.leftRotate(path[i-2:])
				}
				path[i-1].col = black
				path[i-2].col = red
				t.rightRotate(path[i-3:])
				i = 0 // we're done
			}
		} else {
			// parent of x is a right-child grand-parent
			// uncle y is on the left
			y := path[i-2].left
			if y.color() == red {
				// parent and uncle are both red
				// recolor and continue at grand-parent
				y.col = black
				path[i-1].col = black
				path[i-2].col = red
				i -= 2
			} else {
				if path[i-1].left == x {
					t.rightRotate(path[i-2:])
				}
				path[i-1].col = black
				path[i-2].col = red
				t.leftRotate(path[i-3:])
				i = 0 // we're done
			}
		}
	}
	// root may have changed
	t.root = path[1]
	t.root.col = black
}

// balanceDelete re-balances the tree after a black node is removed.  `path` is the
// path from the root down to the child that replaced the removed node.  This
// child may be a leaf, in which case it is represented by a sentinel (either
// leftLeaf or rightLeft, depending on its position wrt the parent node.)
// path[0] is always nil, and represents the parent of the root.
func (t *BST[K, V]) balanceDelete(path nodeList[K, V]) {
	must.BeTrue(len(path) > 1)
	must.BeNil(path[0])
	i := len(path) - 1
	// while not the root and not a red node
	for i > 1 && path[i].color() == black {
		x := path[i]
		// Determine which side of the parent this node is on.
		if (x == &t.leftLeaf) || (x == path[i-1].left) {
			// left side
			w := path[i-1].right
			if w.col == red {
				w.col = black
				path[i-1].col = red
				path[i] = w
				t.leftRotate(path[i-2:])
				path = append(path, x)
				i = len(path) - 1
				w = path[i-1].right
			}
			if w.left.color() == black && w.right.color() == black {
				w.col = red
				i -= 1
			} else {
				if w.right.color() == black {
					w.left.col = black
					w.col = red
					pT := []*node[K, V]{path[i-1], w, w.left}
					t.rightRotate(pT)
					w = path[i-1].right
				}
				w.col = path[i-1].col
				path[i-1].col = black
				w.right.col = black
				path[i] = w
				t.leftRotate(path[i-2:])
				i = 1 // we're done
			}
		} else {
			must.BeTrue((x == &t.rightLeaf) || (x == path[i-1].right))
			// right side
			w := path[i-1].left
			if w.col == red {
				w.col = black
				path[i-1].col = red
				path[i] = w
				t.rightRotate(path[i-2:])
				path = append(path, x)
				i = len(path) - 1
				w = path[i-1].left
			}
			if w.left.color() == black && w.right.color() == black {
				w.col = red
				i -= 1
			} else {
				if w.left.color() == black {
					w.right.col = black
					w.col = red
					pT := []*node[K, V]{path[i-1], w, w.right}
					t.leftRotate(pT)
					w = path[i-1].left
				}
				w.col = path[i-1].col
				path[i-1].col = black
				w.left.col = black
				path[i] = w
				t.rightRotate(path[i-2:])
				i = 1 // we're done
			}
		}
		path = path[:i+1] // trim path to current node
	}
	path[i].col = black
	// root may have changed
	t.root = path[1]
}

// path is the path that is to be rotated:
// path[0] = parent(x)
// path[1] = x - the pivot
// path[2] = right(x)
func (t *BST[K, V]) leftRotate(path nodeList[K, V]) {
	p := path[0] // parent can be nil
	x := path[1]
	y := path[2]
	must.NotBeNil(x)
	must.NotBeNil(y)
	must.BeTrue(y == x.right)
	x.right = y.left
	if p != nil { // pivot is not the root
		// update parent to point to y
		if p.left == x {
			// x is on the left of its parent
			p.left = y
		} else {
			// x is on the right of its parent
			p.right = y
		}
	}
	y.left = x
	path.swap(1, 2)
}

// path is the path that is to be rotated:
// path[0] = parent(x)
// path[1] = x (the pivot)
// path[2] = left(x)
func (t *BST[K, V]) rightRotate(path nodeList[K, V]) {
	xp := path[0] // parent can be nil
	x := path[1]
	must.NotBeNil(x)
	y := path[2]
	must.NotBeNil(y)
	must.BeTrue(y == x.left)
	x.left = y.right
	if xp != nil { // pivot is not the root
		// update parent to point to y
		if xp.left == x {
			// x is on the left of its parent
			xp.left = y
		} else {
			// x is on the right of its parent
			xp.right = y
		}
	}
	y.right = x
	path.swap(1, 2)
}
