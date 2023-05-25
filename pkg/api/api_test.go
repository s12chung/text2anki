package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var server test.Server
var sourcesServer test.Server

func TestMain(m *testing.M) {
	testdb.MustSetupAndSeed("api.TestMain()")

	server = test.Server{Server: httptest.NewServer(DefaultRoutes.Router())}
	sourcesServer = server.WithPathPrefix("/sources")

	code := m.Run()
	server.Close()
	os.Exit(code)
}

func TestRoutes_Router(t *testing.T) {
	testName := "TestRoutes_Router"
	require := require.New(t)

	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))
	r.Mount("/", DefaultRoutes.Router())

	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/sources/1", nil)
	require.NoError(err)

	resp := test.HTTPDo(t, req)
	require.Equal(http.StatusOK, resp.Code)
	jsonBody := test.StaticCopy(t, resp.Body.Bytes(), &db.SourceSerialized{})
	fixture.CompareReadOrUpdate(t, testName+".json", jsonBody)

	req, err = http.NewRequest(http.MethodGet, server.URL+"/healthz", nil)
	require.NoError(err)

	resp = test.HTTPDo(t, req)
	require.Equal(http.StatusOK, resp.Code)
	require.Equal(".", resp.Body.String())
}
