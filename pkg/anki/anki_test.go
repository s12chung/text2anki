package anki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config := DefaultConfig()
	dir, err := os.MkdirTemp("", "text2anki-TestMain-")
	if err != nil {
		log.Fatal(err)
	}
	config.NotesCacheDir = path.Join(dir, "files")
	if err := os.Mkdir(path.Join(dir, "files"), 0750); err != nil {
		log.Fatal(err)
	}
	SetConfig(config)

	exit := m.Run()

	config = DefaultConfig()
	SetConfig(config)
	os.Exit(exit)
}

// NewNoteFromTerm returns a note given the Term and index
func NewNoteFromTerm(term dictionary.Term, translationIndex uint) Note {
	translation := term.Translations[translationIndex]
	return Note{
		Text:         term.Text,
		PartOfSpeech: term.PartOfSpeech,
		Translation:  translation.Text,

		CommonLevel:      term.CommonLevel,
		Explanation:      translation.Explanation,
		DictionarySource: term.DictionarySource,
	}
}

func koreanBasicNotes(t *testing.T) []Note {
	require := require.New(t)

	var notes []Note
	bytes := fixture.Read(t, "notes.json")
	require.Nil(json.Unmarshal(bytes, &notes))

	return notes
}

func koreanBasicNotesWithSounds(t *testing.T) []Note {
	require := require.New(t)

	notes := koreanBasicNotes(t)
	sound := fixture.Read(t, "sound.mp3")
	for _, testIndex := range []uint{1, 3, 7} {
		err := notes[testIndex].SetSound(sound, fmt.Sprintf("Naver CLOVA Speech Synthesis - %v", testIndex))
		require.Nil(err)
	}
	return notes
}

func TestExportFiles(t *testing.T) {
	require := require.New(t)

	exportDir, err := os.MkdirTemp("", "text2anki-TestExportFiles-")
	require.Nil(err)
	err = ExportFiles(koreanBasicNotesWithSounds(t), exportDir)
	require.Nil(err)

	fixture.CompareOrUpdateDir(t, "ExportFiles", exportDir)
}

func TestExportSounds(t *testing.T) {
	require := require.New(t)

	exportDir, err := os.MkdirTemp("", "text2anki-TestExportSounds-")
	require.Nil(err)
	err = ExportSounds(koreanBasicNotesWithSounds(t), exportDir)
	require.Nil(err)

	dirEntries, err := os.ReadDir(exportDir)
	require.Nil(err)
	dirEntryNames := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		dirEntryNames[i] = dirEntry.Name()
	}

	require.Equal([]string{"t2a-가다.mp3", "t2a-가다가.mp3", "t2a-올라가다.mp3"}, dirEntryNames)
}

func TestExportCSVFile(t *testing.T) {
	require := require.New(t)

	dir, err := os.MkdirTemp("", "text2anki-TestExportCSVFile-")
	require.Nil(err)
	dir = path.Join(dir, "TestExportCSVFile.csv")

	err = ExportCSVFile(koreanBasicNotes(t), dir)
	require.Nil(err)
	//nolint:gosec // for tests
	bytes, err := ioutil.ReadFile(dir)
	require.Nil(err)
	fixture.CompareReadOrUpdate(t, "export_csv_expected.csv", bytes)
}

func TestExportCSV(t *testing.T) {
	require := require.New(t)

	buffer := &bytes.Buffer{}
	err := ExportCSV(koreanBasicNotes(t), buffer)
	require.Nil(err)
	fixture.CompareReadOrUpdate(t, "export_csv_expected.csv", buffer.Bytes())
}
