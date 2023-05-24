// Package chiutil contains utils for go-chi/chi
package chiutil

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ParamID retrieves the ID URL param from the request
func ParamID(r *http.Request, key string) (int64, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, fmt.Errorf(key + " not found")
	}
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
