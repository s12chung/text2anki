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

// SourceList returns a list of sources
func (rs Routes) SourceList(r *http.Request) (any, int, error) {
	sourceSerializeds, err := db.Qs().SourceSerializedList(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return sourceSerializeds, 0, nil
}

// SourceGet gets the source
func (rs Routes) SourceGet(r *http.Request) (any, int, error) {
	sourceSerialized, code, err := ctxSourceSerialized(r)
	if err != nil {
		return nil, code, err
	}
	return sourceSerialized, 0, nil
}

// SourcePostRequest represents the SourceCreate request
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

// SourceCreate creates a new source
func (rs Routes) SourceCreate(r *http.Request) (any, int, error) {
	req := SourcePostRequest{}
	if code, err := httputil.BindJSON(r, &req); err != nil {
		return nil, code, err
	}

	tokenizedTexts, err := rs.TextTokenizer.TokenizeTextsFromString(req.TextsString())
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	sourceSerialized, err := db.Qs().SourceSerializedCreate(r.Context(), db.SourceSerialized{
		TokenizedTexts: tokenizedTexts},
	)
	if err != nil {
		return sourceSerialized, http.StatusInternalServerError, err
	}
	return sourceSerialized, 0, nil
}

// SourceDestroy destroys the source
func (rs Routes) SourceDestroy(r *http.Request) (any, int, error) {
	sourceSerialized, code, err := ctxSourceSerialized(r)
	if err != nil {
		return nil, code, err
	}
	if err := db.Qs().SourceDestroy(r.Context(), sourceSerialized.ID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return sourceSerialized, 0, nil
}

func ctxSourceSerialized(r *http.Request) (db.SourceSerialized, int, error) {
	sourceSerialized, ok := r.Context().Value(contextSource).(db.SourceSerialized)
	if !ok {
		return db.SourceSerialized{}, http.StatusInternalServerError, fmt.Errorf("cast to db.SourceSerialized fail")
	}
	return sourceSerialized, 0, nil
}
