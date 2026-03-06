package debounce

import (
	"sync"
	"testing"
	"time"
)

func TestDebounceTriggersOnlyOnce(t *testing.T) {

	var count int
	var mu sync.Mutex

	d := New(200)

	fn := func() {
		mu.Lock()
		count++
		mu.Unlock()
	}

	// simulate rapid events
	d.Trigger(fn)
	d.Trigger(fn)
	d.Trigger(fn)
	d.Trigger(fn)

	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if count != 1 {
		t.Fatalf("expected function to run once, got %d", count)
	}
}