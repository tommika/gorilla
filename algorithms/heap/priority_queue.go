package heap

import (
	"github.com/tommika/gorilla/algorithms/types"
	"github.com/tommika/gorilla/must"
)

type PriorityQueue[T any] struct {
	h Heap[T]
}

func NewPriorityQueue[T any](compFunc types.Compare[T]) *PriorityQueue[T] {
	q := &PriorityQueue[T]{}
	q.h.Init(compFunc)
	return q
}

func (q *PriorityQueue[T]) Size() (len int) {
	return q.h.Size()
}
func (q *PriorityQueue[T]) Enqueue(val T) {
	q.h.Push(val)
}
func (q *PriorityQueue[T]) Head() (head T, ok bool) {
	return q.h.Peek()
}
func (q *PriorityQueue[T]) MustHead() (head T) {
	return must.BeOk(q.h.Peek())
}
func (q *PriorityQueue[T]) Dequeue() (head T, ok bool) {
	return q.h.Pop()
}
func (q *PriorityQueue[T]) MustDequeue() (head T) {
	return must.BeOk(q.Dequeue())
}
