package api

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx/reqtxtest"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func joinPath(elem ...any) string {
	return fmt.Sprintf(strings.Repeat("/%v", len(elem)), elem...)
}

func testIndent(t *testing.T, resp test.Response, testName, name string) {
	jsonBody := test.IndentJSON(t, resp.Body.Bytes())
	fixture.CompareReadOrUpdate(t, fixtureFileName(testName, name), jsonBody)
}

func testModelResponse[T test.StaticCopyable[T]](t *testing.T, resp test.Response, testName, name string, model *T) string {
	jsonBody := test.StaticCopyOrIndent[T](t, resp.Code, resp.Body.Bytes(), model)
	fixtureFile := fixtureFileName(testName, name)
	fixture.CompareReadOrUpdate(t, fixtureFile, jsonBody)
	return fixtureFile
}

func testModelsResponse[T test.StaticCopyable[T]](t *testing.T, resp test.Response, testName, name string, models *[]T) {
	jsonBody := test.StaticCopyOrIndentSlice[T](t, resp.Code, resp.Body.Bytes(), models)
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
	pool reqtxtest.Pool[db.TxQs, config.TxMode]
	test.Server
}

func (s txServer) NewRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	return s.newTxRequest(t, testdb.TxQs(t, nil), txReadOnly, method, path, body)
}
func (s txServer) NewTxRequest(t *testing.T, tx db.TxQs, method, path string, body io.Reader) *http.Request {
	return s.newTxRequest(t, tx, txWritable, method, path, body)
}

func (s txServer) newTxRequest(t *testing.T, tx db.TxQs, mode config.TxMode, method, path string, body io.Reader) *http.Request {
	req := s.Server.NewRequest(t, tx.Ctx(), method, path, body)
	s.pool.SetTx(t, req, tx, mode)
	return req
}

func (s txServer) WithPathPrefix(prefix string) txServer {
	dup := s
	dup.Server = s.Server.WithPathPrefix(prefix, plog)
	return dup
}

type crcTranslator struct{}

func (c crcTranslator) Translate(_ context.Context, s string) (string, error) {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("crc-%x", crc32.ChecksumIEEE([]byte(line)))
	}
	return strings.Join(lines, "\n"), nil
}
