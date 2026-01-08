package validation

import (
	"restaurant/internal/order/models"

	"github.com/google/uuid"
)

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status models.OrderStatus `json:"status" validate:"required,oneof=cart pending preparing served cancelled"`
}

// ListOrdersRequest represents the request to list orders with pagination
type ListOrdersRequest struct {
	Offset    int       `json:"offset" validate:"min=0"`
	Limit     int       `json:"limit" validate:"required,min=1,max=100"`
	SessionID uuid.UUID `json:"session_id" validate:"omitempty"`
}

// CreateOrderItemRequest represents the request to add an item to an order
type CreateOrderItemRequest struct {
	MenuItemID uuid.UUID `json:"menu_item_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"required,gt=0"`
}

// UpdateOrderItemRequest represents the request to update an order item
type UpdateOrderItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

// ValidateCreateOrder validates the create order request
func ValidateCreateOrder(req CreateOrderRequest) error {
	return ValidateStruct(req)
}

// ValidateUpdateOrder validates the update order request
func ValidateUpdateOrder(req UpdateOrderRequest) error {
	return ValidateStruct(req)
}

// ValidateListOrders validates the list orders request
func ValidateListOrders(req ListOrdersRequest) error {
	return ValidateStruct(req)
}

// ValidateCreateOrderItem validates the create order item request
func ValidateCreateOrderItem(req CreateOrderItemRequest) error {
	return ValidateStruct(req)
}

// ValidateUpdateOrderItem validates the update order item request
func ValidateUpdateOrderItem(req UpdateOrderItemRequest) error {
	return ValidateStruct(req)
}
