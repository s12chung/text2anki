package krdict

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/stringclean"
)

// Seed seeds the database from the rscPath XML
func Seed() error {
	lexes, err := unmarshallRscPath()
	if err != nil {
		return err
	}

	ctx := context.Background()
	queries := db.New(db.DB())
	basePopularity := 1
	for _, lex := range lexes {
		for i, entry := range lex.LexicalEntries {
			createParams, err := entry.createParams(basePopularity + i)
			if err != nil {
				if IsNoTranslationsFoundError(err) {
					continue
				}
				return err
			}
			if _, err = queries.TermCreate(ctx, createParams); err != nil {
				return err
			}
		}
	}
	return nil
}

func unmarshallRscPath() ([]*lexicalResource, error) {
	lexes := []*lexicalResource{}
	xmlPaths, err := RscXMLPaths()
	if err != nil {
		return nil, err
	}
	for _, path := range xmlPaths {
		//nolint:gosec // just parsing XML
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		lex, err := unmarshallXML(bytes)
		if err != nil {
			return nil, err
		}
		lexes = append(lexes, lex)
	}
	return lexes, nil
}

func unmarshallXML(bytes []byte) (*lexicalResource, error) {
	lex := &lexicalResource{}
	if err := xml.Unmarshal(bytes, lex); err != nil {
		return nil, err
	}
	return lex, nil
}

type lexicalResource struct {
	LexicalEntries []lexicalEntry `xml:"Lexicon>LexicalEntry"`
}

type lexicalEntry struct {
	ID    uint   `xml:"val,attr"`
	Feats []feat `xml:"feat"`

	Lemmas       []lemma       `xml:"Lemma"`
	RelatedForms []relatedForm `xml:"RelatedForm"`
	Senses       []sense       `xml:"Sense"`
	WordForms    []wordForm    `xml:"WordForm"`
}

func (l *lexicalEntry) createParams(popularity int) (db.TermCreateParams, error) {
	term, err := l.term()
	if err != nil {
		return db.TermCreateParams{}, err
	}
	variants, err := json.Marshal(term.Variants)
	if err != nil {
		return db.TermCreateParams{}, err
	}
	translations, err := json.Marshal(term.Translations)
	if err != nil {
		return db.TermCreateParams{}, err
	}

	return db.TermCreateParams{
		Text:         term.Text,
		Variants:     string(variants),
		PartOfSpeech: string(term.PartOfSpeech),
		CommonLevel:  strconv.Itoa(int(term.CommonLevel)),
		Translations: string(translations),
		Popularity:   strconv.Itoa(popularity),
	}, nil
}

type noTranslationsFoundError struct {
	text string
}

func (e *noTranslationsFoundError) Error() string {
	return fmt.Sprintf("no translations found with text: %v", e.text)
}

// IsNoTranslationsFoundError returns true if the error is a noTranslationsFoundError
func IsNoTranslationsFoundError(err error) bool {
	_, ok := err.(*noTranslationsFoundError)
	return ok
}

func (l *lexicalEntry) term() (dictionary.Term, error) {
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
		return dictionary.Term{}, &noTranslationsFoundError{text: text}
	}

	return dictionary.Term{
		Text:         text,
		Variants:     variants,
		PartOfSpeech: pos,
		CommonLevel:  commonLevel,
		Translations: translations,
	}, nil
}

func (l *lexicalEntry) textVariants() (string, []string, error) {
	text, variants := "", []string{}
	for _, lemma := range l.Lemmas {
		for _, feat := range lemma.Feats {
			switch feat.Att {
			case "writtenForm":
				text = feat.Val
			case "variant":
				variants = stringclean.Split(feat.Val, ",")
			}
		}
	}

	var err error
	if text == "" {
		err = fmt.Errorf("lexicalEntry.writtenForm not found")
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

var vocabularyToCommonLevel = map[string]lang.CommonLevel{
	"":   lang.CommonLevelUnique,
	"없음": lang.CommonLevelUnique,
	"고급": lang.CommonLevelRare,
	"중급": lang.CommonLevelMedium,
	"초급": lang.CommonLevelCommon,
}

func (l *lexicalEntry) posCommonLevel() (lang.PartOfSpeech, lang.CommonLevel, error) {
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

func (l *lexicalEntry) translations() []dictionary.Translation {
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

type lemma struct {
	Feats []feat `xml:"feat"`
}

type relatedForm struct {
	Feats []feat `xml:"feat"`
}

type sense struct {
	ID    uint   `xml:"val,attr"`
	Feats []feat `xml:"feat"`

	Equivalents    []equivalent    `xml:"Equivalent"`
	Multimedias    []multimedia    `xml:"Multimedia"`
	SenseExamples  []senseExample  `xml:"SenseExample"`
	SenseRelations []senseRelation `xml:"SenseRelation"`
}

func (s *sense) translation() (dictionary.Translation, error) {
	for _, equiv := range s.Equivalents {
		translation, err := equiv.translation()
		if err != nil {
			continue
		}
		return translation, nil
	}
	return dictionary.Translation{}, fmt.Errorf("not found")
}

type equivalent struct {
	Feats []feat `xml:"feat"`
}

const engSenseLang = "영어"

func (e *equivalent) translation() (dictionary.Translation, error) {
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
	return dictionary.Translation{
		Text:        text,
		Explanation: explanation,
	}, err
}

type multimedia struct {
	Feats []feat `xml:"feat"`
}

type senseExample struct {
	Feats []feat `xml:"feat"`
}

type senseRelation struct {
	Feats []feat `xml:"feat"`
}

type wordForm struct {
	Feats              []feat             `xml:"feat"`
	FormRepresentation formRepresentation `xml:"FormRepresentation"`
}

type formRepresentation struct {
	Feats []feat `xml:"feat"`
}

type feat struct {
	Att string `xml:"att,attr"`
	Val string `xml:"val,attr"`
}
