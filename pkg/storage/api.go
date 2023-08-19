package storage

import "io"

// API is a wrapper around the API for file storage
type API interface {
	Store(key string, file io.Reader) error
	SignPut(key string) (PreSignedHTTPRequest, error)
	SignGet(key string) (string, error)
	ListKeys(prefix string) ([]string, error)
	KeyFromSignGet(signGet string) (string, error)
}
