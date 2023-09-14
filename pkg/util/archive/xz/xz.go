// Package xz contains utils to help with xz archiving
package xz

import (
	"bytes"
	"io"
	"os/exec"
)

// Read returns the contents of the .xz file at path
func Read(path string) ([]byte, error) {
	cmd := exec.Command("xz", "-dc", path)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return io.ReadAll(out)
}
