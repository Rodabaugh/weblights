-- +goose Up
CREATE TABLE presets (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        name TEXT NOT NULL,
                        colors BIGINT[] NOT NULL,
                        protected BOOLEAN NOT NULL DEFAULT 'false');

-- +goose Down
DROP TABLE presets;