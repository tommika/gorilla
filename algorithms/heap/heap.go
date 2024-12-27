// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
//
// Package heap provides a data structure and algorithm for sorting and priority-queue operations.
// A heap is a _complete_ binary tree: it's perfectly balanced and all levels
// are full (with the possible exception of the lowest level, which may be partially
// filled.) What's particularly interesting is that the nodes of the tree are
// maintained as a array, without the need to have explicit explicit links between nodes.
//
// This is a min-ordered heap, meaning that the "smallest" item is always on the
// top of the heap.  "smallest" is this sense is a logical concept. The actual
// order is defined by a user-defined compare function. As such, the user can
// control the actual ordering of items on the heap (e.g., if a max-ordered heap
// is desired.)
package heap

import (
	"fmt"
	"io"

	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/util"
)

// HeapSort sorts the given array of items using the given compare function.
// The sort is performed in-place, meaning that the given items array is
// modified and no additional memory is allocated.
func HeapSort[T any](items []T, compFunc types.Compare[T]) {
	// The heaps logical order is "smallest on top": i.e., it's a min-heap.
	// However, when we build a heap for sorting, we want to first build a
	// max-heap, and then flatten it into a min-first array.
	// As such, we need to negate the user-defined compare function when we
	// build that heap.
	h := BuildHeap(items, func(a, b T) int {
		return -compFunc(a, b)
	})
	// Swap the "largest" item on the heap (top) with the item at the bottom of
	// the heap. This may temporarily violate the heap property at the root, but
	// the largest item is now in place at the end of the items array. We then
	// shrink the heap by one, and run heapify on the root (thereby ensuring that
	// the heap property is again in place.) We continue this until the entire heap
	// has been flattened-out into a min-order ordered array.
	for i := h.size - 1; i > 0; i-- {
		h.swapItems(0, i)
		h.size--
		h.heapify(0)
	}
}

// Heap is a data structure maintaining a heap of items
type Heap[T any] struct {
	compFunc types.Compare[T]
	items    []T
	size     int
}

func NewHeap[T any](compFunc types.Compare[T]) *Heap[T] {
	h := &Heap[T]{}
	return h.Init(compFunc)
}

func (h *Heap[T]) Init(compFunc types.Compare[T]) *Heap[T] {
	h.compFunc = compFunc
	h.size = 0
	h.items = nil
	return h
}

// Builds a heap from the given items, using the given compare function to order
// items.  Ownership of the items array is transferred to the heap, and will be
// modified by the heap.
func BuildHeap[T any](items []T, compFunc types.Compare[T]) *Heap[T] {
	h := &Heap[T]{
		compFunc: compFunc,
		size:     len(items),
		items:    items,
	}
	// Leaf nodes trivially satisfy the heap property.  As such, we build the
	// heap upwards starting from the deepest non-leaf node, and working
	// backwards through the nodes in the tree.
	for i := h.size/2 - 1; i >= 0; i-- {
		h.heapify(i)
	}
	return h
}

// Size returns the current size of the heap
func (h *Heap[T]) Size() int {
	if h == nil || h.size == 0 {
		return 0
	}
	return h.size
}

// Peek returns the item at the top of the heap.
// This is the smallest item currently on the heap, as
// defined by the compare function.
//
// If the heap is empty, the zero value for the heap
// type is returned, and ok is false.
// Otherwise the top of the heap is returned and ok is true.
func (h *Heap[T]) Peek() (top T, ok bool) {
	if h == nil || h.size == 0 {
		return
	}
	return h.items[0], true
}

// Top removes the item at the top of the heap and returns it.
// This is the smallest item currently on the heap, as
// defined by the compare function.
//
// If the heap is empty, the zero value for the heap
// type is returned, and ok is false.
// Otherwise the top of the heap is returned and ok is true.
func (h *Heap[T]) Pop() (top T, ok bool) {
	if h == nil || h.size == 0 {
		return
	}
	top = h.items[0]
	h.items[0] = h.items[h.size-1]
	h.items[h.size-1] = util.Zero[T]()
	h.size--
	h.heapify(0)
	return top, true
}

// Push pushes an item onto the heap.
func (h *Heap[T]) Push(item T) {
	// Add the item to the bottom of the heap, which
	// may temporarily violate the heap property.
	i := h.size
	h.size++
	// See if we have room in the existing items array
	if h.size > len(h.items) {
		// Need to grow the items array
		// For now, we'll let the Go runtime manage the efficient growth of the array
		h.items = append(h.items, item)
	} else {
		h.items[i] = item
	}
	// Ensure that the heap property is satisfied by sifting-up
	// the new item into its correct position in the heap
	for i > 0 && h.less(i, parent(i)) {
		h.swapItems(i, parent(i))
		i = parent(i)
	}
}

func (h *Heap[T]) Describe(out io.Writer) {
	if h == nil {
		return
	}
	fmt.Fprintf(out, "Heap stats\n")
	fmt.Fprintf(out, "\tsize     : %d\n", h.size)
	fmt.Fprintf(out, "\tcapacity : %d\n", cap(h.items))
	if h.size > 0 {
		fmt.Fprintf(out, "\ttop      : %v\n", h.items[0])
	}
}

// heapify ensures that the subtree rooted at i satisfies the heap property.
// The left and right children of i must already satisfy the heap property.
func (h *Heap[T]) heapify(i int) {
	l := leftChild(i)
	r := rightChild(i)
	// determine which node (i,l,r) is smallest
	min := i
	if l < h.size && h.less(l, i) {
		min = l
	}
	if r < h.size && h.less(r, min) {
		min = r
	}
	if min != i {
		// sift-down
		h.swapItems(i, min)
		h.heapify(min)
	}
}

// less determines if the item i is less than item j.
func (h *Heap[T]) less(i, j int) bool {
	return h.compFunc(h.items[i], h.items[j]) < 0
}

// swapItems swaps the given items
func (h *Heap[T]) swapItems(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Helper functions for navigating the heap

func parent(i int) int {
	return (i - 1) / 2
}
func leftChild(i int) int {
	return (2 * i) + 1
}
func rightChild(i int) int {
	return (2 * i) + 2
}
