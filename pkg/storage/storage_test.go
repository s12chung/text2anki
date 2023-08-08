package storage

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

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

func TestSigner_SignGetByID(t *testing.T) {
	require := require.New(t)
	testName := "TestSigner_SignGetByID"

	prePartList := PrePartList{}
	err := newTestSigner().SignGetByID("sources", "parts", testUUID, &prePartList)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, prePartList))

	err = newTestSigner().SignGetByID("sources", "parts", "some_bad_id", nil)
	require.Error(err)
}

type BasicTree struct {
	Value1Suffix string
	Key1         BasicTreeKey1
	Key2         []BasicTreeKey2
}

type BasicTreeKey1 struct {
	Value2Suffix string
}

type BasicTreeKey2 struct {
	Value3Suffix string
	Value4Suffix string
}

var basicKeys = []string{
	"Value1",
	"Key1.Value2",
	"Key2[0].Value3",
	"Key2[0].Value4",
}

var mixedKeys = []string{
	"Key1.SubKey1.Value1",
	"Key1.SubKey2[0][0].Value2",
	"Key1.SubKey2[0][0].Value3",
	"Key1.SubKey2[0][1].Value2",
	"Key1.SubKey3[0]",
	"Key1.SubKey3[1]",
}

var complexKeys = []string{
	"Key1[0][0]",
	"Key1[0][1]",
	"Key2[0][0].DeepKey1.DeepDeep1",
	"Key2[0][0].DeepKey2[0].SubDeep2[0]",
	"Key2[0][0].DeepKey2[0].SubDeep2[1]",
	"Key2[0][0].DeepKey2[1].SubDeep2[0]",
	"Key2[0][1].DeepKey1.DeepDeep2",
	"Key2[1][0].DeepKey1.DeepDeep1",
	"Key2[1][0].DeepKey2[0].SubDeep2[0]",
	"Key2[1][0].DeepKey2[0].SubDeep3",
}

var stringVSStructKeys = []string{
	"Key1",
	"Key1.Value1",
}
var stringVSAlphaStructKeys = []string{
	"Key1",
	"Key1.zipToAlphaEndToJpg",
}
var stringVSArrayKeys = []string{
	"Key1",
	"Key1[0]",
}
var arrayVsStructKeys = []string{
	"Key1.Value1",
	"Key1[0]",
}

func TestTreeFromKeys(t *testing.T) {
	testName := "TestTreeFromKeys"
	testCases := []struct {
		name string
		keys []string
		err  error
	}{
		{name: "basic", keys: basicKeys},
		{name: "mixed", keys: mixedKeys},
		{name: "complex", keys: complexKeys},
		{name: "string_vs_struct", keys: stringVSStructKeys,
			err: fmt.Errorf("at key: column_name.Key1.jpg, %w",
				fmt.Errorf("unmatched types string and map[string]interface {} at: Key1"))},
		{name: "string_vs_alpha_struct", keys: stringVSAlphaStructKeys,
			err: fmt.Errorf("at key: column_name.Key1.zipToAlphaEndToJpg.jpg, %w",
				fmt.Errorf("expected Map at: zipToAlphaEndToJpg"))},
		{name: "string_vs_array", keys: stringVSArrayKeys,
			err: fmt.Errorf("at key: column_name.Key1[0].jpg, %w",
				fmt.Errorf("expected Slice at: 0"))},
		{name: "array_vs_struct", keys: arrayVsStructKeys,
			err: fmt.Errorf("at key: column_name.Key1[0].jpg, %w",
				fmt.Errorf("expected Slice at: 0"))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			for i, key := range tc.keys {
				tc.keys[i] = "column_name." + key + ".jpg"
			}
			tree, err := treeFromKeys(tc.keys)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, tree))
		})
	}
}

func TestUnmarshallTree(t *testing.T) {
	require := require.New(t)
	testName := "TestUnmarshallTree"

	var tree map[string]any
	err := json.Unmarshal(fixture.Read(t, "TestTreeFromKeys/basic.json"), &tree)
	require.NoError(err)

	obj := BasicTree{}
	err = unmarshallTree(tree, reflect.ValueOf(&obj), "Suffix", func(key string) (string, error) {
		return key, nil
	})
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, obj))
}
