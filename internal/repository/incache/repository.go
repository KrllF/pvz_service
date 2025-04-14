package incache

import (
	"sync"
)

type (
	Repo[T any] struct {
		mu         sync.RWMutex
		Head, Tail *Node[T]
		Storage    map[int64]*Node[T]
		Capacity   int64
	}
	Node[T any] struct {
		Prev, Next *Node[T]
		Key        int64
		Value      T
	}
)

func NewNode[T any](key int64, value T) *Node[T] {
	return &Node[T]{
		Key:   key,
		Value: value,
	}
}

func NewRepository[T any](capacity int64) *Repo[T] {
	head, tail := NewNode(0, *new(T)), NewNode(0, *new(T))
	head.Next = tail
	tail.Prev = head
	return &Repo[T]{
		mu:       sync.RWMutex{},
		Head:     head,
		Tail:     tail,
		Storage:  make(map[int64]*Node[T], capacity),
		Capacity: capacity,
	}
}
