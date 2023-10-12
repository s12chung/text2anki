package firm

import (
	"maps"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

// ErrorMap is a map of TemplateError keys to their respective TemplateError
//
//nolint:errname
type ErrorMap map[ErrorKey]TemplateError

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
func (e ErrorMap) MergeInto(path string, dest ErrorMap) {
	for k, v := range e {
		dest[joinKeys(ErrorKey(path), k)] = v
	}
}

// ToNil returns itself or nil if it's empty
func (e ErrorMap) ToNil() ErrorMap {
	if len(e) == 0 {
		return nil
	}
	return e
}

// Finish finishes the ErrorMap for consumption by filling in the typeName and ValueName
func (e ErrorMap) Finish() ErrorMap {
	for k, v := range e {
		v.ErrorKey = k
		e[k] = v
	}
	return e.ToNil()
}

// TemplateError is an error that contains a key matching a field or top level, a golang template, and template fields
type TemplateError struct {
	Template       string
	TemplateFields map[string]string
	ErrorKey       ErrorKey
}

// Error returns a string for the error
func (t TemplateError) Error() string {
	badTemplateString := t.Template + " (bad format)"
	temp, err := template.New("top").Parse("{{.ValueName}} " + t.Template)
	if err != nil {
		return badTemplateString
	}

	templateDot := map[string]string{}
	if t.TemplateFields != nil {
		templateDot = maps.Clone(t.TemplateFields)
	}
	typeName := t.ErrorKey.RootTypeName()
	if typeName == "" {
		typeName = "NoType"
	}
	templateDot["RootTypeName"] = typeName
	valueName := t.ErrorKey.ValueName()
	if valueName == "" {
		valueName = "value"
	}
	templateDot["ValueName"] = valueName

	var sb strings.Builder
	if err = temp.Execute(&sb, templateDot); err != nil {
		return badTemplateString
	}
	return sb.String()
}

// ErrorKey is a string that has helper functions relating to error keys
type ErrorKey string

// RootTypeName returns the type name of the key
func (e ErrorKey) RootTypeName() string {
	suffix := string(e)
	for i := 0; i < 2; i++ {
		index := strings.Index(suffix, keySeparator)
		if index == -1 {
			return ""
		}
		suffix = suffix[index+1:]
	}
	name := string(e)
	return name[:len(name)-len(suffix)-1]
}

// ValueName returns the value name of the key - the Struct field, array index or value type name
func (e ErrorKey) ValueName() string {
	s := string(e)
	lastIdx := strings.LastIndex(s, keySeparator)
	if lastIdx == -1 {
		return ""
	}
	secLastIdx := strings.LastIndex(s[:lastIdx-1], keySeparator)
	if secLastIdx == -1 {
		return ""
	}
	firstIdx := strings.Index(s, keySeparator)

	start := secLastIdx + 1
	if firstIdx == secLastIdx {
		start = 0
	}
	return s[start:lastIdx]
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
func NewRuleTypeError(ruleName string, typ reflect.Type, badCondition string) *RuleTypeError {
	return &RuleTypeError{RuleName: ruleName, Type: typ, BadCondition: badCondition}
}

// RuleTypeError is an error returned by Rule.TypeCheck
type RuleTypeError struct {
	RuleName     string
	Type         reflect.Type
	BadCondition string
}

// TemplateError returns the TemplateError represented by the RuleTypeError
func (r RuleTypeError) TemplateError() TemplateError {
	valueTypeName := "nil"
	if r.Type != nil {
		valueTypeName = r.Type.String()
	}
	return TemplateError{
		TemplateFields: map[string]string{"ValueTypeName": valueTypeName},
		Template:       r.BadCondition + ", got {{.ValueTypeName}}",
	}
}

// Error returns the error string for the error
func (r RuleTypeError) Error() string { return r.RuleName + ": " + r.TemplateError().Error() }

// TypeCheck is a basic implementation for TypeCheck
func TypeCheck(ruleName string, typ, expectedType reflect.Type, kindString string) *RuleTypeError {
	if typ != expectedType {
		if kindString != "" {
			kindString += " "
		}
		return NewRuleTypeError(ruleName, typ, "is not matching "+kindString+"of type "+expectedType.String())
	}
	return nil
}
