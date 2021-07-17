// Package dictionary contains dictonary related functions
package dictionary

// Dicionary is a interface for a dictionary to search with
type Dicionary interface {
	Search(q string) ([]Term, error)
}

// CommonLevel indicates how common a term is
type CommonLevel uint8

// CommonLevels from 0 to 3, where 3 is the most common
const (
	CommonLevelUnique CommonLevel = iota
	CommonLevelRare
	CommonLevelMedium
	CommonLevelCommon
)

// PartOfSpeech are representations of part of speech
type PartOfSpeech string

// PartOfSpeech types
const (
	PartOfSpeechNoun               PartOfSpeech = "Noun"
	PartOfSpeechPronoun            PartOfSpeech = "Pronoun"
	PartOfSpeechNumeral            PartOfSpeech = "Numberal"
	PartOfSpeechPostposition       PartOfSpeech = "Postposition"
	PartOfSpeechVerb               PartOfSpeech = "Verb"
	PartOfSpeechAdjective          PartOfSpeech = "Adjective"
	PartOfSpeechPrenoun            PartOfSpeech = "Prenoun"
	PartOfSpeechAdverb             PartOfSpeech = "Adverb"
	PartOfSpeechInterjection       PartOfSpeech = "Interjection"
	PartOfSpeechAffix              PartOfSpeech = "Affix"
	PartOfSpeechDependentNoun      PartOfSpeech = "DependentNoun"
	PartOfSpeechAuxiliaryVerb      PartOfSpeech = "AuxiliaryVerb"
	PartOfSpeechAuxiliaryAdjective PartOfSpeech = "AuxiliaryAdjective"
	PartOfSpeechEnding             PartOfSpeech = "Ending"
	PartOfSpeechNone               PartOfSpeech = "None"
	PartOfSpeechInvalid            PartOfSpeech = ""
)

// Term is a word or phrase
type Term struct {
	Text         string
	PartOfSpeech PartOfSpeech
	CommonLevel  CommonLevel
	Translations []Translation

	DictionarySource string
}

// Translation is a translation of a word or phrase
type Translation struct {
	Text        string
	Explanation string
}
