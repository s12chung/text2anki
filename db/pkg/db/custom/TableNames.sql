-- name: TableNames :many
SELECT name FROM sqlite_master WHERE type = 'table';