package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/models"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var sourcesServer test.Server

func init() {
	sourcesServer = server.WithPathPrefix("/sources")
}

func TestRoutes_SourceIndex(t *testing.T) {
	testName := "TestRoutes_SourceIndex"
	resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, "", nil))
	testModelsResponse(t, resp, testName, "", &[]db.SourceSerialized{})
}

func TestRoutes_SourceGet(t *testing.T) {
	testName := "TestRoutes_SourceGet"
	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "normal", path: "/1", expectedCode: http.StatusOK},
		{name: "invalid_id", path: "/9999", expectedCode: http.StatusNotFound},
		{name: "not_a_number", path: "/nan", expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, tc.path, nil))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &db.SourceSerialized{})
		})
	}
}

func TestRoutes_SourceUpdate(t *testing.T) {
	testName := "TestRoutes_SourceUpdate"

	testCases := []struct {
		name         string
		newName      string
		expectedCode int
	}{
		{name: "basic", newName: "new_name.txt", expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
	}

	created, err := db.Qs().SourceCreate(context.Background(), models.SourceSerializedsMust()[1].CreateParams())
	require.NoError(t, err)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			reqBody := test.JSON(t, SourceUpdateRequest{Name: tc.newName})
			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodPatch, idPath("", created.ID), bytes.NewReader(reqBody)))
			require.Equal(tc.expectedCode, resp.Code)

			sourceSerialized := db.SourceSerialized{}
			fixtureFile := testModelResponse(t, resp, testName, tc.name, &sourceSerialized)

			if resp.Code == http.StatusOK {
				source, err := db.Qs().SourceGet(context.Background(), sourceSerialized.ID)
				require.NoError(err)
				fixture.CompareRead(t, fixtureFile, fixture.JSON(t, source.ToSourceSerialized().StaticCopy()))
			}
		})
	}
}

func TestRoutes_SourceCreate(t *testing.T) {
	testName := "TestRoutes_SourceCreate"
	test.CISkip(t, "can't run C environment in CI")

	require.NoError(t, routes.Setup())
	defer func() {
		require.NoError(t, routes.Cleanup())
	}()

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{name: "split", expectedCode: http.StatusOK},
		{name: "weave", expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty_parts", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			body := sourcePostReqFromFile(t, testName, tc.name+".txt")
			if tc.name == "empty_parts" {
				body.Parts = []SourceCreateRequestPart{}
			}
			reqBody := test.JSON(t, body)
			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodPost, "", bytes.NewReader(reqBody)))
			require.Equal(tc.expectedCode, resp.Code)

			sourceSerialized := db.SourceSerialized{}
			fixtureFile := testModelResponse(t, resp, testName, tc.name, &sourceSerialized)

			if resp.Code == http.StatusOK {
				source, err := db.Qs().SourceGet(context.Background(), sourceSerialized.ID)
				require.NoError(err)
				fixture.CompareRead(t, fixtureFile, fixture.JSON(t, source.ToSourceSerialized().StaticCopy()))
			}
		})
	}
}

func sourcePostReqFromFile(t *testing.T, testName, name string) *SourceCreateRequest {
	s := string(test.Read(t, fixture.JoinTestData(testName, name)))
	split := strings.Split(s, "===")
	part := SourceCreateRequestPart{Text: s}
	if len(split) == 2 {
		part = SourceCreateRequestPart{Text: split[0], Translation: split[1]}
	}
	return &SourceCreateRequest{Parts: []SourceCreateRequestPart{part}}
}

func TestRoutes_SourceDestroy(t *testing.T) {
	testName := "TestRoutes_SourceDestroy"

	created, err := db.Qs().SourceCreate(context.Background(), models.SourceSerializedsMust()[1].CreateParams())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "normal", path: idPath("", created.ID), expectedCode: http.StatusOK},
		{name: "invalid_id", path: "/9999", expectedCode: http.StatusNotFound},
		{name: "not_a_number", path: "/nan", expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodDelete, tc.path, nil))
			require.Equal(tc.expectedCode, resp.Code)

			testModelResponse(t, resp, testName, tc.name, &db.SourceSerialized{})
			if resp.Code == http.StatusOK {
				_, err := db.Qs().SourceGet(context.Background(), created.ID)
				require.Equal(fmt.Errorf("sql: no rows in result set"), err)
			}
		})
	}
}

type signPrePartsResponse struct {
	SignPrePartsResponse
}

var testUUID = "123e4567-e89b-12d3-a456-426614174000"
var uuidRegexp = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)

func (s signPrePartsResponse) StaticCopy() any {
	a := s
	a.ID = testUUID
	for i, req := range s.Requests {
		u, err := url.Parse(req.URL)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		u.RawQuery = url.Values{localstore.CipherQueryParam: []string{"testy"}}.Encode()
		req.URL = uuidRegexp.ReplaceAllString(u.String(), testUUID)
		a.Requests[i] = req
	}
	return a
}

func TestRoutes_SignPreParts(t *testing.T) {
	testName := "TestRoutes_SignPreParts"

	testCases := []struct {
		name         string
		queryParams  string
		expectedCode int
	}{
		{name: "one", queryParams: "exts=.png", expectedCode: http.StatusOK},
		{name: "many", queryParams: "exts=.jpg&exts=.png&exts=.jpeg", expectedCode: http.StatusOK},
		{name: "array", queryParams: "exts[0]=.jpg&exts[1]=.png&exts[2]=.jpeg", expectedCode: http.StatusUnprocessableEntity},
		{name: "comma", queryParams: "exts=.jpeg,.png", expectedCode: http.StatusUnprocessableEntity},
		{name: "none", queryParams: "", expectedCode: http.StatusUnprocessableEntity},
		{name: "invalid", queryParams: "exts=.jpg&exts=.png&exts=.waka", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, "/sign_pre_parts?"+tc.queryParams, nil))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &signPrePartsResponse{})
		})
	}
}
