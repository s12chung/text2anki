-- name: SourcesIndex :many
SELECT * FROM sources ORDER BY updated_at DESC;

-- name: SourceGet :one
SELECT * FROM sources WHERE id = ? LIMIT 1;

-- name: SourceUpdate :one
UPDATE sources
SET name = ?,
reference = ?
WHERE id = ? RETURNING *;

-- name: SourceCreate :one
INSERT INTO sources (
    name, reference, parts
) VALUES (?, ?, ?) RETURNING *;

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

-- name: TermsClearAll :exec
DELETE FROM terms;

-- name: NotesIndex :many
SELECT * FROM notes ORDER BY updated_at DESC;

-- name: NotesDownloaded :many
SELECT * FROM notes WHERE downloaded = false ORDER BY updated_at DESC;

-- name: NotesUpdateDownloaded :execrows
UPDATE notes SET downloaded = true WHERE downloaded = false;

-- name: NoteGet :one
SELECT * FROM notes WHERE id = ? LIMIT 1;

-- name: NoteCreate :one
INSERT INTO notes (
    text, part_of_speech, translation, explanation, common_level,
                   usage, usage_translation,
                   source_name, source_reference, dictionary_source, notes
) VALUES (?, ?, ?, ?, ?,
          ?, ?,
          ?, ?, ?, ?) RETURNING *;