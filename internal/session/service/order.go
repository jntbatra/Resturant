package service

import (
	"errors"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
	"time"

	"github.com/google/uuid"
)

// OrderService defines business logic for orders
type OrderService interface {
	CreateOrder(sessionID string) error
	GetOrder(id string) (*models.Order, error)
	ListOrders(limit int, offset int) ([]*models.Order, error)
	UpdateOrder(orderID string, status string) error
	CreateOrderItem(itemID string, quantity int, orderID string) error
	GetOrderItems(orderID string) ([]*models.OrderItems, error)
	GetOrdersBySession(sessionID string) ([]*models.Order, error)
	GetOrderItemsBySessionID(sessionID string) ([]*models.OrderItems, error)
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
func (s *orderService) CreateOrder(sessionID string) error {
	// Validate input: session ID is required
	if sessionID == "" {
		return errors.New("session ID is required")
	}
	// Create new order with generated UUID, initial status 'cart', and current timestamp
	order := &models.Order{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Status:    "cart",
		CreatedAt: time.Now(),
	}
	// Persist the order in the repository
	return s.repo.CreateOrder(order)
}

// GetOrder retrieves an order by ID with validation
func (s *orderService) GetOrder(id string) (*models.Order, error) {
	// Validate input: order ID is required
	if id == "" {
		return nil, errors.New("order ID is required")
	}
	// Retrieve order from repository
	return s.repo.GetOrder(id)
}

// ListOrders lists orders with pagination and validation
func (s *orderService) ListOrders(limit int, offset int) ([]*models.Order, error) {
	// Validate limit: must be between 1 and 100
	if limit <= 0 || limit > 100 {
		return nil, errors.New("limit must be between 1 and 100")
	}
	// Validate offset: cannot be negative
	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}
	// Retrieve paginated orders from repository
	return s.repo.ListOrders(limit, offset)
}

// UpdateOrder updates an order status with validation
func (s *orderService) UpdateOrder(orderID string, status string) error {
	// Validate inputs: order ID and status are required
	if orderID == "" {
		return errors.New("order ID is required")
	}
	if status == "" {
		return errors.New("status is required")
	}
	// Validate status against allowed values
	validStatuses := []string{"cart", "pending", "preparing", "served", "cancelled"}
	valid := false
	for _, s := range validStatuses {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid status")
	}
	// Update order status in repository
	return s.repo.UpdateOrder(orderID, status)
}

// CreateOrderItem creates a new order item with validation
func (s *orderService) CreateOrderItem(itemID string, quantity int, orderID string) error {
	// Validate inputs: item ID, quantity > 0, order ID required
	if itemID == "" {
		return errors.New("item ID is required")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if orderID == "" {
		return errors.New("order ID is required")
	}

	// Validate menu item exists and is available
	menuItem, err := s.menuService.GetMenuItem(itemID)
	if err != nil {
		return errors.New("menu item not found")
	}
	if menuItem.AvalabilityStatus != "available" {
		return errors.New("menu item is unavailable")
	}

	// Create order item with generated UUID
	Item := &models.OrderItems{
		ID:         uuid.New().String(),
		MenuItemID: itemID,
		Quantity:   quantity,
		OrderID:    orderID,
	}
	// Persist order item in repository
	return s.repo.CreateOrderItem(Item)
}

// GetOrderItems retrieves order items by order ID with validation
func (s *orderService) GetOrderItems(orderID string) ([]*models.OrderItems, error) {
	// Validate input: order ID is required
	if orderID == "" {
		return nil, errors.New("order ID is required")
	}
	// Retrieve order items from repository
	return s.repo.GetOrderItems(orderID)
}

// GetOrdersBySession retrieves orders by session ID with validation
func (s *orderService) GetOrdersBySession(sessionID string) ([]*models.Order, error) {
	// Validate input: session ID is required
	if sessionID == "" {
		return nil, errors.New("session ID is required")
	}
	// Retrieve orders from repository
	return s.repo.GetOrdersBySession(sessionID)
}

// GetOrderItemsBySessionID retrieves order items by session ID
// by orchestrating multiple repository calls
func (s *orderService) GetOrderItemsBySessionID(sessionID string) ([]*models.OrderItems, error) {
	// Validate input: session ID is required
	if sessionID == "" {
		return nil, errors.New("session ID is required")
	}

	// Step 1: Get all orders for the session
	orders, err := s.repo.GetOrdersBySession(sessionID)
	if err != nil {
		return nil, err
	}

	// If no orders, return empty slice
	if len(orders) == 0 {
		return []*models.OrderItems{}, nil
	}

	// Step 2: Extract order IDs
	orderIDs := make([]string, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	// Step 3: Get all order items for these orders
	return s.repo.GetOrderItemsByOrderIDs(orderIDs)
}
