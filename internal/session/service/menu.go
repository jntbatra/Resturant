package service

import (
	"context"
	"fmt"
	apperrors "restaurant/internal/errors"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// MenuService defines business logic for menu items
type MenuService interface {
	CreateMenuItem(ctx context.Context, Name string, Description string, Price float64, Category string, AvalabilityStatus models.ItemStatus) (*models.MenuItem, error)
	GetMenuItem(ctx context.Context, id uuid.UUID) (*models.MenuItem, error)
	ListMenuItems(ctx context.Context, offset int, limit int) ([]*models.MenuItem, error)
	UpdateMenuItem(ctx context.Context, id uuid.UUID, name string, desc string, category string, price float64, status models.ItemStatus) error
	DeleteMenuItem(ctx context.Context, id uuid.UUID) error
	GetMenuItemsByCategory(ctx context.Context, category string) ([]*models.MenuItem, error)
	ListCategories(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, name string) (*models.Category, error)
	DeleteCategory(ctx context.Context, name string) error
	UpdateCategory(ctx context.Context, old_name string, new_name string) error
	CategoryIDByNameCreateIfNotPresent(ctx context.Context, name string) (uuid.UUID, error)
	CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error)
}

// menuService implements MenuService
type menuService struct {
	repo repository.MenuRepository
}

// NewMenuService creates a new menu service
func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

// Implementations (wrappers around repository)
func (s *menuService) GetMenuItem(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	item, err := s.repo.GetMenuItem(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve menu item", err)
	}
	return item, nil
}

func (s *menuService) CreateMenuItem(ctx context.Context, Name string, Description string, Price float64, Category string, AvalabilityStatus models.ItemStatus) (*models.MenuItem, error) {
	// Shape validation (name, description, price, category) already done by handler using ValidateStruct

	// Ensure category exists (BUSINESS LOGIC)
	id, err := s.CategoryIDByNameCreateIfNotPresent(ctx, Category)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to ensure category exists", err)
	}

	item := &models.MenuItem{
		ID:                uuid.New(),
		Name:              Name,
		Description:       Description,
		Price:             Price,
		CategoryID:        id,
		AvalabilityStatus: AvalabilityStatus,
		CreatedAt:         time.Now(),
	}
	err = s.repo.CreateMenuItem(ctx, item)
	if err != nil {
		// Check for UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, apperrors.ErrDuplicateMenuItem
		}
		return nil, apperrors.WrapError(500, "failed to create menu item", err)
	}
	return item, nil
}

func (s *menuService) ListMenuItems(ctx context.Context, offset int, limit int) ([]*models.MenuItem, error) {
	items, err := s.repo.ListMenuItems(ctx, offset, limit)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list menu items", err)
	}
	return items, nil
}

func (s *menuService) UpdateMenuItem(ctx context.Context, id uuid.UUID, name string, desc string, category string, price float64, status models.ItemStatus) error {
	// Shape validation (name, description, price, category) already done by handler using ValidateStruct

	// Ensure category exists (BUSINESS LOGIC)
	categoryID, err := s.CategoryIDByNameCreateIfNotPresent(ctx, category)
	if err != nil {
		return apperrors.WrapError(500, "failed to ensure category exists", err)
	}

	item := &models.MenuItem{
		ID:                id,
		Name:              name,
		Description:       desc,
		Price:             price,
		CategoryID:        categoryID,
		AvalabilityStatus: status,
	}
	err = s.repo.UpdateMenuItem(ctx, item)
	if err != nil {
		return apperrors.WrapError(500, "failed to update menu item", err)
	}
	return nil
}

func (s *menuService) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteMenuItem(ctx, id)
	if err != nil {
		// Check for foreign key constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return apperrors.ErrForeignKeyViolation
		}
		return apperrors.WrapError(500, "failed to delete menu item", err)
	}
	return nil
}

func (s *menuService) GetMenuItemsByCategory(ctx context.Context, category string) ([]*models.MenuItem, error) {
	if category == "" {
		return nil, apperrors.NewValidationError("category is required")
	}
	items, err := s.repo.GetMenuItemsByCategory(ctx, category)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to get menu items by category", err)
	}
	return items, nil
}

func (s *menuService) ListCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list categories", err)
	}
	return categories, nil
}

func (s *menuService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	// Shape validation (name) already done by handler using ValidateStruct

	id := uuid.New()
	err := s.repo.CreateCategory(ctx, name, id)
	if err != nil {
		// Handle PostgreSQL UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation error code
				return nil, apperrors.ErrDuplicateCategory
			}
		}
		return nil, apperrors.WrapError(500, "failed to create category", err)
	}

	// Return category with generated ID
	return &models.Category{
		ID:   id,
		Name: name,
	}, nil
}

func (s *menuService) DeleteCategory(ctx context.Context, name string) error {
	err := s.repo.DeleteCategory(ctx, name)
	if err != nil {
		// Check for foreign key constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return apperrors.ErrForeignKeyViolation
		}
		return apperrors.WrapError(500, "failed to delete category", err)
	}
	return nil
}

func (s *menuService) UpdateCategory(ctx context.Context, old_name string, new_name string) error {
	err := s.repo.UpdateCategory(ctx, old_name, new_name)
	if err != nil {
		// Handle PostgreSQL UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation error code
				return apperrors.NewConflictError(fmt.Sprintf("category name '%s' already exists", new_name))
			}
		}
		return apperrors.WrapError(500, "failed to update category", err)
	}
	return nil
}

func (s *menuService) CategoryIDByNameCreateIfNotPresent(ctx context.Context, name string) (uuid.UUID, error) {
	id, err := s.repo.CategoryIDByNameCreateIfNotPresent(ctx, name)
	if err != nil {
		return uuid.Nil, apperrors.WrapError(500, "failed to get or create category", err)
	}
	return id, nil
}

func (s *menuService) CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	id, err := s.repo.CategoryIDByName(ctx, name)
	if err != nil {
		return uuid.Nil, apperrors.WrapError(500, "failed to get category ID", err)
	}
	return id, nil
}
