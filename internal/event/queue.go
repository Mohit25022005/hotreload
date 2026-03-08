package event

import (
	"context"
	"sync"
)

type Queue struct {
	ch     chan struct{}
	closed bool
	mu     sync.Mutex
}

func NewQueue() *Queue {

	return &Queue{
		ch: make(chan struct{}, 1),
	}
}

func (q *Queue) Push() {

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	select {

	case q.ch <- struct{}{}:

	default:
		// queue already has an event
	}
}

func (q *Queue) Events() <-chan struct{} {
	return q.ch
}

func (q *Queue) Wait(ctx context.Context) bool {

	select {

	case <-ctx.Done():
		return false

	case <-q.ch:
		return true
	}
}

func (q *Queue) Close() {

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	close(q.ch)
	q.closed = true
}