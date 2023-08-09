package storage

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func newTestDBStorage() DBStorage {
	return NewDBStorage(testAPI{}, UUIDTest{})
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

func TestDBStorage_SignPut(t *testing.T) {
	require := require.New(t)
	testName := "TestDBStorage_SignPut"

	req, err := newTestDBStorage().SignPut("test_table", "test_column", ".txt")
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, req))
}

func TestDBStorage_SignPutTree(t *testing.T) {
	testName := "TestDBStorage_SignPutTree"

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
			err: InvalidInputError{Message: "empty struct given for DBStorage.SignPutTree() at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts"}},
		{name: "empty_array", req: PrePartListSignRequest{PreParts: []PrePartSignRequest{}},
			err: InvalidInputError{
				Message: "empty slice or array given for DBStorage.SignPutTree() at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts"}},
		{name: "invalid", req: PrePartListSignRequest{PreParts: []PrePartSignRequest{{ImageExt: ".waka"}}},
			err: InvalidInputError{Message: "invalid extension, .waka, at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts[0].Image"}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			resp := PrePartListSignResponse{}

			err := newTestDBStorage().SignPutTree(basicConfig, tc.req, &resp)
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

type PrePartMediaList struct {
	ID       string            `json:"id"`
	PreParts []SourcePartMedia `json:"pre_parts"`
}

type SourcePartMedia struct {
	ImageKey string `json:"image_key,omitempty"`
	AudioKey string `json:"audio_key,omitempty"`
}

func TestDBStorage_KeyTree(t *testing.T) {
	require := require.New(t)
	testName := "TestDBStorage_KeyTree"

	prePartList := PrePartMediaList{}
	err := newTestDBStorage().KeyTree("sources", "parts", testUUID, &prePartList)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, prePartList))

	err = newTestDBStorage().SignGetTree("sources", "parts", "some_bad_id", nil)
	require.Error(err)
}

type PrePartList struct {
	ID       string    `json:"id"`
	PreParts []PrePart `json:"pre_parts"`
}

func (p PrePartList) StaticCopy() any {
	return p
}

type PrePart struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}

func TestDBStorage_SignGetTree(t *testing.T) {
	require := require.New(t)
	testName := "TestDBStorage_SignGetTree"

	prePartList := PrePartList{}
	err := newTestDBStorage().SignGetTree("sources", "parts", testUUID, &prePartList)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, prePartList))

	err = newTestDBStorage().SignGetTree("sources", "parts", "some_bad_id", nil)
	require.Error(err)
}
