-- Revert tables table to original structure
-- Down migration

-- Drop the foreign key constraint
ALTER TABLE sessions DROP CONSTRAINT sessions_table_id_fkey;

-- Drop the simplified tables table
DROP TABLE tables;

-- Recreate the original tables table
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    number INTEGER UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Recreate the foreign key constraint
ALTER TABLE sessions ADD CONSTRAINT sessions_table_id_fkey FOREIGN KEY (table_id) REFERENCES tables(id);

-- Migrate existing data (if any) - this would need to be done manually for existing sessions