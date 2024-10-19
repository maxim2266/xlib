package xlib

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	counter := int32(0)

	errch := Async(func() error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
		return nil
	})

	if err := <-errch; err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if counter != 1 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}

func TestAsyncErr(t *testing.T) {
	counter := int32(0)

	errch := Async(func() error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
		return errors.New("just some error")
	})

	if err := <-errch; err == nil || err.Error() != "just some error" {
		t.Errorf("unexpected error: %q", err)
		return
	}

	if err := <-errch; err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if counter != 1 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}

func TestParallel(t *testing.T) {
	counter := int32(0)

	fn := func() error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
		return nil
	}

	errch := Parallel(fn, fn, fn)

	if err := <-errch; err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if counter != 3 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}

func TestAsyncArg(t *testing.T) {
	counter := int32(0)

	errch := AsyncArg(&counter, func(p *int32) error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(p, 1)
		return nil
	})

	if err := <-errch; err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if counter != 1 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}

func TestParallelArg(t *testing.T) {
	counter := int32(0)

	fn := func(p *int32) error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(p, 1)
		return nil
	}

	errch := ParallelArg(&counter, fn, fn, fn)

	if err := <-errch; err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if counter != 3 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}
