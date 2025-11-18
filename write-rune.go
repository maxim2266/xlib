package xlib

import (
	"io"
	"unicode/utf8"
)

// WriteRune writes the given rune to the given [io.Writer].
func WriteRune(w io.Writer, r rune) (err error) {
	var b [utf8.UTFMax]byte

	_, err = w.Write(b[:utf8.EncodeRune(b[:], r)])
	return
}
