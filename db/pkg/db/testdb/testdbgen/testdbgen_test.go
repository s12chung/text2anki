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
		{Name: "Term", CreateCode: "queries.TermCreate(context.Background(), term.CreateParams())"},
		{Name: "SourceSerialized", CreateCode: "queries.SourceCreate(context.Background(), sourceSerialized.ToSourceCreateParams())"},
		{Name: "Note", CreateCode: "queries.NoteCreate(context.Background(), note.CreateParams())"},
	})
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".go.txt", code)
}
