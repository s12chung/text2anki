// Package lang contains common things for languages
package lang

import (
	"fmt"
)

// CommonLevel indicates how common a term is
type CommonLevel uint8

// CommonLevels from 0 to 3, where 3 is the most common
const (
	CommonLevelUnique CommonLevel = iota
	CommonLevelRare
	CommonLevelMedium
	CommonLevelCommon
)

// ToCommonLevel returns the common level of commonLevel
func ToCommonLevel(commonLevel int64) (CommonLevel, error) {
	if commonLevel < 0 || commonLevel > int64(CommonLevelCommon) {
		return 0, fmt.Errorf("common level not within range %v to %v: %v", 0, CommonLevelCommon, commonLevel)
	}
	return CommonLevel(commonLevel), nil
}

// PartOfSpeech are representations of part of speech
type PartOfSpeech string

// PartOfSpeech types
const (
	// PartOfSpeechCount is the number of Part Of Speech Types - mostly for testing
	PartOfSpeechCount = 26

	//
	// UPDATE PartOfSpeechCount and partOfSpeechTypes WHEN CHANGING
	//
	// Search "PartOfSpeech EQUAL_SIGN" to count
	//
	PartOfSpeechNoun         PartOfSpeech = "Noun"         // 명사
	PartOfSpeechPronoun      PartOfSpeech = "Pronoun"      // 대명사
	PartOfSpeechNumeral      PartOfSpeech = "Numeral"      // 수사
	PartOfSpeechAlphabet     PartOfSpeech = "Alphabet"     // 한글 자소
	PartOfSpeechPostposition PartOfSpeech = "Postposition" // 조사

	PartOfSpeechVerb       PartOfSpeech = "Verb"       // 동사
	PartOfSpeechAdjective  PartOfSpeech = "Adjective"  // 형용사
	PartOfSpeechDeterminer PartOfSpeech = "Determiner" // 보조 용언

	PartOfSpeechAdverb       PartOfSpeech = "Adverb"       // 부사
	PartOfSpeechInterjection PartOfSpeech = "Interjection" // 감탄사

	PartOfSpeechAffix  PartOfSpeech = "Affix"  // General - 접사
	PartOfSpeechPrefix PartOfSpeech = "Prefix" // 체언 접두사
	PartOfSpeechInfix  PartOfSpeech = "Infix"
	PartOfSpeechSuffix PartOfSpeech = "Suffix" // 파생 접미사

	PartOfSpeechRoot PartOfSpeech = "Root" // 어근

	PartOfSpeechDependentNoun PartOfSpeech = "DependentNoun" // 의존 명사

	PartOfSpeechAuxiliaryPredicate PartOfSpeech = "AuxiliaryPredicate" // General - 보조 용언
	PartOfSpeechAuxiliaryVerb      PartOfSpeech = "AuxiliaryVerb"      // 보조 동사
	PartOfSpeechAuxiliaryAdjective PartOfSpeech = "AuxiliaryAdjective" // 보조 형용사

	PartOfSpeechEnding      PartOfSpeech = "Ending" // 어미
	PartOfSpeechCopula      PartOfSpeech = "Copula" // 지정사
	PartOfSpeechPunctuation PartOfSpeech = "Punctuation"

	PartOfSpeechOtherLanguage PartOfSpeech = "OtherLanguage"
	PartOfSpeechOther         PartOfSpeech = "Other"
	PartOfSpeechUnknown       PartOfSpeech = "Unknown"
	PartOfSpeechEmpty         PartOfSpeech = ""
)

var partOfSpeechTypes = []PartOfSpeech{
	PartOfSpeechNoun,
	PartOfSpeechPronoun,
	PartOfSpeechNumeral,
	PartOfSpeechAlphabet,
	PartOfSpeechPostposition,

	PartOfSpeechVerb,
	PartOfSpeechAdjective,
	PartOfSpeechDeterminer,
	PartOfSpeechAdverb,
	PartOfSpeechInterjection,

	PartOfSpeechAffix,
	PartOfSpeechPrefix,
	PartOfSpeechInfix,
	PartOfSpeechSuffix,

	PartOfSpeechRoot,

	PartOfSpeechDependentNoun,

	PartOfSpeechAuxiliaryPredicate,
	PartOfSpeechAuxiliaryVerb,
	PartOfSpeechAuxiliaryAdjective,

	PartOfSpeechEnding,
	PartOfSpeechCopula,
	PartOfSpeechPunctuation,

	PartOfSpeechOtherLanguage,
	PartOfSpeechOther,
	PartOfSpeechUnknown,
	PartOfSpeechEmpty,
}

// PartOfSpeechTypes returns a map of all Part Of Speech Types
func PartOfSpeechTypes() map[string]PartOfSpeech {
	m := map[string]PartOfSpeech{}
	for _, pos := range partOfSpeechTypes {
		m[string(pos)] = pos
	}
	return m
}

var partOfSpeechMap = PartOfSpeechTypes()

// ToPartOfSpeech converts the posString to a PartOfSpeech
func ToPartOfSpeech(posString string) (PartOfSpeech, error) {
	pos, exists := partOfSpeechMap[posString]
	if !exists {
		return PartOfSpeechEmpty, fmt.Errorf("pos not matching lang.PartOfSpeech: %v", posString)
	}
	return pos, nil
}
