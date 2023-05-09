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
		panic(fmt.Sprintf("RegisterType() with type %v already exists", typ.Name()))
	}

	structValidator := NewStructValidator(definition.RuleMap())
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
		validator := NewSliceValidator()
		*rules = append(*rules, &validator)
		r.registerRecursionType(typ.Elem(), &validator.ElementRules)
	}
}

// Validate validates the data with the correct validator
func (r *Registry) Validate(data any) Result {
	return r.DefaultedValidator(reflect.ValueOf(data)).Validate(data)
}

// ValidateValue validates the data value with the correct validator
func (r *Registry) ValidateValue(value reflect.Value) ErrorMap {
	return r.DefaultedValidator(value).ValidateValue(value)
}

// ValidateMerge validates the data value with the correct validator, also doing a merge with the errorMap
func (r *Registry) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	r.DefaultedValidator(value).ValidateMerge(value, key, errorMap)
}

// DefaultedValidator returns the validator for the value, defaulted by r.DefaultValidator, then DefaultValidator
func (r *Registry) DefaultedValidator(value reflect.Value) Validator {
	validator := r.Validator(value)
	if validator != nil {
		return validator
	}
	if r.DefaultValidator != nil {
		return r.DefaultValidator
	}
	return DefaultValidator
}

// Validator returns the validator for the value (not defaulted)
func (r *Registry) Validator(value reflect.Value) Validator {
	value = indirect(value)
	if !value.IsValid() {
		return nil
	}
	return r.ValidatorForType(value.Type())
}

// ValidatorForType returns the validator for the type (not defaulted)
func (r *Registry) ValidatorForType(typ reflect.Type) Validator {
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
