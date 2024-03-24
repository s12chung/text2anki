package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func firstNote(t *testing.T, txQs TxQs) Note {
	require := require.New(t)
	note, err := txQs.NoteGet(txQs.Ctx(), 1)
	require.NoError(err)
	return note
}

func TestNote_StaticCopy(t *testing.T) {
	txQs := TxQsT(t, nil)

	note := firstNote(t, txQs)
	test.EmptyFieldsMatch(t, note, "Downloaded")
	test.EmptyFieldsMatch(t, note.StaticCopy(), "Downloaded", "ID", "UpdatedAt", "CreatedAt")
}

func TestNote_CreateParams(t *testing.T) {
	testName := "TestNote_CreateParams"
	txQs := TxQsT(t, nil)

	createParams := firstNote(t, txQs).CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdateJSON(t, testName, createParams)
}

func TestNote_Anki(t *testing.T) {
	require := require.New(t)
	testName := "TestNote_Anki"
	txQs := TxQsT(t, nil)

	ankiNote, err := firstNote(t, txQs).Anki()
	require.NoError(err)
	test.EmptyFieldsMatch(t, ankiNote, "usageSoundSource")
	fixture.CompareReadOrUpdateJSON(t, testName, ankiNote)

	note := firstNote(t, txQs)
	note.CommonLevel = -1
	_, err = note.Anki()
	require.Error(err)

	note = firstNote(t, txQs)
	note.PartOfSpeech = "not a pos"
	_, err = note.Anki()
	require.Error(err)
}

func TestAnkiNotes(t *testing.T) {
	require := require.New(t)
	testName := "TestAnkiNotes"
	txQs := TxQsT(t, nil)

	notes, err := txQs.NotesIndex(txQs.Ctx())
	require.NoError(err)
	ankiNotes, err := AnkiNotes(notes)
	require.NoError(err)

	for _, note := range ankiNotes {
		test.EmptyFieldsMatch(t, note, "usageSoundSource")
	}
	fixture.CompareReadOrUpdateJSON(t, testName, ankiNotes)
}

func TestQueries_NoteCreate(t *testing.T) {
	require := require.New(t)
	txQs := TxQsT(t, WriteOpts())

	note, err := txQs.NoteCreate(txQs.Ctx(), firstNote(t, txQs).CreateParams())
	require.NoError(err)
	testRecentTimestamps(t, note.CreatedAt, note.UpdatedAt)
}
