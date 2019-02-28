/*
Copyright (c) 2018,2019 Maxim Konakov
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
	"io"
	"io/ioutil"
	"os/exec"
)

// ScanCommandOutput invokes the specified command, reads the command's output (stdout),
// breaks it into tokens using the supplied split function, and calls the given callback function
// once per each token. If the split function is set to nil then the output is split on new-line characters.
// Internally the function is implemented using bufio.Scanner, please refer to its documentation
// for more details on split functions.
func ScanCommandOutput(cmd *exec.Cmd, sf bufio.SplitFunc, fn func([]byte) error) (ret int, err error) {
	ret = -1

	// stdout
	var stdout io.ReadCloser

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}

	// stderr
	var stderr limitedWriter

	stderr.limit = 4 * 1024 // accept only up to 4K
	cmd.Stderr = &stderr

	// start the command
	if err = cmd.Start(); err != nil {
		return
	}

	// stdout reader
	src := bufio.NewScanner(bufio.NewReader(stdout))

	if sf != nil {
		src.Split(sf)
	}

	// iterate
scanner:
	for src.Scan() {
		s := src.Bytes()

		switch err = fn(s[:len(s):len(s)]); err {
		case nil:
			// ok
		case io.EOF: // not an error
			_, err = io.Copy(ioutil.Discard, stdout)
			break scanner
		default:
			io.Copy(ioutil.Discard, stdout)
			break scanner
		}
	}

	if err == nil {
		err = src.Err()
	}

	// error check
	if err != nil {
		cmd.Wait()
	} else if err = cmd.Wait(); err != nil {
		// replace error message with stderr, if any
		if _, ok := err.(*exec.ExitError); ok {
			if s := bytes.TrimSpace(stderr.buff); len(s) > 0 {
				err = errors.New(string(s))
			}
		}
	}

	// all done
	ret = cmd.ProcessState.ExitCode()
	return
}

// limited writer: discards everything beyond the specified number of bytes
type limitedWriter struct {
	buff  []byte
	limit int
}

func (w *limitedWriter) Write(s []byte) (int, error) {
	if n := min(len(s), w.limit-len(w.buff)); n > 0 {
		w.buff = append(w.buff, s[:n]...)
	}

	return len(s), nil
}

// helpers
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
