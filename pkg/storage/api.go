package storage

import (
	"io"
	"io/fs"
)

// API is a wrapper around the API for file storage
type API interface {
	Store(key string, file io.Reader) error
	Get(key string) (fs.File, error)
	ListKeys(prefix string) ([]string, error)

	SignPut(key string) (PreSignedHTTPRequest, error)
	SignGet(key string) (string, error)
	KeyFromSignGet(signGet string) (string, error)
}
