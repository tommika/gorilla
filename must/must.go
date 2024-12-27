// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License

// The must package includes runtime assertions that I use throughout the Gorilla codebase.
// These are often used as pre/post-conditions or other places where some invariant is checked.
package must

import (
	"cmp"
	"fmt"
	"reflect"
)

func BeOk[T any](val T, ok bool) T {
	if !ok {
		panic("it wasn't ok")
	}
	return val
}

func NotBeAnError[T any](val T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: got: %v", err))
	}
	return val
}

func BeEqual[T cmp.Ordered](a, b T) T {
	if cmp.Compare(a, b) != 0 {
		panic(fmt.Sprintf("expected equal: a=%v, b=%v", a, b))
	}
	return a
}

func BeNil[T any](a T) T {
	if !IsNil(a) {
		panic(fmt.Sprintf("expected nil: got: %v", a))
	}
	return a
}

func NotBeNil[T any](a T) T {
	if IsNil(a) {
		panic("it was nil")
	}
	return a
}

func BeTrue(f bool) bool {
	if !f {
		panic("it wasn't true")
	}
	return f
}

func IsNil(v any) bool {
	if v == nil {
		return true
	}
	vT := reflect.ValueOf(v)
	switch vT.Kind() {
	case reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return vT.IsNil()
	}
	return false
}
