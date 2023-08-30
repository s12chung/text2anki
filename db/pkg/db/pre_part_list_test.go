package db

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrePartList_Info(t *testing.T) {
	require := require.New(t)
	testName := "TestPrePartList_Info"

	expectedInfo := PrePartInfo{Name: "info test name", Reference: "https://info.test.ref"}
	b, err := json.Marshal(expectedInfo)
	require.NoError(err)

	key := testName + "/my_key.txt"
	require.NoError(dbStorage.Put(key, bytes.NewReader(b)))

	info, err := PrePartList{InfoKey: key}.Info()
	require.NoError(err)
	require.Equal(expectedInfo, info)
}
