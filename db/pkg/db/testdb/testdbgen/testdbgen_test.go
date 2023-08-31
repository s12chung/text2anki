package testdbgen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestGenerateModelsCodeRaw(t *testing.T) {
	require := require.New(t)
	testName := "TestGenerateModelsCodeRaw"

	code, err := generateModelsCodeRaw([]generateModelsCodeData{
		{Name: "Term", CreateCode: "qs.TermCreate(tx.Ctx(), term.CreateParams())"},
		{Name: "SourceStructured", CreateCode: "qs.SourceCreate(tx.Ctx(), sourceStructured.ToSourceCreateParams())"},
		{Name: "Note", CreateCode: "qs.NoteCreate(tx.Ctx(), note.CreateParams())"},
	})
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".go.txt", code)
}
