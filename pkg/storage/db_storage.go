package storage

import (
	"path"
	"reflect"
)

// DBStorage signs requests
type DBStorage struct {
	api           API
	uuidGenerator UUIDGenerator
}

// NewDBStorage returns a new DBStorage
func NewDBStorage(api API, uuidGenerator UUIDGenerator) DBStorage {
	if uuidGenerator == nil {
		uuidGenerator = UUID7{}
	}
	return DBStorage{api: api, uuidGenerator: uuidGenerator}
}

// SignPutNameToValidExts is a map of short field names and their valid extensions
type SignPutNameToValidExts = map[string]map[string]bool

// SignPutConfig is the config for SignPutTree
type SignPutConfig struct {
	Table           string
	Column          string
	NameToValidExts SignPutNameToValidExts
}

// BaseKey returns the base key for the given args
func BaseKey(table, column, id string) string {
	return path.Join(table, column, id, column)
}

// SignPut signs for the given ext for the able at column
func (d DBStorage) SignPut(table, column, ext string) (PreSignedHTTPRequest, error) {
	id, err := d.uuidGenerator.Generate()
	if err != nil {
		return PreSignedHTTPRequest{}, err
	}
	return d.api.SignPut(BaseKey(table, column, id) + ext)
}

// SignPutTree signs the fields in extTree and fills in the matching signedTree's PreSignedHTTPRequest
func (d DBStorage) SignPutTree(config SignPutConfig, extTree, signedTree any) error {
	id, err := d.uuidGenerator.Generate()
	if err != nil {
		return err
	}
	current := BaseKey(config.Table, config.Column, id)
	signedTreeValue, err := setID(id, signedTree)
	if err != nil {
		return err
	}
	return d.signPutTree(config.NameToValidExts, reflect.ValueOf(extTree), signedTreeValue, current)
}

// KeyTree fills in the matching keyTree's string with storage keys from the key structure
func (d DBStorage) KeyTree(table, column, id string, keyTree any) error {
	tree, objValue, err := d.preUnmarshallTree(table, column, id, keyTree)
	if err != nil {
		return err
	}
	return unmarshallTree(tree, objValue, keySuffix, func(key string) (string, error) {
		return key, nil
	})
}

// SignGetTree fills in the matching signedTree's string with signed GET URLs from the storage key structure
func (d DBStorage) SignGetTree(table, column, id string, signedTree any) error {
	tree, objValue, err := d.preUnmarshallTree(table, column, id, signedTree)
	if err != nil {
		return err
	}
	return unmarshallTree(tree, objValue, urlSuffix, d.api.SignGet)
}

// SignGetKeyTree creates a signedTree from the keyTree
func (d DBStorage) SignGetKeyTree(keyTree any, signedTree any) error {
	tree, err := treeFromKeyTree(keyTree)
	if err != nil {
		return err
	}
	return unmarshallTree(tree, reflect.ValueOf(signedTree), urlSuffix, d.api.SignGet)
}
