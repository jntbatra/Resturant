-- Remove foreign key constraint from menu_items.category
-- Down migration

-- Drop the foreign key constraint
ALTER TABLE menu_items DROP CONSTRAINT menu_items_category_fkey;

-- Drop the index
DROP INDEX IF EXISTS idx_menu_items_category;

-- Add back the old category column as VARCHAR
ALTER TABLE menu_items ADD COLUMN category_name VARCHAR(50);

-- Populate it by casting the UUID back to string
UPDATE menu_items SET category_name = category::text;

-- Drop the category column
ALTER TABLE menu_items DROP COLUMN category;

-- Rename back
ALTER TABLE menu_items RENAME COLUMN category_name TO category;

-- Change categories.id back to VARCHAR
ALTER TABLE categories ALTER COLUMN id TYPE VARCHAR(36) USING id::text;