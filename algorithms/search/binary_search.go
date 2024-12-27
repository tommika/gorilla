// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package search

import "github.com/tommika/gorilla/algorithms/types"

func BinarySearch[T any](items []T, find T, compare types.Compare[T]) (val T, found bool) {
	i, j := 0, len(items)
	for i < j {
		m := i + (j-i)/2
		n := compare(find, items[m])
		if n < 0 {
			j = m
		} else if n > 0 {
			i = m + 1
		} else {
			return items[m], true
		}
	}
	return
}
