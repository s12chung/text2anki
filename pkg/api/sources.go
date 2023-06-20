package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
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

// SourceIndex returns a list of sources
func (rs Routes) SourceIndex(r *http.Request) (any, int, error) {
	return httputil.ReturnModelOr500(func() (any, error) {
		return db.Qs().SourceSerializedIndex(r.Context())
	})
}

// SourceGet gets the source
func (rs Routes) SourceGet(r *http.Request) (any, int, error) {
	return ctxSourceSerialized(r)
}

// SourceUpdateRequest represents the SourceUpdate request
type SourceUpdateRequest struct {
	Name string
}

func init() {
	firm.RegisterType(firm.NewDefinition(SourceUpdateRequest{}).Validates(firm.RuleMap{
		"Name": {rule.Presence{}},
	}))
}

// SourceUpdate updates the source
func (rs Routes) SourceUpdate(r *http.Request) (any, int, error) {
	sourceSerialized, code, err := ctxSourceSerialized(r)
	if err != nil {
		return nil, code, err
	}

	req := SourceUpdateRequest{}
	if code, err = bindAndValidate(r, &req); err != nil {
		return nil, code, err
	}
	sourceSerialized.Name = req.Name

	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := db.Qs().SourceUpdate(r.Context(), sourceSerialized.ToSourceUpdateParams())
		return source.ToSourceSerialized(), err
	})
}

// SourceCreateRequest represents the SourceCreate request
type SourceCreateRequest struct {
	Text        string
	Translation string
}

func init() {
	firm.RegisterType(firm.NewDefinition(SourceCreateRequest{}).Validates(firm.RuleMap{
		"Text": {rule.Presence{}},
	}))
}

// TextsString returns the string for TokenizeTextsFromString
func (s *SourceCreateRequest) TextsString() string {
	if s.Translation == "" {
		return s.Text
	}
	return s.Text + "\n\n" + text.SplitDelimiter + "\n\n" + s.Translation
}

// SourceCreate creates a new source
func (rs Routes) SourceCreate(r *http.Request) (any, int, error) {
	req := SourceCreateRequest{}
	if code, err := bindAndValidate(r, &req); err != nil {
		return nil, code, err
	}

	tokenizedTexts, err := rs.TextTokenizer.TokenizeTextsFromString(req.TextsString())
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := db.Qs().SourceCreate(r.Context(), db.SourceSerialized{TokenizedTexts: tokenizedTexts}.ToSourceCreateParams())
		return source.ToSourceSerialized(), err
	})
}

// SourceDestroy destroys the source
func (rs Routes) SourceDestroy(r *http.Request) (any, int, error) {
	sourceSerialized, code, err := ctxSourceSerialized(r)
	if err != nil {
		return nil, code, err
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return sourceSerialized, db.Qs().SourceDestroy(r.Context(), sourceSerialized.ID)
	})
}

func ctxSourceSerialized(r *http.Request) (db.SourceSerialized, int, error) {
	sourceSerialized, ok := r.Context().Value(contextSource).(db.SourceSerialized)
	if !ok {
		return db.SourceSerialized{}, http.StatusInternalServerError, fmt.Errorf("cast to db.SourceSerialized fail")
	}
	return sourceSerialized, 0, nil
}
