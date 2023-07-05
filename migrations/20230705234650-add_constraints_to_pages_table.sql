
-- +migrate Up
ALTER TABLE pages ALTER COLUMN entity_id SET NOT NULL;
ALTER TABLE pages ADD UNIQUE (entity_id);
-- +migrate Down
ALTER TABLE pages ALTER COLUMN entity_id DROP NOT NULL;
ALTER TABLE pages DROP CONSTRAINT pages_entity_id_key;
