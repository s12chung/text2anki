// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"
)

type Note struct {
	ID               int64  `json:"id"`
	Text             string `json:"text"`
	PartOfSpeech     string `json:"part_of_speech"`
	Translation      string `json:"translation"`
	Explanation      string `json:"explanation"`
	CommonLevel      int64  `json:"common_level"`
	Usage            string `json:"usage"`
	UsageTranslation string `json:"usage_translation"`
	DictionarySource string `json:"dictionary_source"`
	Notes            string `json:"notes"`
	Downloaded       bool   `json:"downloaded"`
}

type Source struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Reference string    `json:"reference"`
	Parts     string    `json:"parts"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Term struct {
	ID           int64  `json:"id"`
	Text         string `json:"text"`
	PartOfSpeech string `json:"part_of_speech"`
	Variants     string `json:"variants"`
	Translations string `json:"translations"`
	CommonLevel  int64  `json:"common_level"`
	Popularity   int64  `json:"popularity"`
}
