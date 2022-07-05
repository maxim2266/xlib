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

import "testing"

func TestStrJoinEx(t *testing.T) {
	sep := [3]string{": ", ", ", " and "}

	cases := []struct {
		res  string
		args []string
	}{
		{"", []string{""}},
		{"AAA", []string{"AAA"}},
		{"AAA: BBB", []string{"AAA", "BBB"}},
		{"AAA: BBB and CCC", []string{"AAA", "BBB", "CCC"}},
		{"AAA: BBB, CCC and DDD", []string{"AAA", "BBB", "CCC", "DDD"}},
		{"AAA: BBB, CCC, DDD and EEE", []string{"AAA", "BBB", "CCC", "DDD", "EEE"}},
		{"AAA: BBB, CCC, DDD, EEE and FFF", []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF"}},
	}

	for i, c := range cases {
		if r := StrJoinEx(sep, c.args...); r != c.res {
			t.Errorf("Unexpected string in case %d: %q instead of %q", i, r, c.res)
			return
		}
	}
}
