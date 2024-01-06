package util

import (
	"context"
)

type Action int

const (
	ActionCreate Action = iota
	ActionUpdate
	ActionDelete
)

type Item struct {
	Action   Action
	ID       int64
	Value    interface{}
	OldValue interface{} //Only in case of update
}

// Queue represents a simple queue structure.
type Queue struct {
	items chan *Item
}

type UpdateQueue interface {
	Enqueue(ctx context.Context, item *Item)
	Dequeue(ctx context.Context) (*Item, bool)
}

func NewQueue() *Queue {
	return &Queue{
		items: make(chan *Item, 100),
	}
}

// Enqueue adds an item to the end of the queue.
func (q *Queue) Enqueue(ctx context.Context, item *Item) {
	select {
	case q.items <- item:
	case <-ctx.Done():
	}
}

// Dequeue removes and returns the item from the front of the queue.
func (q *Queue) Dequeue(ctx context.Context) (*Item, bool) {
	select {
	case v, open := <-q.items:
		return v, open
	case <-ctx.Done():
	}
	return nil, false
}
