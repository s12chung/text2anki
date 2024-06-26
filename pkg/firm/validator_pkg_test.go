package firm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
)

type nonExport struct {
	Public  string
	private string
	privateChild
}

type privateChild struct {
	privateC string
}

var notEmpty = nonExport{Public: "not_empty", private: "not_empty", privateChild: privateChild{privateC: "not_empty"}}

func errorMap(field, name firm.ErrorKey, template string) firm.ErrorMap {
	if name == "" && template == "" {
		name = "TrimPresent"
		template = "is just spaces or empty"
	}
	key := "firm_test.nonExport." + field + "." + name
	return firm.ErrorMap{key: firm.TemplateError{Template: template, ErrorKey: key}}
}

func TestNewStructPkg(t *testing.T) {
	tcs := []struct {
		name    string
		ruleMap firm.RuleMap
		failErr error
	}{
		{name: "exported_field", ruleMap: firm.RuleMap{"Public": {rule.TrimPresent{}}},
			failErr: errorMap("Public", "", "")},
		{name: "non_exported_field", ruleMap: firm.RuleMap{"private": {rule.TrimPresent{}}},
			failErr: errorMap("private", "", "")},
		{name: "non_exported_child", ruleMap: firm.RuleMap{"privateChild": {rule.Present{}}},
			failErr: errorMap("privateChild", "Present", "is not present")},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			validator, err := firm.NewStruct[nonExport](tc.ruleMap)
			require.NoError(err)
			require.Nil(validator.ValidateAny(notEmpty))
			require.Equal(tc.failErr, validator.ValidateAny(nonExport{}))
		})
	}
}
