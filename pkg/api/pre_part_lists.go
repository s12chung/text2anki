package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

func init() {
	httptyped.RegisterType(PrePartListSignResponse{}, db.PrePartListURL{}, PrePartListVerifyResponse{}, PrePartListCreateResponse{})
}

var prePartListPutConfig = storage.SignPutConfig{
	Table:  db.SourcesTable,
	Column: db.PartsColumn,
	NameToValidExts: map[string]map[string]bool{
		"Info": {".json": true},
		"Image": {
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		},
		"Audio": {
			".mp3": true,
		},
	},
}

// PrePartListSignRequest represents the PrePartListSign request
type PrePartListSignRequest struct {
	PreParts []PrePartSignRequest `json:"pre_parts"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(PrePartListSignRequest{}).Validates(firm.RuleMap{
		"PreParts": {rule.Present{}},
	}))
}

// PrePartSignRequest represents a pre_part for PrePartListSign request
type PrePartSignRequest struct {
	ImageExt string `json:"image_ext,omitempty"`
	AudioExt string `json:"audio_ext,omitempty"`
}

// PrePartListSignResponse is the response returned by PrePartListSign
type PrePartListSignResponse struct {
	ID       string                `json:"id"`
	PreParts []PrePartSignResponse `json:"pre_parts"`
}

// PrePartSignResponse represents a pre_part for PrePartListSign response
type PrePartSignResponse struct {
	ImageRequest *storage.PreSignedHTTPRequest `json:"image_request,omitempty"`
	AudioRequest *storage.PreSignedHTTPRequest `json:"audio_request,omitempty"`
}

// PrePartListSign returns signed requests to generate Source Parts
func (rs Routes) PrePartListSign(r *http.Request, _ db.TxQs) (any, *jhttp.HTTPError) {
	req := PrePartListSignRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}

	resp := PrePartListSignResponse{}
	if err := rs.Storage.DBStorage.SignPutTree(prePartListPutConfig, req, &resp); err != nil {
		if storage.IsInvalidInputError(err) {
			return nil, jhttp.Error(http.StatusUnprocessableEntity, err)
		}
		return nil, jhttp.Error(http.StatusInternalServerError, err)
	}
	return resp, nil
}

// PrePartListGet returns the PrePartListURL for a given ID
func (rs Routes) PrePartListGet(r *http.Request, _ db.TxQs) (any, *jhttp.HTTPError) {
	prePartListID := chi.URLParam(r, "id")
	if prePartListID == "" {
		return nil, jhttp.Error(http.StatusNotFound, fmt.Errorf("id not found"))
	}
	prePartList := db.PrePartListURL{}
	err := rs.Storage.DBStorage.SignGetTree(db.SourcesTable, db.PartsColumn, prePartListID, &prePartList)
	if err != nil {
		if storage.IsNotFoundError(err) {
			return nil, jhttp.Error(http.StatusNotFound, err)
		}
		return nil, jhttp.Error(http.StatusInternalServerError, err)
	}
	return prePartList, nil
}

// PrePartListVerifyRequest represents a PrePartListVerify request
type PrePartListVerifyRequest struct {
	Text string `json:"text"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(PrePartListVerifyRequest{}).Validates(firm.RuleMap{
		"Text": {rule.Present{}},
	}))
}

// PrePartListVerifyResponse represents a PrePartListVerify response
type PrePartListVerifyResponse struct {
	ExtractorType string `json:"extractor_type"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartListVerifyResponse) StaticCopy() PrePartListVerifyResponse { return p }

// PrePartListVerify verifies the text whether it fits any extractor and returns the extractor type
func (rs Routes) PrePartListVerify(r *http.Request, _ db.TxQs) (any, *jhttp.HTTPError) {
	req := PrePartListVerifyRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}
	return PrePartListVerifyResponse{ExtractorType: extractor.Verify(req.Text, rs.ExtractorMap)}, nil
}

// PrePartListCreateRequest represents a PrePartListCreate request
type PrePartListCreateRequest struct {
	ExtractorType string `json:"extractor_type"`
	Text          string `json:"text"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(PrePartListCreateRequest{}).Validates(firm.RuleMap{
		"ExtractorType": {rule.Present{}},
		"Text":          {rule.Present{}},
	}))
}

// PrePartListCreateResponse represents a PrePartListCreate response
type PrePartListCreateResponse struct {
	ID string `json:"id"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartListCreateResponse) StaticCopy() PrePartListCreateResponse { return p }

// PrePartListCreate creates PrePartList given the type of extractor and text
func (rs Routes) PrePartListCreate(r *http.Request, _ db.TxQs) (any, *jhttp.HTTPError) {
	req := PrePartListCreateRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}
	ex, exists := rs.ExtractorMap[req.ExtractorType]
	if !exists {
		return nil, jhttp.Error(http.StatusUnprocessableEntity, fmt.Errorf("given type is not valid: %v", req.ExtractorType))
	}
	extraction, err := ex.Extract(req.Text)
	if err != nil {
		return nil, jhttp.Error(http.StatusUnprocessableEntity, err)
	}
	prePartListKey := db.PrePartList{}
	infoFile, err := extraction.InfoFile()
	if err != nil {
		return nil, jhttp.Error(http.StatusUnprocessableEntity, err)
	}
	prePartListFile := db.PrePartListFile{InfoFile: infoFile, PreParts: extraction.Parts}
	if err := rs.Storage.DBStorage.PutTree(prePartListPutConfig, prePartListFile, &prePartListKey); err != nil {
		return nil, jhttp.Error(http.StatusInternalServerError, err)
	}

	return PrePartListCreateResponse{ID: prePartListKey.ID}, nil
}
