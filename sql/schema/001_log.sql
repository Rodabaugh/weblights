-- +goose Up
CREATE TABLE log (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        requester TEXT NOT NULL,
                        request TEXT NOT NULL,
                        result BOOLEAN NOT NULL);

-- +goose Down
DROP TABLE log;