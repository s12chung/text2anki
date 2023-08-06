package storage

import (
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testUUID = "123e4567-e89b-12d3-a456-426614174000"

type UUIDTest struct {
}

func (u UUIDTest) Generate() (string, error) {
	return testUUID, nil
}

type testAPI struct {
}

func keyURL(key string) string {
	return "http://localhost:3000/" + key
}

func (t testAPI) SignPut(key string) (PreSignedHTTPRequest, error) {
	return PreSignedHTTPRequest{
		URL:          keyURL(key) + "?cipher=blah",
		Method:       "PUT",
		SignedHeader: http.Header{},
	}, nil
}

func (t testAPI) SignGet(key string) (string, error) {
	return keyURL(key), nil
}

func (t testAPI) ListKeys(prefix string) ([]string, error) {
	return []string{path.Join(prefix, "a.txt"), path.Join(prefix, "b.txt")}, nil
}

func newTestSigner() Signer {
	return NewSigner(testAPI{}, UUIDTest{})
}

type PrePartListSignRequest struct {
	PreParts []PrePartSignRequest `json:"pre_parts"`
}

type PrePartSignRequest struct {
	ImageExt string `json:"image_ext,omitempty"`
	AudioExt string `json:"audio_ext,omitempty"`
}

type PrePartListSignResponse struct {
	ID       string                `json:"id"`
	PreParts []PrePartSignResponse `json:"pre_parts"`
}

type PrePartSignResponse struct {
	ImageRequest *PreSignedHTTPRequest `json:"image_request,omitempty"`
	AudioRequest *PreSignedHTTPRequest `json:"audio_request,omitempty"`
}

func TestIsExtTreeError(t *testing.T) {
	require := require.New(t)
	require.True(IsInvalidInputError(InvalidInputError{}))
	require.False(IsInvalidInputError(fmt.Errorf("test error")))
}

func TestSigner_SignPut(t *testing.T) {
	require := require.New(t)
	testName := "TestSigner_SignPut"

	req, err := newTestSigner().SignPut("test_table", "test_column", ".txt")
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, req))
}

func TestSigner_SignPutTree(t *testing.T) {
	testName := "TestSigner_SignPutTree"

	basicConfig := SignPutConfig{
		Table:  "sources",
		Column: "parts",
		NameToValidExts: map[string]map[string]bool{
			"Image": {
				".jpg":  true,
				".jpeg": true,
				".png":  true,
			},
			"Audio": {
				".mp3": true,
			},
		},
	}

	testCases := []struct {
		name string
		req  PrePartListSignRequest
		err  error
	}{
		{name: "one", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg"}}},
		},
		{name: "many", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg"}, {ImageExt: ".png"}, {ImageExt: ".jpeg"}}},
		},
		{name: "mixed", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg", AudioExt: ".mp3"}, {AudioExt: ".mp3"}, {ImageExt: ".jpeg"}}},
		},
		{name: "empty", req: PrePartListSignRequest{},
			err: InvalidInputError{Message: "empty struct given for Signer.SignPutTree() at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts"}},
		{name: "empty_array", req: PrePartListSignRequest{PreParts: []PrePartSignRequest{}},
			err: InvalidInputError{
				Message: "empty slice or array given for Signer.SignPutTree() at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts"}},
		{name: "invalid", req: PrePartListSignRequest{PreParts: []PrePartSignRequest{{ImageExt: ".waka"}}},
			err: InvalidInputError{Message: "invalid extension, .waka, at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts[0].Image"}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			resp := PrePartListSignResponse{}

			err := newTestSigner().SignPutTree(basicConfig, tc.req, &resp)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			require.NotEmpty(resp.ID)
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, resp))
		})
	}
}

func TestSigner_SignGetByID(t *testing.T) {
	require := require.New(t)
	urls, err := newTestSigner().SignGetByID("sources", "parts", testUUID)
	require.NoError(err)

	prefix := "http://localhost:3000/sources/parts/123e4567-e89b-12d3-a456-426614174000/"
	require.Equal([]string{prefix + "a.txt", prefix + "b.txt"}, urls)
}
