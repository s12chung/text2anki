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
	httptyped.RegisterType(PrePartListSignResponse{}, PrePartList{})
}

var validSignPartExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// PrePartListSignResponse is the response returned by PrePartListSign
type PrePartListSignResponse struct {
	ID       string                         `json:"id"`
	Requests []storage.PreSignedHTTPRequest `json:"requests"`
}

const sourcesTable = "sources"
const partsColumn = "parts"

// PrePartListSign returns signed requests to generate Source Parts
func (rs Routes) PrePartListSign(r *http.Request) (any, *httputil.HTTPError) {
	exts := r.URL.Query()["exts"]
	if len(exts) == 0 {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("no file extension given"))
	}
	for _, ext := range exts {
		if !validSignPartExts[ext] {
			return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("%v is not a valid file extension", ext))
		}
	}

	reqs, id, err := rs.Storage.Signer.SignPut(sourcesTable, partsColumn, exts)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return PrePartListSignResponse{ID: id, Requests: reqs}, nil
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
