package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestNote_StaticCopy(t *testing.T) {
	require := require.New(t)
	testName := "TestNote_StaticCopy"

	note := Note{}
	err := json.Unmarshal(fixture.Read(t, "NoteSrc.json"), &note)
	require.NoError(err)
	test.EmptyFieldsMatch(t, note, "Downloaded")

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, note.StaticCopy()))
}

func TestNote_CreateParams(t *testing.T) {
	require := require.New(t)
	testName := "TestNote_CreateParams"

	note := Note{}
	err := json.Unmarshal(fixture.Read(t, "NoteSrc.json"), &note)
	require.NoError(err)
	test.EmptyFieldsMatch(t, note, "Downloaded")

	createParams := note.CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}
