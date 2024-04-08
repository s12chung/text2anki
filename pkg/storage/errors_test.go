package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNotFoundError(t *testing.T) {
	require := require.New(t)
	require.True(IsNotFoundError(NotFoundError{}))
	require.False(IsNotFoundError(errors.New("test error")))
}

func TestIsInvalidInputError(t *testing.T) {
	require := require.New(t)
	require.True(IsInvalidInputError(InvalidInputError{}))
	require.False(IsInvalidInputError(errors.New("test error")))
}
