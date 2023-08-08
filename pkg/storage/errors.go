package storage

import (
	"errors"
	"fmt"
)

// NotFoundError is the error returned if the id is not found
type NotFoundError struct {
	ID     string
	IDPath string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("the ID, %v, is not found in the storage for path: %v", n.ID, n.IDPath)
}

// IsNotFoundError returns true if the error is a NotFoundError
func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError
	ok := errors.As(err, &notFoundError)
	return ok
}

// InvalidInputError is the error returned that is an input-related error
type InvalidInputError struct {
	Message string
}

func (i InvalidInputError) Error() string {
	return i.Message
}

// IsInvalidInputError returns true if the error is a InvalidInputError
func IsInvalidInputError(err error) bool {
	var invalidInputError InvalidInputError
	ok := errors.As(err, &invalidInputError)
	return ok
}
