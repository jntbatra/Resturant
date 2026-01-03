-- Simplify tables table: make id the table number (primary key), remove created_at
-- Up migration

-- Drop the foreign key constraint first
ALTER TABLE sessions DROP CONSTRAINT sessions_table_id_fkey;

-- Drop the old tables table
DROP TABLE tables;

-- Create new simplified tables table
CREATE TABLE tables (
    id INTEGER PRIMARY KEY
);

-- Recreate the foreign key constraint
ALTER TABLE sessions ADD CONSTRAINT sessions_table_id_fkey FOREIGN KEY (table_id) REFERENCES tables(id);