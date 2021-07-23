// Package prompt is a temporary prompt CLI interface
package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/s12chung/text2anki/cmd/survey"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/tokenizers"
)

// Revolve revolves through a set of tokenizedTexts to create a set of notes
func Revolve(tokenizedTexts []app.TokenizedText, dict dictionary.Dicionary) ([]anki.Note, error) {
	return (&prompt{
		tokenizedTexts: tokenizedTexts,
		dictionary:     dict,
	}).revolve()
}

type prompt struct {
	tokenizedTexts []app.TokenizedText
	dictionary     dictionary.Dicionary

	tokenizedTextIndex int
	tokenIndex         int

	notes []anki.Note
}

func (p *prompt) revolve() ([]anki.Note, error) {
	for {
		if err := p.promptCurrent(); err != nil {
			return nil, fmt.Errorf("error showing next token: %w", err)
		}
		if !p.increment() {
			break
		}
		time.Sleep(1000)
	}

	return p.notes, nil
}

var ignorePOS = map[lang.PartOfSpeech]bool{
	lang.PartOfSpeechPunctuation: true,
	lang.PartOfSpeechOther:       true,
	lang.PartOfSpeechUnknown:     true,
}

func (p *prompt) promptCurrent() error {
	if ignorePOS[p.currentToken().PartOfSpeech] {
		return nil
	}

	terms, err := p.dictionary.Search(p.currentToken().Text)
	if err != nil {
		return err
	}
	if len(terms) == 0 {
		fmt.Printf("Skipping %v (%v), due to no search results\n", p.currentLabel(), p.currentToken().PartOfSpeech)
		return nil
	}

	var termIndex int
	var keyPress string
	prompt := &survey.Select{
		Message:  p.currentLabel(),
		Options:  itemStringsFromTerms(terms),
		PageSize: 10,
		KeyPressMap: map[rune]string{
			terminal.KeyEscape: "Back to Select Token",
		},
	}
	err = survey.AskOne(prompt, &termIndex)
	if survey.IsKeyPressError(err) {
		key := survey.KeyFromKeyPressError(err)
		keyPress = string(key)
		if _, exists := survey.RuneToKeyString[key]; exists {
			keyPress = survey.RuneToKeyString[key]
		}
	} else if err != nil {
		return err
	}

	if keyPress == "" {
		p.notes = append(p.notes, app.NewNoteFromTerm(terms[termIndex], 0))
	} else {
		fmt.Printf("PRESSED %v\n", keyPress)
	}
	return nil
}

const translationMaxLen = 5

func itemStringsFromTerms(terms []dictionary.Term) []string {
	itemStrings := make([]string, len(terms))
	for i, term := range terms {
		translationTextsMap := map[string]bool{}
		translationTextsA := make([]string, 0, translationMaxLen)
		for j, translation := range term.Translations {
			if j == 0 {
				continue
			}
			text := strings.TrimSpace(translation.Text)
			if text == "" || translationTextsMap[text] {
				continue
			}
			translationTextsMap[text] = true
			translationTextsA = append(translationTextsA, text)
			if len(translationTextsA) == translationMaxLen {
				break
			}
		}

		itemStrings[i] = fmt.Sprintf("%v %v %v: %v - %v\n        %v",
			term.Text,
			term.PartOfSpeech,
			strings.Repeat("*", int(term.CommonLevel)),
			term.Translations[0].Text,
			term.Translations[0].Explanation,
			strings.Join(translationTextsA, "; "),
		)
	}
	return itemStrings
}

func (p *prompt) currentTokenizedText() app.TokenizedText {
	return p.tokenizedTexts[p.tokenizedTextIndex]
}

func (p *prompt) currentToken() tokenizers.Token {
	return p.currentTokenizedText().Tokens[p.tokenIndex]
}

func (p *prompt) currentLabel() string {
	label := []rune(p.currentTokenizedText().Text.Text)

	endIndex := p.currentToken().EndIndex
	label = append(label[:endIndex+1], label[endIndex:]...)
	label[endIndex] = ']'

	beginIndex := p.currentToken().StartIndex
	label = append(label[:beginIndex+1], label[beginIndex:]...)
	label[beginIndex] = '['

	return string(label)
}

func (p *prompt) increment() bool {
	if p.tokenizedTextIndex < len(p.tokenizedTexts) &&
		p.tokenIndex+1 < len(p.currentTokenizedText().Tokens) {
		p.tokenIndex++
		return true
	}
	if p.tokenizedTextIndex+1 < len(p.tokenizedTexts) {
		p.tokenizedTextIndex++
		p.tokenIndex = 0
		return true
	}
	return false
}
