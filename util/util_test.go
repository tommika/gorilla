// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 100, MaxInt(100))
	assert.Equal(t, 100, MaxInt(100, 1))
	assert.Equal(t, 100, MaxInt(1, 100))
	assert.Equal(t, 100, MaxInt(100, 1, 10))
	assert.Equal(t, 100, MaxInt(1, 100, 10))
	assert.Equal(t, 100, MaxInt(1, 10, 100))
}

func TestShuffleSlice(t *testing.T) {
	words, err := ReadFileAsWords("../algorithms/test-data/words.txt")
	assert.Nil(t, err)
	// create a set of all the words
	wordSet := make(map[string]bool)
	for _, word := range words {
		wordSet[word] = true
	}
	assert.Equal(t, len(words), len(wordSet))
	// shuffle the words
	shuffled := ShuffleSlice(CopySlice(words))
	assert.Equal(t, len(words), len(shuffled))
	// check that the shuffled words (and only those words) exist in the set
	// (destroys the contents of the set)
	for _, word := range shuffled {
		val, found := wordSet[word]
		assert.True(t, found && val)
		delete(wordSet, word)
	}
	assert.Equal(t, 0, len(wordSet))
}

func TestReverseSlice(t *testing.T) {
	a := []int{4, 3, 2, 1, 0}
	ReverseSlice(a)
	for i := 0; i < 5; i++ {
		assert.Equal(t, i, a[i])
	}
}

func TestIsNil(t *testing.T) {
	assert.True(t, IsNil(nil))
	assert.False(t, IsNil("wow"))
	var a1 []string
	assert.True(t, IsNil(a1))
	a1 = []string{}
	assert.False(t, IsNil(a1))
}

func TestMakeIntArray(t *testing.T) {
	ints := MakeIntArray(10)
	assert.Equal(t, 10, len(ints))
}

func TestZero(t *testing.T) {
	zero := Zero[int]()
	assert.Equal(t, 0, zero)
}
func okay[T any](val T, ok bool) (T, bool) {
	return val, ok
}

func TestIsOk(t *testing.T) {
	assert.True(t, IsOk(okay(0, true)))
	assert.False(t, IsOk(okay(0, false)))
}

func TestAbsInt(t *testing.T) {
	assert.Equal(t, 2112, AbsInt(2112))
	assert.Equal(t, 2112, AbsInt(-2112))
}
