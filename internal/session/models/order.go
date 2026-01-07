package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID   `json:"id"`         // unique order ID
	SessionID uuid.UUID   `json:"session_id"` // associated session ID
	CreatedAt time.Time   `json:"created_at"` // when the order was created
	Status    OrderStatus `json:"status"`     // e.g., OrderStatusPending, OrderStatusPreparing, etc.
}

type OrderItems struct {
	ID         uuid.UUID `json:"id"`           // unique order ID
	OrderID    uuid.UUID `json:"order_id"`     // associated order ID
	MenuItemID uuid.UUID `json:"menu_item_id"` // associated menu item ID
	Quantity   int       `json:"quantity"`     // quantity of the menu item in the order
}

type OrderStatus string

const (
	OrderStatusCart      OrderStatus = "cart"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCancelled OrderStatus = "cancelled"
)
