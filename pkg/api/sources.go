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
	httptyped.RegisterType(db.SourceSerialized{})
}

const contextSource httputil.ContextKey = "source"
const sourcesTable = "sources"
const partsColumn = "parts"

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
	PrePartListID string                    `json:"pre_part_list_id,omitempty"`
	Parts         []SourceCreateRequestPart `json:"parts"`
}

// SourceCreateRequestPart represents a part of a Source in a SourceCreate request
type SourceCreateRequestPart struct {
	Text        string `json:"text"`
	Translation string `json:"translation,omitempty"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(SourceCreateRequest{}).Validates(firm.RuleMap{
		"Parts": {rule.Presence{}},
	}))
	firm.RegisterType(firm.NewDefinition(SourceCreateRequestPart{}).Validates(firm.RuleMap{
		"Text": {rule.Presence{}},
	}))
}

// PrePartMediaList is the list of media for the SourcePart with an ID
type PrePartMediaList struct {
	ID       string               `json:"id"`
	PreParts []db.SourcePartMedia `json:"pre_parts"`
}

// SourceCreate creates a new source
func (rs Routes) SourceCreate(r *http.Request) (any, *httputil.HTTPError) {
	req := SourceCreateRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}

	prePartList := PrePartMediaList{}
	if req.PrePartListID != "" {
		if err := rs.Storage.DBStorage.KeyTree(sourcesTable, partsColumn, req.PrePartListID, &prePartList); err != nil {
			if storage.IsNotFoundError(err) {
				return nil, httputil.Error(http.StatusNotFound, err)
			}
			return nil, httputil.Error(http.StatusInternalServerError, err)
		}
	}

	parts := make([]db.SourcePart, len(req.Parts))
	for i, part := range req.Parts {
		tokenizedTexts, err := rs.TextTokenizer.TokenizedTexts(part.Text, part.Translation)
		if err != nil {
			return nil, httputil.Error(http.StatusUnprocessableEntity, err)
		}
		part := db.SourcePart{TokenizedTexts: tokenizedTexts}
		if req.PrePartListID != "" {
			part.Media = &prePartList.PreParts[i]
		}
		parts[i] = part
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
