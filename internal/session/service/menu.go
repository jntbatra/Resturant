package service

import (
	"errors"
	"github.com/google/uuid"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
	"time"
)

// MenuService defines business logic for menu items
type MenuService interface {
	CreateMenuItem(Name string, Description string, Price float64, Category string, AvalabilityStatus models.ItemStatus) error
	GetMenuItem(id string) (*models.MenuItem, error)
	ListMenuItems() ([]*models.MenuItem, error)
	UpdateMenuItem(id string, name string, desc string, category string, price float64, status models.ItemStatus) error
	DeleteMenuItem(id string) error
	GetMenuItemsByCategory(category string) ([]*models.MenuItem, error)
	ListCategories() ([]string, error)
	CreateCategory(name string) error
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
func (s *menuService) GetMenuItem(id string) (*models.MenuItem, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetMenuItem(id)
}

func (s *menuService) CreateMenuItem(Name string, Description string, Price float64, Category string, AvalabilityStatus models.ItemStatus) error {
	if Name == "" {
		return errors.New("name is required")
	}
	if len(Name) > 100 {
		return errors.New("name must be less than 100 characters")
	}
	if Description == "" {
		return errors.New("description is required")
	}
	if len(Description) > 500 {
		return errors.New("description must be less than 500 characters")
	}
	if Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if Category == "" {
		return errors.New("category is required")
	}
	if AvalabilityStatus == "" {
		return errors.New("availability status is required")
	}

	// Ensure category exists
	categories, err := s.repo.ListCategories()
	if err != nil {
		return err
	}
	categoryExists := false
	for _, cat := range categories {
		if cat == Category {
			categoryExists = true
			break
		}
	}
	if !categoryExists {
		err = s.repo.CreateCategory(Category)
		if err != nil {
			return err
		}
	}

	item := &models.MenuItem{
		ID:                uuid.New().String(),
		Name:              Name,
		Description:       Description,
		Price:             Price,
		Category:          Category,
		AvalabilityStatus: AvalabilityStatus,
		CreatedAt:         time.Now(),
	}
	return s.repo.CreateMenuItem(item)
}

func (s *menuService) ListMenuItems() ([]*models.MenuItem, error) {
	return s.repo.ListMenuItems()
}

func (s *menuService) UpdateMenuItem(id string, name string, desc string, category string, price float64, status models.ItemStatus) error {
	if id == "" {
		return errors.New("id is required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) > 100 {
		return errors.New("name must be less than 100 characters")
	}
	if desc == "" {
		return errors.New("description is required")
	}
	if len(desc) > 500 {
		return errors.New("description must be less than 500 characters")
	}
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if category == "" {
		return errors.New("category is required")
	}
	if status == "" {
		return errors.New("availability status is required")
	}

	// Ensure category exists
	categories, err := s.repo.ListCategories()
	if err != nil {
		return err
	}
	categoryExists := false
	for _, cat := range categories {
		if cat == category {
			categoryExists = true
			break
		}
	}
	if !categoryExists {
		return errors.New("category does not exist")
	}

	item := &models.MenuItem{
		ID:                id,
		Name:              name,
		Description:       desc,
		Price:             price,
		Category:          category,
		AvalabilityStatus: status,
	}
	return s.repo.UpdateMenuItem(item)
}

func (s *menuService) DeleteMenuItem(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.DeleteMenuItem(id)
}

func (s *menuService) GetMenuItemsByCategory(category string) ([]*models.MenuItem, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}
	return s.repo.GetMenuItemsByCategory(category)
}

func (s *menuService) ListCategories() ([]string, error) {
	return s.repo.ListCategories()
}

func (s *menuService) CreateCategory(name string) error {
	if name == "" {
		return errors.New("category name is required")
	}
	// Check for duplicates
	categories, err := s.repo.ListCategories()
	if err != nil {
		return err
	}
	for _, cat := range categories {
		if cat == name {
			return errors.New("category already exists")
		}
	}
	return s.repo.CreateCategory(name)
}
