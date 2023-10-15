package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFoundRule_ValidateValue(t *testing.T) {
	require := require.New(t)
	errorMap := NotFoundRule{}.ValidateValue(reflect.ValueOf(1))
	require.Equal(NotFoundRule{}.ErrorMap(), errorMap)
	require.NotNil(errorMap)
	require.NotEmpty(errorMap)
	require.Equal("NotFound: value type, NoType, not found in Registry", errorMap.Error())
}
