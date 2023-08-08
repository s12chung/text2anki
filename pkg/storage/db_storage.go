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

func baseKey(table, column, id string) string {
	return path.Join(table, column, id, column)
}

// SignPut signs for the given ext for the able at column
func (d DBStorage) SignPut(table, column, ext string) (PreSignedHTTPRequest, error) {
	id, err := d.uuidGenerator.Generate()
	if err != nil {
		return PreSignedHTTPRequest{}, err
	}
	return d.api.SignPut(baseKey(table, column, id) + ext)
}

// SignPutTree signs the fields in extTree and fills in the matching signedTree's PreSignedHTTPRequest
func (d DBStorage) SignPutTree(config SignPutConfig, extTree, signedTree any) error {
	id, err := d.uuidGenerator.Generate()
	if err != nil {
		return err
	}
	current := baseKey(config.Table, config.Column, id)
	signedTreeValue, err := setID(id, signedTree)
	if err != nil {
		return err
	}
	return d.signPutTree(config.NameToValidExts, reflect.ValueOf(extTree), signedTreeValue, current)
}

// SignGet returns the signed GET URL for the key
func (d DBStorage) SignGet(key string) (string, error) {
	return d.api.SignGet(key)
}

// SignGetByID fills in the matching signedTree's string with signed GET URLs from the storage key structure
func (d DBStorage) SignGetByID(table, column, id string, signedTree any) error {
	objValue, err := setID(id, signedTree)
	if err != nil {
		return err
	}

	idPath := path.Join(table, column, id)
	keys, err := d.api.ListKeys(idPath)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return NotFoundError{ID: id, IDPath: idPath}
	}
	tree, err := treeFromKeys(keys)
	if err != nil {
		return err
	}
	return unmarshallTree(tree, objValue, "URL", d.SignGet)
}
