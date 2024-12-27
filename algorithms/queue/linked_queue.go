// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package queue

import (
	"github.com/tommika/gorilla/must"
)

// LinkedQueue is a first-in/first-out collection of elements implemented as a
// linked-list.
type LinkedQueue[T any] struct {
	head *lqElem[T]
	tail *lqElem[T]
	len  int
}

// lqElem is an element on the linked queue
type lqElem[T any] struct {
	next *lqElem[T]
	val  T
}

// Length returns the current length of the queue
func (q *LinkedQueue[T]) Size() int {
	return q.len
}

// Enqueue adds an element to the end of the queue
func (q *LinkedQueue[T]) Enqueue(val T) {
	e := &lqElem[T]{
		next: nil,
		val:  val,
	}
	if q.tail != nil {
		q.tail.next = e
	} else {
		// was empty; update head
		q.head = e
	}
	// update tail and length in all cases
	q.tail = e
	q.len++
}

// Head returns the element at the front of the queue. If the queue
// is empty, the zero-value fot the queue type is returned, and
// ok is false.
func (q *LinkedQueue[T]) Head() (head T, ok bool) {
	if q.head == nil {
		return
	}
	return q.head.val, true
}

func (q *LinkedQueue[T]) MustHead() T {
	return must.BeOk(q.Head())
}

// Dequeue removes and returns the element at the front of the queue. If the
// queue is empty, the zero-value fot the queue type is returned, and ok is
// false.
func (q *LinkedQueue[T]) Dequeue() (head T, ok bool) {
	if q.head == nil {
		return
	}
	head = q.head.val
	q.head = q.head.next
	if q.head == nil {
		// now empty
		q.tail = nil
	}
	q.len--
	return head, true
}

// Dequeue removes and returns the element at the front of the queue. Panics if
// the queue is empty.
func (q *LinkedQueue[T]) MustDequeue() (head T) {
	return must.BeOk(q.Dequeue())
}
