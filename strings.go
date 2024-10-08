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
	"io"
	"strings"
	"unicode/utf8"
)

// StrJoin is similar to strings.Join(), but more comfortable to use in some scenarios.
func StrJoin(sep string, args ...string) string {
	return strings.Join(args, sep)
}

// WriteString writes the given string to the given io.Writer.
func WriteString(w io.Writer, s string) (err error) {
	if len(s) > 0 {
		_, err = io.WriteString(w, s)
	}

	return
}

// WriteRune writes the given rune to the given io.Writer.
func WriteRune(w io.Writer, r rune) (err error) {
	var b [utf8.UTFMax]byte

	_, err = w.Write(b[:utf8.EncodeRune(b[:], r)])
	return
}

// WriteByte writes the given byte to the given io.Writer.
func WriteByte(w io.Writer, b byte) (err error) {
	m := [1]byte{b}

	_, err = w.Write(m[:])
	return
}
