// Package anki contains Anki export data and functions
package anki

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/stringutil"
)

var config Config

// SetConfig sets the Export and Cache config
func SetConfig(c Config) { config = c }

// Config contains Export and Cache config for Anki
type Config struct {
	ExportPrefix  string
	NotesCacheDir string
}

// SoundFactory generates sounds for SoundSetter
type SoundFactory interface {
	Name() string
	Sound(ctx context.Context, usage string) ([]byte, error)
}

// SoundSetter sets sounds for notes
type SoundSetter struct{ soundFactory SoundFactory }

// NewSoundSetter returns a new SoundSetter
func NewSoundSetter(soundFactory SoundFactory) SoundSetter {
	return SoundSetter{soundFactory: soundFactory}
}

// SetSound sets the sound from the soundFactory
func (s SoundSetter) SetSound(ctx context.Context, notes []Note) error {
	if s.soundFactory == nil {
		return nil
	}
	soundFactoryName := s.soundFactory.Name()
	for i, note := range notes {
		sound, err := s.soundFactory.Sound(ctx, note.Usage)
		if err != nil {
			return err
		}
		if err = notes[i].SetSound(sound, soundFactoryName); err != nil {
			return err
		}
	}
	return nil
}

// Note is an Anki Note, which contains data to create cards from
type Note struct {
	Text              string `json:"text"`
	lang.PartOfSpeech `json:"part_of_speech"`
	Translation       string `json:"translation"`
	Explanation       string `json:"explanation"`
	lang.CommonLevel  `json:"common_level"`

	Usage            string `json:"usage"`
	UsageTranslation string `json:"usage_translation"`

	SourceName       string `json:"source_name"`
	SourceReference  string `json:"source_reference"`
	DictionarySource string `json:"dictionary_source"`
	Notes            string `json:"notes"`

	usageSoundSource string
}

// ID returns the unique id for the note in anki
func (n *Note) ID() string {
	return strings.Join([]string{n.Text, n.SourceName, n.Usage}, "-")
}

// SetSound sets the sound for the note
func (n *Note) SetSound(sound []byte, soundSource string) error {
	err := os.WriteFile(path.Join(config.NotesCacheDir, n.UsageSoundFilename()), sound, ioutil.OwnerRWGroupR)
	if err != nil {
		return err
	}
	n.usageSoundSource = soundSource
	return nil
}

// CSV returns the CSV representation of the Note
func (n *Note) CSV() []string {
	return []string{
		n.ID(),

		n.Text,
		string(n.PartOfSpeech),
		n.Translation,
		n.Explanation,
		strconv.FormatUint(uint64(n.CommonLevel), 10),

		n.Usage,
		n.UsageTranslation,
		n.usageSoundAnkiFormat(),
		n.usageSoundSource,

		n.SourceName,
		n.SourceReference,
		n.DictionarySource,
		n.Notes,
	}
}

// from: https://github.com/ankitects/anki/blob/60748657635cb31c26f9a3deabf7926288471505/rslib/src/media/files.rs#L70
var invalidFilenameRegex = regexp.MustCompile(`/\[|]|<|>|:|"|/|\?|\*|\^|\\|\|`)

const ankiMediaMaxFilenameBytes = 119

// UsageSoundFilename returns the usage sound filename
//
// Inspired by:
// https://github.com/ankitects/anki/blob/60748657635cb31c26f9a3deabf7926288471505/rslib/src/media/files.rs#L248
func (n *Note) UsageSoundFilename() string {
	ext := ".mp3"
	usage := stringutil.TrimBytes(config.ExportPrefix+invalidFilenameRegex.ReplaceAllString(n.Usage, ""), ankiMediaMaxFilenameBytes-len(ext))
	return usage + ext
}

func (n *Note) usageSoundAnkiFormat() string {
	if n.HasUsageSound() {
		return "[sound:" + n.UsageSoundFilename() + "]"
	}
	return ""
}

// HasUsageSound returns true of the Note has a usage sound
func (n *Note) HasUsageSound() bool { return n.usageSoundSource != "" }

// ExportFiles exports all files into the given dst
func ExportFiles(dirPath string, notes []Note) error {
	if err := os.MkdirAll(dirPath, ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	if err := ExportCSVFile(path.Join(dirPath, "text2anki.csv"), notes); err != nil {
		return err
	}

	soundsPath := path.Join(dirPath, "files")
	if err := os.MkdirAll(soundsPath, ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	if err := ExportSounds(soundsPath, notes); err != nil {
		return err
	}
	return nil
}

// ExportSounds exports the sounds of the notes given the dst
func ExportSounds(dirPath string, notes []Note) error {
	for _, note := range notes {
		if !note.HasUsageSound() {
			continue
		}

		src := path.Join(config.NotesCacheDir, note.UsageSoundFilename())
		if err := ioutil.CopyFile(path.Join(dirPath, note.UsageSoundFilename()), src, ioutil.OwnerGroupR); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}
	}
	return nil
}

// ExportCSVFile exports the Note CSV as a file
func ExportCSVFile(filePath string, notes []Note) error {
	f, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	err = ExportCSV(f, notes)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

var csvComments = []string{"#separator:Pipe", "#html:false", "#notetype:Text2Anki", "#deck:Text2Anki"}

// ExportCSV exports the notes as CSV
func ExportCSV(w io.Writer, notes []Note) error {
	for _, comment := range csvComments {
		if _, err := fmt.Fprintf(w, comment+"\n"); err != nil {
			return err
		}
	}
	writer := csv.NewWriter(w)
	writer.Comma = '|'
	for _, note := range notes {
		if err := writer.Write(note.CSV()); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
