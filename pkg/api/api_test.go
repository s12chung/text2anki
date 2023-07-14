package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var routesConfig = Config{
	StorageConfig: StorageConfig{
		LocalStoreConfig: LocalStoreConfig{
			Origin:        "https://test.com",
			KeyBasePath:   path.Join(os.TempDir(), test.GenerateName("filestore-router")),
			EncryptorPath: fixture.TestDataDir,
		},
	}}
var routes = NewRoutes(routesConfig)
var server = test.Server{Server: httptest.NewServer(routes.Router())}

type MustSetupAndSeed struct{}

func TestMain(m *testing.M) {
	testdb.MustSetupAndSeed(MustSetupAndSeed{})
	code := m.Run()
	server.Close()
	os.Exit(code)
}

func TestHttpTypedRegistry(t *testing.T) {
	require := require.New(t)
	testName := "TestHttpTypedRegistry"

	fileNames := make([]string, len(httptyped.Types()))
	for i, typ := range httptyped.Types() {
		fileName := typ.String() + ".json"
		fixture.CompareReadOrUpdate(t, path.Join(testName, fileName), fixture.JSON(t, httptyped.StructureMap(typ)))
		fileNames[i] = fileName
	}

	files, err := os.ReadDir(fixture.JoinTestData(testName))
	require.NoError(err)
	expectedFileNames := make([]string, len(files))
	for i, file := range files {
		expectedFileNames[i] = file.Name()
	}
	require.ElementsMatch(expectedFileNames, fileNames)
}

func idPath(path string, id int64) string {
	return fmt.Sprintf(path+"/%v", id)
}

func testModelResponse(t *testing.T, resp test.Response, testName, name string, model test.StaticCopyable) string {
	jsonBody := test.StaticCopyOrIndent(t, resp.Code, resp.Body.Bytes(), model)
	fixtureFile := testName + ".json"
	if name != "" {
		fixtureFile = path.Join(testName, name+"_response.json")
	}
	fixture.CompareReadOrUpdate(t, fixtureFile, jsonBody)
	return fixtureFile
}

func testModelsResponse(t *testing.T, resp test.Response, testName, name string, models any) {
	jsonBody := test.StaticCopyOrIndentSlice(t, resp.Code, resp.Body.Bytes(), models)
	fixtureFile := testName + ".json"
	if name != "" {
		fixtureFile = path.Join(testName, name+"_response.json")
	}
	fixture.CompareReadOrUpdate(t, fixtureFile, jsonBody)
}

func TestRoutes_Router(t *testing.T) {
	testName := "TestRoutes_Router"
	require := require.New(t)

	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))
	r.Mount("/", routes.Router())

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
