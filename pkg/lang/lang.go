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
	PartOfSpeechNoun         PartOfSpeech = "Noun"
	PartOfSpeechPronoun      PartOfSpeech = "Pronoun"
	PartOfSpeechNumeral      PartOfSpeech = "Numeral"
	PartOfSpeechAlphabet     PartOfSpeech = "Alphabet"
	PartOfSpeechPostposition PartOfSpeech = "Postposition"
	PartOfSpeechVerb         PartOfSpeech = "Verb"
	PartOfSpeechAdjective    PartOfSpeech = "Adjective"
	PartOfSpeechDeterminer   PartOfSpeech = "Determiner"
	PartOfSpeechAdverb       PartOfSpeech = "Adverb"
	PartOfSpeechInterjection PartOfSpeech = "Interjection"

	PartOfSpeechAffix  PartOfSpeech = "Affix" // General
	PartOfSpeechPrefix PartOfSpeech = "Prefix"
	PartOfSpeechSuffix PartOfSpeech = "Suffix"
	PartOfSpeechRoot   PartOfSpeech = "Root"

	PartOfSpeechDependentNoun PartOfSpeech = "DependentNoun"

	PartOfSpeechAuxiliaryPredicate PartOfSpeech = "AuxiliaryPredicate" // General
	PartOfSpeechAuxiliaryVerb      PartOfSpeech = "AuxiliaryVerb"
	PartOfSpeechAuxiliaryAdjective PartOfSpeech = "AuxiliaryAdjective"

	PartOfSpeechEnding      PartOfSpeech = "Ending"
	PartOfSpeechCopula      PartOfSpeech = "Copula"
	PartOfSpeechPunctuation PartOfSpeech = "Punctuation"

	PartOfSpeechOtherLanguage PartOfSpeech = "OtherLanguage"
	PartOfSpeechOther         PartOfSpeech = "Other"
	PartOfSpeechUnknown       PartOfSpeech = "Unknown"
	PartOfSpeechInvalid       PartOfSpeech = ""
)
