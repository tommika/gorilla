// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package search

import (
	"strings"
	"testing"

	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/util"
)

func compIgnoreCase(a, b string) int {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b))
}

func TestBinarySearchWords(t *testing.T) {
	words, err := util.ReadFileAsWords("../test-data/words.txt")
	assert.Nil(t, err)
	testSearchWords(t, words[:1])
	testSearchWords(t, words[:2])
	testSearchWords(t, words[:3])
	testSearchWords(t, words)
}

func testSearchWords(t *testing.T, words []string) {
	for _, w := range words {
		wT, found := BinarySearch(words, w, compIgnoreCase)
		assert.True(t, found)
		assert.Equal(t, w, wT)
	}
	for _, w := range words {
		_, found := BinarySearch(words, w+"-bogus", strings.Compare)
		assert.False(t, found)
	}
}

type kvp[K comparable, V any] struct {
	key K
	val V
}

func TestBinarySearchKVP(t *testing.T) {
	words, err := util.ReadFileAsWords("../test-data/words.txt")
	assert.Nil(t, err)
	data := make([]kvp[string, string], len(words))
	for i, w := range words {
		data[i] = kvp[string, string]{
			key: strings.ToLower(w),
			val: w,
		}
	}
	testSearchKVP(t, data[:1])
	testSearchKVP(t, data[:2])
	testSearchKVP(t, data[:3])
	testSearchKVP(t, data)
}

func compareKey(a, b kvp[string, string]) int {
	return strings.Compare(a.key, b.key)
}

func testSearchKVP(t *testing.T, data []kvp[string, string]) {
	for _, val := range data {
		valT, found := BinarySearch(data, val, compareKey)
		assert.True(t, found)
		assert.Equal(t, val, valT)
	}
	for _, val := range data {
		var valSearch kvp[string, string]
		valSearch.key = val.key + "-bogus"
		_, found := BinarySearch(data, valSearch, compareKey)
		assert.False(t, found)
	}
}
