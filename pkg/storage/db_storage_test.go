package storage

import (
	"fmt"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func newTestDBStorage() DBStorage {
	return NewDBStorage(testAPI{}, uuidTest{})
}

func TestBaseKey(t *testing.T) {
	require := require.New(t)
	require.Equal("my_table/the_column/123e4567-e89b-12d3-a456-426614174000/the_column", BaseKey("my_table", "the_column", testUUID))
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
		{name: "no_pointer", err: fmt.Errorf("storage.PrePartListSignResponse is not a pointer")},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			resp := PrePartListSignResponse{}
			respAny := any(&resp)
			if tc.name == "no_pointer" {
				respAny = resp
			}

			err := newTestDBStorage().SignPutTree(basicConfig, tc.req, respAny)
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
	testName := "TestDBStorage_SignGetTree"

	prePartList := PrePartList{}
	testCases := []struct {
		name       string
		signedTree any
		id         string
		err        error
	}{
		{name: "basic", signedTree: &prePartList},
		{name: "non_pointer", signedTree: prePartList, err: fmt.Errorf("storage.PrePartList is not a pointer")},
		{name: "nil", signedTree: nil, err: fmt.Errorf("passed nil as settable obj")},
		{name: "bad_id", signedTree: &prePartList, id: "bad_id", err: NotFoundError{ID: "bad_id", IDPath: "sources/parts/bad_id"}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			id := tc.id
			if id == "" {
				id = testUUID
			}
			require := require.New(t)
			err := newTestDBStorage().SignGetTree("sources", "parts", id, tc.signedTree)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, prePartList))
		})
	}
}

type SourcePartMediaResponse struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}

func TestDBStorage_SignGetTreeFromKeyTree(t *testing.T) {
	require := require.New(t)
	testName := "TestDBStorage_SignGetTreeFromKeyTree"

	media := SourcePartMedia{ImageKey: "waka.jpg", AudioKey: "haha.mp3"}
	signedTree := SourcePartMediaResponse{}
	err := newTestDBStorage().SignGetTreeFromKeyTree(media, &signedTree)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, signedTree))

	err = newTestDBStorage().SignGetTreeFromKeyTree(media, signedTree)
	require.Equal(fmt.Errorf("signedTree, storage.SourcePartMediaResponse, is not a pointer"), err)
}

func TestDBStorage_KeyTreeFromSignGetTree(t *testing.T) {
	require := require.New(t)
	testName := "TestDBStorage_KeyTreeFromSignGetTree"

	signedTree := SourcePartMediaResponse{ImageURL: keyURL("haha.jpg"), AudioURL: keyURL("me.mp3")}
	media := SourcePartMedia{}

	err := newTestDBStorage().KeyTreeFromSignGetTree(signedTree, &media)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, media))

	err = newTestDBStorage().KeyTreeFromSignGetTree(signedTree, media)
	require.Equal(fmt.Errorf("keyTree, storage.SourcePartMedia, is not a pointer"), err)
}
