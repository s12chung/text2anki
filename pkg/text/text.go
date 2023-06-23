// Package text contains functions to separate text into source and translation
package text

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pemistahl/lingua-go"
)

// Text represents a text line given from the source
type Text struct {
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

// Parser parses text into Text arrays (text and translation)
type Parser struct {
	SourceLanguage      Language
	TranslationLanguage Language
}

// NewParser returns a new parser
func NewParser(sourceLanguage, translationLanguage Language) Parser {
	return Parser{SourceLanguage: sourceLanguage, TranslationLanguage: translationLanguage}
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

var errExtraTextLine = fmt.Errorf("there are more text lines than translation lines")
var errExtraTranslationLine = fmt.Errorf("there are more translation lines than text lines")

// Texts returns an array of Text from the given string
func (p Parser) Texts(s, translation string) ([]Text, error) {
	if translation == "" {
		return p.TextsFromString(s)
	}
	return p.TextsFromTranslation(s, translation)
}

// TextsFromString generates a []Text from a string, can have no translations or weaved
func (p Parser) TextsFromString(s string) ([]Text, error) {
	lines := split(s)
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.Language(p.SourceLanguage), lingua.Language(p.TranslationLanguage)).
		Build()

	texts := make([]Text, 0, len(lines))
	var text Text
	for _, line := range lines {
		language, _ := detector.DetectLanguageOf(line)
		if Language(language) == p.SourceLanguage {
			if text.Text != "" {
				texts = append(texts, text)
				text = Text{}
			}
			text = Text{Text: line}
		} else {
			if text.Text == "" {
				return nil, fmt.Errorf("translation exists for two consecutive non-empty lines: %v", line)
			}
			text.Translation = line
			texts = append(texts, text)
			text = Text{}
		}
	}
	if text.Text != "" {
		texts = append(texts, text)
	}
	return texts, nil
}

// TextsFromTranslation generates a Text[] from two strings (text and translation) that have the same number of lines
func (p Parser) TextsFromTranslation(s, translation string) ([]Text, error) {
	lines := split(s)
	translations := split(translation)
	if len(lines) > len(translations) {
		return nil, errExtraTextLine
	}
	if len(lines) < len(translations) {
		return nil, errExtraTranslationLine
	}

	texts := make([]Text, len(lines))
	for i, line := range lines {
		texts[i] = Text{
			Text:        line,
			Translation: translations[i],
		}
	}
	return texts, nil
}

func split(s string) []string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	clean := make([]string, len(lines))

	i := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		clean[i] = line
		i++
	}
	return clean[:i]
}

// CleanSpeaker removes the CleanSpeakerString names from the next
func CleanSpeaker(texts []Text) []Text {
	cleanedTexts := make([]Text, len(texts))
	for i, t := range texts {
		cleanedTexts[i] = Text{
			Text:        CleanSpeakerString(t.Text),
			Translation: CleanSpeakerString(t.Translation),
		}
	}
	return cleanedTexts
}

var speakerRegex = regexp.MustCompile(`\A[^:\d]{0,25}:`)

// CleanSpeakerString cleans the speaker from the string
func CleanSpeakerString(s string) string {
	s = speakerRegex.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}
