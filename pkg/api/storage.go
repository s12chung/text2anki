package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

func init() {
	httptyped.RegisterType(StoragePutOk{})
}

// StoragePutOk is the Ok response for StoragePut
type StoragePutOk struct {
	Message string `json:"message"`
}

// StaticCopy returns a copy without fields that variate
func (s StoragePutOk) StaticCopy() StoragePutOk { return s }

// StoragePut stores the file with the route's Storer
func (rs Routes) StoragePut(r *http.Request) (any, *jhttp.HTTPError) {
	storer := rs.Storage.Storer
	key := chi.URLParam(r, "*")
	if err := storer.Validate(key, r.URL.Query()); err != nil {
		return nil, jhttp.Error(http.StatusUnprocessableEntity, err)
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		if err := storer.Store(key, r.Body); err != nil {
			return nil, err
		}
		if err := r.Body.Close(); err != nil {
			return nil, err
		}
		return StoragePutOk{Message: "success"}, nil
	})
}

// StorageGet gets the file with the route's Storer
func (rs Routes) StorageGet() http.Handler {
	return rs.Storage.Storer.FileHandler()
}
