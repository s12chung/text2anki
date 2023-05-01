package validates

import (
	"fmt"
	"reflect"
)

func init() {
	AddFieldValidator("presence", &PresenceFieldValidator{})
}

// PresenceFieldValidator validates the presence of the field
type PresenceFieldValidator struct {
}

var errPresenceFieldValidator = fmt.Errorf("value is empty")

// Valid returns true if the value is present
func (v *PresenceFieldValidator) Valid(value reflect.Value) error {
	if !value.IsValid() || value.IsZero() {
		return errPresenceFieldValidator
	}
	return nil
}
