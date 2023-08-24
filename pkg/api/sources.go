package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/chiutil"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(db.SourceStructured{})
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

	r = r.WithContext(context.WithValue(r.Context(), contextSource, source.ToSourceStructured()))
	return r, nil
}

// SourceIndex returns a list of sources
func (rs Routes) SourceIndex(r *http.Request) (any, *httputil.HTTPError) {
	return httputil.ReturnModelOr500(func() (any, error) {
		return db.Qs().SourceStructuredIndex(r.Context())
	})
}

// SourceGet gets the source
func (rs Routes) SourceGet(r *http.Request) (any, *httputil.HTTPError) {
	return ctxSourceStructured(r)
}

// SourceUpdateRequest represents the SourceUpdate request
type SourceUpdateRequest struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(SourceUpdateRequest{}).Validates(firm.RuleMap{
		"Name": {rule.Presence{}},
	}))
}

// SourceUpdate updates the source
func (rs Routes) SourceUpdate(r *http.Request) (any, *httputil.HTTPError) {
	sourceStructured, httpError := ctxSourceStructured(r)
	if httpError != nil {
		return nil, httpError
	}

	req := SourceUpdateRequest{}
	if httpError = extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}
	sourceStructured.Name = req.Name
	sourceStructured.Reference = req.Reference

	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := db.Qs().SourceUpdate(r.Context(), sourceStructured.UpdateParams())
		return source.ToSourceStructured(), err
	})
}

// SourceCreateRequest represents the SourceCreate request
type SourceCreateRequest struct {
	PrePartListID string                    `json:"pre_part_list_id,omitempty"`
	Name          string                    `json:"name"`
	Reference     string                    `json:"reference"`
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

	parts := make([]db.SourcePart, 0, len(req.Parts))
	for i, part := range req.Parts {
		if strings.TrimSpace(part.Text) == "" {
			continue
		}
		tokenizedTexts, err := rs.TextTokenizer.TokenizedTexts(part.Text, part.Translation)
		if err != nil {
			return nil, httputil.Error(http.StatusUnprocessableEntity, err)
		}
		part := db.SourcePart{TokenizedTexts: tokenizedTexts}
		if req.PrePartListID != "" {
			part.Media = &prePartList.PreParts[i]
		}
		parts = append(parts, part)
	}

	if len(parts) == 0 {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("no parts found with text set"))
	}

	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := db.Qs().SourceCreate(r.Context(), db.SourceStructured{Name: req.Name, Reference: req.Reference, Parts: parts}.CreateParams())
		return source.ToSourceStructured(), err
	})
}

// SourceDestroy destroys the source
func (rs Routes) SourceDestroy(r *http.Request) (any, *httputil.HTTPError) {
	sourceStructured, httpError := ctxSourceStructured(r)
	if httpError != nil {
		return nil, httpError
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return sourceStructured, db.Qs().SourceDestroy(r.Context(), sourceStructured.ID)
	})
}

func ctxSourceStructured(r *http.Request) (db.SourceStructured, *httputil.HTTPError) {
	sourceStructured, ok := r.Context().Value(contextSource).(db.SourceStructured)
	if !ok {
		return db.SourceStructured{}, httputil.Error(http.StatusInternalServerError, fmt.Errorf("cast to db.SourceStructured fail"))
	}
	return sourceStructured, nil
}
