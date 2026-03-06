package debounce

import (
	"sync"
	"time"
)

type Debouncer struct {
	delay time.Duration
	timer *time.Timer
	mu    sync.Mutex
}

func New(ms int) *Debouncer {
	return &Debouncer{
		delay: time.Duration(ms) * time.Millisecond,
	}
}

func (d *Debouncer) Trigger(fn func()) {

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.delay, fn)
}