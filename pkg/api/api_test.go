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
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/extractor/extractortest"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
	"github.com/s12chung/text2anki/pkg/util/httputil/reqtx/reqtxtest"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var txPool = reqtxtest.NewPool()

const testUUID = "123e4567-e89b-12d3-a456-426614174000"

type UUIDTest struct{}

func (u UUIDTest) Generate() (string, error) { return testUUID, nil }

var extractorCacheDir = path.Join(os.TempDir(), test.GenerateName("Extractor"))

var routes Routes
var server txServer

func init() {
	testdb.MustSetup()
	routes = NewRoutes(routesConfig)
	server = txServer{pool: txPool, Server: test.Server{Server: httptest.NewServer(routes.Router())}}
	if err := os.MkdirAll(extractorCacheDir, ioutil.OwnerRWXGroupRX); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	server.Close()
	os.Exit(code)
}

const extractorType = "testy"

var routesConfig = config.Config{
	TxPool: txPool,

	StorageConfig: config.StorageConfig{
		LocalStoreConfig: config.LocalStoreConfig{
			Origin:        "https://test.com",
			KeyBasePath:   path.Join(os.TempDir(), test.GenerateName("filestore-router")),
			EncryptorPath: fixture.TestDataDir,
		},
		UUIDGenerator: UUIDTest{},
	},

	ExtractorMap: extractor.Map{
		extractorType: extractor.NewExtractor(extractorCacheDir, extractortest.NewFactory("Extractor")),
	},
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

func TestRoutes_Router(t *testing.T) {
	require := require.New(t)
	testName := "TestRoutes_Router"

	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))
	r.Mount("/", routes.Router())

	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/sources/1", nil)
	require.NoError(err)
	resp := test.HTTPDo(t, txPool.SetTxT(t, req, testdb.TxQs(t)))
	resp.EqualCode(t, http.StatusOK)
	jsonBody := test.StaticCopy(t, resp.Body.Bytes(), &db.SourceStructured{})
	fixture.CompareReadOrUpdate(t, testName+".json", jsonBody)

	req, err = http.NewRequest(http.MethodGet, server.URL+"/healthz", nil)
	require.NoError(err)
	resp = test.HTTPDo(t, req)
	resp.EqualCode(t, http.StatusOK)
	require.Equal(".", resp.Body.String())
}

func TestRoutes_NotFound(t *testing.T) {
	testName := "TestRoutes_NotFound"
	resp := test.HTTPDo(t, server.NewRequest(t, http.MethodGet, "/not_found_me", nil))
	resp.EqualCode(t, http.StatusNotFound)
	testIndent(t, resp, testName, "")
}

func TestRoutes_NotAllowed(t *testing.T) {
	testName := "TestRoutes_NotAllowed"
	resp := test.HTTPDo(t, server.NewRequest(t, http.MethodPost, "/terms/search", nil))
	resp.EqualCode(t, http.StatusMethodNotAllowed)
	testIndent(t, resp, testName, "")
}
