// Package khaiii runs kakao/khaiii, a Korean tokenizer
package khaiii

import (
	"github.com/ebitengine/purego"
)

// DefaultDlPath is the default Dl path relative to the binary
const DefaultDlPath = "lib/libkhaiii.0.5.dylib"

// DefaultRscPath is the default DefaultRscPath relative to the binary
const DefaultRscPath = "lib/rsc"

// Khaiii is a wrapper around the Khaiii C API
type Khaiii struct {
	dlPath string
	libPtr uintptr

	openHandle int
}

// NewKhaiii returns a new Khaiii
func NewKhaiii(dlPath string) (*Khaiii, error) {
	var err error
	libPtr, err := purego.Dlopen(dlPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}

	for functionName, f := range apiFunctions {
		purego.RegisterLibFunc(f, libPtr, functionName)
	}
	return &Khaiii{
		dlPath: dlPath,
		libPtr: libPtr,
	}, nil
}
