package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

func init() {
	httptyped.RegisterType(db.SourceStructured{})
}

// SourcesIndex returns a list of sources
func (rs Routes) SourcesIndex(_ *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	return jhttp.ReturnModelOr500(func() (any, error) {
		return txQs.SourceStructuredIndex(txQs.Ctx())
	})
}

// SourceGet gets the source
func (rs Routes) SourceGet(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	return sourceStructuredFromID(r, txQs)
}

// SourceCreateRequest represents the SourceCreate request
type SourceCreateRequest struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	PartCreateMultiRequest
}

func init() {
	firm.MustRegisterType(firm.NewDefinition[SourceCreateRequest]().Validates(firm.RuleMap{
		"PartCreateMultiRequest": {},
	}))
}

// SourceCreate creates a new source
func (rs Routes) SourceCreate(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	req := SourceCreateRequest{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	source, err := rs.sourceStructuredFromBodyReq(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		source, err := txQs.SourceCreate(txQs.Ctx(), source.CreateParams())
		return source.ToSourceStructured(), err
	})
}

// SourceUpdateRequest represents the SourceUpdate request
type SourceUpdateRequest struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

func init() {
	firm.MustRegisterType(firm.NewDefinition[SourceUpdateRequest]().Validates(firm.RuleMap{
		"Name": {rule.Present{}},
	}))
}

// SourceUpdate updates the source
func (rs Routes) SourceUpdate(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	sourceStructured, httpErr := sourceStructuredFromID(r, txQs)
	if httpErr != nil {
		return nil, httpErr
	}

	req := SourceUpdateRequest{}
	if httpErr = extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	sourceStructured.Name = req.Name
	sourceStructured.Reference = req.Reference

	return jhttp.ReturnModelOr500(func() (any, error) {
		source, err := txQs.SourceUpdate(txQs.Ctx(), sourceStructured.UpdateParams())
		return source.ToSourceStructured(), err
	})
}

// SourceDestroy destroys the source
func (rs Routes) SourceDestroy(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	sourceStructured, httpErr := sourceStructuredFromID(r, txQs)
	if httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		return sourceStructured, txQs.SourceDestroy(txQs.Ctx(), sourceStructured.ID)
	})
}

func (rs Routes) prePartListFromID(prePartListID string) (*db.PrePartList, *jhttp.HTTPError) {
	if prePartListID == "" {
		return nil, nil
	}
	prePartList := &db.PrePartList{}
	if err := rs.Storage.DBStorage.KeyTree(db.SourcesTable, db.PartsColumn, prePartListID, prePartList); err != nil {
		if storage.IsNotFoundError(err) {
			return nil, jhttp.Error(http.StatusNotFound, err)
		}
		return nil, jhttp.Error(http.StatusInternalServerError, err)
	}
	return prePartList, nil
}

func sourceStructuredFromID(r *http.Request, txQs db.TxQs) (db.SourceStructured, *jhttp.HTTPError) {
	id, httpErr := idFromRequest(r)
	if httpErr != nil {
		return db.SourceStructured{}, httpErr
	}
	source, err := txQs.SourceGet(r.Context(), id)
	if err != nil {
		return db.SourceStructured{}, jhttp.Error(http.StatusNotFound, err)
	}
	return source.ToSourceStructured(), nil
}

func (rs Routes) sourceStructuredFromBodyReq(ctx context.Context, req SourceCreateRequest) (*db.SourceStructured, *jhttp.HTTPError) {
	prePartList, httpErr := rs.prePartListFromID(req.PrePartListID)
	if httpErr != nil {
		return nil, httpErr
	}
	name, reference, httpErr := rs.chooseSourceNameRef(req.Name, req.Reference, prePartList)
	if httpErr != nil {
		return nil, httpErr
	}
	parts, httpErr := rs.requestPartsToDBParts(ctx, req.Parts, prePartList)
	if httpErr != nil {
		return nil, httpErr
	}
	return &db.SourceStructured{Name: name, Reference: reference, Parts: parts}, nil
}

func (rs Routes) chooseSourceNameRef(name, ref string, prePartList *db.PrePartList) (string, string, *jhttp.HTTPError) {
	if prePartList == nil || (name != "" && ref != "") {
		return name, ref, nil
	}

	info, err := prePartList.Info()
	if err != nil {
		return "", "", jhttp.Error(http.StatusInternalServerError, err)
	}
	if name == "" {
		name = info.Name
	}
	if ref == "" {
		ref = info.Reference
	}
	return name, ref, nil
}

func (rs Routes) requestPartsToDBParts(ctx context.Context, reqParts []PartCreateMultiRequestPart,
	prePartList *db.PrePartList) ([]db.SourcePart, *jhttp.HTTPError) {
	parts := make([]db.SourcePart, 0, len(reqParts))
	for i, part := range reqParts {
		if strings.TrimSpace(part.Text) == "" {
			continue
		}
		tokenizedTexts, err := rs.TextTokenizer.TokenizedTexts(ctx, part.Text, part.Translation)
		if err != nil {
			return nil, jhttp.Error(http.StatusUnprocessableEntity, err)
		}
		sourcePart := db.SourcePart{TokenizedTexts: tokenizedTexts}
		if prePartList != nil {
			sourcePart.Media = &prePartList.PreParts[i]
		}
		parts = append(parts, sourcePart)
	}
	if len(parts) == 0 {
		return nil, jhttp.Error(http.StatusUnprocessableEntity, errors.New("no parts found with text set"))
	}
	return parts, nil
}
