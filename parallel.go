/*
Copyright (c) 2018,2019,2022,2023 Maxim Konakov
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
  this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.
3. Neither the name of the copyright holder nor the names of its contributors
  may be used to endorse or promote products derived from this software without
  specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY
OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package xlib

import (
	"sync/atomic"
)

// Parallel starts the given functions each in a separate goroutine, and returns
// a channel where errors from the functions are posted. The channel is closed
// when all goroutines have completed.
func Parallel(tasks ...func() error) <-chan error {
	if len(tasks) == 0 {
		panic("xlib.Parallel: no tasks to run")
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

// Async starts the given function in a dedicated goroutine and returns
// a channel where any error from the function will be posted. The channel is
// closed when the function has completed.
func Async(fn func() error) <-chan error {
	return Parallel(fn)
}

// ParallelArg starts the given functions each in a separate goroutine, and returns
// a channel where errors from the functions are posted. The channel is closed
// when all goroutines have completed. The given argument is passed over to all
// functions. Hint: the argument may be a context that all the functions check
// time to time, and terminate when it's cancelled.
func ParallelArg[T any](arg T, tasks ...func(T) error) <-chan error {
	if len(tasks) == 0 {
		panic("xlib.ParallelArg: no tasks to run")
	}

	errch := make(chan error, len(tasks))
	count := int32(len(tasks))

	for _, fn := range tasks {
		go func(fn func(T) error) {
			defer func() {
				if atomic.AddInt32(&count, -1) == 0 {
					close(errch)
				}
			}()

			if err := fn(arg); err != nil {
				errch <- err
			}
		}(fn)
	}

	return errch
}

// AsyncArg starts the given function in a dedicated goroutine and returns
// a channel where any error from the function will be posted. The channel is
// closed when the function has completed. The given argument is passed over to all
// functions. Hint: the argument may be a context that the function checks
// time to time, and terminates when it's cancelled.
func AsyncArg[T any](arg T, fn func(T) error) <-chan error {
	return ParallelArg(arg, fn)
}
