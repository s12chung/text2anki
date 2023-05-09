package firm

// Result is the result of a validation
type Result interface {
	IsValid() bool
	ErrorMap() ErrorMap
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
