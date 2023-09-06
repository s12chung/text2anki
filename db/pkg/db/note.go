package db

import (
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
)

func init() {
	firm.RegisterType(firm.NewDefinition(NoteCreateParams{}).Validates(firm.RuleMap{
		"Text":        {rule.Presence{}},
		"Translation": {rule.Presence{}},

		"Explanation":      {rule.Presence{}},
		"Usage":            {rule.Presence{}},
		"UsageTranslation": {rule.Presence{}},

		"DictionarySource": {rule.Presence{}},
	}))
}

// StaticCopy returns a copy without fields that variate
func (n Note) StaticCopy() any {
	c := n
	c.ID = 0
	return c
}

// CreateParams converts the note to a NoteCreateParams
func (n Note) CreateParams() NoteCreateParams {
	return NoteCreateParams{
		Text:         n.Text,
		PartOfSpeech: n.PartOfSpeech,
		Translation:  n.Translation,

		CommonLevel:      n.CommonLevel,
		Explanation:      n.Explanation,
		Usage:            n.Usage,
		UsageTranslation: n.UsageTranslation,

		SourceName:       n.SourceName,
		SourceReference:  n.SourceReference,
		DictionarySource: n.DictionarySource,
		Notes:            n.Notes,
	}
}
