package seedkrdict

import (
	"encoding/xml"
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/stringutil"
)

// Seed seeds the database from the rscPath XML
func Seed(tx db.Tx, rscPath string) error {
	lexes, err := UnmarshallRscPath(rscPath)
	if err != nil {
		return err
	}

	basePopularity := 1
	for _, lex := range lexes {
		if err := seedLex(tx, lex, basePopularity); err != nil {
			return err
		}
		basePopularity += len(lex.LexicalEntries)
	}
	return nil
}

// SeedFile seeds a rscPath XML file to the database
func SeedFile(tx db.Tx, file []byte) error {
	lex, err := UnmarshallRscXML(file)
	if err != nil {
		return err
	}
	return seedLex(tx, lex, 1)
}

func seedLex(tx db.Tx, lex *LexicalResource, basePopularity int) error {
	// default to 1
	if basePopularity == 0 {
		basePopularity = 1
	}
	qs := db.New(tx)
	for i, entry := range lex.LexicalEntries {
		createParams, err := entry.CreateParams(basePopularity + i)
		if err != nil {
			if IsNoTranslationsFoundError(err) {
				continue
			}
			return err
		}
		if _, err = qs.TermCreate(tx.Ctx(), createParams); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshallRscPath unmarshalls the rsc path of XML files to LexicalResources
func UnmarshallRscPath(rscPath string) ([]*LexicalResource, error) {
	lexes := []*LexicalResource{}
	xmlPaths, err := RscXMLPaths(rscPath)
	if err != nil {
		return nil, err
	}
	for _, path := range xmlPaths {
		//nolint:gosec // just parsing XML
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		lex, err := UnmarshallRscXML(bytes)
		if err != nil {
			return nil, err
		}
		lexes = append(lexes, lex)
	}
	return lexes, nil
}

// UnmarshallRscXML unmarshalls the rsc XML files to a LexicalResource
func UnmarshallRscXML(bytes []byte) (*LexicalResource, error) {
	lex := &LexicalResource{}
	if err := xml.Unmarshal(bytes, lex); err != nil {
		return nil, err
	}
	return lex, nil
}

// LexicalResource represents a rsc XML file
type LexicalResource struct {
	LexicalEntries []LexicalEntry `xml:"Lexicon>LexicalEntry"`
}

// LexicalEntry represents a entry in the dictionary
type LexicalEntry struct {
	ID    uint   `xml:"val,attr"`
	Feats []Feat `xml:"feat"`

	Lemmas       []Lemma       `xml:"Lemma"`
	RelatedForms []RelatedForm `xml:"RelatedForm"`
	Senses       []Sense       `xml:"Sense"`
	WordForms    []WordForm    `xml:"WordForm"`
}

// CreateParams returns the db create params from the LexicalEntry
func (l *LexicalEntry) CreateParams(popularity int) (db.TermCreateParams, error) {
	term, err := l.Term()
	if err != nil {
		return db.TermCreateParams{}, err
	}
	dbTerm, err := db.ToDBTerm(term, popularity)
	if err != nil {
		return db.TermCreateParams{}, err
	}
	return dbTerm.CreateParams(), nil
}

// NoTranslationsFoundError is returned if no translation is found within a LexicalEntry
type NoTranslationsFoundError struct {
	text string
}

func (e NoTranslationsFoundError) Error() string {
	return fmt.Sprintf("no translations found with text: %v", e.text)
}

// IsNoTranslationsFoundError returns true if the error is a NoTranslationsFoundError
func IsNoTranslationsFoundError(err error) bool {
	ok := errors.As(err, &NoTranslationsFoundError{})
	return ok
}

// Term returns the dictionary.Term from the LexicalEntry
func (l *LexicalEntry) Term() (dictionary.Term, error) {
	text, variants, err := l.textVariants()
	if err != nil {
		return dictionary.Term{}, err
	}
	pos, commonLevel, err := l.posCommonLevel()
	if err != nil {
		return dictionary.Term{}, fmt.Errorf("%w with text: %v", err, text)
	}
	translations := l.translations()
	if len(translations) == 0 {
		return dictionary.Term{}, NoTranslationsFoundError{text: text}
	}
	return dictionary.Term{
		ID:           int64(l.ID),
		Text:         text,
		Variants:     variants,
		PartOfSpeech: pos,
		CommonLevel:  commonLevel,
		Translations: translations,
	}, nil
}

func (l *LexicalEntry) textVariants() (string, []string, error) {
	text, variants := "", []string{}
	for _, lemma := range l.Lemmas {
		innerText, innerVariants := lemma.textOrVariants()
		if innerText != "" {
			text = innerText
		}
		if len(innerVariants) != 0 {
			variants = innerVariants
		}
	}

	var err error
	if text == "" {
		err = fmt.Errorf("LexicalEntry.writtenForm not found")
	}
	return text, variants, err
}

var partOfSpeechMap = map[string]lang.PartOfSpeech{
	"명사":  lang.PartOfSpeechNoun,
	"대명사": lang.PartOfSpeechPronoun,
	"수사":  lang.PartOfSpeechNumeral,
	"조사":  lang.PartOfSpeechPostposition,

	"동사":  lang.PartOfSpeechVerb,
	"형용사": lang.PartOfSpeechAdjective,
	"관형사": lang.PartOfSpeechDeterminer,

	"부사":  lang.PartOfSpeechAdverb,
	"감탄사": lang.PartOfSpeechInterjection,

	"접사": lang.PartOfSpeechAffix,

	"의존 명사": lang.PartOfSpeechDependentNoun,

	"보조 동사":  lang.PartOfSpeechAuxiliaryVerb,
	"보조 형용사": lang.PartOfSpeechAuxiliaryAdjective,

	"어미": lang.PartOfSpeechEnding,

	"품사 없음": lang.PartOfSpeechUnknown,
	"":      lang.PartOfSpeechUnknown,
}

// PartOfSpeechMap returns a copy of the partOfSpeechMap
func PartOfSpeechMap() map[string]lang.PartOfSpeech {
	return maps.Clone(partOfSpeechMap)
}

var vocabularyToCommonLevel = map[string]lang.CommonLevel{
	"":   lang.CommonLevelUnique,
	"없음": lang.CommonLevelUnique,
	"고급": lang.CommonLevelRare,
	"중급": lang.CommonLevelMedium,
	"초급": lang.CommonLevelCommon,
}

func (l *LexicalEntry) posCommonLevel() (lang.PartOfSpeech, lang.CommonLevel, error) {
	posSrc, vocabularyLevel := "", ""
	for _, feat := range l.Feats {
		switch feat.Att {
		case "partOfSpeech":
			posSrc = feat.Val
		case "vocabularyLevel":
			vocabularyLevel = feat.Val
		}
	}

	var err error
	pos, exists := partOfSpeechMap[posSrc]
	if !exists {
		err = fmt.Errorf("part of speech not found: %v", posSrc)
	}
	commonLevel, exists := vocabularyToCommonLevel[vocabularyLevel]
	if !exists {
		err = fmt.Errorf("common level not found: %v", vocabularyLevel)
	}
	return pos, commonLevel, err
}

func (l *LexicalEntry) translations() []dictionary.Translation {
	translations := []dictionary.Translation{}
	for _, sense := range l.Senses {
		translation, err := sense.translation()
		if err != nil {
			continue
		}
		translations = append(translations, translation)
	}
	return translations
}

// Lemma contains the written forms and variants of the entry
type Lemma struct {
	Feats []Feat `xml:"feat"`
}

func (l *Lemma) textOrVariants() (string, []string) {
	text, variants := "", []string{}
	for _, feat := range l.Feats {
		switch feat.Att {
		case "writtenForm":
			text = feat.Val
		case "variant":
			variants = stringutil.SplitClean(feat.Val, ",")
		}
	}
	return text, variants
}

// RelatedForm contains the related words of the entry
type RelatedForm struct {
	Feats []Feat `xml:"feat"`
}

// Sense represents the meaning of the entry for all languages
type Sense struct {
	ID    uint   `xml:"val,attr"`
	Feats []Feat `xml:"feat"`

	Equivalents    []Equivalent    `xml:"Equivalent"`
	Multimedias    []Multimedia    `xml:"Multimedia"`
	SenseExamples  []SenseExample  `xml:"SenseExample"`
	SenseRelations []SenseRelation `xml:"SenseRelation"`
}

func (s *Sense) translation() (dictionary.Translation, error) {
	for _, equiv := range s.Equivalents {
		translation, err := equiv.translation()
		if err != nil {
			continue
		}
		return translation, nil
	}
	return dictionary.Translation{}, fmt.Errorf("not found")
}

// Equivalent represents the translation of the entry given a special language
type Equivalent struct {
	Feats []Feat `xml:"feat"`
}

const engSenseLang = "영어"

var cleanTranslationMap = map[string]string{
	"&quot;": "\"",
}

func (e *Equivalent) translation() (dictionary.Translation, error) {
	isEng := false
	for _, feat := range e.Feats {
		if feat.Val == engSenseLang && feat.Att == "language" {
			isEng = true
			break
		}
	}
	if !isEng {
		return dictionary.Translation{}, fmt.Errorf("not found")
	}

	text, explanation := "", ""
	for _, feat := range e.Feats {
		if feat.Att == "lemma" {
			text = feat.Val
		}
		if feat.Att == "definition" {
			explanation = feat.Val
		}
	}

	var err error
	if text == "" {
		err = fmt.Errorf("text is empty")
	}
	if explanation == "" {
		err = fmt.Errorf("explanation is empty")
	}
	for k, v := range cleanTranslationMap {
		explanation = strings.ReplaceAll(explanation, k, v)
	}
	return dictionary.Translation{
		Text:        text,
		Explanation: explanation,
	}, err
}

// Multimedia represents media related to the Sense
type Multimedia struct {
	Feats []Feat `xml:"feat"`
}

// SenseExample gives examples for the Sense
type SenseExample struct {
	Feats []Feat `xml:"feat"`
}

// SenseRelation gives relations to the Sense
type SenseRelation struct {
	Feats []Feat `xml:"feat"`
}

// WordForm contains other word forms of the entry
type WordForm struct {
	Feats              []Feat             `xml:"feat"`
	FormRepresentation FormRepresentation `xml:"FormRepresentation"`
}

// FormRepresentation contains the word form representation
type FormRepresentation struct {
	Feats []Feat `xml:"feat"`
}

// Feat represents any feature
type Feat struct {
	Att string `xml:"att,attr"`
	Val string `xml:"val,attr"`
}
