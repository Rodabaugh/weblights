-- name: CreatePreset :one
INSERT INTO presets (id, created_at, updated_at, name, color1, color2)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetAllPresets :many
SELECT * FROM presets;

-- name: GetPresetByID :one
SELECT * FROM presets WHERE id = $1;

-- name: DeletePresetByID :exec
DELETE FROM presets WHERE id = $1;