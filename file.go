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
