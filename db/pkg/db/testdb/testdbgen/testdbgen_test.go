package testdbgen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestGenerateModelsCode(t *testing.T) {
	require := require.New(t)
	testName := "TestGenerateModelsCode"

	code, err := GenerateModelsCode()
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".go.txt", code)
}
