// Package krdict contains functions for a Korean dictionary connected to a database
package krdict

import (
	"context"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
)

// New returns a new KrDict
func New() KrDict { return KrDict{} }

// KrDict is a Korean dictionary connected to a database
type KrDict struct{}

var mergePosMap = map[lang.PartOfSpeech]lang.PartOfSpeech{
	lang.PartOfSpeechNoun:         lang.PartOfSpeechNoun,
	lang.PartOfSpeechPronoun:      lang.PartOfSpeechPronoun,
	lang.PartOfSpeechNumeral:      lang.PartOfSpeechNumeral,
	lang.PartOfSpeechAlphabet:     lang.PartOfSpeechNoun, // Dictionary Encoding
	lang.PartOfSpeechPostposition: lang.PartOfSpeechPostposition,

	lang.PartOfSpeechVerb:       lang.PartOfSpeechVerb,
	lang.PartOfSpeechAdjective:  lang.PartOfSpeechAdjective,
	lang.PartOfSpeechDeterminer: lang.PartOfSpeechDeterminer,

	lang.PartOfSpeechAdverb:       lang.PartOfSpeechAdverb,
	lang.PartOfSpeechInterjection: lang.PartOfSpeechInterjection,

	lang.PartOfSpeechAffix:  lang.PartOfSpeechAffix,
	lang.PartOfSpeechPrefix: lang.PartOfSpeechAffix, // Make them the same
	lang.PartOfSpeechInfix:  lang.PartOfSpeechAffix, // Make them the same
	lang.PartOfSpeechSuffix: lang.PartOfSpeechAffix, // Make them the same

	lang.PartOfSpeechRoot: lang.PartOfSpeechEnding, // Convert, but untested

	lang.PartOfSpeechDependentNoun: lang.PartOfSpeechDependentNoun,

	lang.PartOfSpeechAuxiliaryPredicate: lang.PartOfSpeechUnknown, // Convert
	lang.PartOfSpeechAuxiliaryVerb:      lang.PartOfSpeechAuxiliaryVerb,
	lang.PartOfSpeechAuxiliaryAdjective: lang.PartOfSpeechAuxiliaryAdjective,

	lang.PartOfSpeechEnding:      lang.PartOfSpeechEnding,
	lang.PartOfSpeechCopula:      lang.PartOfSpeechPostposition, // Convert
	lang.PartOfSpeechPunctuation: lang.PartOfSpeechEmpty,        // Skip

	lang.PartOfSpeechOtherLanguage: lang.PartOfSpeechEmpty, // Skip
	lang.PartOfSpeechOther:         lang.PartOfSpeechEmpty, // Skip
	lang.PartOfSpeechUnknown:       lang.PartOfSpeechUnknown,
	lang.PartOfSpeechEmpty:         lang.PartOfSpeechEmpty,
}

// Search searches for the query inside the dictionary
func (k KrDict) Search(q string, pos lang.PartOfSpeech) ([]dictionary.Term, error) {
	pos = mergePosMap[pos]

	txQs, err := db.NewTxQs(context.Background())
	if err != nil {
		return nil, err
	}
	rows, err := txQs.TermsSearch(txQs.Ctx(), q, pos)
	if err != nil {
		return nil, err
	}
	terms := make([]dictionary.Term, len(rows))
	for i, row := range rows {
		term, err := row.Term.DictionaryTerm()
		if err != nil {
			return nil, err
		}
		terms[i] = term
	}
	return terms, txQs.Commit()
}
