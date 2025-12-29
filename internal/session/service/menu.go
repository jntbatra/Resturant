package service

import (
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
)

// MenuService defines business logic for menu items
type MenuService interface {
	CreateMenuItem(item *models.MenuItem) error
	GetMenuItem(id string) (*models.MenuItem, error)
	ListMenuItems() ([]*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) error
	DeleteMenuItem(id string) error
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
func (s *menuService) CreateMenuItem(item *models.MenuItem) error {
	return s.repo.CreateMenuItem(item)
}

func (s *menuService) GetMenuItem(id string) (*models.MenuItem, error) {
	return s.repo.GetMenuItem(id)
}

func (s *menuService) ListMenuItems() ([]*models.MenuItem, error) {
	return s.repo.ListMenuItems()
}

func (s *menuService) UpdateMenuItem(item *models.MenuItem) error {
	return s.repo.UpdateMenuItem(item)
}

func (s *menuService) DeleteMenuItem(id string) error {
	return s.repo.DeleteMenuItem(id)
}