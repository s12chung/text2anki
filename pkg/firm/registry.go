package firm

import (
	"fmt"
	"reflect"
)

// Registry registers types find the right validator to validate with
type Registry struct {
	typeToValidator      map[reflect.Type]*ValueAny
	unregisteredTypeRefs map[reflect.Type][]*[]Rule
	DefaultValidator     Validator
}

// MustRegisterType registers the Definition to validate the type, panics if there is an error
func (r *Registry) MustRegisterType(definition *Definition) {
	if err := r.RegisterType(definition); err != nil {
		panic(err.Error())
	}
}

// RegisterType registers the Definition to validate the type
func (r *Registry) RegisterType(definition *Definition) error {
	if r.typeToValidator == nil {
		r.typeToValidator = map[reflect.Type]*ValueAny{}
		r.unregisteredTypeRefs = map[reflect.Type][]*[]Rule{}
	}

	typ := definition.typ
	if _, exists := r.typeToValidator[typ]; exists {
		return fmt.Errorf("RegisterType() with type %v already exists", typ.String())
	}

	r.typeToValidator[typ] = r.registeredStruct(definition)
	for _, rules := range r.unregisteredTypeRefs[typ] {
		*rules = append(*rules, r.typeToValidator[typ])
	}
	delete(r.unregisteredTypeRefs, typ)
	return nil
}

func (r *Registry) registeredStruct(definition *Definition) *ValueAny {
	typ := definition.Type()
	valueValidatorRules := definition.TopLevelRules()
	if len(definition.RuleMap()) > 0 {
		structValidator := mustNewValidator(func() (StructAny, error) { return NewStructAny(definition.typ, definition.RuleMap()) })
		for fieldName := range structValidator.ruleMap {
			field, _ := typ.FieldByName(fieldName)
			r.registerRecursionType(field.Type, structValidator.ruleMap[fieldName])
		}
		valueValidatorRules = append(valueValidatorRules, structValidator)
	}
	v := mustNewValidator(func() (ValueAny, error) { return NewValueAny(typ, valueValidatorRules...) })
	return &v
}

func (r *Registry) registerRecursionType(typ reflect.Type, rules *[]Rule) {
	typ = indirectType(typ)

	//nolint:exhaustive // just need these cases
	switch typ.Kind() {
	case reflect.Struct:
		validator := r.typeToValidator[typ]
		if validator == nil {
			// when type is registered, appends to the unregisteredTypeRef, similar to inside the else statement
			r.unregisteredTypeRefs[typ] = append(r.unregisteredTypeRefs[typ], rules)
		} else {
			*rules = append(*rules, validator.Rules()...) // add existing type rules
		}
	case reflect.Slice, reflect.Array:
		// No access to field type via generics
		validator := mustNewValidator(func() (SliceAny, error) { return NewSliceAny(typ) })
		*rules = append(*rules, &validator)
		r.registerRecursionType(typ.Elem(), &validator.elementRules)
	}
}

// Type returns the Type the Registry handles
func (r *Registry) Type() reflect.Type { return r.Validator(nil).Type() }

// Validate validates the data with the correct validator
func (r *Registry) Validate(data any) ErrorMap {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return errInvalidValue
	}
	return validateValueResult(r.DefaultedValidator(value.Type()), value)
}

// ValidateValue validates the data value with the correct validator (assumes TypeCheck is called)
func (r *Registry) ValidateValue(value reflect.Value) ErrorMap {
	return r.DefaultedValidator(value.Type()).ValidateValue(value)
}

// TypeCheck checks whether the type is valid for the Rule
func (r *Registry) TypeCheck(typ reflect.Type) *RuleTypeError {
	return r.DefaultedValidator(typ).TypeCheck(typ)
}

// ValidateMerge validates the data value with the correct validator, also doing a merge with the errorMap (assumes TypeCheck is called)
func (r *Registry) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
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
