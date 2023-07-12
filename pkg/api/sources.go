package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/chiutil"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(db.SourceSerialized{}, storage.PresignedHTTPRequest{})
}

const contextSource httputil.ContextKey = "source"

// SourceCtx sets the source context from the sourceID
func SourceCtx(r *http.Request) (*http.Request, *httputil.HTTPError) {
	sourceID, err := chiutil.ParamID(r, "sourceID")
	if err != nil {
		return nil, httputil.Error(http.StatusNotFound, err)
	}

	source, err := db.Qs().SourceGet(r.Context(), sourceID)
	if err != nil {
		return nil, httputil.Error(http.StatusNotFound, err)
	}

	r = r.WithContext(context.WithValue(r.Context(), contextSource, source.ToSourceSerialized()))
	return r, nil
}

// SourceIndex returns a list of sources
func (rs Routes) SourceIndex(r *http.Request) (any, *httputil.HTTPError) {
	return httputil.ReturnModelOr500(func() (any, error) {
		return db.Qs().SourceSerializedIndex(r.Context())
	})
}

// SourceGet gets the source
func (rs Routes) SourceGet(r *http.Request) (any, *httputil.HTTPError) {
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
func (rs Routes) SourceUpdate(r *http.Request) (any, *httputil.HTTPError) {
	sourceSerialized, httpError := ctxSourceSerialized(r)
	if httpError != nil {
		return nil, httpError
	}

	req := SourceUpdateRequest{}
	if httpError = extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}
	sourceSerialized.Name = req.Name

	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := db.Qs().SourceUpdate(r.Context(), sourceSerialized.UpdateParams())
		return source.ToSourceSerialized(), err
	})
}

// SourceCreateRequest represents the SourceCreate request
type SourceCreateRequest struct {
	Parts []SourceCreateRequestPart
}

// SourceCreateRequestPart represents a part of a Source in a SourceCreate request
type SourceCreateRequestPart struct {
	Text        string
	Translation string
}

func init() {
	firm.RegisterType(firm.NewDefinition(SourceCreateRequest{}).Validates(firm.RuleMap{
		"Parts": {rule.Presence{}},
	}))
	firm.RegisterType(firm.NewDefinition(SourceCreateRequestPart{}).Validates(firm.RuleMap{
		"Text": {rule.Presence{}},
	}))
}

// SourceCreate creates a new source
func (rs Routes) SourceCreate(r *http.Request) (any, *httputil.HTTPError) {
	req := SourceCreateRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}

	parts := make([]db.SourcePart, len(req.Parts))
	for i, part := range req.Parts {
		tokenizedTexts, err := rs.TextTokenizer.TokenizedTexts(part.Text, part.Translation)
		if err != nil {
			return nil, httputil.Error(http.StatusUnprocessableEntity, err)
		}
		parts[i] = db.SourcePart{TokenizedTexts: tokenizedTexts}
	}

	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := db.Qs().SourceCreate(r.Context(), db.SourceSerialized{Parts: parts}.CreateParams())
		return source.ToSourceSerialized(), err
	})
}

// SourceDestroy destroys the source
func (rs Routes) SourceDestroy(r *http.Request) (any, *httputil.HTTPError) {
	sourceSerialized, httpError := ctxSourceSerialized(r)
	if httpError != nil {
		return nil, httpError
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return sourceSerialized, db.Qs().SourceDestroy(r.Context(), sourceSerialized.ID)
	})
}

func ctxSourceSerialized(r *http.Request) (db.SourceSerialized, *httputil.HTTPError) {
	sourceSerialized, ok := r.Context().Value(contextSource).(db.SourceSerialized)
	if !ok {
		return db.SourceSerialized{}, httputil.Error(http.StatusInternalServerError, fmt.Errorf("cast to db.SourceSerialized fail"))
	}
	return sourceSerialized, nil
}

var validSignPartsExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// SignParts returns signed requests to generate Source Parts
func (rs Routes) SignParts(r *http.Request) (any, *httputil.HTTPError) {
	exts := r.URL.Query()["exts"]
	if len(exts) == 0 {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("no file extension given"))
	}
	for _, ext := range exts {
		if !validSignPartsExts[ext] {
			return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("%v is not a valid file extension", ext))
		}
	}

	reqs, err := rs.Signer.Sign("sources", "parts", exts)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return reqs, nil
}
