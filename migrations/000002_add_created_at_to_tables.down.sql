-- Remove created_at column from tables table
-- Down migration

ALTER TABLE tables DROP COLUMN created_at;