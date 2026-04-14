package xlib

import "iter"

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

// PipeIt creates a new iterator that runs the original iterator in
// a dedicated goroutine.
func PipeIt[T any](src iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		// channels
		pipe, done := make(chan T, 10), make(chan struct{})
		stopped := false

		defer func() {
			close(done)

			if stopped {
				// drain the pipe
				for range pipe {
					// do nothing
				}
			}
		}()

		// feeder
		go func() {
		loop:
			for v := range src {
				select {
				case pipe <- v:
				case <-done:
					break loop
				}
			}

			// we don't want to close the channel on panic, because
			// that would allow the reader loop to exit while the panic
			// is still in progress
			close(pipe)
		}()

		// reader loop
		for v := range pipe {
			if stopped = !yield(v); stopped {
				break
			}
		}
	}
}
