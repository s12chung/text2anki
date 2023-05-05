package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Presence checks if value is Zero (or present)
type Presence struct {
}

// ValidateValue returns true if the value is present
func (p Presence) ValidateValue(value reflect.Value) firm.ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return errorMapPresence
	}
	return nil
}

var errorMapPresence = firm.ErrorMap{
	"Presence": &firm.TemplatedError{
		Template: "value is not present",
	},
}
