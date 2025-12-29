package service

import (
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
)

// OrderService defines business logic for orders
type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	ListOrders() ([]*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id string) error
	CreateOrderItem(item *models.OrderItems) error
	GetOrderItems(orderID string) ([]*models.OrderItems, error)
}

// orderService implements OrderService
type orderService struct {
	repo repository.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

// Implementations (wrappers around repository)
func (s *orderService) CreateOrder(order *models.Order) error {
	return s.repo.CreateOrder(order)
}

func (s *orderService) GetOrder(id string) (*models.Order, error) {
	return s.repo.GetOrder(id)
}

func (s *orderService) ListOrders() ([]*models.Order, error) {
	return s.repo.ListOrders()
}

func (s *orderService) UpdateOrder(order *models.Order) error {
	return s.repo.UpdateOrder(order)
}

func (s *orderService) DeleteOrder(id string) error {
	return s.repo.DeleteOrder(id)
}

func (s *orderService) CreateOrderItem(item *models.OrderItems) error {
	return s.repo.CreateOrderItem(item)
}

func (s *orderService) GetOrderItems(orderID string) ([]*models.OrderItems, error) {
	return s.repo.GetOrderItems(orderID)
}