package queue

import "github.com/tommika/gorilla/must"

// SliceQueue is a first-in/first-out collection of elements implemented directly using
// Go's underlying slice and automatic garbage collection. This is by far the
// simplest of the implementations here, with performance comparable to a
// dynamic-circular-array.
type SliceQueue[T any] struct {
	items []T
}

func (q *SliceQueue[T]) Size() (size int) {
	return len(q.items)
}

func (q *SliceQueue[T]) Enqueue(val T) {
	q.items = append(q.items, val)
}

func (q *SliceQueue[T]) Head() (head T, ok bool) {
	if q == nil || len(q.items) == 0 {
		return
	}
	return q.items[0], true
}

func (q *SliceQueue[T]) MustHead() (head T) {
	return must.BeOk(q.Head())
}

func (q *SliceQueue[T]) Dequeue() (head T, ok bool) {
	if q == nil || len(q.items) == 0 {
		return
	}
	head, q.items = q.items[0], q.items[1:]
	return head, true
}

func (q *SliceQueue[T]) MustDequeue() (head T) {
	return must.BeOk(q.Dequeue())
}
