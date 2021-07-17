// Package iotools contains iotools
package iotools

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CopyFile copies the contents from src to dst atomically.
// If dst does not exist, CopyFile creates it with permissions perm.
// If the copy fails, CopyFile aborts and dst is preserved.
//
// https://go.googlesource.com/go/+/c766fc4dc357a1fed47ba31212fb2a12d0e050d6%5E%21
func CopyFile(dst, src string, perm os.FileMode) error {
	//nolint:gosec // generic library
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	//nolint:errcheck,gosec // just closing file
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		//nolint:errcheck,gosec // just closing file
		tmp.Close()
		//nolint:errcheck,gosec // just removing temp file
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		//nolint:errcheck,gosec // just removing temp file
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		//nolint:errcheck,gosec // just removing temp file
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dst)
}
