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
	Sign(key string) (PresignedHTTPRequest, error)
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

// Sign signs the files for a table's field
func (s Signer) Sign(table, field string, exts []string) ([]PresignedHTTPRequest, error) {
	reqs := make([]PresignedHTTPRequest, len(exts))
	for i, ext := range exts {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		req, err := s.api.Sign(path.Join(table, field, id.String(), strconv.Itoa(i)+ext))
		if err != nil {
			return nil, err
		}
		reqs[i] = req
	}
	return reqs, nil
}
