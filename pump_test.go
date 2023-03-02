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
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestPump(t *testing.T) {
	const N = 1000_000

	count := 0

	err := Pump(makeCountingSource(N), func(i int) error {
		if i != count {
			return fmt.Errorf("unexpected parameter: %d instead of %d", i, count)
		}

		count++
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	if count != N {
		t.Errorf("unexpected final value: %d instead of %d", count, N)
		return
	}
}

func TestPumpError(t *testing.T) {
	const N = 1000

	count := 0

	err := Pump(makeCountingSource(N), func(i int) error {
		if count >= N/2 {
			return fmt.Errorf("unexpected call with value %d", i)
		}

		if i != count {
			return fmt.Errorf("unexpected parameter: %d instead of %d", i, count)
		}

		if count++; count == N/2 {
			return fmt.Errorf("expected error: reached value %d", count)
		}

		return nil
	})

	if err == nil {
		t.Error("missing expected error")
		return
	}

	t.Log(err)

	if count != N/2 {
		t.Errorf("unexpected final value: %d instead of %d", count, N)
		return
	}
}

func TestPumpSourceError(t *testing.T) {
	const N = 1000

	src := makeCountingSourceErr(N, errors.New("expected error"))

	count := 0

	err := Pump(src, func(i int) error {
		if i != count {
			return fmt.Errorf("unexpected parameter: %d instead of %d", i, count)
		}

		count++
		return nil
	})

	if err == nil {
		t.Error("missing expected error")
		return
	}

	t.Log(err)

	if count != N {
		t.Errorf("unexpected final value: %d instead of %d", count, N)
		return
	}
}

func makeCountingSource(N int) func() (int, error) {
	return makeCountingSourceErr(N, io.EOF)
}

func makeCountingSourceErr(N int, stop error) func() (int, error) {
	counter := 0

	return func() (int, error) {
		if counter < N {
			ret := counter

			counter++
			return ret, nil
		}

		return 0, stop
	}
}
