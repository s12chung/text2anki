package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func (d DBStorage) putTreeSetup(config SignPutConfig, tree any) (reflect.Value, string, error) {
	id, err := d.uuidGenerator.Generate()
	if err != nil {
		return reflect.Value{}, "", err
	}
	treeValue, err := setID(id, tree)
	if err != nil {
		return reflect.Value{}, "", err
	}
	current := BaseKey(config.Table, config.Column, id)
	return treeValue, current, nil
}

type putTreeFunc = func(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error

var preSignedRequestType = reflect.TypeOf(&PreSignedHTTPRequest{})

func (d DBStorage) signPutTree(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	extTree = indirect(extTree)
	signedTree = indirect(signedTree)

	//nolint:exhaustive // default will handle the rest
	switch extTree.Kind() {
	case reflect.String:
		if extTree.IsZero() {
			return nil
		}
		if !signedTree.IsValid() || signedTree.Type() != preSignedRequestType {
			return fmt.Errorf("not valid PreSignedHTTPRequest at %v", current)
		}
		ext := extTree.String()
		if err := fieldExtError(nameToValidExts, ext, current); err != nil {
			return err
		}

		req, err := d.api.SignPut(current + ext)
		if err != nil {
			return err
		}
		signedTree.Set(reflect.ValueOf(&req))
		return nil
	case reflect.Slice, reflect.Array:
		return putTreeSlice(nameToValidExts, extTree, signedTree, current, d.signPutTree)
	case reflect.Struct:
		return putTreeStruct(nameToValidExts, extTree, signedTree, signExtSuffix, signRequestSuffix, current, d.signPutTree)
	default:
		return fmt.Errorf("invalid type for DBStorage.SignPutTree(): %v", extTree.Kind())
	}
}

var fileType = reflect.TypeOf((*fs.File)(nil)).Elem()

func (d DBStorage) putTree(nameToValidExts SignPutNameToValidExts, fileTree, keyTree reflect.Value, current string) error {
	fileTree = indirect(fileTree)
	keyTree = indirect(keyTree)

	if fileTree.Type().Implements(fileType) {
		if fileTree.IsZero() {
			return nil
		}
		file, ok := fileTree.Interface().(fs.File)
		if !ok {
			return fmt.Errorf("expecting fs.File for DBStorage.SignPutTree() at %v", current)
		}
		if !keyTree.IsValid() || keyTree.Kind() != reflect.String {
			return fmt.Errorf("not valid String at %v", current)
		}

		key, err := d.storeFile(file, nameToValidExts, current)
		if err != nil {
			return err
		}
		keyTree.SetString(key)
		return nil
	}

	//nolint:exhaustive // default will handle the rest
	switch fileTree.Kind() {
	case reflect.Slice, reflect.Array:
		return putTreeSlice(nameToValidExts, fileTree, keyTree, current, d.putTree)
	case reflect.Struct:
		return putTreeStruct(nameToValidExts, fileTree, keyTree, fileSuffix, keySuffix, current, d.putTree)
	default:
		return fmt.Errorf("invalid type for DBStorage.PutTree(): %v", fileTree.Kind())
	}
}

func (d DBStorage) storeFile(file fs.File, nameToValidExts SignPutNameToValidExts, current string) (string, error) {
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	ext := path.Ext(info.Name())
	if err := fieldExtError(nameToValidExts, ext, current); err != nil {
		return "", err
	}
	key := current + ext
	if err := d.api.Store(key, file); err != nil {
		return "", err
	}
	return key, file.Close()
}

func fieldExtError(nameToValidExts SignPutNameToValidExts, ext, current string) error {
	fieldName := current[strings.LastIndex(current, ".")+1:]
	if !nameToValidExts[fieldName][ext] {
		return InvalidInputError{Message: fmt.Sprintf("invalid extension, %v, at %v", ext, current)}
	}
	return nil
}

func putTreeSlice(nameToValidExts SignPutNameToValidExts, srcTree, destTree reflect.Value, current string, treeFunc putTreeFunc) error {
	if srcTree.IsZero() || srcTree.Len() == 0 {
		return InvalidInputError{Message: fmt.Sprintf("srcTree empty slice or array given at %v", current)}
	}
	if !destTree.IsValid() || (destTree.Kind() != reflect.Slice && destTree.Kind() != reflect.Array) {
		return fmt.Errorf("destTree not valid Slice or Array at %v", current)
	}

	destTree.Set(reflect.MakeSlice(destTree.Type(), srcTree.Len(), srcTree.Len()))
	for i := range srcTree.Len() {
		if err := treeFunc(nameToValidExts, srcTree.Index(i), destTree.Index(i), current+"["+strconv.Itoa(i)+"]"); err != nil {
			return err
		}
	}
	return nil
}

func putTreeStruct(nameToValidExts SignPutNameToValidExts, srcTree, destTree reflect.Value,
	srcSuffix, destSuffix, current string, treeFunc putTreeFunc) error {
	if srcTree.IsZero() {
		return InvalidInputError{Message: fmt.Sprintf("srcTree empty struct given at %v", current)}
	}
	if !destTree.IsValid() || destTree.Kind() != reflect.Struct {
		return fmt.Errorf("destTree not valid Struct at %v", current)
	}
	for i := range srcTree.NumField() {
		shortName := srcTree.Type().Field(i).Name
		requestName := shortName //nolint:copyloopvar // an actual working copy
		if strings.HasSuffix(shortName, srcSuffix) {
			shortName = shortName[:len(shortName)-len(srcSuffix)]
			requestName = shortName + destSuffix
		}
		signedTreeField := destTree.FieldByName(requestName)
		if !signedTreeField.IsValid() {
			return fmt.Errorf("destTree does not have matching field name, %v, at %v", requestName, current)
		}
		if err := treeFunc(nameToValidExts, srcTree.Field(i), signedTreeField, current+"."+shortName); err != nil {
			return err
		}
	}
	return nil
}

var indexRegex = regexp.MustCompile(`\[\d+]$`)

func treeFromKeys(keys []string) (map[string]any, error) {
	sort.Strings(keys)
	tree := map[string]any{}

	for _, key := range keys {
		treeParts, err := treePartsFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("at key: %v, %w", key, err)
		}
		_, err = setTree(tree, treeParts, key)
		if err != nil {
			return nil, fmt.Errorf("at key: %v, %w", key, err)
		}
	}
	return tree, nil
}

type treePartType int

const (
	treePartStruct treePartType = iota
	treePartSlice
)

type treePart struct {
	typ   treePartType
	key   string
	index int
}

func treePartsFromKey(key string) ([]treePart, error) {
	key = key[0 : len(key)-len(filepath.Ext(key))]
	parts := strings.Split(key, ".")[1:]

	//nolint:prealloc // can't calculate upfront
	var treeParts []treePart
	for _, part := range parts {
		var arrayParts []string
		part, arrayParts = splitArrayParts(part)
		treeParts = append(treeParts, treePart{typ: treePartStruct, key: part})
		for _, arrayPart := range arrayParts {
			index, err := strconv.Atoi(arrayPart[1 : len(arrayPart)-1])
			if err != nil {
				return nil, err
			}
			treeParts = append(treeParts, treePart{typ: treePartSlice, index: index})
		}
	}
	return treeParts, nil
}

func setTree(treeObj any, treeParts []treePart, finalValue string) (any, error) {
	if len(treeParts) == 0 {
		return finalValue, nil
	}

	tPart := treeParts[0]
	treeParts = treeParts[1:]

	switch tPart.typ {
	case treePartSlice:
		return setTreeSlice(tPart, treeObj, treeParts, finalValue)
	case treePartStruct:
		return setTreeStruct(tPart, treeObj, treeParts, finalValue)
	default:
		return nil, fmt.Errorf("got invalid treePart.typ: %v", tPart.typ)
	}
}

func setTreeSlice(part treePart, treeObj any, treeParts []treePart, finalValue string) (any, error) {
	if treeObj == nil {
		treeObj = []any{}
	}
	currentSlice, ok := treeObj.([]any)
	if !ok {
		return nil, fmt.Errorf("expected Slice at: %v", part.index)
	}
	if len(currentSlice) < part.index {
		return nil, fmt.Errorf("slice index (%v) lower than: %v", len(currentSlice), part.index)
	}
	if len(currentSlice) == part.index {
		currentSlice = append(currentSlice, nil)
		treeObj = currentSlice
	}

	value, err := setTree(currentSlice[part.index], treeParts, finalValue)
	if err != nil {
		return nil, err
	}
	if currentSlice[part.index] != nil {
		valueType := reflect.TypeOf(value)
		existingType := reflect.TypeOf(currentSlice[part.index])
		if valueType != existingType {
			return nil, fmt.Errorf("unmatched types %v and %v at: %v", valueType.String(), existingType.String(), part.key)
		}
	}
	currentSlice[part.index] = value
	return treeObj, nil
}

func setTreeStruct(part treePart, treeObj any, treeParts []treePart, finalValue string) (any, error) {
	if treeObj == nil {
		treeObj = map[string]any{}
	}
	currentMap, ok := treeObj.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected Map at: %v", part.key)
	}
	value, err := setTree(currentMap[part.key], treeParts, finalValue)
	if err != nil {
		return nil, err
	}
	if currentMap[part.key] != nil {
		valueType := reflect.TypeOf(value)
		existingType := reflect.TypeOf(currentMap[part.key])
		if valueType != existingType {
			return nil, fmt.Errorf("unmatched types %v and %v at: %v", valueType.String(), existingType.String(), part.key)
		}
	}
	currentMap[part.key] = value
	return treeObj, nil
}

func splitArrayParts(part string) (string, []string) {
	var aParts []string

	matchIndices := indexRegex.FindStringIndex(part)
	for matchIndices != nil {
		aParts = append([]string{part[matchIndices[0]:]}, aParts...)
		part = part[:matchIndices[0]]
		matchIndices = indexRegex.FindStringIndex(part)
	}
	return part, aParts
}

func (d DBStorage) preUnmarshallTree(table, column, id string, obj any) (map[string]any, reflect.Value, error) {
	objValue, err := setID(id, obj)
	if err != nil {
		return nil, reflect.Value{}, err
	}

	idPath := path.Join(table, column, id)
	keys, err := d.api.ListKeys(idPath)
	if err != nil {
		return nil, reflect.Value{}, err
	}
	if len(keys) == 0 {
		return nil, reflect.Value{}, NotFoundError{ID: id, IDPath: idPath}
	}
	tree, err := treeFromKeys(keys)
	if err != nil {
		return nil, reflect.Value{}, err
	}
	return tree, objValue, nil
}

func mapTree(keyTree any, fieldSuffix string) (map[string]any, error) {
	keyTreeValue := reflect.ValueOf(keyTree)
	if keyTreeValue.IsZero() {
		return map[string]any{}, nil
	}
	current := indirect(keyTreeValue)
	if current.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%v is not a Struct", keyTreeValue.Type().String())
	}
	tree, err := mapTreeFromValue(current, fieldSuffix)
	if err != nil {
		return nil, err
	}
	treeMap, ok := tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%v is not a Struct", keyTreeValue.Type().String())
	}
	return treeMap, nil
}

func mapTreeFromValue(current reflect.Value, fieldSuffix string) (any, error) {
	current = indirect(current)

	//nolint:exhaustive // default will handle the rest
	switch current.Kind() {
	case reflect.String:
		return current.String(), nil
	case reflect.Slice, reflect.Array:
		treeObjSlice := make([]any, current.Len())
		for i := range current.Len() {
			value, err := mapTreeFromValue(current.Index(i), fieldSuffix)
			if err != nil {
				return nil, err
			}
			treeObjSlice[i] = value
		}
		return treeObjSlice, nil
	case reflect.Struct:
		treeObjMap := map[string]any{}
		for i := range current.NumField() {
			fieldType := current.Type().Field(i)
			if !fieldType.IsExported() {
				continue
			}
			value, err := mapTreeFromValue(current.Field(i), fieldSuffix)
			if err != nil {
				return nil, err
			}
			keyName := fieldType.Name
			if _, isString := value.(string); isString {
				keyName = keyName[:len(keyName)-len(fieldSuffix)]
			}
			treeObjMap[keyName] = value
		}
		return treeObjMap, nil
	default:
		return nil, fmt.Errorf("invalid type for mapTree(): %v", current.Kind())
	}
}

type treeValueFunc = func(key string) (string, error)

func unmarshallTree(tree map[string]any, obj reflect.Value, suffix string, valueFunc treeValueFunc) error {
	_, err := unmarshallTreeValue(tree, "", obj, suffix, valueFunc)
	return err
}

func unmarshallTreeValue(current any, currentKey string, obj reflect.Value, suffix string, valueFunc treeValueFunc) (reflect.Value, error) {
	switch currentTyped := current.(type) {
	case string:
		if currentTyped == "" {
			return reflect.ValueOf(currentTyped), nil
		}
		v, err := valueFunc(currentTyped)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v), nil
	case []any:
		return unmarshallTreeSlice(currentTyped, currentKey, obj, suffix, valueFunc)
	case map[string]any:
		return unmarshallTreeStruct(currentTyped, currentKey, obj, suffix, valueFunc)
	default:
		return reflect.Value{}, fmt.Errorf("at key: %v, unexcepted tree type: %T", currentKey, currentTyped)
	}
}

func unmarshallTreeSlice(current []any, currentKey string, obj reflect.Value, suffix string, valueFunc treeValueFunc) (reflect.Value, error) {
	if obj.IsNil() {
		obj = reflect.MakeSlice(obj.Type(), len(current), len(current))
	}
	if obj.Kind() != reflect.Slice && obj.Kind() != reflect.Array {
		return reflect.Value{}, fmt.Errorf("at key: %v, expected Nil or Slice/Array, but got %v", currentKey, obj.Type().String())
	}

	for i, value := range current {
		indexObj := obj.Index(i)
		v, err := unmarshallTreeValue(value, fmt.Sprintf("%v[%v]", currentKey, i), indexObj, suffix, valueFunc)
		if err != nil {
			return reflect.Value{}, err
		}
		indexObj.Set(v)
	}
	return obj, nil
}

func unmarshallTreeStruct(current map[string]any, currentKey string, obj reflect.Value, suffix string,
	valueFunc treeValueFunc) (reflect.Value, error) {
	obj = indirect(obj)
	if obj.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("at key: %v, expected Struct, but got %v", currentKey, obj.Type().String())
	}

	for key, value := range current {
		_, isString := value.(string)
		if isString {
			key += suffix
		}
		fieldObj := obj.FieldByName(key)
		if !fieldObj.IsValid() {
			return reflect.Value{}, fmt.Errorf("at key: %v.%v, not valid field name", currentKey, key)
		}

		v, err := unmarshallTreeValue(value, fmt.Sprintf("%v.%v", currentKey, key), fieldObj, suffix, valueFunc)
		if err != nil {
			return reflect.Value{}, err
		}
		fieldObj.Set(v)
	}
	return obj, nil
}

func indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Pointer && !value.IsZero() {
		value = value.Elem()
	}
	return value
}

// IDFieldName is the field name for SignPutTree()'s signTree's ID
const IDFieldName = "ID"

func setID(id string, obj any) (reflect.Value, error) {
	value := reflect.ValueOf(obj)
	if !value.IsValid() {
		return reflect.Value{}, errors.New("passed nil as settable obj")
	}
	if value.Kind() != reflect.Pointer {
		return reflect.Value{}, fmt.Errorf("%v is not a pointer", value.Type().String())
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%v is not a struct", value.Type().String())
	}
	idField := value.FieldByName(IDFieldName)
	if !idField.IsValid() || idField.Kind() != reflect.String {
		return reflect.Value{}, fmt.Errorf("%v field, %v is not a valid String", value.Type().String(), IDFieldName)
	}
	idField.SetString(id)
	return value, nil
}
