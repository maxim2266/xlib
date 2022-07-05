/*
Copyright (c) 2018,2019,2022 Maxim Konakov
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
	"unsafe"
)

// StrJoin is similar to strings.Join(), but more comfortable to use in some scenarios.
func StrJoin(sep string, args ...string) string {
	switch len(args) {
	case 0:
		return ""
	case 1:
		return args[0]
	}

	// total length
	n := len(sep) * (len(args) - 1)

	for _, s := range args {
		n += len(s)
	}

	// compose
	b := make([]byte, 0, n)

	if len(sep) == 0 {
		for _, s := range args {
			b = append(b, s...)
		}
	} else {
		b = append(b, args[0]...)

		for _, s := range args[1:] {
			b = append(append(b, sep...), s...)
		}
	}

	// all done
	return *(*string)(unsafe.Pointer(&b))
}

// StrJoinEx joins the given arguments using the supplied separators, for example,
// given the separator list of [": ", ", ", ", and "] and the arguments
// ["AAA", "BBB", "CCC", "DDD"], the resulting string will be "AAA: BBB, CCC, and DDD".
func StrJoinEx(sep [3]string, args ...string) string {
	// process short argument lists
	switch len(args) {
	case 0:
		return ""
	case 1:
		return args[0]
	case 2:
		return args[0] + sep[0] + args[1]
	case 3:
		return args[0] + sep[0] + args[1] + sep[2] + args[2]
	case 4:
		return args[0] + sep[0] + args[1] + sep[1] + args[2] + sep[2] + args[3]
	}

	// total length
	n := len(sep[0]) + len(sep[1])*(len(args)-2) + len(sep[2])

	for _, s := range args {
		n += len(s)
	}

	// compose
	buff := append(append(append(make([]byte, 0, n), args[0]...), sep[0]...), args[1]...)

	for _, s := range args[2 : len(args)-1] {
		buff = append(append(buff, sep[1]...), s...)
	}

	buff = append(append(buff, sep[2]...), args[len(args)-1]...)

	// all done
	return *(*string)(unsafe.Pointer(&buff))
}
