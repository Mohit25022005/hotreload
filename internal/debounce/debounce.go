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

	// stop previous timer safely
	if d.timer != nil {
		if !d.timer.Stop() {
			select {
			case <-d.timer.C:
			default:
			}
		}
	}

	// create new timer
	d.timer = time.AfterFunc(d.delay, func() {
		fn()
	})
}

func (d *Debouncer) Cancel() {

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}