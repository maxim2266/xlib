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

// Package xlib is an ever growing collection of useful Go functions.
package xlib

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

/*
WriteFile safely replaces the content of the given file.  First, it creates a temporary file,
then it calls the supplied function to actually write to the file, and in the end it moves
the temporary to the given target file. In case of any error or a panic the temporary file
is always removed. The target pathname must either not exist, or refer to an existing regular
file, in which case it will be replaced. To avoid copying files across different filesystems
the temporary file is created in the same directory as the target.
*/
func WriteFile(pathname string, fn func(*bufio.Writer) error) (err error) {
	// check target and copy permission bits
	perm := fs.FileMode(0600)

	var info fs.FileInfo

	if info, err = os.Lstat(pathname); err == nil {
		if !info.Mode().IsRegular() {
			return errors.New(strconv.Quote(pathname) + " is not a regular file")
		}

		if perm = info.Mode().Perm(); perm&0200 == 0 {
			return errors.New(strconv.Quote(pathname) + " is not writable")
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return
	}

	// create temporary file
	var fd *os.File

	if fd, err = os.CreateTemp(filepath.Dir(pathname), "tmp-"); err != nil {
		return
	}

	temp := fd.Name()

	// make sure the temporary is always deleted
	defer func() {
		if p := recover(); p != nil {
			os.Remove(temp)
			panic(p)
		}

		if err != nil {
			os.Remove(temp)
		}
	}()

	// copy permission bits
	if err = fd.Chmod(perm); err != nil {
		return
	}

	// write and move file
	if err = writeFile(fd, fn); err == nil {
		err = os.Rename(temp, pathname) // usually, an atomic operation
	}

	return
}

func writeFile(fd *os.File, fn func(*bufio.Writer) error) (err error) {
	// make sure the file gets closed afterwards
	defer func() {
		if e := fd.Close(); e != nil && err == nil {
			err = e
		}
	}()

	// add buffer
	file := bufio.NewWriter(fd)

	// write and flush
	if err = fn(file); err == nil {
		err = file.Flush()
	}

	return
}
