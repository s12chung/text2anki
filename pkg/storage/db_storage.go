package storage

import (
	"fmt"
	"io/fs"
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

// Get returns the file at key
func (d DBStorage) Get(key string) (fs.File, error) {
	return d.api.Get(key)
}

// SignPutTree signs the fields in extTree and fills in the matching signedTree's PreSignedHTTPRequest
func (d DBStorage) SignPutTree(config SignPutConfig, extTree, signedTree any) error {
	signedTreeValue, current, err := d.putTreeSetup(config, signedTree)
	if err != nil {
		return err
	}
	return d.signPutTree(config.NameToValidExts, reflect.ValueOf(extTree), signedTreeValue, current)
}

// PutTree puts the files from fileTree and sets the keyTree
func (d DBStorage) PutTree(config SignPutConfig, fileTree, keyTree any) error {
	keyTreeValue, current, err := d.putTreeSetup(config, keyTree)
	if err != nil {
		return err
	}
	return d.putTree(config.NameToValidExts, reflect.ValueOf(fileTree), keyTreeValue, current)
}

const signExtSuffix = "Ext"
const signRequestSuffix = "Request"
const keySuffix = "Key"
const urlSuffix = "URL"
const fileSuffix = "File"

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

// SignGetTree fills in the matching signedTree's strings with signed GET URLs from the storage key structure
func (d DBStorage) SignGetTree(table, column, id string, signedTree any) error {
	tree, objValue, err := d.preUnmarshallTree(table, column, id, signedTree)
	if err != nil {
		return err
	}
	return unmarshallTree(tree, objValue, urlSuffix, d.api.SignGet)
}

// SignGetTreeFromKeyTree creates a signedTree from the keyTree
func (d DBStorage) SignGetTreeFromKeyTree(keyTree, signedTree any) error {
	tree, err := mapTree(keyTree, keySuffix)
	if err != nil {
		return err
	}
	signedTreeValue := reflect.ValueOf(signedTree)
	if signedTreeValue.Kind() != reflect.Pointer {
		return fmt.Errorf("signedTree, %v, is not a pointer", signedTreeValue.Type().String())
	}
	return unmarshallTree(tree, signedTreeValue, urlSuffix, d.api.SignGet)
}

// KeyTreeFromSignGetTree creates a keyTree from a signedTree
func (d DBStorage) KeyTreeFromSignGetTree(signedTree, keyTree any) error {
	tree, err := mapTree(signedTree, urlSuffix)
	if err != nil {
		return err
	}
	keyTreeValue := reflect.ValueOf(keyTree)
	if keyTreeValue.Kind() != reflect.Pointer {
		return fmt.Errorf("keyTree, %v, is not a pointer", keyTreeValue.Type().String())
	}
	return unmarshallTree(tree, keyTreeValue, keySuffix, d.api.KeyFromSignGet)
}
