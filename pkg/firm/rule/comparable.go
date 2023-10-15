package rule

import (
	"cmp"
	"fmt"
	"reflect"
	"unicode"

	"github.com/s12chung/text2anki/pkg/firm"
)

const equalName = "Equal"
const lessName = "Less"
const greaterName = "Greater"

// Equal checks if data is equal to .To
type Equal[T comparable] struct{ To T }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (e Equal[T]) ValidateValue(value reflect.Value) firm.ErrorMap {
	return comparableValidateValue[T](e, value)
}

// Validate validates the data value
func (e Equal[T]) Validate(data T) firm.ErrorMap { return comparableValidate[T](e, data) }

// Compare returns true if the data is valid
func (e Equal[T]) Compare(data T) bool { return data == e.To }

// ErrorMap returns the ErrorMap returned from ValidateValue
func (e Equal[T]) ErrorMap() firm.ErrorMap {
	return firm.ErrorMap{equalName: firm.TemplateError{
		TemplateFields: map[string]string{"To": fmt.Sprintf("%v", e.To)},
		Template:       "is not equal to {{.To}}",
	}}
}

// TypeCheck checks whether the type is valid for the Rule
func (e Equal[T]) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	return comparableTypeCheck(equalName, e.To, typ)
}

// Name returns the name of the Rule
func (e Equal[T]) Name() string { return equalName }

// Less checks if data is less (or equal to) .To
type Less[T cmp.Ordered] struct {
	OrEqual bool
	To      T
}

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (l Less[T]) ValidateValue(value reflect.Value) firm.ErrorMap {
	return comparableValidateValue[T](l, value)
}

// Validate validates the data value
func (l Less[T]) Validate(data T) firm.ErrorMap { return comparableValidate[T](l, data) }

// Compare returns true if the data is valid
func (l Less[T]) Compare(data T) bool { return less(l.OrEqual, data, l.To) }

// TypeCheck checks whether the type is valid for the Rule
func (l Less[T]) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	return comparableTypeCheck("Less", l.To, typ)
}

// ErrorMap returns the ErrorMap returned from ValidateValue
func (l Less[T]) ErrorMap() firm.ErrorMap {
	return orderedErrorMap(l.Name(), lessName, l.To, l.OrEqual)
}

// Name returns the name of the Rule
func (l Less[T]) Name() string { return orderedName(lessName, l.OrEqual) }

// Greater checks if data is greater (or equal to) .To
type Greater[T cmp.Ordered] struct {
	OrEqual bool
	To      T
}

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (g Greater[T]) ValidateValue(value reflect.Value) firm.ErrorMap {
	return comparableValidateValue[T](g, value)
}

// Validate validates the data value
func (g Greater[T]) Validate(data T) firm.ErrorMap { return comparableValidate[T](g, data) }

// Compare returns true if the data is valid
func (g Greater[T]) Compare(data T) bool { return !less(!g.OrEqual, data, g.To) }

// TypeCheck checks whether the type is valid for the Rule
func (g Greater[T]) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	return comparableTypeCheck("Greater", g.To, typ)
}

// ErrorMap returns the ErrorMap returned from ValidateValue
func (g Greater[T]) ErrorMap() firm.ErrorMap {
	return orderedErrorMap(g.Name(), greaterName, g.To, g.OrEqual)
}

// Name returns the name of the Rule
func (g Greater[T]) Name() string { return orderedName(greaterName, g.OrEqual) }

type comparableRule[T comparable] interface {
	firm.RuleTyped[T]
	Compare(data T) bool
	Name() string
}

func comparableValidateValue[T comparable](rule comparableRule[T], value reflect.Value) firm.ErrorMap {
	data, ok := value.Interface().(T)
	if !ok {
		panic("comparable ValidateValue type not matching type--called before TypeCheck?")
	}
	return comparableValidate(rule, data)
}
func comparableValidate[T comparable](rule comparableRule[T], data T) firm.ErrorMap {
	if rule.Compare(data) {
		return nil
	}
	return rule.ErrorMap()
}
func comparableTypeCheck[T comparable](ruleName string, to T, typ reflect.Type) *firm.RuleTypeError {
	//nolint:godox // want the comment
	toType := reflect.TypeOf(to) // TODO: cache in struct?, or use cast switch to package level cache?
	if toType == typ {
		return nil
	}
	return firm.NewRuleTypeError(ruleName, typ, "is not a "+toType.String())
}

func less[T cmp.Ordered](orEqual bool, data, to T) bool {
	return cmp.Less(data, to) || (orEqual && data == to)
}

func orderedErrorMap[T cmp.Ordered](name, baseName string, to T, orEqual bool) firm.ErrorMap {
	baseName = string(unicode.ToLower(rune(baseName[0]))) + baseName[1:]
	orEqualTemplate := ""
	if orEqual {
		orEqualTemplate = "or equal to "
	}
	return firm.ErrorMap{firm.ErrorKey(name): firm.TemplateError{
		TemplateFields: map[string]string{"To": fmt.Sprintf("%v", to)},
		Template:       fmt.Sprintf("is not %v than %v{{.To}}", baseName, orEqualTemplate),
	}}
}
func orderedName(name string, orEqual bool) string {
	if orEqual {
		name += "OrEqual"
	}
	return name
}
