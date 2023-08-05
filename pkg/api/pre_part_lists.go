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

var signedImageConfig = signFieldConfig{
	Name: "image",
	ValidExts: map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	},
}

var signedAudioConfig = signFieldConfig{
	Name: "audio",
	ValidExts: map[string]bool{
		".mp3": true,
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

const sourcesTable = "sources"
const partsColumn = "parts"

// PrePartListSign returns signed requests to generate Source Parts
func (rs Routes) PrePartListSign(r *http.Request) (any, *httputil.HTTPError) {
	req := PrePartListSignRequest{}
	if httpError := extractAndValidate(r, &req); httpError != nil {
		return nil, httpError
	}

	builder, err := rs.Storage.Signer.SignPutBuilder(sourcesTable, partsColumn)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}

	response := PrePartListSignResponse{ID: builder.ID(), PreParts: make([]PrePartSignResponse, len(req.PreParts))}

	for i, prePart := range req.PreParts {
		b := builder.Index(i)

		respPrePart := PrePartSignResponse{}

		var httpError *httputil.HTTPError
		respPrePart.ImageRequest, httpError = signFieldIfExists(b, signedImageConfig, prePart.ImageExt)
		if httpError != nil {
			return nil, httpError
		}
		respPrePart.AudioRequest, httpError = signFieldIfExists(b, signedAudioConfig, prePart.AudioExt)
		if httpError != nil {
			return nil, httpError
		}

		response.PreParts[i] = respPrePart
	}
	return response, nil
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
	URL string `json:"url"`
}

// PrePartListGet returns the PrePartList for a given ID
func (rs Routes) PrePartListGet(r *http.Request) (any, *httputil.HTTPError) {
	prePartListID := chi.URLParam(r, "prePartListID")
	if prePartListID == "" {
		return nil, httputil.Error(http.StatusNotFound, fmt.Errorf("prePartListID not found"))
	}
	urls, err := rs.Storage.Signer.SignGetByID(sourcesTable, partsColumn, prePartListID)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	if len(urls) == 0 {
		return nil, httputil.Error(http.StatusNotFound, fmt.Errorf("no pre-parts found"))
	}

	preParts := make([]PrePart, len(urls))
	for i, u := range urls {
		preParts[i] = PrePart{URL: u}
	}
	return PrePartList{ID: prePartListID, PreParts: preParts}, nil
}
