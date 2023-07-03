-- name: SourceIndex :many
SELECT * FROM sources ORDER BY created_at DESC;

-- name: SourceGet :one
SELECT * FROM sources WHERE id = ? LIMIT 1;

-- name: SourceUpdate :one
UPDATE sources
SET name = ?
WHERE id = ? RETURNING *;

-- name: SourceCreate :one
INSERT INTO sources (
    name, tokenized_texts
) VALUES (?, ?) RETURNING *;

-- name: SourceDestroy :exec
DELETE FROM sources WHERE id = ?;

-- name: TermsCount :one
SELECT COUNT(*) FROM terms;

-- name: TermsPopular :many
SELECT * FROM terms ORDER BY CAST(popularity AS INT) LIMIT 100;

-- name: TermCreate :one
INSERT INTO terms (
    text, variants, part_of_speech, common_level, translations, popularity
) VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: NoteGet :one
SELECT * FROM notes WHERE id = ? LIMIT 1;

-- name: NoteCreate :one
INSERT INTO notes (
    text, part_of_speech, translation, common_level, explanation, usage, usage_translation, dictionary_source, notes
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;