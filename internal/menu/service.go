package menu

import (
	"context"
	"fmt"
	apperrors "restaurant/internal/errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// MenuService defines business logic for menu items
type MenuService interface {
	CreateMenuItem(ctx context.Context, Name string, Description string, Price float64, Category string, AvalabilityStatus ItemStatus) (*MenuItem, error)
	GetMenuItem(ctx context.Context, id uuid.UUID) (*MenuItem, error)
	ListMenuItems(ctx context.Context, offset int, limit int) ([]*MenuItem, error)
	UpdateMenuItem(ctx context.Context, id uuid.UUID, name string, desc string, category string, price float64, avalabilityStatus ItemStatus) error
	DeleteMenuItem(ctx context.Context, id uuid.UUID) error
	GetMenuItemsByCategory(ctx context.Context, category string) ([]*MenuItem, error)
	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, name string) (*Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
	DeleteCategory(ctx context.Context, name string) error
	UpdateCategory(ctx context.Context, old_name string, new_name string) (*Category, error)
	CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error)
}

// menuService implements MenuService
type menuService struct {
	repo MenuRepository
}

// NewMenuService creates a new menu service
func NewMenuService(repo MenuRepository) MenuService {
	return &menuService{repo: repo}
}

// Implementations (wrappers around repository)
func (s *menuService) GetMenuItem(ctx context.Context, id uuid.UUID) (*MenuItem, error) {
	item, err := s.repo.GetMenuItem(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve menu item", err)
	}
	return item, nil
}

func (s *menuService) CreateMenuItem(ctx context.Context, Name string, Description string, Price float64, Category string, AvalabilityStatus ItemStatus) (*MenuItem, error) {
	// Shape validation (name, description, price, category) already done by handler using ValidateStruct

	// Ensure category exists (BUSINESS LOGIC)
	id, err := s.CategoryIDByName(ctx, Category)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to ensure category exists", err)
	}

	item := &MenuItem{
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

func (s *menuService) ListMenuItems(ctx context.Context, offset int, limit int) ([]*MenuItem, error) {
	items, err := s.repo.ListMenuItems(ctx, offset, limit)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list menu items", err)
	}
	return items, nil
}

func (s *menuService) UpdateMenuItem(ctx context.Context, id uuid.UUID, name string, desc string, category string, price float64, avalabilityStatus ItemStatus) error {
	// Shape validation (name, description, price, category) already done by handler using ValidateStruct

	// Ensure category exists (BUSINESS LOGIC)
	categoryID, err := s.CategoryIDByName(ctx, category)
	if err != nil {
		return apperrors.WrapError(500, "failed to ensure category exists", err)
	}

	item := &MenuItem{
		ID:                id,
		Name:              name,
		Description:       desc,
		Price:             price,
		CategoryID:        categoryID,
		AvalabilityStatus: avalabilityStatus,
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

func (s *menuService) GetMenuItemsByCategory(ctx context.Context, category string) ([]*MenuItem, error) {
	if category == "" {
		return nil, apperrors.NewValidationError("category is required")
	}
	items, err := s.repo.GetMenuItemsByCategory(ctx, category)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to get menu items by category", err)
	}
	return items, nil
}

func (s *menuService) ListCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list categories", err)
	}
	return categories, nil
}

func (s *menuService) CreateCategory(ctx context.Context, name string) (*Category, error) {
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
	return &Category{
		ID:   id,
		Name: name,
	}, nil
}

func (s *menuService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to get category", err)
	}
	return category, nil
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

func (s *menuService) UpdateCategory(ctx context.Context, old_name string, new_name string) (*Category, error) {
	err := s.repo.UpdateCategory(ctx, old_name, new_name)
	if err != nil {
		// Handle PostgreSQL UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation error code
				return nil, apperrors.NewConflictError(fmt.Sprintf("category name '%s' already exists", new_name))
			}
		}
		return nil, apperrors.WrapError(500, "failed to update category", err)
	}

	// Retrieve the updated category
	id, err := s.repo.CategoryIDByName(ctx, new_name)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve updated category ID", err)
	}

	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve updated category", err)
	}

	return category, nil
}

func (s *menuService) CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	id, err := s.repo.CategoryIDByName(ctx, name)
	if err != nil {
		return uuid.Nil, apperrors.WrapError(500, "failed to get category ID", err)
	}
	return id, nil
}
