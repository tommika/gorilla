// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package types

// Queue is the interface to a first-in/first-out container of items.
type Queue[T any] interface {
	Size() (len int)
	Enqueue(val T)
	Head() (head T, ok bool)
	MustHead() (head T)
	Dequeue() (head T, ok bool)
	MustDequeue() (head T)
}
