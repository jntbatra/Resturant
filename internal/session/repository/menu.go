package repository

import (
	"database/sql"
	"errors"
	"restaurant/internal/session/models"
)

// MenuRepository defines methods for menu item database operations
type MenuRepository interface {
	// CreateMenuItem creates a new menu item
	CreateMenuItem(item *models.MenuItem) error

	// GetMenuItem retrieves a menu item by ID
	GetMenuItem(id string) (*models.MenuItem, error)

	// ListMenuItems lists all menu items
	ListMenuItems() ([]*models.MenuItem, error)

	// GetMenuItemsByCategory retrieves menu items by category
	GetMenuItemsByCategory(category string) ([]*models.MenuItem, error)

	// UpdateMenuItem updates a menu item
	UpdateMenuItem(item *models.MenuItem) error

	// DeleteMenuItem deletes a menu item by ID
	DeleteMenuItem(id string) error
}

// ErrMenuItemNotFound is returned when a menu item is not found
var ErrMenuItemNotFound = errors.New("menu item not found")

// postgresMenuRepository implements MenuRepository using PostgreSQL
type postgresMenuRepository struct {
	db *sql.DB
}

// NewMenuRepository creates a new PostgreSQL-based menu repository
func NewMenuRepository(db *sql.DB) MenuRepository {
	return &postgresMenuRepository{db: db}
}

// Implementations (stubs for now)
func (r *postgresMenuRepository) CreateMenuItem(item *models.MenuItem) error {
	_, err := r.db.Exec("INSERT INTO menu_items (id, name, description, price, avalability_status, category, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", 
		item.ID, item.Name, item.Description, item.Price, item.AvalabilityStatus, item.Category, item.CreatedAt)
	return err
}

func (r *postgresMenuRepository) GetMenuItem(id string) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.QueryRow("SELECT id, name, description, price, avalability_status, category, created_at FROM menu_items WHERE id = $1", id).Scan(
		&item.ID, &item.Name, &item.Description, &item.Price, &item.AvalabilityStatus, &item.Category, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMenuItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *postgresMenuRepository) ListMenuItems() ([]*models.MenuItem, error) {
	// TODO: Consider adding pagination (limit, offset) and availability filter for production use
	rows, err := r.db.Query("SELECT id, name, description, price, avalability_status, category, created_at FROM menu_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.AvalabilityStatus, &item.Category, &item.CreatedAt)
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

func (r *postgresMenuRepository) GetMenuItemsByCategory(category string) ([]*models.MenuItem, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, avalability_status, category, created_at FROM menu_items WHERE category = $1", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.AvalabilityStatus, &item.Category, &item.CreatedAt)
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

func (r *postgresMenuRepository) UpdateMenuItem(item *models.MenuItem) error {
	_, err := r.db.Exec("UPDATE menu_items SET name = $1, description = $2, price = $3, avalability_status = $4, category = $5 WHERE id = $6",
		item.Name, item.Description, item.Price, item.AvalabilityStatus, item.Category, item.ID)
	return err
}

func (r *postgresMenuRepository) DeleteMenuItem(id string) error {
	_, err := r.db.Exec("DELETE FROM menu_items WHERE id = $1", id)
	return err
}