package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"time"

	"golang.org/x/exp/slog"

	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizer"
	"github.com/s12chung/text2anki/pkg/util/stringutil"
)

// SourceStructured is a copy of Source for with JSON columns structured
type SourceStructured struct {
	ID        int64        `json:"id,omitempty"`
	Name      string       `json:"name"`
	Reference string       `json:"reference"`
	Parts     []SourcePart `json:"parts"`
	UpdatedAt time.Time    `json:"updated_at"`
	CreatedAt time.Time    `json:"created_at"`
}

// PrepareSerialize prepares the model for Serializing for API endpoints
func (s SourceStructured) PrepareSerialize() {
	for _, part := range s.Parts {
		if part.Media == nil {
			continue
		}
		part.Media.toSerialize = true
	}
}

// SourcePart is a part of the Source that contains text
type SourcePart struct {
	Media          *SourcePartMedia `json:"media,omitempty"`
	TokenizedTexts []TokenizedText  `json:"tokenized_texts"`
}

// SourcePartMedia is the media of the SourcePart
type SourcePartMedia struct {
	toSerialize bool
	ImageKey    string `json:"image_key,omitempty"`
	AudioKey    string `json:"audio_key,omitempty"`
}

type sourcePartMediaAlias SourcePartMedia

// SourcePartMediaSerialized is the API endpoint version of SourcePartMedia
type SourcePartMediaSerialized struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}

func (s SourcePartMediaSerialized) toDB() (SourcePartMedia, error) {
	db := SourcePartMedia{}
	return db, dbStorage.KeyTreeFromSignGetTree(s, &db)
}

// SourcePartMediaFile is the File version of SourcePartMedia
type SourcePartMediaFile struct {
	ImageFile fs.File `json:"image_file,omitempty"`
	AudioFile fs.File `json:"audio_file,omitempty"`
}

// SerializedEmpty returns an empty model for Serializing for API endpoints
func (s *SourcePartMedia) SerializedEmpty() any {
	return SourcePartMediaSerialized{}
}

func (s *SourcePartMedia) toSerialized() (SourcePartMediaSerialized, error) {
	serialized := SourcePartMediaSerialized{}
	return serialized, dbStorage.SignGetTreeFromKeyTree(s, &serialized)
}

// MarshalJSON returns the JSON representation
func (s *SourcePartMedia) MarshalJSON() ([]byte, error) {
	if !s.toSerialize {
		return json.Marshal(sourcePartMediaAlias(*s))
	}
	serialized, err := s.toSerialized()
	if err != nil {
		return nil, err
	}
	return json.Marshal(serialized)
}

// UnmarshalJSON sets the data based on the data JSON
//
// Should be removed if tests are improved such that they're not relying on Unmarshalling requests
func (s *SourcePartMedia) UnmarshalJSON(data []byte) error {
	alias := sourcePartMediaAlias(*s)
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	if alias != (sourcePartMediaAlias{}) {
		*s = SourcePartMedia(alias)
		return nil
	}

	serialized := SourcePartMediaSerialized{}
	if err := json.Unmarshal(data, &serialized); err != nil {
		return err
	}
	db, err := serialized.toDB()
	if err != nil {
		return err
	}
	*s = db
	s.toSerialize = true
	return nil
}

// DefaultedName returns the Default name
func (s SourceStructured) DefaultedName() string {
	if s.Name != "" {
		return s.Name
	}
	if len(s.Parts) == 0 && len(s.Parts[0].TokenizedTexts) == 0 {
		return ""
	}
	return stringutil.FirstUnbrokenSubstring(s.Parts[0].TokenizedTexts[0].Text.Text, 30)
}

// StaticCopy returns a copy without fields that variate
func (s SourceStructured) StaticCopy() any {
	c := s
	c.ID = 0
	c.UpdatedAt = time.Time{}
	c.CreatedAt = time.Time{}
	return c
}

// UpdateParams returns the SourceUpdateParams for the SourceStructured
func (s SourceStructured) UpdateParams() SourceUpdateParams {
	return SourceUpdateParams{
		Name:      s.Name,
		Reference: s.Reference,
		ID:        s.ID,
	}
}

// CreateParams returns the SourceCreateParams for the SourceStructured
func (s SourceStructured) CreateParams() SourceCreateParams {
	return SourceCreateParams{
		Name:      s.DefaultedName(),
		Reference: s.Reference,
		Parts:     s.ToSource().Parts,
	}
}

// ToSource returns the Source of the SourceStructured
func (s SourceStructured) ToSource() Source {
	bytes, err := json.Marshal(s.Parts)
	if err != nil {
		slog.Error(err.Error())
		panic(-1)
	}
	return Source{
		ID:        s.ID,
		Name:      s.Name,
		Reference: s.Reference,
		Parts:     string(bytes),
		UpdatedAt: s.UpdatedAt,
		CreatedAt: s.CreatedAt,
	}
}

// ToSourceStructured returns the SourceStructured of the Source
func (s Source) ToSourceStructured() SourceStructured {
	var parts []SourcePart
	if err := json.Unmarshal([]byte(s.Parts), &parts); err != nil {
		slog.Error(err.Error())
		panic(-1)
	}
	return SourceStructured{
		ID:        s.ID,
		Name:      s.Name,
		Reference: s.Reference,
		Parts:     parts,
		UpdatedAt: s.UpdatedAt,
		CreatedAt: s.CreatedAt,
	}
}

// TextTokenizer is used to generate TokenizedText
type TextTokenizer struct {
	Parser       text.Parser
	Tokenizer    tokenizer.Tokenizer
	CleanSpeaker bool
}

// TokenizedText is the text grouped with its tokens
type TokenizedText struct {
	text.Text
	Tokens []tokenizer.Token `json:"tokens"`
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
		var tokens []tokenizer.Token
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

// SourceStructuredIndex returns a SourceStructured from the DB
func (q *Queries) SourceStructuredIndex(ctx context.Context) ([]SourceStructured, error) {
	sources, err := q.SourceIndex(ctx)
	if err != nil {
		return nil, err
	}

	sourceStructureds := make([]SourceStructured, len(sources))
	for i, source := range sources {
		sourceStructureds[i] = source.ToSourceStructured()
	}
	return sourceStructureds, nil
}
