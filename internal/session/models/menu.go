package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID                uuid.UUID  // unique menu item ID
	Name              string     // name of the menu item
	Description       string     // description of the menu item
	Price             float64    // price of the menu item
	CategoryID        uuid.UUID  // category of the menu item
	AvalabilityStatus ItemStatus // status of the menu item in stock (e.g., "in_stock", "out_of_stock")
	CreatedAt         time.Time  // when the menu item was created
}

type ItemStatus string

const (
	ItemStatusInStock    ItemStatus = "in_stock"
	ItemStatusOutOfStock ItemStatus = "out_of_stock"
)

type category struct {
	ID   uuid.UUID // unique category ID
	Name string    // name of the category
}

type Category struct {
	ID   uuid.UUID `json:"id"`   // unique category ID
	Name string    `json:"name"` // name of the category
}
