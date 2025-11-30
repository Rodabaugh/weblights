-- name: CreatePreset :one
INSERT INTO presets (id, created_at, updated_at, name, colors)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetAllPresets :many
SELECT *
FROM presets
ORDER BY protected DESC;

-- name: GetPresetByID :one
SELECT * FROM presets WHERE id = $1;

-- name: DeletePresetByID :exec
DELETE FROM presets WHERE id = $1;