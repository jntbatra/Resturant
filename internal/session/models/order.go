package models

import (
	"time"
)

type Order struct {
	ID        string      // unique order ID
	SessionID string      // associated session ID
	CreatedAt time.Time   // when the order was created
	Status    OrderStatus // e.g., OrderStatusPending, OrderStatusPreparing, etc.
}

type OrderItems struct {
	ID         string // unique order ID
	OrderID    string // associated order ID
	MenuItemID string // associated menu item ID
	Quantity   int    // quantity of the menu item in the order
}

type OrderStatus string

const (
	OrderStatusCart      OrderStatus = "cart"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCancelled OrderStatus = "cancelled"
)
