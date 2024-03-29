package db

import (
	"time"

	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/lang"
)

func init() {
	firm.MustRegisterType(firm.NewDefinition[NoteCreateParams]().Validates(firm.RuleMap{
		"Text":        {rule.Present{}},
		"Translation": {rule.Present{}},
		"Explanation": {rule.Present{}},

		"Usage":            {rule.Present{}},
		"UsageTranslation": {rule.Present{}},

		"SourceName":       {rule.Present{}},
		"DictionarySource": {rule.Present{}},
	}))
}

// StaticCopy returns a copy without fields that variate
func (n Note) StaticCopy() Note {
	c := n
	c.ID = 0
	c.UpdatedAt = time.Time{}
	c.CreatedAt = time.Time{}
	return c
}

// CreateParams converts the note to a NoteCreateParams
func (n Note) CreateParams() NoteCreateParams {
	return NoteCreateParams{
		Text:         n.Text,
		PartOfSpeech: n.PartOfSpeech,
		Translation:  n.Translation,
		Explanation:  n.Explanation,
		CommonLevel:  n.CommonLevel,

		Usage:            n.Usage,
		UsageTranslation: n.UsageTranslation,

		SourceName:       n.SourceName,
		SourceReference:  n.SourceReference,
		DictionarySource: n.DictionarySource,
		Notes:            n.Notes,
	}
}

// Anki returns the anki.Note representation
func (n Note) Anki() (anki.Note, error) {
	pos, err := lang.ToPartOfSpeech(n.PartOfSpeech)
	if err != nil {
		return anki.Note{}, err
	}
	commonLevel, err := lang.ToCommonLevel(n.CommonLevel)
	if err != nil {
		return anki.Note{}, err
	}
	return anki.Note{
		Text:         n.Text,
		PartOfSpeech: pos,
		Translation:  n.Translation,
		Explanation:  n.Explanation,
		CommonLevel:  commonLevel,

		Usage:            n.Usage,
		UsageTranslation: n.UsageTranslation,

		SourceName:       n.SourceName,
		SourceReference:  n.SourceReference,
		DictionarySource: n.DictionarySource,
		Notes:            n.Notes,
	}, nil
}

// AnkiNotes converts the []Note to []anki.Note
func AnkiNotes(notes []Note) ([]anki.Note, error) {
	ankiNotes := make([]anki.Note, len(notes))
	var err error
	for i, note := range notes {
		ankiNotes[i], err = note.Anki()
		if err != nil {
			return nil, err
		}
	}
	return ankiNotes, nil
}
