// Package ioutil contains utils for io
package ioutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// OwnerGroupR is the file mode for Owners and Group read
const OwnerGroupR os.FileMode = 0440

// OwnerRWGroupR is the file mode for Owners read/write and Group read
const OwnerRWGroupR os.FileMode = 0640

// OwnerRWXGroupRX is the mode for Owners read/write/execute + Group read/write, commonly used for directories
const OwnerRWXGroupRX os.FileMode = 0750

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
	tmp, err := os.CreateTemp(filepath.Dir(dst), "")
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

// FilenamesWithExtensions returns the file names with the given extensions in the directory (non-recursive)
func FilenamesWithExtensions(entries []os.DirEntry, extensions []string) []string {
	filenames := make([]string, 0, len(entries))
	for _, file := range entries {
		if file.IsDir() {
			continue
		}
		for _, ext := range extensions {
			if strings.HasSuffix(file.Name(), ext) {
				filenames = append(filenames, file.Name())
				break
			}
		}
	}
	return filenames
}
