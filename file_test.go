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
