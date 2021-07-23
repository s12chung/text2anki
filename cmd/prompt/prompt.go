// Package prompt is a temporary prompt CLI interface
package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/dictionary"
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

// TODO: clean up POS types from tokenizer
var ignorePOS = map[string]bool{
	"SF": true,
	"SP": true,
	"SL": true,
}

func (p *prompt) promptCurrent() error {
	if ignorePOS[p.currentToken().POS] {
		return nil
	}

	terms, err := p.dictionary.Search(p.currentToken().Morph)
	if err != nil {
		return err
	}
	if len(terms) == 0 {
		fmt.Printf("Skipping %v (%v), due to no search results\n", p.currentLabel(), p.currentToken().POS)
		return nil
	}
	prompt := promptui.Select{
		Label:  p.currentLabel(),
		Items:  itemStringsFromTerms(terms),
		Size:   10,
		Stdout: &bellSkipper{},
	}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	p.notes = append(p.notes, app.NewNoteFromTerm(terms[i], 0))
	return nil
}

const translationMaxLen = 5

func itemStringsFromTerms(terms []dictionary.Term) []string {
	itemStrings := make([]string, len(terms))
	for i, term := range terms {
		translationTextsMap := map[string]bool{}
		translationTextsA := make([]string, 0, translationMaxLen)
		for _, translation := range term.Translations {
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

		itemStrings[i] = fmt.Sprintf("%v %v %v - %v",
			term.Text,
			term.PartOfSpeech,
			strings.Repeat("*", int(term.CommonLevel)),
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

	beginIndex := p.currentToken().BeginIndex
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
