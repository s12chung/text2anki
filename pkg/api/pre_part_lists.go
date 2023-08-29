package api

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(PrePartListSignResponse{}, PrePartListURL{}, PrePartListVerifyResponse{}, PrePartListCreateResponse{})
}

// PrePartListFile is the fileTree for PreParts
type PrePartListFile struct {
	InfoFile fs.File                  `json:"info_file,omitempty"`
	PreParts []db.SourcePartMediaFile `json:"pre_parts"`
}

// PrePartList is a KeyTree for PreParts
type PrePartList struct {
	ID       string               `json:"id"`
	InfoKey  string               `json:"info_key,omitempty"`
	PreParts []db.SourcePartMedia `json:"pre_parts"`
}

var prePartListPutConfig = storage.SignPutConfig{
	Table:  sourcesTable,
	Column: partsColumn,
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
		"PreParts": {rule.Presence{}},
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
func (rs Routes) PrePartListSign(r *http.Request) (any, *httputil.HTTPError) {
	req := PrePartListSignRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}

	resp := PrePartListSignResponse{}
	if err := rs.Storage.DBStorage.SignPutTree(prePartListPutConfig, req, &resp); err != nil {
		if storage.IsInvalidInputError(err) {
			return nil, httputil.Error(http.StatusUnprocessableEntity, err)
		}
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return resp, nil
}

// PrePartListURL represents all the Source parts together for a given id
type PrePartListURL struct {
	ID       string            `json:"id"`
	PreParts []PrePartMediaURL `json:"pre_parts"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartListURL) StaticCopy() any {
	return p
}

// PrePartMediaURL represents a SourcePartMedia before it is created, only stored via. Routes.Storage.Storer
type PrePartMediaURL struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}

// PrePartListGet returns the PrePartListURL for a given ID
func (rs Routes) PrePartListGet(r *http.Request) (any, *httputil.HTTPError) {
	prePartListID := chi.URLParam(r, "prePartListID")
	if prePartListID == "" {
		return nil, httputil.Error(http.StatusNotFound, fmt.Errorf("prePartListID not found"))
	}
	prePartList := PrePartListURL{}
	err := rs.Storage.DBStorage.SignGetTree(sourcesTable, partsColumn, prePartListID, &prePartList)
	if err != nil {
		if storage.IsNotFoundError(err) {
			return nil, httputil.Error(http.StatusNotFound, err)
		}
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return prePartList, nil
}

// PrePartListVerifyRequest represents a PrePartListVerify request
type PrePartListVerifyRequest struct {
	Text string `json:"text"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(PrePartListVerifyRequest{}).Validates(firm.RuleMap{
		"Text": {rule.Presence{}},
	}))
}

// PrePartListVerifyResponse represents a PrePartListVerify response
type PrePartListVerifyResponse struct {
	ExtractorType string `json:"extractor_type"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartListVerifyResponse) StaticCopy() any {
	return p
}

// PrePartListVerify verifies the text whether it fits any extractor and returns the extractor type
func (rs Routes) PrePartListVerify(r *http.Request) (any, *httputil.HTTPError) {
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
		"ExtractorType": {rule.Presence{}},
		"Text":          {rule.Presence{}},
	}))
}

// PrePartListCreateResponse represents a PrePartListCreate response
type PrePartListCreateResponse struct {
	ID string `json:"id"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartListCreateResponse) StaticCopy() any {
	return p
}

// PrePartListCreate creates PrePartList given the type of extractor and text
func (rs Routes) PrePartListCreate(r *http.Request) (any, *httputil.HTTPError) {
	req := PrePartListCreateRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}
	ex, exists := rs.ExtractorMap[req.ExtractorType]
	if !exists {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("given type is not valid: %v", req.ExtractorType))
	}
	extraction, err := ex.Extract(req.Text)
	if err != nil {
		return nil, httputil.Error(http.StatusUnprocessableEntity, err)
	}
	prePartListKey := PrePartList{}
	infoFile, err := extraction.InfoFile()
	if err != nil {
		return nil, httputil.Error(http.StatusUnprocessableEntity, err)
	}
	prePartListFile := PrePartListFile{InfoFile: infoFile, PreParts: extraction.Parts}
	if err := rs.Storage.DBStorage.PutTree(prePartListPutConfig, prePartListFile, &prePartListKey); err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}

	return PrePartListCreateResponse{ID: prePartListKey.ID}, nil
}
