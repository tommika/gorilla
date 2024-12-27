// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package types

import "cmp"

// Map is the interface to an associative container of key-value pairs with
// unique keys.
type Map[K cmp.Ordered, V any] interface {
	Size() (size int)
	Get(key K) (val V, found bool)
	MustGet(key K) (val V)
	Put(key K, val V)
	Delete(key K) (deleted bool)
}
