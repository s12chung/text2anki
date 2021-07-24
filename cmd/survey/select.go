package survey

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

/*
Select is a prompt that presents a list of various options to the user
for them to select using the arrow keys and enter. Response type is a string.

	color := ""
	prompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{"red", "blue", "green"},
	}
	survey.AskOne(prompt, &color)
*/
type Select struct {
	Renderer
	Message           string
	Options           []string
	Default           interface{}
	PageSize          int
	KeyPressMap       map[rune]string
	selectedIndex     int
	useDefault        bool
	resultingKeyPress rune
}

// SelectTemplateData is the data available to the templates when processing
type SelectTemplateData struct {
	Select
	PageEntries   []core.OptionAnswer
	SelectedIndex int
	Answer        string
	ShowAnswer    bool
	Config        *PromptConfig
	KeyPressText  string
}

var SelectQuestionTemplate = `
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- "  "}}{{- color "cyan"}}[Arrows=Move{{ .KeyPressText }}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex }}{{color $.Config.Icons.SelectFocus.Format }}{{ $.Config.Icons.SelectFocus.Text }} {{else}}{{color "default"}}  {{end}}
    {{- $choice.Value}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (s *Select) OnChange(key rune, config *PromptConfig) bool {
	options := s.filterOptions(config)

	// if the user pressed the enter key and the index is a valid option
	if key == terminal.KeyEnter || key == '\n' {
		// if the selected index is a valid option
		if len(options) > 0 && s.selectedIndex < len(options) {

			// we're done (stop prompting the user)
			return true
		}

		// we're not done (keep prompting)
		return false

		// if the user pressed the up arrow or 'k' to emulate vim
	} else if key == terminal.KeyArrowUp && len(options) > 0 {
		s.useDefault = false

		// if we are at the top of the list
		if s.selectedIndex == 0 {
			// start from the button
			s.selectedIndex = len(options) - 1
		} else {
			// otherwise we are not at the top of the list so decrement the selected index
			s.selectedIndex--
		}

		// if the user pressed down or 'j' to emulate vim
	} else if (key == terminal.KeyTab || key == terminal.KeyArrowDown) && len(options) > 0 {
		s.useDefault = false
		// if we are at the bottom of the list
		if s.selectedIndex == len(options)-1 {
			// start from the top
			s.selectedIndex = 0
		} else {
			// increment the selected index
			s.selectedIndex++
		}
		// only show the help message if we have one
	} else if _, exists := s.KeyPressMap[unicode.ToUpper(key)]; exists {
		s.resultingKeyPress = unicode.ToUpper(key)
		return true
	}

	// figure out the options and index to render
	// figure out the page size
	pageSize := s.PageSize
	// if we dont have a specific one
	if pageSize == 0 {
		// grab the global value
		pageSize = config.PageSize
	}

	// TODO if we have started filtering and were looking at the end of a list
	// and we have modified the filter then we should move the page back!
	opts, idx := paginate(pageSize, options, s.selectedIndex)

	// render the options
	s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			SelectedIndex: idx,
			PageEntries:   opts,
			Config:        config,
			KeyPressText:  s.keyPressText(),
		},
	)

	// keep prompting
	return false
}

var RuneToKeyString = map[rune]string{
	terminal.KeyEscape:     "ESC",
	terminal.KeyArrowLeft:  "←",
	terminal.KeyArrowRight: "→",
}

func (s *Select) keyPressText() string {
	keyPressText := make([]string, len(s.KeyPressMap))
	i := 0
	for r, help := range s.KeyPressMap {
		k := string(r)
		if _, exists := RuneToKeyString[r]; exists {
			k = RuneToKeyString[r]
		}
		keyPressText[i] = fmt.Sprintf("%v=%v", k, help)
		i++
	}
	sort.Strings(keyPressText)
	return ", " + strings.Join(keyPressText, ", ")
}

func (s *Select) filterOptions(config *PromptConfig) []core.OptionAnswer {
	return core.OptionAnswerList(s.Options)
}

func (s *Select) Prompt(config *PromptConfig) (interface{}, error) {
	// if there are no options to render
	if len(s.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// start off with the first option selected
	sel := 0
	// if there is a default
	if s.Default != "" {
		// find the choice
		for i, opt := range s.Options {
			// if the option corresponds to the default
			if opt == s.Default {
				// we found our initial value
				sel = i
				// stop looking
				break
			}
		}
	}
	// save the selected index
	s.selectedIndex = sel

	// figure out the page size
	pageSize := s.PageSize
	// if we dont have a specific one
	if pageSize == 0 {
		// grab the global value
		pageSize = config.PageSize
	}

	// figure out the options and index to render
	opts, idx := paginate(pageSize, core.OptionAnswerList(s.Options), sel)

	// ask the question
	err := s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			PageEntries:   opts,
			SelectedIndex: idx,
			Config:        config,
			KeyPressText:  s.keyPressText(),
		},
	)
	if err != nil {
		return "", err
	}

	// by default, use the default value
	s.useDefault = true

	rr := s.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	cursor := s.NewCursor()
	cursor.Hide()       // hide the cursor
	defer cursor.Show() // show the cursor when we're done

	// start waiting for input
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		if s.OnChange(r, config) {
			break
		}
	}
	options := s.filterOptions(config)

	// the index to report
	var val string
	// if we are supposed to use the default value
	if s.useDefault || s.selectedIndex >= len(options) {
		// if there is a default value
		if s.Default != nil {
			// if the default is a string
			if defaultString, ok := s.Default.(string); ok {
				// use the default value
				val = defaultString
				// the default value could also be an interpret which is interpretted as the index
			} else if defaultIndex, ok := s.Default.(int); ok {
				val = s.Options[defaultIndex]
			} else {
				return val, errors.New("default value of select must be an int or string")
			}
		} else if len(options) > 0 {
			// there is no default value so use the first
			val = options[0].Value
		}
		// otherwise the selected index points to the value
	} else if s.selectedIndex < len(options) {
		// the
		val = options[s.selectedIndex].Value
	}

	// now that we have the value lets go hunt down the right index to return
	idx = -1
	for i, optionValue := range s.Options {
		if optionValue == val {
			idx = i
		}
	}

	if s.resultingKeyPress != rune(0) {
		err = &KeyPressError{Key: s.resultingKeyPress}
	}
	return core.OptionAnswer{Value: val, Index: idx}, err
}

type KeyPressError struct {
	Key rune
}

func (k *KeyPressError) Error() string {
	return fmt.Sprintf("user pressed %v, select exited", k.Key)
}

func IsKeyPressError(err error) bool {
	_, ok := err.(*KeyPressError)
	return ok
}

func KeyFromKeyPressError(err error) rune {
	if !IsKeyPressError(err) {
		return rune(0)
	}
	e := err.(*KeyPressError)
	return e.Key
}

func (s *Select) Cleanup(config *PromptConfig, val interface{}) error {
	return s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:       *s,
			Answer:       val.(core.OptionAnswer).Value,
			ShowAnswer:   true,
			Config:       config,
			KeyPressText: s.keyPressText(),
		},
	)
}
