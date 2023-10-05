package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFoundRule_ValidateValue(t *testing.T) {
	require := require.New(t)
	errorMap := NotFoundRule{}.ValidateValue(reflect.ValueOf(1))
	require.NotNil(errorMap)
	require.NotEmpty(errorMap)
}
