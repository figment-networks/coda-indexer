-- +goose Up
ALTER TABLE blocks ADD COLUMN reward_calculated BOOLEAN NOT NULL;

-- +goose Down
ALTER TABLE blocks DROP COLUMN reward_calculated;