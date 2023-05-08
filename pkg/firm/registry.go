package firm

import (
	"fmt"
	"reflect"
)

// Registry registers types find the right validator to validate with
type Registry struct {
	typeToDefinition map[reflect.Type]StructuredDefinition
	DefaultValidator Validator
}

// RegisterType registers the TypedDefinition
func (r *Registry) RegisterType(definition *TypedDefinition) {
	if r.typeToDefinition == nil {
		r.typeToDefinition = map[reflect.Type]StructuredDefinition{}
	}

	typ := definition.typ
	if _, exists := r.typeToDefinition[typ]; exists {
		panic(fmt.Sprintf("RegisterType() with type %v already exists", typ.Name()))
	}
	// not more need for .typ, so take StructDefinition
	r.typeToDefinition[typ] = definition
}

// Validate validates the data with the correct validator
func (r *Registry) Validate(data any) Result {
	return MapResult{errorMap: r.ValidateValue(reflect.ValueOf(data))}
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
	definition := r.Definition(value)
	if definition != nil {
		return definition.Validator(r)
	}
	return nil
}

// Definition returns the definition for the value
func (r *Registry) Definition(value reflect.Value) StructuredDefinition {
	value = indirect(value)
	if !value.IsValid() {
		return nil
	}
	return r.DefinitionForType(value.Type())
}

// DefinitionForType returns the definition for the type
func (r *Registry) DefinitionForType(typ reflect.Type) StructuredDefinition {
	if typ == nil || r.typeToDefinition == nil {
		return nil
	}
	definition, exists := r.typeToDefinition[typ]
	if !exists {
		return nil
	}
	return definition
}
