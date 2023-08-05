// Package storage provides interfaces to abstract out file storage services like s3
package storage

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
)

// PreSignedHTTPRequest is the data to do the signed request
type PreSignedHTTPRequest struct {
	URL          string      `json:"url"`
	Method       string      `json:"method"`
	SignedHeader http.Header `json:"signed_header"`
}

// API is a wrapper around the API for file storage
type API interface {
	SignPut(key string) (PreSignedHTTPRequest, error)
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

// SignPutBuilder returns a new SignPutBuilder
func (s Signer) SignPutBuilder(table, column string) (SignPutBuilder, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return newSignPutBuilder(table, column, id.String(), s.api), nil
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

// SignPutBuilder builds paths to sign Put requests to the storage
type SignPutBuilder interface {
	ID() string
	Index(index int) SignPutBuilder
	Field(field string) SignPutBuilder
	Sign(ext string) (PreSignedHTTPRequest, error)
}

type signPutBuilder struct {
	id     string
	prefix string
	api    API
}

func newSignPutBuilder(table, column, id string, api API) signPutBuilder {
	prefix := path.Join(table, column, id)
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	prefix += column
	return signPutBuilder{id: id, prefix: prefix, api: api}
}

// ID returns the ID for the files signed by the builder
func (s signPutBuilder) ID() string {
	return s.id
}

// Index sets the index for the new builder
func (s signPutBuilder) Index(index int) SignPutBuilder {
	s.prefix += "[" + strconv.Itoa(index) + "]"
	return s
}

// Field sets the field name for the new builder
func (s signPutBuilder) Field(field string) SignPutBuilder {
	s.prefix += "." + field
	return s
}

// Sign signs a new Put request for the builder and the extension
func (s signPutBuilder) Sign(ext string) (PreSignedHTTPRequest, error) {
	return s.api.SignPut(s.prefix + ext)
}
