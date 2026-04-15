package xlib

import (
	"context"
	"iter"
	"sync"
)

// MapIt creates a new iterator that converts each item from the original
// iterator via the given function.
func MapIt[T, U any](src iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range src {
			if !yield(fn(v)) {
				break
			}
		}
	}
}

// FilterIt creates a new iterator that filters items from the original
// iterator using the supplied predicate.
func FilterIt[T any](src iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range src {
			if pred(v) && !yield(v) {
				break
			}
		}
	}
}

// TakeWhileIt creates a new iterator that yields values while the given
// predicate returns true.
func TakeWhileIt[T any](src iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range src {
			if !pred(v) || !yield(v) {
				break
			}
		}
	}
}

// PipeIt creates a new iterator that runs the original iterator in
// a dedicated goroutine.
func PipeIt[T any](src iter.Seq[T]) iter.Seq[T] {
	return PipeItCtx(context.Background(), src)
}

// PipeIt creates a new iterator that runs the original iterator in
// a dedicated goroutine, with lifetime-controlling context.
func PipeItCtx[T any](ctx context.Context, src iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		ctx, cancel := context.WithCancel(ctx)

		// channel
		pipe := make(chan T, 10)

		// waiter
		var wg sync.WaitGroup

		defer func() {
			cancel()
			wg.Wait()
		}()

		// feeder
		wg.Go(func() {
		loop:
			for v := range src {
				select {
				case pipe <- v:
				case <-ctx.Done():
					break loop
				}
			}

			// closing the channel on panic would allow the reader loop to
			// exit while the panic is still in progress, so we only close
			// on normal exit from the function (ignoring runtime.Goexit)
			close(pipe)
		})

		// reader loop
		for v := range pipe {
			if !yield(v) {
				break
			}
		}
	}
}
