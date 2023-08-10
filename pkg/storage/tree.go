package storage

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// SignExtSuffix is the suffix of SignPutTree's extTree for the extensions
const SignExtSuffix = "Ext"

// SignRequestSuffix is the suffix of SignPutTree's signedTree for the requests
const SignRequestSuffix = "Request"

var preSignedRequestType = reflect.TypeOf(&PreSignedHTTPRequest{})

func (d DBStorage) signPutTree(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	extTree = indirect(extTree)
	signedTree = indirect(signedTree)

	//nolint:exhaustive // default will handle the rest
	switch extTree.Kind() {
	case reflect.String:
		return d.signPutTreeString(nameToValidExts, extTree, signedTree, current)
	case reflect.Slice, reflect.Array:
		return d.signPutTreeSlice(nameToValidExts, extTree, signedTree, current)
	case reflect.Struct:
		return d.signPutTreeStruct(nameToValidExts, extTree, signedTree, current)
	default:
		return fmt.Errorf("invalid type for DBStorage.SignPutTree(): %v", extTree.Kind())
	}
}

func (d DBStorage) signPutTreeString(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	if extTree.IsZero() {
		return nil
	}
	if !signedTree.IsValid() || !signedTree.CanSet() || signedTree.Type() != preSignedRequestType {
		return fmt.Errorf("not valid settable for PreSignedHTTPRequest at %v", current)
	}
	ext := extTree.String()
	fieldName := current[strings.LastIndex(current, ".")+1:]
	if !nameToValidExts[fieldName][ext] {
		return InvalidInputError{Message: fmt.Sprintf("invalid extension, %v, at %v", ext, current)}
	}

	req, err := d.api.SignPut(current + ext)
	if err != nil {
		return err
	}
	signedTree.Set(reflect.ValueOf(&req))
	return nil
}

func (d DBStorage) signPutTreeSlice(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	if !signedTree.IsValid() || !signedTree.CanSet() || (signedTree.Kind() != reflect.Slice && signedTree.Kind() != reflect.Array) {
		return fmt.Errorf("signedTree not valid settable Slice or Array at %v", current)
	}
	if extTree.IsZero() || extTree.Len() == 0 {
		return InvalidInputError{Message: fmt.Sprintf("empty slice or array given for DBStorage.SignPutTree() at %v", current)}
	}

	signedTree.Set(reflect.MakeSlice(signedTree.Type(), extTree.Len(), extTree.Len()))
	for i := 0; i < extTree.Len(); i++ {
		if err := d.signPutTree(nameToValidExts, extTree.Index(i), signedTree.Index(i), current+"["+strconv.Itoa(i)+"]"); err != nil {
			return err
		}
	}
	return nil
}

func (d DBStorage) signPutTreeStruct(nameToValidExts SignPutNameToValidExts, extTree, signedTree reflect.Value, current string) error {
	if !signedTree.IsValid() || !signedTree.CanSet() || signedTree.Kind() != reflect.Struct {
		return fmt.Errorf("signedTree not valid settable Struct at %v", current)
	}
	if extTree.IsZero() {
		return InvalidInputError{Message: fmt.Sprintf("empty struct given for DBStorage.SignPutTree() at %v", current)}
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
		if err := d.signPutTree(nameToValidExts, extTree.Field(i), signedTreeField, current+"."+shortName); err != nil {
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

type treeValueFunc = func(key string) (string, error)

func unmarshallTree(tree map[string]any, obj reflect.Value, suffix string, valueFunc treeValueFunc) error {
	_, err := unmarshallTreeValue(tree, "", obj, suffix, valueFunc)
	return err
}

func unmarshallTreeValue(current any, currentKey string, obj reflect.Value, suffix string, valueFunc treeValueFunc) (reflect.Value, error) {
	switch currentTyped := current.(type) {
	case string:
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
		return reflect.Value{}, fmt.Errorf("at key: %v, expected Nil or Slice/Array", currentKey)
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
		return reflect.Value{}, fmt.Errorf("at key: %v, expected Struct", currentKey)
	}

	for key, value := range current {
		_, isString := value.(string)
		if isString {
			key += suffix
		}
		fieldObj := obj.FieldByName(key)
		if !fieldObj.IsValid() || !fieldObj.CanSet() {
			return reflect.Value{}, fmt.Errorf("at key: %v.%v, not valid settable field name", currentKey, key)
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
	if value.Kind() != reflect.Pointer {
		return reflect.Value{}, fmt.Errorf("obj is not a pointer")
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("obj is not a struct")
	}
	idField := value.FieldByName(IDFieldName)
	if !idField.IsValid() || !idField.CanSet() || idField.Kind() != reflect.String {
		return reflect.Value{}, fmt.Errorf("obj field, %v is not a valid settable String", IDFieldName)
	}
	idField.SetString(id)
	return value, nil
}
