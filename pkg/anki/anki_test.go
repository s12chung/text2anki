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
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestMain(m *testing.M) {
	config, err := DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	dir := path.Join(os.TempDir(), test.GenerateName("anki.TestMain"))
	config.NotesCacheDir = path.Join(dir, "files")
	if err = os.MkdirAll(path.Join(dir, "files"), ioutil.OwnerRWXGroupRX); err != nil {
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

func koreanBasicNotesWithSounds(t *testing.T) []Note {
	require := require.New(t)

	notes := notesFromTerms(t)
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
	testName := "TestExportFiles"

	exportDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, exportDir)
	require.NoError(ExportFiles(koreanBasicNotesWithSounds(t), exportDir))

	fixture.CompareOrUpdateDir(t, "ExportFiles", exportDir)
}

func TestExportSounds(t *testing.T) {
	require := require.New(t)
	testName := "TestExportSounds"

	exportDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, exportDir)
	require.NoError(ExportSounds(koreanBasicNotesWithSounds(t), exportDir))

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
	testName := "TestExportCSVFile"

	dir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, dir)

	dir = path.Join(dir, "basic.csv")
	require.NoError(ExportCSVFile(notesFromTerms(t), dir))
	//nolint:gosec // for tests
	bytes, err := os.ReadFile(dir)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "export_csv_expected.csv", bytes)
}

func TestExportCSV(t *testing.T) {
	require := require.New(t)

	buffer := &bytes.Buffer{}
	err := ExportCSV(notesFromTerms(t), buffer)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "export_csv_expected.csv", buffer.Bytes())
}

func notesFromTerms(t *testing.T) []Note {
	require := require.New(t)

	var terms []dictionary.Term
	// from .../TestKoreanBasic_Search/basic.json
	err := json.Unmarshal(fixture.Read(t, "terms.json"), &terms)
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
	return notes
}
