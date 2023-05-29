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
	ID             int64           `json:"id,omitempty"`
	Name           string          `json:"name"`
	TokenizedTexts []TokenizedText `json:"tokenized_texts,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
}

// DefaultedName returns the Default name
func (s SourceSerialized) DefaultedName() string {
	if s.Name != "" {
		return s.Name
	}
	if len(s.TokenizedTexts) == 0 {
		return ""
	}
	return stringutil.FirstUnbrokenSubstring(s.TokenizedTexts[0].Text.Text, 25)
}

// StaticCopy returns a copy with fields that variate
func (s SourceSerialized) StaticCopy() any {
	c := s
	c.ID = 0
	c.UpdatedAt = time.Time{}
	c.CreatedAt = time.Time{}
	return c
}

// ToSource returns the Source of the SourceSerialized
func (s SourceSerialized) ToSource() Source {
	bytes, err := json.Marshal(s.TokenizedTexts)
	if err != nil {
		slog.Error(err.Error())
		panic(-1)
	}
	return Source{
		ID:             s.ID,
		Name:           s.Name,
		TokenizedTexts: string(bytes),
		UpdatedAt:      s.UpdatedAt,
		CreatedAt:      s.CreatedAt,
	}
}

// ToSourceSerialized returns the SourceSerialized of the Source
func (s Source) ToSourceSerialized() SourceSerialized {
	var tokenizedTexts []TokenizedText
	if err := json.Unmarshal([]byte(s.TokenizedTexts), &tokenizedTexts); err != nil {
		slog.Error(err.Error())
		panic(-1)
	}
	return SourceSerialized{
		ID:             s.ID,
		Name:           s.Name,
		TokenizedTexts: tokenizedTexts,
		UpdatedAt:      s.UpdatedAt,
		CreatedAt:      s.CreatedAt,
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
	Tokens []tokenizers.Token `json:"tokens,omitempty"`
}

// Setup sets up the TextTokenizer
func (t TextTokenizer) Setup() error {
	return t.Tokenizer.Setup()
}

// Cleanup cleans up the TextTokenizer
func (t TextTokenizer) Cleanup() error {
	return t.Tokenizer.Cleanup()
}

// TokenizeTextsFromString converts a string to TokenizedText
func (t TextTokenizer) TokenizeTextsFromString(s string) ([]TokenizedText, error) {
	texts, err := t.Parser.TextsFromString(s)
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

// SourceSerializedList returns a SourceSerialized from the DB
func (q *Queries) SourceSerializedList(ctx context.Context) ([]SourceSerialized, error) {
	sources, err := q.SourceList(ctx)
	if err != nil {
		return nil, err
	}

	sourceSerializeds := make([]SourceSerialized, len(sources))
	for i, source := range sources {
		sourceSerializeds[i] = source.ToSourceSerialized()
	}
	return sourceSerializeds, nil
}

// SourceSerializedCreate creates a source in the DB
func (q *Queries) SourceSerializedCreate(ctx context.Context, sourceSerialized SourceSerialized) (SourceSerialized, error) {
	source, err := q.SourceCreate(ctx, SourceCreateParams{
		Name:           sourceSerialized.DefaultedName(),
		TokenizedTexts: sourceSerialized.ToSource().TokenizedTexts,
	})
	if err != nil {
		return SourceSerialized{}, err
	}
	return source.ToSourceSerialized(), nil
}
