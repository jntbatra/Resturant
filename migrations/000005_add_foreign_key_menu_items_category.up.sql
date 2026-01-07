-- Add foreign key constraint to menu_items.category referencing categories.id
-- Up migration

-- First, change categories.id to UUID type
ALTER TABLE categories ALTER COLUMN id TYPE UUID USING id::uuid;

-- Add a temporary column for the new category ID in menu_items
ALTER TABLE menu_items ADD COLUMN category_id UUID;

-- Populate the new column by casting the existing category (which contains UUID strings)
UPDATE menu_items SET category_id = category::uuid;

-- Drop the old category column
ALTER TABLE menu_items DROP COLUMN category;

-- Rename the new column to category
ALTER TABLE menu_items RENAME COLUMN category_id TO category;

-- Add the foreign key constraint
ALTER TABLE menu_items ADD CONSTRAINT menu_items_category_fkey FOREIGN KEY (category) REFERENCES categories(id);

-- Add index for performance
CREATE INDEX idx_menu_items_category ON menu_items(category);