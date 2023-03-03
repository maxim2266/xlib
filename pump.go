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

import "io"

/*
Pump repeatedly calls `src` function to obtain items for processing until the function returns
an error, and forwards each item to the `dest` function running in a separate goroutine.
This forms a processing conveyor where the input data are retrieved in one goroutine and
processed in another. In case of an error in either source or destination the processing stops
at the first error encountered and the error gets returned to the caller. Error io.EOF from
the source is treated as a signal to stop the iteration and returned as `nil`. Since `dest`
function is invoked from a different goroutine any data shared between `src` and `dest` should
be protected by a mutex. Upon return from Pump it is guaranteed, that the processing goroutine
has fully completed its job.
*/
func Pump[T any](src func() (T, error), dest func(T) error) (err error) {
	// error channel
	errch := make(chan error, 1)

	// pump
	if err = pump(src, dest, errch); err == nil {
		err = <-errch
	} else {
		<-errch
	}

	return
}

func pump[T any](src func() (T, error), dest func(T) error, errch chan error) (err error) {
	// start `dest` pump and get the write end of the work queue
	queue := startPump(dest, errch)

	defer close(queue)

	// feed the pump
	var item T

	for item, err = src(); err == nil; item, err = src() {
		select {
		case queue <- item:
			// ok
		case err = <-errch:
			return
		}
	}

	// io.EOF from source means end of input
	if err == io.EOF {
		err = nil
	}

	return
}

func startPump[T any](dest func(T) error, errch chan<- error) chan<- T {
	// work queue
	queue := make(chan T, 20)

	// processor
	go func() {
		defer close(errch)

		for item := range queue {
			if err := dest(item); err != nil {
				errch <- err
				break
			}
		}
	}()

	return queue
}
