package api

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"testing"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx/reqtxtest"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func idPath(path string, id int64) string {
	return fmt.Sprintf(path+"/%v", id)
}

func testIndent(t *testing.T, resp test.Response, testName, name string) {
	jsonBody := test.IndentJSON(t, resp.Body.Bytes())
	fixture.CompareReadOrUpdate(t, fixtureFileName(testName, name), jsonBody)
}

func testModelResponse(t *testing.T, resp test.Response, testName, name string, model test.StaticCopyable) string {
	jsonBody := test.StaticCopyOrIndent(t, resp.Code, resp.Body.Bytes(), model)
	fixtureFile := fixtureFileName(testName, name)
	fixture.CompareReadOrUpdate(t, fixtureFile, jsonBody)
	return fixtureFile
}

func testModelsResponse(t *testing.T, resp test.Response, testName, name string, models any) {
	jsonBody := test.StaticCopyOrIndentSlice(t, resp.Code, resp.Body.Bytes(), models)
	fixture.CompareReadOrUpdate(t, fixtureFileName(testName, name), jsonBody)
}

func fixtureFileName(testName, name string) string {
	fixtureFile := testName + ".json"
	if name != "" {
		fixtureFile = path.Join(testName, name+"_response.json")
	}
	return fixtureFile
}

type txServer struct {
	pool reqtxtest.Pool[db.TxQs]
	test.Server
}

func (s txServer) NewRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	return s.NewTxRequest(t, testdb.TxQs(t, nil), method, path, body)
}

func (s txServer) NewTxRequest(t *testing.T, tx db.TxQs, method, path string, body io.Reader) *http.Request {
	req := s.Server.NewRequest(t, tx.Ctx(), method, path, body)
	s.pool.SetTxT(t, req, tx)
	return req
}

func (s txServer) WithPathPrefix(prefix string) txServer {
	dup := s
	dup.Server = s.Server.WithPathPrefix(prefix)
	return dup
}
