package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(PrePartsSignResponse{}, PreParts{})
}

var validSignPartsExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// PrePartsSignResponse is the response returned by PrePartsSign
type PrePartsSignResponse struct {
	ID       string                         `json:"id"`
	Requests []storage.PreSignedHTTPRequest `json:"requests"`
}

const sourcesTable = "sources"
const partsColumn = "parts"

// PrePartsSign returns signed requests to generate Source Parts
func (rs Routes) PrePartsSign(r *http.Request) (any, *httputil.HTTPError) {
	exts := r.URL.Query()["exts"]
	if len(exts) == 0 {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("no file extension given"))
	}
	for _, ext := range exts {
		if !validSignPartsExts[ext] {
			return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("%v is not a valid file extension", ext))
		}
	}

	reqs, id, err := rs.Storage.Signer.SignPut(sourcesTable, partsColumn, exts)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return PrePartsSignResponse{ID: id, Requests: reqs}, nil
}

// PreParts represents all the Source parts together for a given id
type PreParts struct {
	ID       string    `json:"id"`
	PreParts []PrePart `json:"pre_parts"`
}

// StaticCopy returns a copy without fields that variate
func (p PreParts) StaticCopy() any {
	return p
}

// PrePart represents a Source part before it is created, only stored via. Routes.Storage.Storer
type PrePart struct {
	URL string `json:"url"`
}

// PrePartsGet returns the PreParts for a given ID
func (rs Routes) PrePartsGet(r *http.Request) (any, *httputil.HTTPError) {
	prePartsID := chi.URLParam(r, "prePartsID")
	if prePartsID == "" {
		return nil, httputil.Error(http.StatusNotFound, fmt.Errorf("prePartsID not found"))
	}
	urls, err := rs.Storage.Signer.SignGetByID(sourcesTable, partsColumn, prePartsID)
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
	return PreParts{ID: prePartsID, PreParts: preParts}, nil
}
