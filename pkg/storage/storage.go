// Package storage provides interfaces to abstract out file storage services like s3
package storage

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gofrs/uuid"
)

// PreSignedHTTPRequest is the data to do the signed request
type PreSignedHTTPRequest struct {
	URL          string      `json:"url"`
	Method       string      `json:"method"`
	SignedHeader http.Header `json:"signed_header"`
}

// Storer is a wrapper around the storage API
type Storer interface {
	Validate(key string, values url.Values) error
	Store(key string, file io.Reader) error
	FileHandler() http.Handler
}

// UUIDGenerator generates UUID for signer
type UUIDGenerator interface {
	Generate() (string, error)
}

// UUID7 generates UUID v7 uuids
type UUID7 struct {
}

// Generate generates a UUId
func (u UUID7) Generate() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
