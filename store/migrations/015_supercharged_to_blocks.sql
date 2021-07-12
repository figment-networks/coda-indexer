-- +goose Up
ALTER TABLE blocks ADD COLUMN supercharged BOOLEAN NOT NULL;

-- +goose Down
ALTER TABLE blocks DROP COLUMN supercharged;