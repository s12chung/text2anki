package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/util/stringutil"
)

// SourceSerialized is a copy of Source for Serializing
type SourceSerialized struct {
	ID        int64        `json:"id,omitempty"`
	Name      string       `json:"name"`
	Parts     []SourcePart `json:"parts"`
	UpdatedAt time.Time    `json:"updated_at"`
	CreatedAt time.Time    `json:"created_at"`
}

// SourcePart is a part of the Source that contains text
type SourcePart struct {
	TokenizedTexts []TokenizedText `json:"tokenized_texts"`
}

// DefaultedName returns the Default name
func (s SourceSerialized) DefaultedName() string {
	if s.Name != "" {
		return s.Name
	}
	if len(s.Parts) == 0 && len(s.Parts[0].TokenizedTexts) == 0 {
		return ""
	}
	return stringutil.FirstUnbrokenSubstring(s.Parts[0].TokenizedTexts[0].Text.Text, 25)
}

// StaticCopy returns a copy without fields that variate
func (s SourceSerialized) StaticCopy() any {
	c := s
	c.ID = 0
	c.UpdatedAt = time.Time{}
	c.CreatedAt = time.Time{}
	return c
}

// UpdateParams returns the SourceUpdateParams for the SourceSerialized
func (s SourceSerialized) UpdateParams() SourceUpdateParams {
	return SourceUpdateParams{
		Name: s.Name,
		ID:   s.ID,
	}
}

// CreateParams returns the SourceCreateParams for the SourceSerialized
func (s SourceSerialized) CreateParams() SourceCreateParams {
	return SourceCreateParams{
		Name:  s.DefaultedName(),
		Parts: s.ToSource().Parts,
	}
}

// ToSource returns the Source of the SourceSerialized
func (s SourceSerialized) ToSource() Source {
	bytes, err := json.Marshal(s.Parts)
	if err != nil {
		slog.Error(err.Error())
		panic(-1)
	}
	return Source{
		ID:        s.ID,
		Name:      s.Name,
		Parts:     string(bytes),
		UpdatedAt: s.UpdatedAt,
		CreatedAt: s.CreatedAt,
	}
}

// ToSourceSerialized returns the SourceSerialized of the Source
func (s Source) ToSourceSerialized() SourceSerialized {
	var parts []SourcePart
	if err := json.Unmarshal([]byte(s.Parts), &parts); err != nil {
		slog.Error(err.Error())
		panic(-1)
	}
	return SourceSerialized{
		ID:        s.ID,
		Name:      s.Name,
		Parts:     parts,
		UpdatedAt: s.UpdatedAt,
		CreatedAt: s.CreatedAt,
	}
}

// TextTokenizer is used to generate TokenizedText
type TextTokenizer struct {
	Parser       text.Parser
	Tokenizer    tokenizers.Tokenizer
	CleanSpeaker bool
}

// TokenizedText is the text grouped with its tokens
type TokenizedText struct {
	text.Text
	Tokens []tokenizers.Token `json:"tokens"`
}

// Setup sets up the TextTokenizer
func (t TextTokenizer) Setup() error {
	return t.Tokenizer.Setup()
}

// Cleanup cleans up the TextTokenizer
func (t TextTokenizer) Cleanup() error {
	return t.Tokenizer.Cleanup()
}

// TokenizedTexts converts a string to TokenizedText
func (t TextTokenizer) TokenizedTexts(s, translation string) ([]TokenizedText, error) {
	texts, err := t.Parser.Texts(s, translation)
	if err != nil {
		return nil, err
	}
	if t.CleanSpeaker {
		texts = text.CleanSpeaker(texts)
	}
	return t.TokenizeTexts(texts)
}

// TokenizeTexts takes the texts and tokenizes them
func (t TextTokenizer) TokenizeTexts(texts []text.Text) (tokenizedTexts []TokenizedText, err error) {
	if !t.Tokenizer.IsSetup() {
		return nil, fmt.Errorf("TextTokenizer not set up")
	}

	tokenizedTexts = make([]TokenizedText, len(texts))
	for i, text := range texts {
		var tokens []tokenizers.Token
		tokens, err = t.Tokenizer.Tokenize(text.Text)
		if err != nil {
			return nil, err
		}
		tokenizedTexts[i] = TokenizedText{
			Text:   text,
			Tokens: tokens,
		}
	}

	return tokenizedTexts, nil
}

// SourceSerializedIndex returns a SourceSerialized from the DB
func (q *Queries) SourceSerializedIndex(ctx context.Context) ([]SourceSerialized, error) {
	sources, err := q.SourceIndex(ctx)
	if err != nil {
		return nil, err
	}

	sourceSerializeds := make([]SourceSerialized, len(sources))
	for i, source := range sources {
		sourceSerializeds[i] = source.ToSourceSerialized()
	}
	return sourceSerializeds, nil
}
