-- Add created_at column to tables table
-- Up migration

ALTER TABLE tables ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();