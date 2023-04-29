-- name: TermCreate :one
INSERT INTO terms (
    text, variants, part_of_speech, common_level, translations, popularity
)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: TermsCount :one
SELECT COUNT(*) FROM terms;

-- name: TermsPopular :many
SELECT * FROM terms ORDER BY CAST(popularity AS INT) LIMIT 100;