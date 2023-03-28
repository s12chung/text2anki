package anki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/test/fixture"
)

func TestMain(m *testing.M) {
	config, err := DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.MkdirTemp("", "text2anki-anki.TestMain-")
	if err != nil {
		log.Fatal(err)
	}
	config.NotesCacheDir = path.Join(dir, "files")
	if err = os.Mkdir(path.Join(dir, "files"), 0750); err != nil {
		log.Fatal(err)
	}
	SetConfig(config)

	exit := m.Run()

	config, err = DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	SetConfig(config)
	os.Exit(exit)
}

const notesFixture = "notes.json"

func koreanBasicNotes(t *testing.T) []Note {
	require := require.New(t)

	var notes []Note
	if fixture.WillUpdate() {
		fixture.Update(t, notesFixture, dictionaryNotes(t))
	}
	require.NoError(json.Unmarshal(fixture.Read(t, notesFixture), &notes))

	return notes
}

func koreanBasicNotesWithSounds(t *testing.T) []Note {
	require := require.New(t)

	notes := koreanBasicNotes(t)
	sound := fixture.Read(t, "sound.mp3")
	for i, note := range notes {
		if note.Usage != "" {
			err := notes[i].SetSound(sound, fmt.Sprintf("Naver CLOVA Speech Synthesis - %v", i))
			require.NoError(err)
		}
	}
	return notes
}

func TestExportFiles(t *testing.T) {
	require := require.New(t)

	exportDir, err := os.MkdirTemp("", "text2anki-TestExportFiles-")
	require.NoError(err)
	err = ExportFiles(koreanBasicNotesWithSounds(t), exportDir)
	require.NoError(err)

	fixture.CompareOrUpdateDir(t, "ExportFiles", exportDir)
}

func TestExportSounds(t *testing.T) {
	require := require.New(t)

	exportDir, err := os.MkdirTemp("", "text2anki-TestExportSounds-")
	require.NoError(err)
	err = ExportSounds(koreanBasicNotesWithSounds(t), exportDir)
	require.NoError(err)

	dirEntries, err := os.ReadDir(exportDir)
	require.NoError(err)
	dirEntryNames := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		dirEntryNames[i] = dirEntry.Name()
	}

	require.Equal([]string{"t2a-소풍: usage0.mp3", "t2a-소풍: usage2.mp3", "t2a-소풍: usage4.mp3"}, dirEntryNames)
}

func TestExportCSVFile(t *testing.T) {
	require := require.New(t)

	dir, err := os.MkdirTemp("", "text2anki-TestExportCSVFile-")
	require.NoError(err)
	dir = path.Join(dir, "TestExportCSVFile.csv")

	err = ExportCSVFile(koreanBasicNotes(t), dir)
	require.NoError(err)
	//nolint:gosec // for tests
	bytes, err := os.ReadFile(dir)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "export_csv_expected.csv", bytes)
}

func TestExportCSV(t *testing.T) {
	require := require.New(t)

	buffer := &bytes.Buffer{}
	err := ExportCSV(koreanBasicNotes(t), buffer)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "export_csv_expected.csv", buffer.Bytes())
}

func dictionaryNotes(t *testing.T) []byte {
	require := require.New(t)

	sourcePath := path.Join("..", "dictionary", "koreanbasic", fixture.TestDataDir, "search_expected.json")
	//nolint:gosec // for tests
	sourceBytes, err := os.ReadFile(sourcePath)
	require.NoError(err)

	var terms []dictionary.Term
	err = json.Unmarshal(sourceBytes, &terms)
	require.NoError(err)

	notes := make([]Note, len(terms))
	for i, term := range terms {
		notes[i] = NewNoteFromTerm(term, 0)
	}
	for _, testIndex := range []uint{0, 2, 4} {
		notes[testIndex].Usage = fmt.Sprintf("소풍: /\\usage%v", testIndex)
	}
	for _, testIndex := range []uint{0, 2, 4} {
		notes[testIndex].UsageTranslation = fmt.Sprintf("Test usage translation, index: %v", testIndex)
	}

	bytes, err := json.MarshalIndent(notes, "", "  ")
	require.NoError(err)
	return bytes
}
