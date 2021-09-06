// Package anki contains Anki export data and functions
package anki

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/s12chung/text2anki/pkg/lang"

	"github.com/s12chung/text2anki/pkg/iotools"
)

// Note is a Anki Note, which contains data to create cards from
type Note struct {
	Text string
	lang.PartOfSpeech
	Translation string

	lang.CommonLevel
	Explanation      string
	Usage            string
	UsageTranslation string
	DictionarySource string

	hasSound    bool
	soundSource string

	Notes string
}

// Valid returns true when the Note is valid
func (n *Note) Valid() bool {
	return n.Text != "" && n.PartOfSpeech != lang.PartOfSpeechInvalid && n.Translation != ""
}

// SetSound sets the sound for the note
func (n *Note) SetSound(sound []byte, soundSource string) error {
	err := ioutil.WriteFile(path.Join(config.NotesCacheDir, n.soundFilename()), sound, 0600)
	if err != nil {
		return err
	}
	n.hasSound = true
	n.soundSource = soundSource
	return nil
}

// CSV returns the CSV representation of the Note
func (n *Note) CSV() []string {
	soundAnkiFormat := ""
	if n.hasSound {
		soundAnkiFormat = "[sound:" + n.soundFilename() + "]"
	}
	return []string{
		n.Text,
		string(n.PartOfSpeech),
		n.Translation,

		strconv.FormatUint(uint64(n.CommonLevel), 10),
		n.Explanation,
		n.Usage,
		n.UsageTranslation,
		n.DictionarySource,

		soundAnkiFormat,
		n.soundSource,

		n.Notes,
	}
}

var invalidFilenameRegex = regexp.MustCompile("/\\\\")

func (n *Note) soundFilename() string {
	return config.ExportPrefix + invalidFilenameRegex.ReplaceAllString(n.Usage, "") + ".mp3"
}

// ExportFiles exports all files into the given dst
func ExportFiles(notes []Note, dst string) error {
	if err := ExportCSVFile(notes, path.Join(dst, "text2anki.csv")); err != nil {
		return err
	}

	if err := os.Mkdir(path.Join(dst, "files"), 0750); err != nil {
		return nil
	}
	if err := ExportSounds(notes, path.Join(dst, "files")); err != nil {
		return err
	}
	return nil
}

// ExportSounds exports the sounds of the notes given the dst
func ExportSounds(notes []Note, dst string) error {
	for _, note := range notes {
		if !note.hasSound {
			continue
		}

		src := path.Join(config.NotesCacheDir, note.soundFilename())
		if err := iotools.CopyFile(path.Join(dst, note.soundFilename()), src, 0400); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}
	}
	return nil
}

// ExportCSVFile exports the Note CSV as a file
func ExportCSVFile(notes []Note, dst string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	err = ExportCSV(notes, f)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

// ExportCSV exports the notes as CSV
func ExportCSV(notes []Note, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, note := range notes {
		if err := writer.Write(note.CSV()); err != nil {
			return err
		}
	}
	return nil
}
