package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(PrePartListSignResponse{}, PrePartList{})
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

var prePartListSignPutConfig = storage.SignPutConfig{
	Table:  sourcesTable,
	Column: partsColumn,
	NameToValidExts: map[string]map[string]bool{
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

// PrePartListSign returns signed requests to generate Source Parts
func (rs Routes) PrePartListSign(r *http.Request) (any, *httputil.HTTPError) {
	req := PrePartListSignRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}

	resp := PrePartListSignResponse{}
	if err := rs.Storage.DBStorage.SignPutTree(prePartListSignPutConfig, req, &resp); err != nil {
		if storage.IsInvalidInputError(err) {
			return nil, httputil.Error(http.StatusUnprocessableEntity, err)
		}
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return resp, nil
}

// PrePartList represents all the Source parts together for a given id
type PrePartList struct {
	ID       string    `json:"id"`
	PreParts []PrePart `json:"pre_parts"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartList) StaticCopy() any {
	return p
}

// PrePart represents a Source part before it is created, only stored via. Routes.Storage.Storer
type PrePart struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}

// PrePartListGet returns the PrePartList for a given ID
func (rs Routes) PrePartListGet(r *http.Request) (any, *httputil.HTTPError) {
	prePartListID := chi.URLParam(r, "prePartListID")
	if prePartListID == "" {
		return nil, httputil.Error(http.StatusNotFound, fmt.Errorf("prePartListID not found"))
	}
	prePartList := PrePartList{}
	err := rs.Storage.DBStorage.SignGetTree(sourcesTable, partsColumn, prePartListID, &prePartList)
	if err != nil {
		if storage.IsNotFoundError(err) {
			return nil, httputil.Error(http.StatusNotFound, err)
		}
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return prePartList, nil
}
