package rule

import (
	"cmp"
	"fmt"
	"reflect"
	"unicode"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Equal checks if data is equal to .To
type Equal[T comparable] struct{ To T }

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (e Equal[T]) ValidateValue(value reflect.Value) firm.ErrorMap {
	return comparableValidateValue[T](e, value)
}

// Compare returns true if the data is valid
func (e Equal[T]) Compare(data T) bool { return data == e.To }

// ErrorMap returns the ErrorMap returned from ValidateValue
func (e Equal[T]) ErrorMap() firm.ErrorMap {
	return firm.ErrorMap{"Equal": firm.TemplateError{
		TemplateFields: map[string]string{"To": fmt.Sprintf("%v", e.To)},
		Template:       "is not equal to {{.To}}",
	}}
}

// TypeCheck checks whether the type is valid for the Rule
func (e Equal[T]) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	return comparableTypeCheck(e.To, typ)
}

// Less checks if data is less (or equal to) .To
type Less[T cmp.Ordered] struct {
	OrEqual bool
	To      T
}

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (l Less[T]) ValidateValue(value reflect.Value) firm.ErrorMap {
	return comparableValidateValue[T](l, value)
}

// Compare returns true if the data is valid
func (l Less[T]) Compare(data T) bool { return less(l.OrEqual, data, l.To) }

// ErrorMap returns the ErrorMap returned from ValidateValue
func (l Less[T]) ErrorMap() firm.ErrorMap { return orderedErrorMap("less", l.To, l.OrEqual) }

// TypeCheck checks whether the type is valid for the Rule
func (l Less[T]) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	return comparableTypeCheck(l.To, typ)
}

// Greater checks if data is greater (or equal to) .To
type Greater[T cmp.Ordered] struct {
	OrEqual bool
	To      T
}

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (g Greater[T]) ValidateValue(value reflect.Value) firm.ErrorMap {
	return comparableValidateValue[T](g, value)
}

// Compare returns true if the data is valid
func (g Greater[T]) Compare(data T) bool { return !less(!g.OrEqual, data, g.To) }

// ErrorMap returns the ErrorMap returned from ValidateValue
func (g Greater[T]) ErrorMap() firm.ErrorMap { return orderedErrorMap("greater", g.To, g.OrEqual) }

// TypeCheck checks whether the type is valid for the Rule
func (g Greater[T]) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	return comparableTypeCheck(g.To, typ)
}

type comparableRule[T comparable] interface {
	firm.RuleBasic
	Compare(data T) bool
}

func comparableValidateValue[T comparable](rule comparableRule[T], value reflect.Value) firm.ErrorMap {
	data, ok := value.Interface().(T)
	if !ok {
		panic("comparable ValidateValue type not matching type--called before TypeCheck?")
	}
	if rule.Compare(data) {
		return nil
	}
	return rule.ErrorMap()
}
func comparableTypeCheck[T comparable](to T, typ reflect.Type) *firm.RuleTypeError {
	//nolint:godox // want the comment
	toType := reflect.TypeOf(to) // TODO: cache in struct?, or use cast switch to package level cache?
	if toType == typ {
		return nil
	}
	return firm.NewRuleTypeError(typ, "is not a "+toType.String())
}

func less[T cmp.Ordered](orEqual bool, data, to T) bool {
	return cmp.Less(data, to) || (orEqual && data == to)
}
func orderedErrorMap[T cmp.Ordered](name string, to T, orEqual bool) firm.ErrorMap {
	fullName := string(unicode.ToUpper(rune(name[0]))) + name[1:]
	orEqualTemplate := ""
	if orEqual {
		fullName += "OrEqual"
		orEqualTemplate = "or equal to "
	}
	return firm.ErrorMap{firm.ErrorKey(fullName): firm.TemplateError{
		TemplateFields: map[string]string{"To": fmt.Sprintf("%v", to)},
		Template:       fmt.Sprintf("is not %v than %v{{.To}}", name, orEqualTemplate),
	}}
}
