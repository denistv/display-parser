
-- +migrate Up
ALTER TABLE models ADD COLUMN panel_bit_depth int DEFAULT NULL;
-- +migrate Down
ALTER TABLE models DROP COLUMN panel_bit_depth;
