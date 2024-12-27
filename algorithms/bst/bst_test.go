// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package bst

import (
	"cmp"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/util"
)

// keyFunc is a function type that determines the key for the given value
type keyFunc[K cmp.Ordered, V comparable] func(val V) K

func TestWithNilRoot(t *testing.T) {
	bst := NewBST[string, string](false)
	assert.Nil(t, bst.root)
	assert.False(t, util.IsOk(bst.Get("fred")))
	assert.False(t, util.IsOk(bst.Min()))
	assert.False(t, util.IsOk(bst.Max()))
	assert.Equal(t, 0, bst.count())
	bst.Put("fred", "Fred")
	assert.NotNil(t, bst.root)
	assert.True(t, util.IsOk(bst.Get("fred")))
}

func TestCase2AtRoot(t *testing.T) {
	bst := NewBST[string, types.Unit](false)
	bst.Put("A", types.Nothing)
	bst.Put("B", types.Nothing)
	deleted := bst.Delete("A")
	assert.True(t, deleted)
	assert.True(t, util.IsOk(bst.Get("B")))
}

func TestWithWords(t *testing.T) {
	// Read test words; guaranteed to be sorted
	words, err := util.ReadFileAsWords("../test-data/words.txt")
	assert.Nil(t, err)
	testBST(t, false, words, func(val string) string {
		return strings.ToLower(val)
	})
	testBST(t, true, words, func(val string) string {
		return strings.ToLower(val)
	})
}

func TestWithInts(t *testing.T) {
	ints := util.MakeIntArray(10000)
	testBST(t, false, ints, func(val int) int {
		return val
	})
	testBST(t, true, ints, func(val int) int {
		return val
	})
}

func testBST[K cmp.Ordered, V comparable](t *testing.T, balance bool, vals []V, keyFunc keyFunc[K, V]) {
	var bst = NewBST[K, V](balance)
	// Shuffle values and insert into the tree.
	// We shuffle the values so that we get a more balanced tree, something with
	// a height close to lg N (where N is the number of nodes in the tree.)
	shuffled := util.ShuffleSlice(util.CopySlice(vals))
	fmt.Printf("inserting values into tree\n")
	for _, v := range shuffled {
		bst.Put(keyFunc(v), v)
	}
	assert.Nil(t, bst.validate())
	fmt.Printf("tree is valid after insertions\n")
	// Print some stats about the tree
	count := bst.count()
	height := bst.height()
	fmt.Printf("count=%d, height=%d, log2(count)=%d\n", count, height, int(math.Log2(float64(count))))

	// Test that the counts match
	assert.Equal(t, len(vals), bst.Size())
	assert.Equal(t, len(vals), bst.count())

	// Test Min
	min, ok := bst.Min()
	assert.True(t, ok)
	assert.Equal(t, keyFunc(vals[0]), min)

	// Test Max
	max, ok := bst.Max()
	assert.True(t, ok)
	assert.Equal(t, keyFunc(vals[len(vals)-1]), max)

	// Test that we can find every key and that the value is preserved
	for _, v := range vals {
		vT := bst.MustGet(keyFunc(v))
		assert.Equal(t, v, vT)
	}
	// Test that the tree is sorted by visiting every node in order,
	// and comparing to the original word list.
	myCount := 0
	bst.VisitInOrder(func(k K, v V) {
		assert.Equal(t, k, keyFunc(v))
		assert.Equal(t, vals[myCount], v)
		myCount++
	})
	assert.Equal(t, len(vals), myCount)
	assert.Equal(t, len(vals), bst.Size())

	// Test that we can update values without changing the shape of the tree
	oldRoot := bst.root
	for _, v := range vals {
		k := keyFunc(v)
		bst.Put(k, v)
		// The root should not change
		assert.Equal(t, oldRoot, bst.root)
	}
	// Should still be a valid true
	assert.Nil(t, bst.validate())
	// Count should not have changed
	assert.Equal(t, count, bst.count())
	assert.Equal(t, count, bst.Size())
	// Height should not have changed
	assert.Equal(t, height, bst.height())

	// Delete half of the items
	var deleted bool
	for _, v := range shuffled[count/2:] {
		k := keyFunc(v)
		// Deleting from tree; root might change
		deleted = bst.Delete(k)
		assert.True(t, deleted)
		// Already deleted; should return false
		deleted = bst.Delete(k)
		assert.False(t, deleted)
	}

	// Add them back
	for _, v := range shuffled[count/2:] {
		bst.Put(keyFunc(v), v)
	}

	// Test that we can find every key and that the value is preserved
	for _, v := range vals {
		vT := bst.MustGet(keyFunc(v))
		assert.Equal(t, v, vT)
	}

	// Delete all of the items
	for _, v := range shuffled {
		k := keyFunc(v)
		// Deleting from tree; root might change
		deleted = bst.Delete(k)
		assert.True(t, deleted)
		// Already deleted; should return false
		deleted = bst.Delete(k)
		assert.False(t, deleted)
	}
	// Tree should be empty
	assert.Equal(t, 0, bst.Size())
	assert.Equal(t, 0, bst.count())
	assert.Nil(t, bst.root)
}

func TestBSTContainerInterface(t *testing.T) {
	var bst = NewBST[int, int](true)
	assertIsMap(t, bst)
	assertIsOrdered(t, bst)
}

func assertIsMap[K cmp.Ordered, V any](t *testing.T, a types.Map[K, V]) {
	assert.NotNil(t, a)
}

func assertIsOrdered[K cmp.Ordered, V any](t *testing.T, a types.Ordered[K, V]) {
	assert.NotNil(t, a)
}

func TestInvalid(t *testing.T) {
	bst := NewBST[int, int](true)
	bst.Put(1, 1)
	bst.root.col = red
	assert.NotNil(t, bst.validate())

	n := newNode(1, 1)
	assert.NotNil(t, n.validateRBNode())

	_, err := n.blackHeight()
	assert.NotNil(t, err)

	n0 := newNode(1, 0)
	n0.col = red
	n1 := newNode(0, 1)
	n1.col = red
	n2 := newNode(2, 2)
	n2.col = red
	n0.left = n1
	n0.right = n2
	_, err = n0.blackHeight()
	assert.NotNil(t, err)

	n0.col = black
	n1.col = black
	_, err = n0.blackHeight()
	assert.NotNil(t, err)
}

func TestStringer(t *testing.T) {
	path := nodeList[int, int]{}
	path = append(path, newNode(1, 1))
	path = append(path, newNode(2, 2))
	t.Logf("path: %v", path)
}
