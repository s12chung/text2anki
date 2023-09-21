package api

import (
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/util/chiutil"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

// PartCreateMultiRequest represents a PartCreateMulti request
type PartCreateMultiRequest struct {
	PrePartListID string                       `json:"pre_part_list_id,omitempty"`
	Parts         []PartCreateMultiRequestPart `json:"parts"`
}

// PartCreateMultiRequestPart represents a SourcePart for PartCreateMultiRequest
type PartCreateMultiRequestPart struct {
	Text        string `json:"text"`
	Translation string `json:"translation,omitempty"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(PartCreateMultiRequest{}).Validates(firm.RuleMap{
		"Parts": {rule.Present{}},
	}))
}

// PartCreateMulti creates multiple Source parts
func (rs Routes) PartCreateMulti(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	req := PartCreateMultiRequest{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	prePartList, httpErr := rs.prePartListFromID(req.PrePartListID)
	if httpErr != nil {
		return nil, httpErr
	}
	parts, httpErr := rs.requestPartsToDBParts(r.Context(), req.Parts, prePartList)
	if httpErr != nil {
		return nil, httpErr
	}

	sourceStructured, httpErr := sourceStructuredFromID(r, txQs)
	if httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		sourceStructured.Parts = append(sourceStructured.Parts, parts...)
		source, err := txQs.SourcePartsUpdate(txQs.Ctx(), sourceStructured.UpdatePartsParams())
		return source.ToSourceStructured(), err
	})
}

// PartCreateOrUpdateRequest represents a PartCreate or PartUpdate request
type PartCreateOrUpdateRequest PartCreateMultiRequestPart

func init() {
	firm.RegisterType(firm.NewDefinition(PartCreateOrUpdateRequest{}).Validates(firm.RuleMap{
		"Text": {rule.TrimPresent{}},
	}))
}

// PartCreate creates a single Source part
func (rs Routes) PartCreate(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	return rs.sourceUpdatePart(r, txQs, func(sourceStructured *db.SourceStructured, tokenizedTexts []db.TokenizedText) *jhttp.HTTPError {
		sourceStructured.Parts = append(sourceStructured.Parts, db.SourcePart{TokenizedTexts: tokenizedTexts})
		return nil
	})
}

// PartUpdate updates a single Source part
func (rs Routes) PartUpdate(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	return rs.sourceUpdatePart(r, txQs, func(sourceStructured *db.SourceStructured, tokenizedTexts []db.TokenizedText) *jhttp.HTTPError {
		partIndex, httpErr := sourcePartIndex(r, *sourceStructured)
		if httpErr != nil {
			return httpErr
		}
		sourceStructured.Parts[partIndex].TokenizedTexts = tokenizedTexts
		return nil
	})
}

func (rs Routes) sourceUpdatePart(r *http.Request, txQs db.TxQs,
	changeFunc func(sourceStructured *db.SourceStructured, tokenizedTexts []db.TokenizedText) *jhttp.HTTPError) (any, *jhttp.HTTPError) {
	req := PartCreateOrUpdateRequest{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	sourceStructured, httpErr := sourceStructuredFromID(r, txQs)
	if httpErr != nil {
		return nil, httpErr
	}

	tokenizedTexts, err := rs.TextTokenizer.TokenizedTexts(r.Context(), req.Text, req.Translation)
	if err != nil {
		return nil, jhttp.Error(http.StatusInternalServerError, err)
	}
	if httpErr := changeFunc(&sourceStructured, tokenizedTexts); httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		source, err := txQs.SourcePartsUpdate(txQs.Ctx(), sourceStructured.UpdatePartsParams())
		return source.ToSourceStructured(), err
	})
}

// PartDestroy destroys a Source part
func (rs Routes) PartDestroy(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	sourceStructured, partIndex, httpErr := sourceStructuredAndPartIndex(r, txQs)
	if httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		sourceStructured.Parts = append(sourceStructured.Parts[:partIndex], sourceStructured.Parts[partIndex+1:]...)
		source, err := txQs.SourcePartsUpdate(txQs.Ctx(), sourceStructured.UpdatePartsParams())
		return source.ToSourceStructured(), err
	})
}

func sourceStructuredAndPartIndex(r *http.Request, txQs db.TxQs) (db.SourceStructured, int, *jhttp.HTTPError) {
	sourceStructured, httpErr := sourceStructuredFromID(r, txQs)
	if httpErr != nil {
		return db.SourceStructured{}, 0, httpErr
	}
	partIndex, httpErr := sourcePartIndex(r, sourceStructured)
	if httpErr != nil {
		return db.SourceStructured{}, 0, httpErr
	}
	return sourceStructured, partIndex, nil
}

func sourcePartIndex(r *http.Request, sourceStructured db.SourceStructured) (int, *jhttp.HTTPError) {
	index, err := chiutil.ParamID(r, "partIndex")
	if err != nil {
		return 0, jhttp.Error(http.StatusUnprocessableEntity, err)
	}
	partIndex := int(index)

	partsLen := len(sourceStructured.Parts)
	if partIndex < 0 {
		return 0, jhttp.Error(http.StatusUnprocessableEntity,
			fmt.Errorf("partIndex (%v) is less than zero", partIndex))
	}
	if partIndex >= partsLen {
		return 0, jhttp.Error(http.StatusUnprocessableEntity,
			fmt.Errorf("partIndex (%v) is greater than existing parts len (%v)", partIndex, partsLen))
	}
	return partIndex, nil
}
