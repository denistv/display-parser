
-- +migrate Up
ALTER TABLE pages ADD COLUMN entity_id TEXT;
ALTER TABLE models ADD COLUMN entity_id TEXT UNIQUE NOT NULL;
-- +migrate Down
ALTER TABLE pages DROP COLUMN entity_id;
ALTER TABLE models DROP COLUMN entity_id;
