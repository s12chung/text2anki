// Package lang contains common things for languages
package lang

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
