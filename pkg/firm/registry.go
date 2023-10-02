package firm

import (
	"fmt"
	"reflect"
)

// Registry registers types find the right validator to validate with
type Registry struct {
	typeToValidator            map[reflect.Type]*ValueValidator
	unregisteredTypeReferences map[reflect.Type][]*[]Rule
	DefaultValidator           Validator
}

// RegisterType registers the Definition to validate the type
func (r *Registry) RegisterType(definition *Definition) {
	if r.typeToValidator == nil {
		r.typeToValidator = map[reflect.Type]*ValueValidator{}
		r.unregisteredTypeReferences = map[reflect.Type][]*[]Rule{}
	}

	typ := definition.typ
	if _, exists := r.typeToValidator[typ]; exists {
		panic(fmt.Sprintf("RegisterType() with type %v already exists", typ.String()))
	}

	structValidator := NewStructValidator(typ, definition.RuleMap())
	for fieldName := range structValidator.RuleMap {
		field, _ := typ.FieldByName(fieldName)
		r.registerRecursionType(field.Type, structValidator.RuleMap[fieldName])
	}

	validator := NewValueValidator(append(definition.TopLevelRules(), &structValidator)...)
	r.typeToValidator[typ] = &validator

	for _, rules := range r.unregisteredTypeReferences[typ] {
		*rules = append(*rules, r.typeToValidator[typ])
	}
	delete(r.unregisteredTypeReferences, typ)
}

func (r *Registry) registerRecursionType(typ reflect.Type, rules *[]Rule) {
	typ = indirectType(typ)

	//nolint:exhaustive // just need these cases
	switch typ.Kind() {
	case reflect.Struct:
		validator := r.typeToValidator[typ]
		if validator != nil {
			*rules = append(*rules, validator.Rules...)
		} else {
			references, exists := r.unregisteredTypeReferences[typ]
			if !exists {
				references = []*[]Rule{}
			}
			r.unregisteredTypeReferences[typ] = append(references, rules)
		}
	case reflect.Slice, reflect.Array:
		validator := NewSliceValidator(typ)
		*rules = append(*rules, &validator)
		r.registerRecursionType(typ.Elem(), &validator.ElementRules)
	}
}

// Validate validates the data with the correct validator
func (r *Registry) Validate(data any) ErrorMap {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return errInvalidValue
	}
	return validateValueResult(value, r.DefaultedValidator(value.Type()))
}

// ValidateValue validates the data value with the correct validator (assumes ValidateType is called)
func (r *Registry) ValidateValue(value reflect.Value) ErrorMap {
	return r.DefaultedValidator(value.Type()).ValidateValue(value)
}

// ValidateType checks whether the type is valid for the Rule
func (r *Registry) ValidateType(typ reflect.Type) *RuleTypeError {
	return r.DefaultedValidator(typ).ValidateType(typ)
}

// ValidateMerge validates the data value with the correct validator, also doing a merge with the errorMap (assumes ValidateType is called)
func (r *Registry) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	r.DefaultedValidator(value.Type()).ValidateMerge(value, key, errorMap)
}

// DefaultedValidator returns the validator for the value, defaulted by r.DefaultValidator, then DefaultValidator
func (r *Registry) DefaultedValidator(typ reflect.Type) Validator {
	validator := r.Validator(typ)
	if validator != nil {
		return validator
	}
	if r.DefaultValidator != nil {
		return r.DefaultValidator
	}
	return DefaultValidator
}

// Validator returns the validator for the type (not defaulted)
func (r *Registry) Validator(typ reflect.Type) Validator {
	if typ == nil || r.typeToValidator == nil {
		return nil
	}
	typ = indirectType(typ)
	validator := r.typeToValidator[typ]
	if validator == nil {
		return nil
	}
	return validator
}
