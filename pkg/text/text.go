// Package text contains functions to separate text into source and translation
package text

import (
	"fmt"
	"strings"

	lingua "github.com/pemistahl/lingua-go"
)

// Text represents a text line given from the source
type Text struct {
	Text        string
	Translation string
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
	splittingTranslation
	weaveTranslation
)

// TextsFromString returns an array of Texts from the given string
func (p *Parser) TextsFromString(s string) []Text {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.Language(p.SourceLanguage), lingua.Language(p.TranslationLanguage)).
		Build()

	texts := make([]Text, 0, len(lines))
	nonEmptyIndex := 0
	mode := noTranslation
	var text Text
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		if mode == noTranslation {
			if l == "===" {
				nonEmptyIndex = 0
				mode = splittingTranslation
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
		case splittingTranslation:
			fmt.Println(nonEmptyIndex)
			fmt.Println(texts[nonEmptyIndex].Text)
			fmt.Println(line)
			texts[nonEmptyIndex].Translation = line
		}
		nonEmptyIndex++
	}
	return texts
}
