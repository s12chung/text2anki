// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: queries.sql

package db

import (
	"context"
)

const sourceCreate = `-- name: SourceCreate :one
INSERT INTO sources (
    name, tokenized_texts
) VALUES (?, ?) RETURNING id, name, tokenized_texts, updated_at, created_at
`

type SourceCreateParams struct {
	Name           string `json:"name"`
	TokenizedTexts string `json:"tokenized_texts"`
}

func (q *Queries) SourceCreate(ctx context.Context, arg SourceCreateParams) (Source, error) {
	row := q.db.QueryRowContext(ctx, sourceCreate, arg.Name, arg.TokenizedTexts)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.TokenizedTexts,
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
SELECT id, name, tokenized_texts, updated_at, created_at FROM sources WHERE id = ? LIMIT 1
`

func (q *Queries) SourceGet(ctx context.Context, id int64) (Source, error) {
	row := q.db.QueryRowContext(ctx, sourceGet, id)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.TokenizedTexts,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const sourceList = `-- name: SourceList :many
SELECT id, name, tokenized_texts, updated_at, created_at FROM sources ORDER BY created_at DESC
`

func (q *Queries) SourceList(ctx context.Context) ([]Source, error) {
	rows, err := q.db.QueryContext(ctx, sourceList)
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
			&i.TokenizedTexts,
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
SET name = ?
WHERE id = ? RETURNING id, name, tokenized_texts, updated_at, created_at
`

type SourceUpdateParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (q *Queries) SourceUpdate(ctx context.Context, arg SourceUpdateParams) (Source, error) {
	row := q.db.QueryRowContext(ctx, sourceUpdate, arg.Name, arg.ID)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.TokenizedTexts,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const termCreate = `-- name: TermCreate :one
INSERT INTO terms (
    text, variants, part_of_speech, common_level, translations, popularity
) VALUES (?, ?, ?, ?, ?, ?) RETURNING text, variants, part_of_speech, common_level, translations, popularity
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
		&i.Text,
		&i.Variants,
		&i.PartOfSpeech,
		&i.CommonLevel,
		&i.Translations,
		&i.Popularity,
	)
	return i, err
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
SELECT text, variants, part_of_speech, common_level, translations, popularity FROM terms ORDER BY CAST(popularity AS INT) LIMIT 100
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
			&i.Text,
			&i.Variants,
			&i.PartOfSpeech,
			&i.CommonLevel,
			&i.Translations,
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
