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

func TestDebounceDelayedExecution(t *testing.T) {

	var called bool
	d := New(200)

	d.Trigger(func() {
		called = true
	})

	time.Sleep(100 * time.Millisecond)

	if called {
		t.Fatal("function executed too early")
	}

	time.Sleep(150 * time.Millisecond)

	if !called {
		t.Fatal("function did not execute after delay")
	}
}

func TestDebounceCancel(t *testing.T) {

	var called bool

	d := New(200)

	d.Trigger(func() {
		called = true
	})

	d.Cancel()

	time.Sleep(300 * time.Millisecond)

	if called {
		t.Fatal("function should not run after cancel")
	}
}

func TestDebounceMultipleTriggers(t *testing.T) {

	var count int
	var mu sync.Mutex

	d := New(100)

	fn := func() {
		mu.Lock()
		count++
		mu.Unlock()
	}

	for i := 0; i < 10; i++ {
		d.Trigger(fn)
	}

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if count != 1 {
		t.Fatalf("expected single execution, got %d", count)
	}
}