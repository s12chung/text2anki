package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Presence checks if value is non-Zero and non-Empty
type Presence struct {
}

// ValidateValue returns true if the value is present
func (p Presence) ValidateValue(value reflect.Value) firm.ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return errorMapPresence
	}
	//nolint:exhaustive // only checking against .Len() kinds
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map, reflect.String:
		if value.Len() == 0 {
			return errorMapPresence
		}
	case reflect.Ptr:
		elem := value.Type().Elem()
		if elem.Kind() == reflect.Array {
			if elem.Len() == 0 {
				return errorMapPresence
			}
		}
	}
	return nil
}

var errorMapPresence = firm.ErrorMap{
	"Presence": &firm.TemplatedError{
		Template: "value is not present",
	},
}
