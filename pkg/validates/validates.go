// Package validates contains functions to validate structs
package validates

import (
	"fmt"
	"reflect"
	"strings"
)

// New returns a new Validator
func New(data any) *Validator {
	return &Validator{data: data}
}

// Validator validates the data
type Validator struct {
	data any

	// hidden untested API
	Key string
}

// IsValid returns true if the data is valid
func (v *Validator) IsValid() bool {
	var err error
	v.Key, err = v.validate(reflect.Indirect(reflect.ValueOf(v.data)))
	return err == nil
}

func (v *Validator) validate(dataValue reflect.Value) (string, error) {
	if dataValue.Kind() != reflect.Struct {
		return "", fmt.Errorf("passed in data is not a struct")
	}
	for i := 0; i < dataValue.NumField(); i++ {
		field := dataValue.Type().Field(i)
		tag, exists := field.Tag.Lookup("validates")
		if !exists {
			return "", nil
		}

		fieldValue := dataValue.Field(i)
		structValue := reflect.Indirect(fieldValue)
		if structValue.Kind() == reflect.Struct {
			key, err := v.validate(structValue)
			if err != nil {
				return field.Name + "." + key, err
			}
		}

		if tag == "" {
			continue
		}
		for _, fieldValidation := range strings.Split(tag, ",") {
			fieldValidator, exists := fieldValidatorMap[fieldValidation]
			if !exists {
				return field.Name, fmt.Errorf("fieldValidator does not exist: %v", fieldValidation)
			}
			err := fieldValidator.Valid(fieldValue)
			if err != nil {
				return field.Name, err
			}
		}
	}
	return "", nil
}
