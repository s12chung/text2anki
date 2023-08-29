// Package archive contains utils to help with archiving
package archive

import (
	"bytes"
	"io"
	"os/exec"
)

// XZBytes returns the contents of the .xz file at path
func XZBytes(path string) ([]byte, error) {
	cmd := exec.Command("xz", "-dc", path)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return io.ReadAll(out)
}
