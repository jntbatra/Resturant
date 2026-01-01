package service

import (
	"context"
	"errors"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
	"time"

	"github.com/google/uuid"
)

// OrderService defines business logic for orders
type OrderService interface {
	CreateOrder(ctx context.Context, sessionID uuid.UUID) (*models.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error)
	ListOrders(ctx context.Context, limit int, offset int) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, orderID uuid.UUID, status string) error
	CreateOrderItem(ctx context.Context, itemID uuid.UUID, quantity int, orderID uuid.UUID) (*models.OrderItems, error)
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItems, error)
	GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.Order, error)
	GetOrderItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.OrderItems, error)
}

// orderService implements OrderService
type orderService struct {
	repo        repository.OrderRepository
	menuService MenuService
}

// NewOrderService creates a new order service
func NewOrderService(repo repository.OrderRepository, menuService MenuService) OrderService {
	return &orderService{
		repo:        repo,
		menuService: menuService,
	}
}

// Implementations (wrappers around repository)

// CreateOrder creates a new order for the given session ID with validation
func (s *orderService) CreateOrder(ctx context.Context, sessionID uuid.UUID) (*models.Order, error) {
	// Create new order with generated UUID, initial status 'cart', and current timestamp
	order := &models.Order{
		ID:        uuid.New(),
		SessionID: sessionID,
		Status:    "cart",
		CreatedAt: time.Now(),
	}
	// Persist the order in the repository
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetOrder retrieves an order by ID with validation
func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	// Retrieve order from repository
	return s.repo.GetOrder(ctx, id)
}

// ListOrders lists orders with pagination and validation
func (s *orderService) ListOrders(ctx context.Context, limit int, offset int) ([]*models.Order, error) {
	// Shape validation (limit, offset ranges) already done by handler using ValidateStruct
	// Retrieve paginated orders from repository
	return s.repo.ListOrders(ctx, limit, offset)
}

// UpdateOrder updates an order status with validation
func (s *orderService) UpdateOrder(ctx context.Context, orderID uuid.UUID, status string) error {
	// Shape validation (status oneof) already done by handler using ValidateStruct
	// Update order status in repository
	return s.repo.UpdateOrder(ctx, orderID, status)
}

// CreateOrderItem creates a new order item with validation
func (s *orderService) CreateOrderItem(ctx context.Context, itemID uuid.UUID, quantity int, orderID uuid.UUID) (*models.OrderItems, error) {
	// Shape validation (quantity > 0) already done by handler using ValidateStruct

	// Validate menu item exists and is available (BUSINESS LOGIC)
	menuItem, err := s.menuService.GetMenuItem(ctx, itemID)
	if err != nil {
		return nil, errors.New("menu item not found")
	}
	if menuItem.AvalabilityStatus != "in_stock" {
		return nil, errors.New("menu item is out of stock")
	}

	// Create order item with generated UUID
	Item := &models.OrderItems{
		ID:         uuid.New(),
		MenuItemID: itemID,
		Quantity:   quantity,
		OrderID:    orderID,
	}
	// Persist order item in repository
	err = s.repo.CreateOrderItem(ctx, Item)
	if err != nil {
		return nil, err
	}
	return Item, nil
}

// GetOrderItems retrieves order items by order ID with validation
func (s *orderService) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItems, error) {
	// Retrieve order items from repository
	return s.repo.GetOrderItems(ctx, orderID)
}

// GetOrdersBySession retrieves orders by session ID with validation
func (s *orderService) GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.Order, error) {
	// Retrieve orders from repository
	return s.repo.GetOrdersBySession(ctx, sessionID)
}

// GetOrderItemsBySessionID retrieves order items by session ID
// by orchestrating multiple repository calls
func (s *orderService) GetOrderItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.OrderItems, error) {
	// Step 1: Get all orders for the session
	orders, err := s.repo.GetOrdersBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// If no orders, return empty slice
	if len(orders) == 0 {
		return []*models.OrderItems{}, nil
	}

	// Step 2: Extract order IDs
	orderIDs := make([]uuid.UUID, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	// Step 3: Get all order items for these orders
	return s.repo.GetOrderItemsByOrderIDs(ctx, orderIDs)
}
