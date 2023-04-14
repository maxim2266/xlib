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
	"context"
	"errors"
)

// Pump is a function that iterates or generates a sequence of items of type T, and
// calls the given callback once per each item from the sequence. It is expected
// to stop on the first error encountered either from the iteration itself, or returned
// from the callback.
type Pump[T any] func(func(T) error) error

// PipelinedPump creates a pipelined version of the given pump, where the pump itself
// is running in a separate goroutine, while the calling goroutine only does the callback
// invocations.
func PipelinedPump[T any](pump Pump[T]) Pump[T] {
	return func(fn func(T) error) (err error) {
		queue := make(chan T, 20)
		errch := make(chan error, 1)
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			defer func() {
				close(queue)
				close(errch)
			}()

			err := pump(func(item T) error {
				select {
				case queue <- item:
					return nil
				case <-ctx.Done():
					return cancelledPumpError // just to stop the pump
				}
			})

			if err != nil && ctx.Err() == nil {
				errch <- err
			}
		}()

		for item := range queue {
			if err := fn(item); err != nil {
				cancel()
				<-errch // wait for the pump to stop
				return err
			}
		}

		return <-errch
	}
}

var cancelledPumpError = errors.New("pump cancelled")
