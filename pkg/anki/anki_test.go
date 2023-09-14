package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/logg"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func init() {
	dir := path.Join(os.TempDir(), test.GenerateName("anki.TestMain"))
	c := Config{ExportPrefix: "t2a-", NotesCacheDir: dir}
	if err := os.MkdirAll(dir, ioutil.OwnerRWXGroupRX); err != nil {
		slog.Error("anki.init()", logg.Err(err)) //nolint:forbidigo // used in init only
		os.Exit(-1)
	}
	SetConfig(c)
}

func notesFromFixture(t *testing.T) []Note {
	require := require.New(t)

	bytes, err := os.ReadFile(fixture.JoinTestData("Notes.json"))
	require.NoError(err)
	var notes []Note
	require.NoError(json.Unmarshal(bytes, &notes))
	for _, note := range notes {
		test.EmptyFieldsMatch(t, note, "usageSoundSource")
	}
	return notes
}

func notesWithSounds(t *testing.T) []Note {
	require := require.New(t)

	notes := notesFromFixture(t)
	sound := fixture.Read(t, "sound.mp3")
	for i := range notes {
		if i%2 == 1 {
			continue
		}
		err := notes[i].SetSound(sound, fmt.Sprintf("Naver CLOVA Speech Synthesis - %v", i))
		require.NoError(err)
	}
	return notes
}

type soundFactory struct{}

func (s soundFactory) Name() string { return "soundFactory name" }
func (s soundFactory) Sound(_ context.Context, usage string) ([]byte, error) {
	return []byte(usage), nil
}

func TestSoundSetter_SetSound(t *testing.T) {
	require := require.New(t)

	soundSetter := NewSoundSetter(soundFactory{})
	notes := notesFromFixture(t)
	require.NoError(soundSetter.SetSound(context.Background(), notes))
	for _, note := range notes {
		require.Equal(soundSetter.soundFactory.Name(), note.usageSoundSource)
		require.Equal(note.Usage, string(test.Read(t, path.Join(config.NotesCacheDir, note.soundFilename()))))
	}
}

func TestNote_ID(t *testing.T) {
	require := require.New(t)
	require.Equal("어른-Flower Road-모자람 없이 주신 사랑이 과분하다 느낄 때쯤 난 어른이 됐죠", notesFromFixture(t)[0].ID())
}

func TestNote_SetSound(t *testing.T) {
	require := require.New(t)

	note := notesFromFixture(t)[0]
	require.Empty(note.usageSoundSource)

	soundSource := "the source"
	soundContents := []byte("my_test")
	require.NoError(note.SetSound(soundContents, soundSource))

	require.Equal(soundSource, note.usageSoundSource)
	require.Equal(soundContents, test.Read(t, path.Join(config.NotesCacheDir, note.soundFilename())))
}

func TestNote_CSV(t *testing.T) {
	testName := "TestNote_CSV"
	fixture.CompareReadOrUpdate(t, testName+".json", test.JSON(t, notesFromFixture(t)[0].CSV()))
}

func TestExportFiles(t *testing.T) {
	require := require.New(t)
	testName := "TestExportFiles"

	exportDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, exportDir)
	require.NoError(ExportFiles(exportDir, notesWithSounds(t)))

	fixture.CompareOrUpdateDir(t, "ExportFiles", exportDir)
}

func TestExportSounds(t *testing.T) {
	require := require.New(t)
	testName := "TestExportSounds"

	exportDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, exportDir)
	require.NoError(ExportSounds(exportDir, notesWithSounds(t)))

	dirEntries, err := os.ReadDir(exportDir)
	require.NoError(err)
	dirEntryNames := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		dirEntryNames[i] = dirEntry.Name()
	}
	require.Equal([]string{"t2a-꽃길만 걷게 해줄게요.mp3", "t2a-모자람 없이 주신 사랑이 과분하다 느낄 때쯤 난 어른이 됐죠.mp3"}, dirEntryNames)
	for _, entry := range dirEntries {
		require.Equal("sound.mp3 fake", string(test.Read(t, path.Join(exportDir, entry.Name()))))
	}
}

func TestExportCSVFile(t *testing.T) {
	require := require.New(t)
	testName := "TestExportCSVFile"

	dir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, dir)

	dir = path.Join(dir, "basic.csv")
	require.NoError(ExportCSVFile(dir, notesFromFixture(t)))
	fixture.CompareReadOrUpdate(t, testName+".csv", test.Read(t, dir))
}

func TestExportCSV(t *testing.T) {
	require := require.New(t)

	buffer := &bytes.Buffer{}
	err := ExportCSV(buffer, notesFromFixture(t))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestExportCSVFile.csv", buffer.Bytes())
}
