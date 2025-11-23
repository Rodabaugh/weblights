-- name: CreateLogEntry :one
INSERT INTO log (id, created_at, updated_at, requester, request, result)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetAllLogEntries :many
SELECT * FROM log;

-- name: GetLogEntryByID :one
SELECT * FROM log WHERE id = $1;

-- name: GetLogEntryByRequester :one
SELECT * FROM log WHERE requester = $1;

-- name: GetLogEntryByResult :one
SELECT * FROM log WHERE result = $1;

-- name: DeleteLogEntryByID :exec
DELETE FROM log WHERE id = $1;