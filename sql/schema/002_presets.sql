-- +goose Up
CREATE TABLE presets (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        name TEXT NOT NULL,
                        color1 BIGINT NOT NULL,
                        color2 BIGINT NOT NULL);

-- +goose Down
DROP TABLE presets;