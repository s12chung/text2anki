// Package storage provides interfaces to abstract out file storage services like s3
package storage

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
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

// Signer signs requests
type Signer struct {
	api           API
	uuidGenerator UUIDGenerator
}

// NewSigner returns a new Signer
func NewSigner(api API, uuidGenerator UUIDGenerator) Signer {
	if uuidGenerator == nil {
		uuidGenerator = UUID7{}
	}
	return Signer{api: api, uuidGenerator: uuidGenerator}
}

// SignIDFieldName is the field name for SignPutTree()'s signTree's ID
const SignIDFieldName = "ID"

// SignPutNameToValidExts is a map of short field names and their valid extensions
type SignPutNameToValidExts = map[string]map[string]bool

// SignPutConfig is the config for SignPutTree
type SignPutConfig struct {
	Table           string
	Column          string
	NameToValidExts SignPutNameToValidExts
}

func baseKey(table, column, id string) string {
	return path.Join(table, column, id, column)
}

// InvalidInputError is the error returned that is an input (extTree) related error
type InvalidInputError struct {
	Message string
}

func (e InvalidInputError) Error() string {
	return e.Message
}

// IsInvalidInputError returns true if the error is a InvalidInputError
func IsInvalidInputError(err error) bool {
	var errExtTree InvalidInputError
	ok := errors.As(err, &errExtTree)
	return ok
}

// SignPut signs for the given ext for the able at column
func (s Signer) SignPut(table, column, ext string) (PreSignedHTTPRequest, error) {
	id, err := s.uuidGenerator.Generate()
	if err != nil {
		return PreSignedHTTPRequest{}, err
	}
	return s.api.SignPut(baseKey(table, column, id) + ext)
}

// SignPutTree signs the fields in extTree and fills in the matching signedTree's PreSignedHTTPRequest
func (s Signer) SignPutTree(config SignPutConfig, extTree, signedTree any) error {
	id, err := s.uuidGenerator.Generate()
	if err != nil {
		return err
	}
	current := baseKey(config.Table, config.Column, id)

	signedTreeValue := reflect.ValueOf(signedTree)
	if signedTreeValue.Kind() != reflect.Pointer {
		return fmt.Errorf("signedTree is not a pointer")
	}
	signedTreeValue = signedTreeValue.Elem()
	if signedTreeValue.Kind() != reflect.Struct {
		return fmt.Errorf("signedTree is not a struct")
	}
	idField := signedTreeValue.FieldByName(SignIDFieldName)
	if !idField.IsValid() {
		return fmt.Errorf("signedTree does not have matching field name, %v, at %v", SignIDFieldName, current)
	}
	if !idField.CanSet() || idField.Kind() != reflect.String {
		return fmt.Errorf("signedTree field, %v, at %v is not a settable String", SignIDFieldName, current)
	}
	idField.SetString(id)

	return s.signPutTree(config.NameToValidExts, reflect.ValueOf(extTree), signedTreeValue, current)
}

// SignExtSuffix is the suffix of SignPutTree's extTree for the extensions
const SignExtSuffix = "Ext"

// SignRequestSuffix is the suffix of SignPutTree's signedTree for the requests
const SignRequestSuffix = "Request"

var preSignedRequestType = reflect.TypeOf(&PreSignedHTTPRequest{})

func (s Signer) signPutTree(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	extTree = indirect(extTree)
	signedTree = indirect(signedTree)

	//nolint:exhaustive // default will handle the rest
	switch extTree.Kind() {
	case reflect.String:
		return s.signPutTreeString(nameToValidExts, extTree, signedTree, current)
	case reflect.Slice, reflect.Array:
		return s.signPutTreeSlice(nameToValidExts, extTree, signedTree, current)
	case reflect.Struct:
		return s.signPutTreeStruct(nameToValidExts, extTree, signedTree, current)
	default:
		return fmt.Errorf("invalid type for Signer.SignPutTree(): %v", extTree.Kind())
	}
}

func (s Signer) signPutTreeString(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	if extTree.IsZero() {
		return nil
	}
	if !signedTree.CanSet() || signedTree.Type() != preSignedRequestType {
		return fmt.Errorf("signedTree not settable PreSignedHTTPRequest at %v", current)
	}
	ext := extTree.String()
	fieldName := current[strings.LastIndex(current, ".")+1:]
	if !nameToValidExts[fieldName][ext] {
		return InvalidInputError{Message: fmt.Sprintf("invalid extension, %v, at %v", ext, current)}
	}

	req, err := s.api.SignPut(current + ext)
	if err != nil {
		return err
	}
	signedTree.Set(reflect.ValueOf(&req))
	return nil
}

func (s Signer) signPutTreeSlice(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	if !signedTree.CanSet() || (signedTree.Kind() != reflect.Slice && signedTree.Kind() != reflect.Array) {
		return fmt.Errorf("signedTree not settable Slice or Array at %v", current)
	}
	if extTree.IsZero() || extTree.Len() == 0 {
		return InvalidInputError{Message: fmt.Sprintf("empty slice or array given for Signer.SignPutTree() at %v", current)}
	}

	signedTree.Set(reflect.MakeSlice(signedTree.Type(), extTree.Len(), extTree.Len()))
	for i := 0; i < extTree.Len(); i++ {
		if err := s.signPutTree(nameToValidExts, extTree.Index(i), signedTree.Index(i), current+"["+strconv.Itoa(i)+"]"); err != nil {
			return err
		}
	}
	return nil
}

func (s Signer) signPutTreeStruct(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	if signedTree.Kind() != reflect.Struct {
		return fmt.Errorf("signedTree not Struct at %v", current)
	}
	if extTree.IsZero() {
		return InvalidInputError{Message: fmt.Sprintf("empty struct given for Signer.SignPutTree() at %v", current)}
	}
	for i := 0; i < extTree.NumField(); i++ {
		shortName := extTree.Type().Field(i).Name
		requestName := shortName
		if strings.HasSuffix(shortName, SignExtSuffix) {
			shortName = shortName[:len(shortName)-len(SignExtSuffix)]
			requestName = shortName + SignRequestSuffix
		}
		signedTreeField := signedTree.FieldByName(requestName)
		if !signedTreeField.IsValid() {
			return fmt.Errorf("signedTree does not have matching field name, %v, at %v", requestName, current)
		}
		if err := s.signPutTree(nameToValidExts, extTree.Field(i), signedTreeField, current+"."+shortName); err != nil {
			return err
		}
	}
	return nil
}

func indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Pointer && !value.IsZero() {
		value = value.Elem()
	}
	return value
}

// SignGet returns the signed GET URL for the key
func (s Signer) SignGet(key string) (string, error) {
	return s.api.SignGet(key)
}

// SignGetByID returns the signed GET URLs for the given table, column, and ID
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
