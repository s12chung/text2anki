package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNotFoundError(t *testing.T) {
	require := require.New(t)
	require.True(IsNotFoundError(NotFoundError{}))
	require.False(IsNotFoundError(fmt.Errorf("test error")))
}

func TestIsInvalidInputError(t *testing.T) {
	require := require.New(t)
	require.True(IsInvalidInputError(InvalidInputError{}))
	require.False(IsInvalidInputError(fmt.Errorf("test error")))
}
