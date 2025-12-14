-- +goose Up
ALTER TABLE feeds
  ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE feeds
  DROP COLUMN updated_at,
  DROP COLUMN created_at;
