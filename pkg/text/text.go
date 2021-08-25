// Package text contains functions to separate text into source and translation
package text

import (
	"fmt"
	"strings"
	"time"

	astisub "github.com/asticode/go-astisub"
	lingua "github.com/pemistahl/lingua-go"
)

// Text represents a text line given from the source
type Text struct {
	Text        string
	Translation string
}

// ParseSubtitles parses the subtitles to return an array of Text
func ParseSubtitles(sourceFile, translationFile string) ([]Text, error) {
	sourceSub, err := astisub.OpenFile(sourceFile)
	if err != nil {
		return nil, err
	}
	translationSub, err := astisub.OpenFile(translationFile)
	if err != nil {
		return nil, err
	}

	texts := make([]Text, 0, len(sourceSub.Items))
	translationIndex := 0
	for i, item := range sourceSub.Items {
		var nextItem *astisub.Item
		if i+1 < len(sourceSub.Items) {
			nextItem = sourceSub.Items[i+1]
		}

		var translation []string
		for ; translationIndex < len(translationSub.Items); translationIndex++ {
			translationItem := translationSub.Items[translationIndex]
			if !isRelatedToItem(item, nextItem, translationItem) {
				break
			}
			translation = append(translation, itemString(translationItem))
		}

		texts = append(texts, Text{
			Text:        itemString(item),
			Translation: strings.Join(translation, " "),
		})
	}
	return texts, nil
}

func isRelatedToItem(item, nextItem, translationItem *astisub.Item) bool {
	if nextItem == nil {
		return true
	}
	iOverlap, nextOverlap := itemOverlap(item, translationItem), itemOverlap(nextItem, translationItem)
	if iOverlap != nextOverlap {
		return iOverlap > nextOverlap
	}
	// equal non-zero overlap
	if iOverlap != 0 {
		return true
	}
	// no overlap, so calculate distance
	return translationItem.StartAt-item.EndAt < nextItem.StartAt-translationItem.EndAt
}

func itemOverlap(a, b *astisub.Item) time.Duration {
	if b.StartAt < a.StartAt {
		b, a = a, b
	}
	if a.EndAt < b.StartAt {
		return 0
	}
	return a.EndAt - b.StartAt
}

func itemString(i *astisub.Item) string {
	os := make([]string, len(i.Lines))
	for i, l := range i.Lines {
		os[i] = strings.TrimSpace(l.String())
	}
	return strings.Join(os, " ")
}

// Parser parses text into Text arrays (text and translation)
type Parser struct {
	SourceLanguage      Language
	TranslationLanguage Language
}

// NewParser returns a new parser
func NewParser(sourceLanguage, translationLanguage Language) *Parser {
	return &Parser{SourceLanguage: sourceLanguage, TranslationLanguage: translationLanguage}
}

// Language represents a language to parse translations with
type Language int

// Language values
const (
	Afrikaans Language = iota
	Albanian
	Arabic
	Armenian
	Azerbaijani
	Basque
	Belarusian
	Bengali
	Bokmal
	Bosnian
	Bulgarian
	Catalan
	Chinese
	Croatian
	Czech
	Danish
	Dutch
	English
	Esperanto
	Estonian
	Finnish
	French
	Ganda
	Georgian
	German
	Greek
	Gujarati
	Hebrew
	Hindi
	Hungarian
	Icelandic
	Indonesian
	Irish
	Italian
	Japanese
	Kazakh
	Korean
	Latin
	Latvian
	Lithuanian
	Macedonian
	Malay
	Maori
	Marathi
	Mongolian
	Nynorsk
	Persian
	Polish
	Portuguese
	Punjabi
	Romanian
	Russian
	Serbian
	Shona
	Slovak
	Slovene
	Somali
	Sotho
	Spanish
	Swahili
	Swedish
	Tagalog
	Tamil
	Telugu
	Thai
	Tsonga
	Tswana
	Turkish
	Ukrainian
	Urdu
	Vietnamese
	Welsh
	Xhosa
	Yoruba
	Zulu
	Unknown
)

type parseMode int

const (
	noTranslation parseMode = iota
	splitTranslation
	weaveTranslation
)

var errExtraTextLine = fmt.Errorf("there are more text lines than translation lines")
var errExtraTranslationLine = fmt.Errorf("there are more translation lines than text lines")

//nolint:gocognit // too many states for simplification at O(n) time
// TextsFromString returns an array of Text from the given string
func (p *Parser) TextsFromString(s string) ([]Text, error) {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.Language(p.SourceLanguage), lingua.Language(p.TranslationLanguage)).
		Build()

	texts := make([]Text, 0, len(lines))
	nonEmptyIndex := 0
	mode := noTranslation
	var text Text
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if mode == noTranslation {
			if line == "===" {
				if nonEmptyIndex == 1 {
					texts = append(texts, text)
				}
				nonEmptyIndex = 0
				mode = splitTranslation
				continue
			}
			if nonEmptyIndex == 1 {
				language, _ := detector.DetectLanguageOf(line)
				if Language(language) == p.TranslationLanguage {
					mode = weaveTranslation
				} else {
					texts = append(texts, text)
				}
			}
		}

		switch mode {
		case weaveTranslation:
			if nonEmptyIndex%2 == 0 {
				text = Text{Text: line}
			} else {
				text.Translation = line
				texts = append(texts, text)
			}
		case noTranslation:
			text = Text{Text: line}
			if nonEmptyIndex != 0 {
				texts = append(texts, text)
			}
		case splitTranslation:
			if len(texts) <= nonEmptyIndex {
				return texts, errExtraTranslationLine
			}
			texts[nonEmptyIndex].Translation = line
		}
		nonEmptyIndex++
	}
	if mode == splitTranslation && len(texts) != nonEmptyIndex {
		return texts, errExtraTextLine
	}

	return texts, nil
}
