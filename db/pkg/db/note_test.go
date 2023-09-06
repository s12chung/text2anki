package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func noteSrc(t *testing.T) Note {
	require := require.New(t)

	note := Note{}
	err := json.Unmarshal(fixture.Read(t, "NoteSrc.json"), &note)
	require.NoError(err)
	test.EmptyFieldsMatch(t, note, "Downloaded")
	return note
}

func TestNote_StaticCopy(t *testing.T) {
	require := require.New(t)

	note := noteSrc(t)
	noteCopy := note
	noteCopy.ID = 0
	require.Equal(noteCopy, note.StaticCopy())
}

func TestNote_CreateParams(t *testing.T) {
	testName := "TestNote_CreateParams"

	createParams := noteSrc(t).CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}
