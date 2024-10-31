package xlib

import (
	"errors"
	"sync/atomic"
)

// Async starts each of the given functions in a separate goroutine, and returns
// a channel where errors from the functions are posted. The channel is closed
// when all goroutines have completed.
func Async(tasks ...func() error) <-chan error {
	if len(tasks) == 0 {
		panic("xlib.Async: no task to run")
	}

	errch := make(chan error, len(tasks))
	count := int32(len(tasks))

	for _, fn := range tasks {
		go func(fn func() error) {
			defer func() {
				if atomic.AddInt32(&count, -1) == 0 {
					close(errch)
				}
			}()

			if err := fn(); err != nil {
				errch <- err
			}
		}(fn)
	}

	return errch
}

// Await takes the error channel returned from Async() function, and waits for all the tasks
// to complete, collecting errors via errors.Join().
func Await(errch <-chan error) error {
	var errs []error

	for err := range errch {
		errs = append(errs, err)
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errors.Join(errs...)
	}
}
