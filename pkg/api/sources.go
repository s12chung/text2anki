package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/util/chiutil"
	"github.com/s12chung/text2anki/pkg/util/httputil"
)

const contextSource httputil.ContextKey = "source"

// SourceCtx sets the source context from the sourceID
func SourceCtx(r *http.Request) (*http.Request, int, error) {
	sourceID, err := chiutil.ParamID(r, "sourceID")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	source, err := db.Qs().SourceGet(r.Context(), sourceID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	r = r.WithContext(context.WithValue(r.Context(), contextSource, source.ToSourceSerialized()))
	return r, 0, nil
}

// SourceGet gets the source
func (rs Routes) SourceGet(r *http.Request) (any, int, error) {
	sourceSerialized, ok := r.Context().Value(contextSource).(db.SourceSerialized)
	if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("cast to db.Source fail")
	}
	return sourceSerialized, 0, nil
}

// SourcePostRequest represents the SourcePost request
type SourcePostRequest struct {
	Text        string
	Translation string
}

// TextsString returns the string for TokenizeTextsFromString
func (s *SourcePostRequest) TextsString() string {
	if s.Translation == "" {
		return s.Text
	}
	return s.Text + "\n\n" + text.SplitDelimiter + "\n\n" + s.Translation
}

// SourcePost creates a new source
func (rs Routes) SourcePost(r *http.Request) (any, int, error) {
	req := SourcePostRequest{}
	if code, err := httputil.BindJSON(r, &req); err != nil {
		return nil, code, err
	}

	tokenizedTexts, err := rs.TextTokenizer.TokenizeTextsFromString(req.TextsString())
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	sourceSerialized, err := db.Qs().SourceSerializedCreate(r.Context(), tokenizedTexts)
	return sourceSerialized, http.StatusInternalServerError, err
}
