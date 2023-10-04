package firm

import (
	"reflect"
	"sort"
	"strings"
	"text/template"
)

// ErrorMap is a map of TemplatedError keys to their respective TemplatedError
//
//nolint:errname
type ErrorMap map[ErrorKey]*TemplatedError

func (e ErrorMap) Error() string {
	errors := make([]string, len(e))
	keys := e.sortedKeys()
	for i, k := range keys {
		errors[i] = k + ": " + e[ErrorKey(k)].Error()
	}
	return strings.Join(errors, ", ")
}

func (e ErrorMap) sortedKeys() []string {
	keys := make([]string, len(e))
	i := 0
	for k := range e {
		keys[i] = string(k)
		i++
	}
	sort.Strings(keys)
	return keys
}

// MergeInto merges into dest, given appending path to the src keys
func (e ErrorMap) MergeInto(path ErrorKey, dest ErrorMap) {
	for k, v := range e {
		dest[joinKeys(path, k)] = v
	}
}

// ToNil returns itself or nil if it's empty
func (e ErrorMap) ToNil() ErrorMap {
	if len(e) == 0 {
		return nil
	}
	return e
}

// TemplatedError is an error that contains a key matching a field or top level, a golang template, and template fields
type TemplatedError struct {
	Template       string
	TemplateFields map[string]string
}

// Error returns a string for the error
func (t *TemplatedError) Error() string {
	badTemplateString := t.Template + " (bad format)"
	temp, err := template.New("top").Parse(t.Template)
	if err != nil {
		return badTemplateString
	}
	var sb strings.Builder
	if err = temp.Execute(&sb, t.TemplateFields); err != nil {
		return badTemplateString
	}
	return sb.String()
}

// ErrorKey is a string that has helper functions relating to error keys
type ErrorKey string

// TypeName returns the type name of the key
func (e ErrorKey) TypeName() string {
	s := string(e)
	firstIdx := strings.Index(s, keySeparator)
	if firstIdx == -1 {
		return ""
	}
	return s[:firstIdx]
}

// ErrorName returns the error name of the key
func (e ErrorKey) ErrorName() string {
	s := string(e)
	lastIdx := strings.LastIndex(s, keySeparator)
	if lastIdx == -1 {
		return ""
	}
	return s[lastIdx+len(keySeparator):]
}

// NewRuleTypeError returns a new RuleTypeError
func NewRuleTypeError(typ reflect.Type, badCondition string) *RuleTypeError {
	return &RuleTypeError{Type: typ, BadCondition: badCondition}
}

// RuleTypeError is an error returned by Rule.ValidateType
type RuleTypeError struct {
	Type         reflect.Type
	BadCondition string
}

// TemplatedError returns the TemplatedError represented by the RuleTypeError
func (r RuleTypeError) TemplatedError() *TemplatedError {
	typeString := "nil"
	if r.Type != nil {
		typeString = r.Type.String()
	}
	return &TemplatedError{
		TemplateFields: map[string]string{"Type": typeString},
		Template:       "value to validate " + r.BadCondition + ", got {{.Type}}",
	}
}

// Error returns the error string for the error
func (r RuleTypeError) Error() string { return r.TemplatedError().Error() }
