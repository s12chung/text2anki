package storage

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func newTestDBStorage() DBStorage {
	return NewDBStorage(newTestAPI(), uuidTest{})
}

func TestSignPutConfig_KeyFor(t *testing.T) {
	require := require.New(t)

	key := SignPutConfig{
		Table:  "my_table",
		Column: "sources",
	}.KeyFor(testUUID, "output.json")
	require.Equal("my_table/sources/123e4567-e89b-12d3-a456-426614174000/output.json", key)
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

var basicConfig = SignPutConfig{
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

func TestDBStorage_SignPutTree(t *testing.T) {
	testName := "TestDBStorage_SignPutTree"

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
			err: InvalidInputError{Message: "srcTree empty struct given at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts"}},
		{name: "empty_array", req: PrePartListSignRequest{PreParts: []PrePartSignRequest{}},
			err: InvalidInputError{
				Message: "srcTree empty slice or array given at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts"}},
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

type PrePartListFileTree struct {
	PreParts []PrePartFileTree `json:"pre_parts"`
}

type PrePartFileTree struct {
	ImageFile fs.File `json:"image_file,omitempty"`
	AudioFile fs.File `json:"audio_file,omitempty"`
}

type PrePartListKeyTree struct {
	ID       string           `json:"id"`
	PreParts []PrePartKeyTree `json:"pre_parts"`
}

type PrePartKeyTree struct {
	ImageKey string `json:"image_key,omitempty"`
	AudioKey string `json:"audio_key,omitempty"`
}

func TestDBStorage_PutTree(t *testing.T) {
	testName := "TestDBStorage_PutTree"

	testCases := []struct {
		name      string
		partCount int
		err       error
	}{
		{name: "basic"},
		{name: "many", partCount: 3},
		{name: "empty",
			err: InvalidInputError{Message: "srcTree empty struct given at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts[0]"}},
		{name: "empty_parts",
			err: InvalidInputError{Message: "srcTree empty slice or array given at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts"}},
		{name: "invalid",
			err: InvalidInputError{Message: "invalid extension, .txt, at sources/parts/123e4567-e89b-12d3-a456-426614174000/parts.PreParts[0].Image"}},
		{name: "no_pointer", err: fmt.Errorf("storage.PrePartListKeyTree is not a pointer")},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			fileTree := PrePartListFileTree{PreParts: fileTreeParts(t, tc.name, testName, tc.partCount)}
			keyTree := PrePartListKeyTree{}
			keyTreeAny := any(&keyTree)
			if tc.name == "no_pointer" {
				keyTreeAny = keyTree
			}

			api := newTestAPI()
			err := NewDBStorage(api, uuidTest{}).PutTree(basicConfig, fileTree, keyTreeAny)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, keyTree))
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+"_storeMap.json"), fixture.JSON(t, api.storeMap))
		})
	}
}

func fileTreeParts(t *testing.T, caseName, testName string, partCount int) []PrePartFileTree {
	if partCount == 0 {
		partCount = 1
	}
	parts := make([]PrePartFileTree, partCount)
	switch caseName {
	case "empty":
		parts[0] = PrePartFileTree{}
	case "empty_parts":
		parts = []PrePartFileTree{}
	default:
		for i := 0; i < partCount; i++ {
			parts[i] = fileTreePartsFromFile(t, testName, caseName+strconv.Itoa(i))
		}
	}
	return parts
}

func fileTreePartsFromFile(t *testing.T, testName, name string) PrePartFileTree {
	require := require.New(t)

	basePath := fixture.JoinTestData(testName, name)
	part := PrePartFileTree{
		ImageFile: fileFromFile(t, basePath+"_image"),
		AudioFile: fileFromFile(t, basePath+"_audio"),
	}
	require.NotEmpty(part, "PrePartFileTree is empty")
	return part
}

func fileFromFile(t *testing.T, glob string) fs.File {
	require := require.New(t)

	files, err := filepath.Glob(glob + "*")
	require.NoError(err)
	require.LessOrEqual(len(files), 1, "found more than 1 file with glob: %v", glob)

	if len(files) == 0 {
		return nil
	}
	file, err := os.Open(files[0])
	require.NoError(err)
	return file
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
