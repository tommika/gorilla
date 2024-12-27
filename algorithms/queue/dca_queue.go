// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package queue

import (
	"github.com/tommika/gorilla/must"
)

// DynamicCircularArrayQueue is a first-in/first-out collection of elements implemented
// using an underlying circular array/slice.
type DynamicCircularArrayQueue[T any] struct {
	items []T // storage for items on the queue
	len   int // length of the queue
	head  int // index of item at head of queue
	free  int // index to next free slot on queue
}

// NewDynamicCircularArrayQueue creates a new queue with the given initial capacity
func NewDynamicCircularArrayQueue[T any](initCapacity int) *DynamicCircularArrayQueue[T] {
	return &DynamicCircularArrayQueue[T]{
		items: make([]T, initCapacity),
	}
}

func (q *DynamicCircularArrayQueue[T]) Size() (len int) {
	if q != nil {
		len = q.len
	}
	return
}

func (q *DynamicCircularArrayQueue[T]) Enqueue(val T) {
	if q.items == nil {
		// auto initialize
		q.items = make([]T, 1)
	} else if q.len == len(q.items) {
		// reached capacity; need to expand
		must.BeEqual(q.head, q.free)
		must.BeTrue(q.head < q.len)
		must.BeEqual(len(q.items), cap(q.items))
		// double capacity
		items := make([]T, q.len*2)
		must.BeEqual(len(items), cap(items))
		// copy existing items into new array
		if q.head > 0 {
			// copy tail items to new array (if needed)
			n := copy(items[q.len:q.len+q.head], q.items[:q.head])
			must.BeEqual(q.head, n)
		}
		// copy head items to new array
		n := copy(items[q.head:q.len], q.items[q.head:q.len])
		must.BeEqual(q.len-q.head, n)
		// update state
		q.items = items
		q.free += q.len
	}
	q.items[q.free] = val
	q.free = (q.free + 1) % len(q.items)
	q.len++
}

func (q *DynamicCircularArrayQueue[T]) Head() (head T, ok bool) {
	if q == nil || q.len == 0 {
		return
	}
	return q.items[q.head], true
}

func (q *DynamicCircularArrayQueue[T]) MustHead() T {
	return must.BeOk(q.Head())
}

func (q *DynamicCircularArrayQueue[T]) Dequeue() (head T, ok bool) {
	if q == nil || q.len == 0 {
		return
	}
	head = q.items[q.head]
	q.head = (q.head + 1) % len(q.items)
	q.len--
	return head, true
}

// MustDequeue removes and returns the element at the front of the queue. Panics if
// the queue is empty.
func (q *DynamicCircularArrayQueue[T]) MustDequeue() (head T) {
	return must.BeOk(q.Dequeue())
}
