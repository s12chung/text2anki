package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/attr"
)

func intEqual(i int) Equal[int] { return Equal[int]{To: i} }

func TestAttr_ErrorMap(t *testing.T) {
	rule := Attr{Of: attr.Len{}, Rule: intEqual(1)}
	testErrorMap(t, rule, "Len-Equal: value attribute, Len, is not equal to 1")
	require.Equal(t, rule.ValidateValue(reflect.ValueOf("")), rule.ErrorMap())
}

func TestAttr_ValidateValue(t *testing.T) {
	tcs := []struct {
		name string
		attr Attribute
		rule firm.RuleBasic

		data     any
		errorMap firm.ErrorMap
	}{
		{name: "normal", data: " ", rule: intEqual(1)},
		{name: "multi", data: " ", rule: intEqual(1)},
		{name: "invalid", data: " ", rule: intEqual(2), errorMap: intEqual(2).ErrorMap()},
		{name: "invalid_with_empty_template_fields", data: "", rule: Present{}, errorMap: errorMapPresent},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			attribute := tc.attr
			if attribute == nil {
				attribute = attr.Len{}
			}
			var expected firm.ErrorMap
			for k, v := range tc.errorMap {
				err := v
				if err.TemplateFields == nil {
					err.TemplateFields = map[string]string{}
				}
				err.TemplateFields["AttrName"] = "Len"
				err.Template = "attribute, {{.AttrName}}, " + err.Template

				if expected == nil {
					expected = firm.ErrorMap{}
				}
				expected["Len-"+k] = err
			}
			require.Equal(expected, Attr{Of: attribute, Rule: tc.rule}.ValidateValue(reflect.ValueOf(tc.data)))
		})
	}
}

func TestAttr_TypeCheck(t *testing.T) {
	tcs := []struct {
		name string
		attr Attribute
		rule firm.RuleBasic

		data         any
		errData      any
		badCondition string
	}{
		{name: "normal", data: " ", rule: intEqual(1)},
		{name: "invalid_data", data: 0, rule: intEqual(2),
			badCondition: "does not have a length (not a Slice, Array, Array pointer, Channel, Map or String)"},
		{name: "invalid_rule", data: " ", rule: Less[string]{To: ""},
			errData: 0, badCondition: "has Attr, Len, which is not a string"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			attribute := tc.attr
			if attribute == nil {
				attribute = attr.Len{}
			}

			typ := reflect.TypeOf(tc.data)
			var err *firm.RuleTypeError
			if tc.badCondition != "" {
				errData := typ
				if tc.errData != nil {
					errData = reflect.TypeOf(tc.errData)
				}
				err = firm.NewRuleTypeError(errData, tc.badCondition)
			}
			require.Equal(err, Attr{Of: attribute, Rule: tc.rule}.TypeCheck(typ))
		})
	}
}
