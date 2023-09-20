package rule

import (
	"reflect"
	"strings"

	"github.com/s12chung/text2anki/pkg/firm"
)

// TrimPresent checks if value is not "" when strings.TrimSpace is applied
type TrimPresent struct{}

// ValidateValue returns true if the value is present
func (t TrimPresent) ValidateValue(value reflect.Value) firm.ErrorMap {
	if value.Kind() != reflect.String {
		return errorMapNotSpace("value is not a string")
	}
	if strings.TrimSpace(value.String()) == "" {
		return errorMapNotSpace("value is just spaces or empty")
	}
	return nil
}

func errorMapNotSpace(template string) firm.ErrorMap {
	return firm.ErrorMap{"TrimPresent": &firm.TemplatedError{Template: template}}
}
