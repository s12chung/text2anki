-- name: TermCreate :one
INSERT INTO terms (
  text, variants, part_of_speech, common_level, translations, popularity
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;