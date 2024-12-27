// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package types

import "cmp"

// Visitor defines a function to be applied when visiting the key-value pairs
// of a container.
type Visitor[K cmp.Ordered, V any] func(key K, val V)

// Ordered is the interface to a container of key-value pairs sorted by key
type Ordered[K cmp.Ordered, V any] interface {
	Size() (size int)
	// Min returns the smallest key in the container, or false if the container is empty
	Min() (minKey K, ok bool)
	// Max returns the largest key in the container, or false if the container is empty
	Max() (maxKey K, ok bool)
	// VisitInOrder visits all key,value pairs in order
	VisitInOrder(v Visitor[K, V])
}
