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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFileWrite(t *testing.T) {
	tests := []func(string) error{
		testSimpleWrite,
		testOverWrite,
		testPanic,
		testError,
	}

	for i, test := range tests {
		if err := run(test); err != nil {
			t.Error(i, err)
			return
		}
	}
}

func TestFileWriteAsync(t *testing.T) {
	errch := Async(
		func() error { return run(testSimpleWrite) },
		func() error { return run(testOverWrite) },
		func() error { return run(testPanic) },
		func() error { return run(testError) },
	)

	for err := range errch {
		t.Error(err)
	}
}

func testSimpleWrite(dir string) error {
	fname := filepath.Join(dir, "zzz")

	err := WriteFile(fname, func(w *bufio.Writer) error {
		_, e := w.WriteString(fileContent)
		return e
	})

	if err != nil {
		return err
	}

	if err = assertOneFile(dir, "zzz"); err != nil {
		return err
	}

	s, err := os.ReadFile(fname)

	if err != nil {
		return err
	}

	if !bytes.Equal(s, []byte(fileContent)) {
		return fmt.Errorf("Unexpected file content: %q", string(s))
	}

	return nil
}

func testOverWrite(dir string) error {
	fname := filepath.Join(dir, "zzz")

	if err := os.WriteFile(fname, []byte("XXX"), 0666); err != nil {
		return err
	}

	return testSimpleWrite(dir)
}

const fileContent = "Hello, world!"

func testPanic(dir string) (err error) {
	type myPanic int

	fname := filepath.Join(dir, "zzz")

	defer func() {
		p := recover()

		if p == nil {
			err = errors.New("No panic detected")
			return
		}

		n, ok := p.(myPanic)

		if !ok {
			err = errors.New("Unknown panic type")
			return
		}

		if n != 123 {
			err = errors.New("Unknown panic value")
			return
		}

		err = assertNoFiles(dir)
	}()

	err = WriteFile(fname, func(_ *bufio.Writer) error {
		panic(myPanic(123))
	})

	return
}

func testError(dir string) error {
	const errMsg = "Test error"

	fname := filepath.Join(dir, "zzz")

	err := WriteFile(fname, func(_ *bufio.Writer) error {
		return errors.New(errMsg)
	})

	if err == nil {
		return errors.New("Missing error")
	}

	if err.Error() != errMsg {
		return fmt.Errorf("Unexpected error: %q", err)
	}

	return assertNoFiles(dir)
}

func run(fn func(string) error) error {
	dir, err := os.MkdirTemp("", "xlib-")

	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	return fn(dir)
}

func assertNoFiles(dir string) error {
	files, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	if len(files) != 0 {
		return fmt.Errorf("Unexpected file found: %q", files[0])
	}

	return nil
}

func assertOneFile(dir, fname string) error {
	files, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("No files found")
	}

	for _, file := range files {
		if file.Name() != fname {
			return fmt.Errorf("Unexpected file: %q", file.Name())
		}
	}

	return nil
}
