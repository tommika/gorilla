// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
// The assert package includes the assertion functions that I use in unit tests.
package assert

import (
	"math"
	"reflect"
	"testing"

	"github.com/tommika/gorilla/must"
)

func Equal[T comparable](t testing.TB, expected, actual T) {
	if expected != actual {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected %v; got %v", expected, actual)
	}
}

func EqualEpsilon(t testing.TB, expected, actual, epsilon float64) {
	e := math.Abs(expected - actual)
	if e > epsilon {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected %f; got %f (epsilon allowed=%f, actual=%f)", expected, actual, epsilon, e)
	}
}

func DeepEqual[T any](t testing.TB, expected, actual T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected %v; got %v", expected, actual)
	}
}

func NotEqual[T comparable](t testing.TB, notExpected, actual T) {
	if notExpected == actual {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected %v; got %v", notExpected, actual)
	}
}

func DeepNotEqual[T comparable](t testing.TB, notExpected, actual T) {
	if reflect.DeepEqual(notExpected, actual) {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected %v; got %v", notExpected, actual)
	}
}

func Nil(t testing.TB, actual any) {
	if !must.IsNil(actual) {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected nil; got %v", actual)
	}
}
func NotNil(t testing.TB, actual any) {
	if must.IsNil(actual) {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected non-nil")
	}
}

func True(t testing.TB, actual bool) {
	if !actual {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected true")
	}
}

func False(t testing.TB, actual bool) {
	if actual {
		t.Helper() // tells t to ignore this depth of the stack when reporting file/line info
		t.Fatalf("expected false")
	}
}
