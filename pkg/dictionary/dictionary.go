// Package dictionary contains dictonary related functions
package dictionary

// Dicionary is a interface for a dictionary to search with
type Dicionary interface {
	Search(q string) ([]Term, error)
}

// CommonLevel indicates how common a term is
type CommonLevel uint8

// CommonLevels from 0 to 3, where 3 is the most common
const (
	CommonLevelUnique CommonLevel = iota
	CommonLevelRare
	CommonLevelMedium
	CommonLevelCommon
)

// Term is a word or phrase
type Term struct {
	Text         string
	CommonLevel  CommonLevel
	Translations []Translation
}

// Translation is a translation of a word or phrase
type Translation struct {
	Text        string
	Explanation string
}
