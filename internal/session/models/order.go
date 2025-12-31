package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID   // unique order ID
	SessionID uuid.UUID   // associated session ID
	CreatedAt time.Time   // when the order was created
	Status    OrderStatus // e.g., OrderStatusPending, OrderStatusPreparing, etc.
}

type OrderItems struct {
	ID         uuid.UUID // unique order ID
	OrderID    uuid.UUID // associated order ID
	MenuItemID uuid.UUID // associated menu item ID
	Quantity   int       // quantity of the menu item in the order
}

type OrderStatus string

const (
	OrderStatusCart      OrderStatus = "cart"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCancelled OrderStatus = "cancelled"
)
