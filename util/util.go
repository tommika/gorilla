// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"math/rand"

	"github.com/tommika/gorilla/must"
)

// CopySlice returns a copy of the given slice.
func CopySlice[T any](items []T) []T {
	copyOfItems := make([]T, len(items))
	copy(copyOfItems, items)
	return copyOfItems
}

// ShuffleSlice randomizes the order of elements in the given array.
// This is done in-place. The slice is returned for convenience,
// for example to use in chaining.
func ShuffleSlice[T any](items []T) []T {
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	return items
}

// ReverseSlice reverses a slice in-place.
// The slice is returned as a convenience.
func ReverseSlice[T any](items []T) []T {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return items
}

func IsNil(v any) bool {
	return must.IsNil(v)
}

// MaxInt returns largest of the given integer arguments
func MaxInt(first int, rest ...int) int {
	max := first
	for _, n := range rest {
		if n > max {
			max = n
		}
	}
	return max
}

func MakeIntArray(size int) []int {
	ints := make([]int, size)
	for i := 0; i < size; i++ {
		ints[i] = i
	}
	return ints
}

// Zero returns the zero value of the given type T
func Zero[T any]() (t T) {
	return
}

func IsOk(_ any, ok bool) bool {
	return ok
}

func AbsInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](val T) T {
	if val < 0 {
		val = -val
	}
	return val
}
