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

	if err := Await(errch); err == nil || err.Error() != "just some error" {
		t.Errorf("unexpected error: %q", err)
		return
	}

	if counter != 1 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}

func TestAsyncN(t *testing.T) {
	counter := int32(0)

	fn := func() error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
		return nil
	}

	if err := Await(Async(fn, fn, fn)); err != nil {
		t.Errorf("unexpected error(s): %s", err)
		return
	}

	if counter != 3 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}

func TestAsyncErrN(t *testing.T) {
	counter := int32(0)

	fn := func() error {
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
		return errors.New("X")
	}

	err := Await(Async(fn, fn, fn))

	if err == nil {
		t.Error("missing error")
		return
	}

	if msg := err.Error(); msg != "X\nX\nX" {
		t.Errorf("unexpected error message: %q", msg)
		return
	}

	if counter != 3 {
		t.Errorf("unexpected counter: %d", counter)
		return
	}
}
