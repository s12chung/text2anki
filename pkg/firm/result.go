package firm

import (
	"golang.org/x/exp/maps"
)

// Result is the result of a validation
type Result interface {
	IsValid() bool
	ErrorMap() ErrorMap
	Errors() []*TemplatedError
	Error() *TemplatedError
}

// MapResult is a result with a map of errors
type MapResult struct {
	errorMap ErrorMap
}

// IsValid returns true if the data is valid from the validator validation
func (s MapResult) IsValid() bool {
	return len(s.errorMap) == 0
}

// ErrorMap returns the map of errors from the validator validation
func (s MapResult) ErrorMap() ErrorMap {
	return s.errorMap
}

// Errors returns an array of errors from the validator validation
func (s MapResult) Errors() []*TemplatedError {
	return maps.Values(s.ErrorMap())
}

// Error the first error (if any) from the validator validation
func (s MapResult) Error() *TemplatedError {
	for _, v := range s.errorMap {
		return v
	}
	return nil
}
