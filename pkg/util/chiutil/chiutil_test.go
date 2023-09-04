package chiutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestParamID(t *testing.T) {
	r := chi.NewRouter()
	r.Route("/{articleId}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			id, err := ParamID(r, "articleId")
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			if _, err = fmt.Fprintf(w, "%v", id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	testCases := []struct {
		name         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{name: "normal", path: "/123", expectedCode: http.StatusOK, expectedBody: "123"},
		{name: "not a number", path: "/nan", expectedCode: http.StatusNotFound, expectedBody: "strconv.ParseInt: parsing \"nan\": invalid syntax\n0"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+tc.path, nil)
			require.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(err)
			defer func() { require.NoError(resp.Body.Close()) }()
			require.Equal(tc.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(err)
			require.Equal(tc.expectedBody, string(body))
		})
	}
}
