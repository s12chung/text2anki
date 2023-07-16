// Package storage provides interfaces to abstract out file storage services like s3
package storage

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/gofrs/uuid"
)

// PresignedHTTPRequest is the data to do the signed request
type PresignedHTTPRequest struct {
	URL          string
	Method       string
	SignedHeader http.Header
}

// API is a wrapper around the API for file storage
type API interface {
	SignPut(key string) (PresignedHTTPRequest, error)
	SignGet(key string) (string, error)
	ListKeys(prefix string) ([]string, error)
}

// Storer is a wrapper around the storage API
type Storer interface {
	Validate(key string, values url.Values) error
	Store(key string, file io.Reader) error
	FileHandler() http.Handler
}

// Signer signs requests
type Signer struct {
	api API
}

// NewSigner returns a new Signer
func NewSigner(api API) Signer {
	return Signer{api: api}
}

// SignPut signs the files for a table's field
func (s Signer) SignPut(table, column string, exts []string) ([]PresignedHTTPRequest, string, error) {
	reqs := make([]PresignedHTTPRequest, len(exts))
	id, err := uuid.NewV7()
	if err != nil {
		return nil, "", err
	}
	stringID := id.String()
	for i, ext := range exts {
		req, err := s.api.SignPut(path.Join(table, column, stringID, strconv.Itoa(i)+ext))
		if err != nil {
			return nil, "", err
		}
		reqs[i] = req
	}
	return reqs, stringID, nil
}

// SignGet returns the signed URL for the key
func (s Signer) SignGet(key string) (string, error) {
	return s.api.SignGet(key)
}

// SignGetByID returns the signed URLs for the given table, column, and ID
func (s Signer) SignGetByID(table, column, id string) ([]string, error) {
	keys, err := s.api.ListKeys(path.Join(table, column, id))
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(keys))
	for i, key := range keys {
		u, err := s.SignGet(key)
		if err != nil {
			return nil, err
		}
		urls[i] = u
	}
	return urls, nil
}
