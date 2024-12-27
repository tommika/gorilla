// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package types

// Compare is a function used to compare items during a search
//
//	-1 if a is less than b,
//	 0 if a equals b,
//	+1 if a is greater than b.
type Compare[T any] func(a, b T) int
