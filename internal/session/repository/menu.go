package repository

import (
	"context"
	"database/sql"
	"restaurant/internal/errors"
	"restaurant/internal/session/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// MenuRepository defines methods for menu item database operations
type MenuRepository interface {
	// CreateMenuItem creates a new menu item
	CreateMenuItem(ctx context.Context, item *models.MenuItem) error

	// GetMenuItem retrieves a menu item by ID
	GetMenuItem(ctx context.Context, id uuid.UUID) (*models.MenuItem, error)

	// ListMenuItems lists all menu items
	ListMenuItems(ctx context.Context, offset int, limit int) ([]*models.MenuItem, error)

	// GetMenuItemsByCategory retrieves menu items by category
	GetMenuItemsByCategory(ctx context.Context, category string) ([]*models.MenuItem, error)

	// UpdateMenuItem updates a menu item
	UpdateMenuItem(ctx context.Context, item *models.MenuItem) error

	// DeleteMenuItem deletes a menu item by ID
	DeleteMenuItem(ctx context.Context, id uuid.UUID) error

	// ListCategories lists all unique categories
	ListCategories(ctx context.Context) ([]models.Category, error)

	// CreateCategory creates a new category
	CreateCategory(ctx context.Context, name string, id uuid.UUID) error

	// GetCategoryByID retrieves a category by ID
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error)

	// DeleteCategory deletes a category
	DeleteCategory(ctx context.Context, name string) error

	UpdateCategory(ctx context.Context, old_name string, new_name string) error

	CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error)

	CategoryIDByNameCreateIfNotPresent(ctx context.Context, name string) (uuid.UUID, error)
}

// postgresMenuRepository implements MenuRepository using PostgreSQL
type postgresMenuRepository struct {
	db *sql.DB
}

// NewMenuRepository creates a new PostgreSQL-based menu repository
func NewMenuRepository(db *sql.DB) MenuRepository {
	return &postgresMenuRepository{db: db}
}

// Implementations (stubs for now)
func (r *postgresMenuRepository) CreateMenuItem(ctx context.Context, item *models.MenuItem) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO menu_items (id, name, description, price, avalability_status, category, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		item.ID, item.Name, item.Description, item.Price, item.AvalabilityStatus, item.CategoryID, item.CreatedAt)
	return err
}

func (r *postgresMenuRepository) GetMenuItem(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.QueryRowContext(ctx, "SELECT id, name, description, price, avalability_status, category, created_at FROM menu_items WHERE id = $1", id).Scan(
		&item.ID, &item.Name, &item.Description, &item.Price, &item.AvalabilityStatus, &item.CategoryID, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrMenuItemNotFound
		}
		return nil, errors.NewInternalError("failed to get menu item", err)
	}
	return &item, nil
}

func (r *postgresMenuRepository) ListMenuItems(ctx context.Context, offset int, limit int) ([]*models.MenuItem, error) {
	// TODO: Consider adding pagination (limit, offset) and availability filter for production use
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, price, avalability_status, category, created_at FROM menu_items ORDER BY created_at DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.AvalabilityStatus, &item.CategoryID, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresMenuRepository) GetMenuItemsByCategory(ctx context.Context, category string) ([]*models.MenuItem, error) {
	// First get the category ID by name
	categoryID, err := r.CategoryIDByName(ctx, category)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, price, avalability_status, category, created_at FROM menu_items WHERE category = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.AvalabilityStatus, &item.CategoryID, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresMenuRepository) UpdateMenuItem(ctx context.Context, item *models.MenuItem) error {
	_, err := r.db.ExecContext(ctx, "UPDATE menu_items SET name = $1, description = $2, price = $3, avalability_status = $4, category = $5 WHERE id = $6",
		item.Name, item.Description, item.Price, item.AvalabilityStatus, item.CategoryID, item.ID)
	return err
}

func (r *postgresMenuRepository) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM menu_items WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.ErrMenuItemNotFound
	}
	return nil
}

func (r *postgresMenuRepository) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *postgresMenuRepository) CreateCategory(ctx context.Context, name string, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO categories (id, name) VALUES ($1, $2)", id, name)
	return err
}

func (r *postgresMenuRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.QueryRowContext(ctx, "SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrCategoryNotFound
		}
		return nil, errors.NewInternalError("failed to get category", err)
	}
	return &category, nil
}

// UpdateCategory updates an existing category name
func (r *postgresMenuRepository) UpdateCategory(ctx context.Context, old_name string, new_name string) error {
	result, err := r.db.ExecContext(ctx, "UPDATE categories SET name = $1 WHERE LOWER(name) = LOWER($2)", new_name, old_name)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.ErrCategoryNotFound
	}
	return nil
}

func (r *postgresMenuRepository) DeleteCategory(ctx context.Context, name string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE LOWER(name) = LOWER($1)", name)
	if err != nil {
		// Check if it's a foreign key violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return errors.ErrForeignKeyViolation
		}
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.ErrCategoryNotFound
	}
	return nil
}

func (r *postgresMenuRepository) CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, "SELECT id FROM categories WHERE LOWER(name) = LOWER($1)", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, errors.ErrCategoryNotFound
		}
		return uuid.Nil, errors.NewInternalError("failed to get category", err)
	}
	return id, nil
}

func (r *postgresMenuRepository) CategoryIDByNameCreateIfNotPresent(ctx context.Context, name string) (uuid.UUID, error) {
	id, err := r.CategoryIDByName(ctx, name)
	if err == errors.ErrCategoryNotFound {
		// Create new category
		newID := uuid.New()
		err = r.CreateCategory(ctx, name, newID)
		if err != nil {
			// Check if it's a UNIQUE constraint violation
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				// Category was created by another request, try to get it
				return r.CategoryIDByName(ctx, name)
			}
			return uuid.Nil, errors.NewInternalError("failed to create category", err)
		}
		return newID, nil
	} else if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
