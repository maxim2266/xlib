// Package xlib is an ever growing collection of useful Go functions.
package xlib

import (
	"bufio"
	"errors"
	"io"
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

	// cleanup
	defer func() {
		temp := fd.Name()

		// close and delete on panic
		if p := recover(); p != nil {
			fd.Close()
			os.Remove(temp)
			panic(p)
		}

		// close file
		if e := fd.Close(); e != nil && err == nil {
			err = e
		}

		// move file
		if err == nil {
			err = os.Rename(temp, pathname) // an atomic operation on most filesystems
		}

		// delete file on error
		if err != nil {
			os.Remove(temp)
		}
	}()

	// copy permission bits
	if err = fd.Chmod(perm); err != nil {
		return
	}

	// add buffer
	file := bufio.NewWriter(fd)

	// write and flush
	if err = fn(file); err == nil {
		err = file.Flush()
	}

	// wrap error
	if err == io.ErrShortWrite {
		err = errors.New("writing " + strconv.Quote(fd.Name()) + ": " + err.Error())
	}

	return
}
