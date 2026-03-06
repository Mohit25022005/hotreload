package event

type Queue struct {
	ch chan struct{}
}

func NewQueue() *Queue {

	return &Queue{
		ch: make(chan struct{}, 1),
	}
}

func (q *Queue) Push() {

	select {

	case q.ch <- struct{}{}:

	default:
	}
}

func (q *Queue) Events() <-chan struct{} {
	return q.ch
}