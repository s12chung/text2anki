package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test"
)

var termsServer txServer

func init() {
	termsServer = server.WithPathPrefix("/terms")
}

func TestRoutes_TermsSearch(t *testing.T) {
	testName := "TestRoutes_TermsSearch"
	testCases := []struct {
		name   string
		values url.Values
	}{
		{name: "normal", values: map[string][]string{"query": {testdb.SearchTerm}, "pos": {string(testdb.SearchPOS)}}},
		{name: "empty_pos", values: map[string][]string{"query": {testdb.SearchTerm}, "pos": {string(lang.PartOfSpeechUnknown)}}},
		{name: "bad_pos", values: map[string][]string{"query": {testdb.SearchTerm}, "pos": {"waka"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resp test.Response
			db.WithTermsSearchConfig(testdb.SearchConfig, func() {
				resp = test.HTTPDo(t, termsServer.NewRequest(t, http.MethodGet, "/search?"+tc.values.Encode(), nil))
			})
			testModelsResponse[dictionary.Term](t, resp, testName, tc.name, nil)
		})
	}
}
