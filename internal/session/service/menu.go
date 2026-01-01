package service

import (
	"context"
	"errors"
	"fmt"
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
	CategoryIDByNameCareateIfNotPresent(ctx context.Context, name string) (uuid.UUID, error)
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
	return s.repo.GetMenuItem(ctx, id)
}

func (s *menuService) CreateMenuItem(ctx context.Context, Name string, Description string, Price float64, Category string, AvalabilityStatus models.ItemStatus) (*models.MenuItem, error) {
	// Shape validation (name, description, price, category) already done by handler using ValidateStruct

	// Ensure category exists (BUSINESS LOGIC)
	id, err := s.CategoryIDByNameCareateIfNotPresent(ctx, Category)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return item, nil
}

func (s *menuService) ListMenuItems(ctx context.Context, offset int, limit int) ([]*models.MenuItem, error) {
	return s.repo.ListMenuItems(ctx, offset, limit)
}

func (s *menuService) UpdateMenuItem(ctx context.Context, id uuid.UUID, name string, desc string, category string, price float64, status models.ItemStatus) error {
	// Shape validation (name, description, price, category) already done by handler using ValidateStruct

	// Ensure category exists (BUSINESS LOGIC)
	id, err := s.CategoryIDByNameCareateIfNotPresent(ctx, category)
	if err != nil {
		return err
	}

	item := &models.MenuItem{
		ID:                id,
		Name:              name,
		Description:       desc,
		Price:             price,
		CategoryID:        id,
		AvalabilityStatus: status,
	}
	return s.repo.UpdateMenuItem(ctx, item)
}

func (s *menuService) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMenuItem(ctx, id)
}

func (s *menuService) GetMenuItemsByCategory(ctx context.Context, category string) ([]*models.MenuItem, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}
	return s.repo.GetMenuItemsByCategory(ctx, category)
}

func (s *menuService) ListCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *menuService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	// Shape validation (name) already done by handler using ValidateStruct

	id := uuid.New()
	err := s.repo.CreateCategory(ctx, name, id)
	if err != nil {
		// Handle PostgreSQL UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation error code
				return nil, errors.New("category already exists")
			}
		}
		return nil, err
	}

	// Return category with generated ID
	return &models.Category{
		ID:   id,
		Name: name,
	}, nil
}

func (s *menuService) DeleteCategory(ctx context.Context, name string) error {
	return s.repo.DeleteCategory(ctx, name)
}

func (s *menuService) UpdateCategory(ctx context.Context, old_name string, new_name string) error {
	err := s.repo.UpdateCategory(ctx, old_name, new_name)
	if err != nil {
		// Handle PostgreSQL UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation error code
				return fmt.Errorf("category name '%s' already exists", new_name)
			}
		}
	}
	return err
}

func (s *menuService) CategoryIDByNameCareateIfNotPresent(ctx context.Context, name string) (uuid.UUID, error) {
	id, err := s.repo.CategoryIDByNameCareateIfNotPresent(ctx, name)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *menuService) CategoryIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	return s.repo.CategoryIDByName(ctx, name)
}
