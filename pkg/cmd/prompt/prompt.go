// Package prompt is a temporary prompt CLI interface
package prompt

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	yaml "gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/cmd/survey"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/text"
)

// CreateCards initializes the create card UI prompt
func CreateCards(tokenizedTexts []app.TokenizedText, dict dictionary.Dicionary) ([]anki.Note, error) {
	return (&createCards{
		tokenizedTexts: tokenizedTexts,
		dictionary:     dict,
	}).start()
}

type createCards struct {
	tokenizedTexts []app.TokenizedText
	dictionary     dictionary.Dicionary

	tokenizedTextIndex int

	notes []anki.Note
}

type transition string

const previousTransition transition = "Previous"
const nextTransition transition = "Next"
const finishedTransition transition = "Finished"
const noneTransition transition = "None"

const searchToken = 'S'

func (c *createCards) start() ([]anki.Note, error) {
	for {
		trans, err := c.showTokenizedText(c.tokenizedTexts[c.tokenizedTextIndex])
		if err != nil {
			return nil, err
		}
		switch trans {
		case previousTransition:
			c.tokenizedTextIndex = (c.tokenizedTextIndex + len(c.tokenizedTexts) - 1) % len(c.tokenizedTexts)
		case nextTransition:
			c.tokenizedTextIndex = (c.tokenizedTextIndex + 1) % len(c.tokenizedTexts)
		case finishedTransition:
			return c.notes, nil
		}
	}
}

func (c *createCards) showTokenizedText(tokenizedText app.TokenizedText) (transition, error) {
	for {
		context := tokenizedText.Text
		selectText := fmt.Sprintf("%v/%v: %v\n%v\n",
			c.tokenizedTextIndex+1, len(c.tokenizedTexts), context.Text, context.Translation)
		options, noValidTokens := tokenOptions(tokenizedText)

		var token string
		keyPress, err := showSelect(selectText, options, &token, map[rune]string{
			terminal.KeyEscape:     "Finish and Export",
			terminal.KeyArrowLeft:  "Prev Text",
			terminal.KeyArrowRight: "Next Text",
			searchToken:            "Search Dictionary",
		})
		if !survey.IsKeyPressError(err) && err != nil {
			return noneTransition, err
		}

		switch keyPress {
		case terminal.KeyEscape:
			return finishedTransition, nil
		case terminal.KeyArrowLeft:
			return previousTransition, nil
		case terminal.KeyArrowRight:
			return nextTransition, nil
		case searchToken:
			if err := c.showSearchInput(context); err != nil {
				return noneTransition, err
			}
		default:
			if noValidTokens {
				if err := c.showCreateNote(nil); err != nil {
					return noneTransition, err
				}
			} else {
				if err := c.showSearch(context, token); err != nil {
					return noneTransition, err
				}
			}
		}
	}
}

var ignorePOS = map[lang.PartOfSpeech]bool{
	lang.PartOfSpeechPunctuation: true,
	lang.PartOfSpeechOther:       true,
	lang.PartOfSpeechUnknown:     true,
}

func tokenOptions(tokenizedText app.TokenizedText) ([]string, bool) {
	options := make([]string, 0, len(tokenizedText.Tokens))
	for _, token := range tokenizedText.Tokens {
		if ignorePOS[token.PartOfSpeech] {
			continue
		}
		options = append(options, token.Text)
	}
	if len(options) == 0 {
		return []string{"No valid token, create from empty card"}, true
	}
	return options, false
}

func (c *createCards) showSearchInput(context text.Text) error {
	query := ""
	prompt := &survey.Input{
		Message: "Search Dictionary:",
	}
	if err := survey.AskOne(prompt, &query); err != nil {
		return err
	}
	return c.showSearch(context, query)
}

func (c *createCards) showSearch(context text.Text, query string) error {
	for {
		terms, err := c.dictionary.Search(query)
		if err != nil {
			return err
		}

		options, noSearchResults := itemStringsFromTerms(terms)
		var termIndex int
		selectText := context.Text + "\n" + context.Translation
		keyPress, err := showSelect(selectText, options, &termIndex, map[rune]string{
			terminal.KeyEscape: "Back to Select Token",
			searchToken:        "Search Dictionary",
		})
		if !survey.IsKeyPressError(err) && err != nil {
			return err
		}

		switch keyPress {
		case terminal.KeyEscape:
			return nil
		case searchToken:
			if err := c.showSearchInput(context); err != nil {
				return err
			}
		default:
			if noSearchResults {
				if err := c.showCreateNote(nil); err != nil {
					return err
				}
			} else if err := c.showCreateNote(&terms[termIndex]); err != nil {
				return err
			}
		}
	}
}

var posTypes = []lang.PartOfSpeech{
	lang.PartOfSpeechNoun,
	lang.PartOfSpeechPronoun,
	lang.PartOfSpeechNumeral,
	lang.PartOfSpeechPostposition,
	lang.PartOfSpeechVerb,
	lang.PartOfSpeechAdjective,
	lang.PartOfSpeechDeterminer,
	lang.PartOfSpeechAdverb,
	lang.PartOfSpeechInterjection,

	lang.PartOfSpeechAffix,
	lang.PartOfSpeechPrefix,
	lang.PartOfSpeechInfix,
	lang.PartOfSpeechSuffix,

	lang.PartOfSpeechDependentNoun,

	lang.PartOfSpeechAuxiliaryPredicate,
	lang.PartOfSpeechAuxiliaryVerb,
	lang.PartOfSpeechAuxiliaryAdjective,

	lang.PartOfSpeechEnding,
	lang.PartOfSpeechCopula,
	lang.PartOfSpeechPunctuation,

	lang.PartOfSpeechOther,
	lang.PartOfSpeechUnknown,
	lang.PartOfSpeechInvalid,
}

func (c *createCards) showCreateNote(term *dictionary.Term) error {
	filename, err := createNoteTempfile(term, c.tokenizedTexts[c.tokenizedTextIndex].Text)
	if err != nil {
		return err
	}
	if err = openEditor(filename); err != nil {
		return err
	}
	note, err := noteFromFile(filename)
	if err != nil {
		return err
	}
	c.notes = append(c.notes, *note)
	return nil
}

func createNoteTempfile(term *dictionary.Term, context text.Text) (s string, err error) {
	f, err := ioutil.TempFile("", "text2anki-showCreateNote-*.yaml")
	if err != nil {
		return "", err
	}
	defer func() {
		err2 := f.Close()
		if err == nil {
			err = err2
		}
	}()
	if err = addCreateNoteHeaders(f, term); err != nil {
		return "", err
	}
	if err = addNote(f, term, context); err != nil {
		return "", err
	}
	return f.Name(), err
}

func openEditor(filename string) error {
	//nolint:gosec // can't get around it for now
	cmd := exec.Command("subl", "-w", filename)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func noteFromFile(filename string) (*anki.Note, error) {
	//nolint:gosec // always writing to tempfile
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	note := &anki.Note{}
	if err = yaml.Unmarshal(bytes, note); err != nil {
		return nil, err
	}
	return note, err
}

func addCreateNoteHeaders(f io.Writer, term *dictionary.Term) error {
	if term != nil {
		bytes, err := yaml.Marshal(term)
		if err != nil {
			return err
		}
		termString := string(bytes)
		termString = "# " + strings.ReplaceAll(termString, "\n", "\n# ") + "\n"
		if _, err := f.Write([]byte(termString)); err != nil {
			return err
		}
	}

	a := make([]string, len(posTypes))
	for i, posType := range posTypes {
		a[i] = string(posType)
	}
	postTypesString := fmt.Sprintf("#\n# %v\n# \n# \n", strings.Join(a, ", "))
	if _, err := f.Write([]byte(postTypesString)); err != nil {
		return err
	}
	return nil
}

func addNote(f io.Writer, term *dictionary.Term, context text.Text) error {
	note := anki.Note{}
	if term != nil {
		note = app.NewNoteFromTerm(*term, 0)
	}
	note.Usage = context.Text
	note.UsageTranslation = context.Translation
	bytes, err := yaml.Marshal(&note)
	if err != nil {
		return err
	}
	if _, err := f.Write(bytes); err != nil {
		return err
	}
	return nil
}

const translationMaxLen = 5

func itemStringsFromTerms(terms []dictionary.Term) ([]string, bool) {
	if len(terms) == 0 {
		return []string{"No Search Results, create from empty card"}, true
	}

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
	return itemStrings, false
}

func showSelect(message string, options []string, resp interface{}, keyPressMap map[rune]string) (rune, error) {
	sel := &survey.Select{
		Message:     message,
		Options:     options,
		PageSize:    10,
		KeyPressMap: keyPressMap,
	}
	err := survey.AskOne(sel, resp)
	return survey.KeyFromKeyPressError(err), err
}
