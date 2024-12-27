// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package queue

import (
	"testing"

	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/assert"
)

func TestQueue(t *testing.T) {
	testQueue(t, 16, 3, &LinkedQueue[int]{})
	testQueue(t, 16, 3, &DynamicCircularArrayQueue[int]{})
	testQueue(t, 16, 3, &SliceQueue[int]{})
}

func BenchmarkLinkedQueue(b *testing.B) {
	benchQueue(b, 10000, 1000, &LinkedQueue[int]{})
}

func BenchmarkDynamicCircularArrayQueue(b *testing.B) {
	benchQueue(b, 10000, 1000, &DynamicCircularArrayQueue[int]{})
}

func BenchmarkSliceQueue(b *testing.B) {
	benchQueue(b, 10000, 1000, &DynamicCircularArrayQueue[int]{})
}

func TestSliceQueueExpansion(t *testing.T) {
	N := 4
	q := NewDynamicCircularArrayQueue[int](N * 2)
	for i := 0; i < N; i++ {
		q.Enqueue(i)
	}
	for i := 0; i < N; i++ {
		val, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}
	// Test expansion when head index non-zero
	assert.True(t, q.head != 0)
	testQueue(t, 32, 32, q)
}

func benchQueue(b *testing.B, n, m int, q types.Queue[int]) {
	for i := 0; i < b.N; i++ {
		testQueue(b, n, m, q)
	}
}

// testQueue tests a queue implementation.
// The given queue must be empty.
// * n is the total number of items to enqueue / dequeue
// * m is the max allowable queue length
func testQueue(t testing.TB, n, m int, q types.Queue[int]) {
	assert.Equal(t, 0, q.Size())
	_, ok := q.Dequeue()
	assert.False(t, ok)
	_, ok = q.Head()
	assert.False(t, ok)

	assert.Equal(t, 0, q.Size())

	assert.Equal(t, 0, q.Size())

	for i, j := 0, 0; i < n || j < n; {
		if i < n {
			q.Enqueue(i)
			i++
		}
		if i >= m && j < n {
			val := q.MustHead()
			assert.Equal(t, j, val)
			val = q.MustDequeue()
			assert.Equal(t, j, val)
			j++
		}
		assert.True(t, q.Size() <= m)
	}
	assert.Equal(t, 0, q.Size())
}
