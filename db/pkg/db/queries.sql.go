// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: queries.sql

package db

import (
	"context"
)

const noteCreate = `-- name: NoteCreate :one
INSERT INTO notes (
    text, part_of_speech, translation,
                   common_level, explanation, usage, usage_translation,
                   source_name, source_reference, dictionary_source, notes
) VALUES (?, ?, ?,
          ?, ?, ?, ?,
          ?, ?, ?, ?) RETURNING id, text, part_of_speech, translation, explanation, common_level, usage, usage_translation, source_name, source_reference, dictionary_source, notes, downloaded
`

type NoteCreateParams struct {
	Text             string `json:"text"`
	PartOfSpeech     string `json:"part_of_speech"`
	Translation      string `json:"translation"`
	CommonLevel      int64  `json:"common_level"`
	Explanation      string `json:"explanation"`
	Usage            string `json:"usage"`
	UsageTranslation string `json:"usage_translation"`
	SourceName       string `json:"source_name"`
	SourceReference  string `json:"source_reference"`
	DictionarySource string `json:"dictionary_source"`
	Notes            string `json:"notes"`
}

func (q *Queries) NoteCreate(ctx context.Context, arg NoteCreateParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, noteCreate,
		arg.Text,
		arg.PartOfSpeech,
		arg.Translation,
		arg.CommonLevel,
		arg.Explanation,
		arg.Usage,
		arg.UsageTranslation,
		arg.SourceName,
		arg.SourceReference,
		arg.DictionarySource,
		arg.Notes,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Text,
		&i.PartOfSpeech,
		&i.Translation,
		&i.Explanation,
		&i.CommonLevel,
		&i.Usage,
		&i.UsageTranslation,
		&i.SourceName,
		&i.SourceReference,
		&i.DictionarySource,
		&i.Notes,
		&i.Downloaded,
	)
	return i, err
}

const noteGet = `-- name: NoteGet :one
SELECT id, text, part_of_speech, translation, explanation, common_level, usage, usage_translation, source_name, source_reference, dictionary_source, notes, downloaded FROM notes WHERE id = ? LIMIT 1
`

func (q *Queries) NoteGet(ctx context.Context, id int64) (Note, error) {
	row := q.db.QueryRowContext(ctx, noteGet, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Text,
		&i.PartOfSpeech,
		&i.Translation,
		&i.Explanation,
		&i.CommonLevel,
		&i.Usage,
		&i.UsageTranslation,
		&i.SourceName,
		&i.SourceReference,
		&i.DictionarySource,
		&i.Notes,
		&i.Downloaded,
	)
	return i, err
}

const sourceCreate = `-- name: SourceCreate :one
INSERT INTO sources (
    name, reference, parts
) VALUES (?, ?, ?) RETURNING id, name, reference, parts, updated_at, created_at
`

type SourceCreateParams struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Parts     string `json:"parts"`
}

func (q *Queries) SourceCreate(ctx context.Context, arg SourceCreateParams) (Source, error) {
	row := q.db.QueryRowContext(ctx, sourceCreate, arg.Name, arg.Reference, arg.Parts)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Reference,
		&i.Parts,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const sourceDestroy = `-- name: SourceDestroy :exec
DELETE FROM sources WHERE id = ?
`

func (q *Queries) SourceDestroy(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, sourceDestroy, id)
	return err
}

const sourceGet = `-- name: SourceGet :one
SELECT id, name, reference, parts, updated_at, created_at FROM sources WHERE id = ? LIMIT 1
`

func (q *Queries) SourceGet(ctx context.Context, id int64) (Source, error) {
	row := q.db.QueryRowContext(ctx, sourceGet, id)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Reference,
		&i.Parts,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const sourceIndex = `-- name: SourceIndex :many
SELECT id, name, reference, parts, updated_at, created_at FROM sources ORDER BY created_at DESC
`

func (q *Queries) SourceIndex(ctx context.Context) ([]Source, error) {
	rows, err := q.db.QueryContext(ctx, sourceIndex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Source
	for rows.Next() {
		var i Source
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Reference,
			&i.Parts,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const sourceUpdate = `-- name: SourceUpdate :one
UPDATE sources
SET name = ?,
reference = ?
WHERE id = ? RETURNING id, name, reference, parts, updated_at, created_at
`

type SourceUpdateParams struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	ID        int64  `json:"id"`
}

func (q *Queries) SourceUpdate(ctx context.Context, arg SourceUpdateParams) (Source, error) {
	row := q.db.QueryRowContext(ctx, sourceUpdate, arg.Name, arg.Reference, arg.ID)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Reference,
		&i.Parts,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const termCreate = `-- name: TermCreate :one
INSERT INTO terms (
    text, variants, part_of_speech, common_level, translations, popularity
) VALUES (?, ?, ?, ?, ?, ?) RETURNING id, text, part_of_speech, variants, translations, common_level, popularity
`

type TermCreateParams struct {
	Text         string `json:"text"`
	Variants     string `json:"variants"`
	PartOfSpeech string `json:"part_of_speech"`
	CommonLevel  int64  `json:"common_level"`
	Translations string `json:"translations"`
	Popularity   int64  `json:"popularity"`
}

func (q *Queries) TermCreate(ctx context.Context, arg TermCreateParams) (Term, error) {
	row := q.db.QueryRowContext(ctx, termCreate,
		arg.Text,
		arg.Variants,
		arg.PartOfSpeech,
		arg.CommonLevel,
		arg.Translations,
		arg.Popularity,
	)
	var i Term
	err := row.Scan(
		&i.ID,
		&i.Text,
		&i.PartOfSpeech,
		&i.Variants,
		&i.Translations,
		&i.CommonLevel,
		&i.Popularity,
	)
	return i, err
}

const termsClearAll = `-- name: TermsClearAll :exec
DELETE FROM terms
`

func (q *Queries) TermsClearAll(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, termsClearAll)
	return err
}

const termsCount = `-- name: TermsCount :one
SELECT COUNT(*) FROM terms
`

func (q *Queries) TermsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, termsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const termsPopular = `-- name: TermsPopular :many
SELECT id, text, part_of_speech, variants, translations, common_level, popularity FROM terms ORDER BY CAST(popularity AS INT) LIMIT 100
`

func (q *Queries) TermsPopular(ctx context.Context) ([]Term, error) {
	rows, err := q.db.QueryContext(ctx, termsPopular)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Term
	for rows.Next() {
		var i Term
		if err := rows.Scan(
			&i.ID,
			&i.Text,
			&i.PartOfSpeech,
			&i.Variants,
			&i.Translations,
			&i.CommonLevel,
			&i.Popularity,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
