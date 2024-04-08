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

type BasicTree struct {
	//nolint:unused // used to test if notExported fields are ignored
	notExported bool
	Value1Key   string          `json:"value_1_key,omitempty"`
	EmptyKey    string          `json:"empty_key,omitempty"`
	Key1        BasicTreeKey1   `json:"key_1"`
	Key2        []BasicTreeKey2 `json:"key_2,omitempty"`
}

type BasicTreeKey1 struct {
	Value2Key string `json:"value_2_key,omitempty"`
}

type BasicTreeKey2 struct {
	Value3Key string `json:"value_3_key,omitempty"`
	Value4Key string `json:"value_4_key,omitempty"`
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
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, tc.name), tree)
		})
	}
}

func TestUnmarshallTree(t *testing.T) {
	require := require.New(t)
	testName := "TestUnmarshallTree"

	var tree map[string]any
	err := json.Unmarshal(fixture.Read(t, "TestTreeFromKeys/basic.json"), &tree)
	tree["Empty"] = ""
	require.NoError(err)

	obj := BasicTree{}
	err = unmarshallTree(tree, reflect.ValueOf(&obj), keySuffix, func(key string) (string, error) {
		return path.Join("testPrefix", key), nil
	})
	require.NoError(err)
	fixture.CompareReadOrUpdateJSON(t, testName, obj)
}

func TestMapTree(t *testing.T) {
	require := require.New(t)
	testName := "TestMapTree"

	basicTree := &BasicTree{}
	err := json.Unmarshal(fixture.Read(t, "TestUnmarshallTree.json"), basicTree)
	require.NoError(err)

	tree, err := mapTree(basicTree, keySuffix)
	require.NoError(err)
	fixture.CompareReadOrUpdateJSON(t, testName, tree)

	basicTree = nil
	tree, err = mapTree(basicTree, keySuffix)
	require.NoError(err)
	require.Equal(map[string]any{}, tree)
}
