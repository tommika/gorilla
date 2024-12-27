// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package heap

import (
	"cmp"
	"os"
	"strings"
	"testing"

	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/util"
)

func orderIgnoreCase(a, b string) int {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b))
}

func orderNatural[T cmp.Ordered](a, b T) int {
	return cmp.Compare(a, b)
}

func isOk(_ any, ok bool) bool {
	return ok
}

func TestNilHeap(t *testing.T) {
	var h *Heap[string]
	h.Describe(os.Stderr)
	assert.Equal(t, 0, h.Size())
	assert.False(t, isOk(h.Peek()))
	assert.False(t, isOk(h.Pop()))
}

func TestEmptyHeap(t *testing.T) {
	h := NewHeap(orderIgnoreCase)
	h.Describe(os.Stderr)
	assert.Equal(t, 0, h.Size())
	assert.False(t, isOk(h.Peek()))
	assert.False(t, isOk(h.Pop()))
}

func TestBuildHeapWithInts(t *testing.T) {
	items := util.MakeIntArray(1000)
	testBuildHeap(t, items, orderNatural)
}

func TestBuildHeapWithWords(t *testing.T) {
	items, err := util.ReadFileAsWords("../test-data/words.txt")
	assert.Nil(t, err)
	testBuildHeap(t, items, orderIgnoreCase)
}

func TestHeapSortWithInts(t *testing.T) {
	items := util.MakeIntArray(1000)
	testHeapSort(t, items, orderNatural)
}

func TestHeapSortWithWords(t *testing.T) {
	items, err := util.ReadFileAsWords("../test-data/words.txt")
	assert.Nil(t, err)
	testHeapSort(t, items, orderIgnoreCase)
}

func TestPriorityQueueWithInts(t *testing.T) {
	items := util.MakeIntArray(1000)
	testPriorityQueue(t, items, orderNatural)
}

func TestPriorityQueueWithWords(t *testing.T) {
	items, err := util.ReadFileAsWords("../test-data/words.txt")
	assert.Nil(t, err)
	testPriorityQueue(t, items, orderIgnoreCase)
}

func testBuildHeap[T comparable](t *testing.T, items []T, compFunc types.Compare[T]) {
	h := BuildHeap(util.ShuffleSlice(util.CopySlice(items)), compFunc)
	h.Describe(os.Stderr)
	assert.Equal(t, len(items), h.Size())
	top, ok := h.Peek()
	assert.True(t, ok)
	assert.Equal(t, items[0], top)
	for i := 0; i < len(items); i++ {
		assert.True(t, h.Size() > 0)
		top, ok := h.Pop()
		assert.True(t, ok)
		assert.Equal(t, items[i], top)
	}
	h.Describe(os.Stderr)
	assert.Equal(t, 0, h.Size())
}

func testHeapSort[T comparable](t *testing.T, items []T, compFunc types.Compare[T]) {
	itemsToSort := util.ShuffleSlice(util.CopySlice(items))
	HeapSort(itemsToSort, compFunc)
	for i := 0; i < len(items); i++ {
		assert.Equal(t, items[i], itemsToSort[i])
	}
}

func testPriorityQueue[T comparable](t *testing.T, items []T, compFunc types.Compare[T]) {
	shuffled := util.ShuffleSlice(util.CopySlice(items))
	q := NewPriorityQueue[T](compFunc)
	assertIsQueue(t, q)
	q.h.Describe(os.Stderr)
	assert.Equal(t, 0, q.Size())
	_, ok := q.Head()
	assert.False(t, ok)
	for i := 0; i < len(shuffled); i++ {
		q.Enqueue(shuffled[i])
	}
	q.MustHead()
	q.h.Describe(os.Stderr)
	assert.Equal(t, len(shuffled), q.Size())
	for i := range items {
		top := q.MustDequeue()
		assert.Equal(t, items[i], top)
	}
	assert.Equal(t, 0, q.Size())
	q.h.Describe(os.Stderr)
	for i := 0; i < len(shuffled); i++ {
		q.Enqueue(shuffled[i])
	}
	q.h.Describe(os.Stderr)
	assert.Equal(t, len(shuffled), q.Size())
}

func assertIsQueue[T any](t *testing.T, q types.Queue[T]) {
	assert.NotNil(t, q)
}
