package xlib

import (
	"io"
	"strings"
)

// JoinStrings is similar to [strings.Join], but more comfortable to use in some scenarios.
func JoinStrings(sep string, args ...string) string {
	return strings.Join(args, sep)
}

// WriteString writes the given string to the given [io.Writer].
func WriteString(w io.Writer, s string) (err error) {
	if len(s) > 0 {
		_, err = io.WriteString(w, s)
	}

	return
}

// WriteByte writes the given byte to the given [io.Writer].
func WriteByte(w io.Writer, b byte) (err error) {
	m := [1]byte{b}

	_, err = w.Write(m[:])
	return
}
