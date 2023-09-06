// Package models contains the models used by testdb
package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
)

const modelsDir = "modeldata"

var seederMap = map[string]seeder{}
var callerPath string

func init() {
	_, callerFilePath, _, ok := runtime.Caller(0)
	if !ok {
		slog.Error("runtime.Caller not ok for models package")
		os.Exit(-1)
	}
	callerPath = path.Dir(callerFilePath)
	setSeederMap(notes, sourceStructureds, terms)
}

// SeedList seeds the models for the testdb
func SeedList(txQs db.TxQs, list map[string]bool) error {
	isWhiteList, isBlacklist := true, true
	for k, v := range list {
		if v {
			isBlacklist = false
		} else {
			isWhiteList = false
			continue
		}
		s, exists := seederMap[k]
		if !exists {
			return fmt.Errorf("seedFunc for '%v' doesn't exist", k)
		}
		if err := s.Seed(txQs); err != nil {
			return err
		}
	}
	for k, s := range seederMap {
		if _, exists := list[k]; !(isBlacklist && isWhiteList) &&
			(isBlacklist && exists || isWhiteList && !exists) {
			continue
		}
		if err := s.Seed(txQs); err != nil {
			return err
		}
	}
	return nil
}

// Notes returns the Notes fixture
func Notes() Fixture[db.Note] { return notes }

// SourceStructureds returns the SourceStructureds fixture
func SourceStructureds() Fixture[db.SourceStructured] { return sourceStructureds }

// Terms returns the Terms fixture
func Terms() Fixture[db.Term] { return terms }

var notes = newFixture[db.Note]("Notes", func(txQs db.TxQs, model db.Note) error {
	_, err := txQs.NoteCreate(txQs.Ctx(), model.CreateParams())
	return err
})
var sourceStructureds = newFixture[db.SourceStructured]("SourceStructureds", func(txQs db.TxQs, model db.SourceStructured) error {
	_, err := txQs.SourceCreate(txQs.Ctx(), model.CreateParams())
	return err
})
var terms = newFixture[db.Term]("Terms", func(txQs db.TxQs, model db.Term) error {
	_, err := txQs.TermCreate(txQs.Ctx(), model.CreateParams())
	return err
})

// Fixture is the interface to the fixture data
type Fixture[T any] struct {
	name   string
	create func(txQs db.TxQs, model T) error
}

func newFixture[T any](name string, create func(txQs db.TxQs, model T) error) Fixture[T] {
	return Fixture[T]{name: name, create: create}
}

// Name returns the name of the fixtures
func (f Fixture[T]) Name() string { return f.name }

// Filename returns the filename of the fixture
func (f Fixture[T]) Filename() string { return f.name + "Seed.json" }

// Models returns the models of the fixture
func (f Fixture[T]) Models() ([]T, error) {
	var models []T
	if err := unmarshall(f.Filename(), &models); err != nil {
		return nil, err
	}
	return models, nil
}

// ModelsT returns the models of the fixture
func (f Fixture[T]) ModelsT(t *testing.T) []T {
	models, err := f.Models()
	require.NoError(t, err)
	return models
}

// Seed seeds the Models
func (f Fixture[T]) Seed(txQs db.TxQs) error {
	models, err := f.Models()
	if err != nil {
		return err
	}
	for _, model := range models {
		if err := f.create(txQs, model); err != nil {
			return err
		}
	}
	return nil
}

type seeder interface {
	Name() string
	Seed(txQs db.TxQs) error
}

func setSeederMap(seeders ...seeder) {
	for _, s := range seeders {
		seederMap[s.Name()] = s
	}
}

func unmarshall(filename string, models any) error {
	bytes, err := os.ReadFile(path.Join(callerPath, modelsDir, filename)) //nolint:gosec // for testing
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, models)
}
