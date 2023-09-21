package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Present checks if value is non-Zero and non-Empty
type Present struct{}

// ValidateValue returns true if the value is present
func (p Present) ValidateValue(value reflect.Value) firm.ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return errorMapPresent
	}
	//nolint:exhaustive // only checking against .Len() kinds
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map, reflect.String:
		if value.Len() == 0 {
			return errorMapPresent
		}
	case reflect.Ptr:
		elem := value.Type().Elem()
		if elem.Kind() == reflect.Array {
			if elem.Len() == 0 {
				return errorMapPresent
			}
		}
	}
	return nil
}

var errorMapPresent = firm.ErrorMap{"Present": &firm.TemplatedError{Template: "value is not present"}}
