-- +goose Up
CREATE TYPE OWNER_TYPE AS ENUM ('validator', 'delegator', 'unknown');

-- +goose Down
DROP TYPE IF EXISTS OWNER_TYPE;