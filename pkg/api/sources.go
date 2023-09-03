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

const sourceContextKey httputil.ContextKey = "source"

// SourceCtx sets the source context from the sourceID
func (rs Routes) SourceCtx(r *http.Request) (*http.Request, *httputil.HTTPError) {
	sourceID, err := chiutil.ParamID(r, "sourceID")
	if err != nil {
		return nil, httputil.Error(http.StatusNotFound, err)
	}

	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	source, err := txQs.SourceGet(r.Context(), sourceID)
	if err != nil {
		return nil, httputil.Error(http.StatusNotFound, err)
	}

	r = r.WithContext(context.WithValue(r.Context(), sourceContextKey, source.ToSourceStructured()))
	return r, nil
}

// SourceIndex returns a list of sources
func (rs Routes) SourceIndex(r *http.Request) (any, *httputil.HTTPError) {
	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return txQs.SourceStructuredIndex(txQs.Ctx())
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
	sourceStructured, httpErr := ctxSourceStructured(r)
	if httpErr != nil {
		return nil, httpErr
	}

	req := SourceUpdateRequest{}
	if httpErr = extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	sourceStructured.Name = req.Name
	sourceStructured.Reference = req.Reference

	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := txQs.SourceUpdate(txQs.Ctx(), sourceStructured.UpdateParams())
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

// SourceCreate creates a new source
func (rs Routes) SourceCreate(r *http.Request) (any, *httputil.HTTPError) {
	req := SourceCreateRequest{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}

	var prePartList *db.PrePartList
	if req.PrePartListID != "" {
		prePartList = &db.PrePartList{}
		if err := rs.Storage.DBStorage.KeyTree(db.SourcesTable, db.PartsColumn, req.PrePartListID, prePartList); err != nil {
			if storage.IsNotFoundError(err) {
				return nil, httputil.Error(http.StatusNotFound, err)
			}
			return nil, httputil.Error(http.StatusInternalServerError, err)
		}
	}

	source, err := rs.sourceCreateSource(req, prePartList)
	if err != nil {
		return nil, err
	}

	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		source, err := txQs.SourceCreate(txQs.Ctx(), source.CreateParams())
		return source.ToSourceStructured(), err
	})
}

func (rs Routes) sourceCreateSource(req SourceCreateRequest, prePartList *db.PrePartList) (*db.SourceStructured, *httputil.HTTPError) {
	name, reference, err := rs.sourceCreateSourceNameRef(req.Name, req.Reference, prePartList)
	if err != nil {
		return nil, err
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
		if prePartList != nil {
			part.Media = &prePartList.PreParts[i]
		}
		parts = append(parts, part)
	}

	if len(parts) == 0 {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("no parts found with text set"))
	}
	return &db.SourceStructured{Name: name, Reference: reference, Parts: parts}, nil
}

func (rs Routes) sourceCreateSourceNameRef(name, ref string, prePartList *db.PrePartList) (string, string, *httputil.HTTPError) {
	if prePartList == nil || (name != "" && ref != "") {
		return name, ref, nil
	}

	info, err := prePartList.Info()
	if err != nil {
		return "", "", httputil.Error(http.StatusInternalServerError, err)
	}
	if name == "" {
		name = info.Name
	}
	if ref == "" {
		ref = info.Reference
	}
	return name, ref, nil
}

// SourceDestroy destroys the source
func (rs Routes) SourceDestroy(r *http.Request) (any, *httputil.HTTPError) {
	sourceStructured, httpErr := ctxSourceStructured(r)
	if httpErr != nil {
		return nil, httpErr
	}
	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return sourceStructured, txQs.SourceDestroy(txQs.Ctx(), sourceStructured.ID)
	})
}

func ctxSourceStructured(r *http.Request) (db.SourceStructured, *httputil.HTTPError) {
	sourceStructured, ok := r.Context().Value(sourceContextKey).(db.SourceStructured)
	if !ok {
		return db.SourceStructured{}, httputil.Error(http.StatusInternalServerError, fmt.Errorf("cast to db.SourceStructured fail"))
	}
	return sourceStructured, nil
}
