package service

import (
	"context"
	apperrors "restaurant/internal/errors"
	menuService "restaurant/internal/menu/service"
	"restaurant/internal/order/models"
	"restaurant/internal/order/repository"
	sessionService "restaurant/internal/session/service"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// OrderService defines business logic for orders
type OrderService interface {
	CreateOrder(ctx context.Context, sessionID uuid.UUID) (*models.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error)
	ListOrders(ctx context.Context, limit int, offset int) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, orderID uuid.UUID, status string) (*models.Order, error)
	CreateOrderItem(ctx context.Context, itemID uuid.UUID, quantity int, orderID uuid.UUID) (*models.OrderItems, error)
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItems, error)
	GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.Order, error)
	GetOrderItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.OrderItems, error)
}

// orderService implements OrderService
type orderService struct {
	repo           repository.OrderRepository
	menuService    menuService.MenuService
	sessionService sessionService.SessionService
}

// NewOrderService creates a new order service
func NewOrderService(repo repository.OrderRepository, menuService menuService.MenuService, sessionService sessionService.SessionService) OrderService {
	return &orderService{
		repo:           repo,
		menuService:    menuService,
		sessionService: sessionService,
	}
}

// Implementations (wrappers around repository)

// CreateOrder creates a new order for the given session ID with validation
func (s *orderService) CreateOrder(ctx context.Context, sessionID uuid.UUID) (*models.Order, error) {
	// Validate that the session exists
	_, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("session not found")
	}

	// Create new order with generated UUID, initial status 'cart', and current timestamp
	order := &models.Order{
		ID:        uuid.New(),
		SessionID: sessionID,
		Status:    "cart",
		CreatedAt: time.Now(),
	}
	// Persist the order in the repository
	err = s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to create order", err)
	}
	return order, nil
}

// GetOrder retrieves an order by ID with validation
func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	// Retrieve order from repository
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve order", err)
	}
	return order, nil
}

// ListOrders lists orders with pagination and validation
func (s *orderService) ListOrders(ctx context.Context, limit int, offset int) ([]*models.Order, error) {
	// Shape validation (limit, offset ranges) already done by handler using ValidateStruct
	// Retrieve paginated orders from repository
	orders, err := s.repo.ListOrders(ctx, limit, offset)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list orders", err)
	}
	return orders, nil
}

// UpdateOrder updates an order status with validation
func (s *orderService) UpdateOrder(ctx context.Context, orderID uuid.UUID, status string) (*models.Order, error) {
	// Get current order to validate state transition
	currentOrder, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve order", err)
	}

	// Validate state transition
	if err := s.validateOrderStatusTransition(currentOrder, models.OrderStatus(status)); err != nil {
		return nil, err
	}

	// Update order status in repository
	err = s.repo.UpdateOrder(ctx, orderID, status)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to update order status", err)
	}

	// Retrieve the updated order
	updatedOrder, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve updated order", err)
	}

	return updatedOrder, nil
}

// validateOrderStatusTransition validates order status transitions
func (s *orderService) validateOrderStatusTransition(order *models.Order, newStatus models.OrderStatus) error {
	currentStatus := order.Status

	// Allow transition to cancelled from any state, but check time limit
	if newStatus == models.OrderStatusCancelled {
		if time.Since(order.CreatedAt) > 30*time.Second {
			return apperrors.NewValidationError("orders can only be cancelled within 30 seconds of creation")
		}
		return nil
	}

	// Define valid forward transitions
	switch currentStatus {
	case models.OrderStatusCart:
		if newStatus != models.OrderStatusPending {
			return apperrors.NewValidationError("cart orders can only transition to pending")
		}
	case models.OrderStatusPending:
		if newStatus != models.OrderStatusPreparing {
			return apperrors.NewValidationError("pending orders can only transition to preparing")
		}
	case models.OrderStatusPreparing:
		if newStatus != models.OrderStatusServed {
			return apperrors.NewValidationError("preparing orders can only transition to served")
		}
	case models.OrderStatusServed:
		return apperrors.NewValidationError("served orders cannot be updated")
	case models.OrderStatusCancelled:
		return apperrors.NewValidationError("cancelled orders cannot be updated")
	default:
		return apperrors.NewValidationError("invalid current order status")
	}

	return nil
}

// CreateOrderItem creates a new order item with validation or updates quantity if item already exists
func (s *orderService) CreateOrderItem(ctx context.Context, itemID uuid.UUID, quantity int, orderID uuid.UUID) (*models.OrderItems, error) {
	// Shape validation (quantity > 0) already done by handler using ValidateStruct

	// Get the order to check its status
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve order", err)
	}

	// Only allow adding items if order is in cart status
	if order.Status != models.OrderStatusCart {
		return nil, apperrors.NewValidationError("can only add items to orders in cart status")
	}

	// Validate menu item exists and is available (BUSINESS LOGIC)
	menuItem, err := s.menuService.GetMenuItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if menuItem.AvalabilityStatus != "in_stock" {
		return nil, apperrors.ErrOutOfStock
	}

	// Check if this menu item already exists in the order
	existingItems, err := s.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to check existing order items", err)
	}

	// Look for existing item with same menu_item_id
	for _, item := range existingItems {
		if item.MenuItemID == itemID {
			// Update existing item's quantity
			newQuantity := item.Quantity + quantity
			err = s.repo.UpdateOrderItemQuantity(ctx, item.ID, newQuantity)
			if err != nil {
				return nil, apperrors.WrapError(500, "failed to update order item quantity", err)
			}
			// Return the updated item
			item.Quantity = newQuantity
			return item, nil
		}
	}

	// Create new order item if it doesn't exist
	Item := &models.OrderItems{
		ID:         uuid.New(),
		MenuItemID: itemID,
		Quantity:   quantity,
		OrderID:    orderID,
	}
	// Persist order item in repository
	err = s.repo.CreateOrderItem(ctx, Item)
	if err != nil {
		// Check for foreign key constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, apperrors.ErrForeignKeyViolation
		}
		return nil, apperrors.WrapError(500, "failed to create order item", err)
	}

	return Item, nil
}

// GetOrderItems retrieves order items by order ID with validation
func (s *orderService) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItems, error) {
	// Retrieve order items from repository
	items, err := s.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve order items", err)
	}
	return items, nil
}

// GetOrdersBySession retrieves orders by session ID with validation
func (s *orderService) GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.Order, error) {
	// Retrieve orders from repository
	orders, err := s.repo.GetOrdersBySession(ctx, sessionID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve orders for session", err)
	}
	return orders, nil
}

// GetOrderItemsBySessionID retrieves order items by session ID
// by orchestrating multiple repository calls
func (s *orderService) GetOrderItemsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.OrderItems, error) {
	// Step 1: Get all orders for the session
	orders, err := s.repo.GetOrdersBySession(ctx, sessionID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve orders for session", err)
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
	items, err := s.repo.GetOrderItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve order items", err)
	}
	return items, nil
}
