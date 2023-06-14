-- +migrate Up
ALTER TABLE models ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();

-- +migrate Down
ALTER TABLE models DROP COLUMN updated_at;
