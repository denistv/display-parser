
-- +migrate Up
ALTER TABLE models
ALTER COLUMN size TYPE decimal;
-- +migrate Down
ALTER TABLE models
ALTER COLUMN size TYPE bigint;