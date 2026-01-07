package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID                uuid.UUID  `json:"id"`                  // unique menu item ID
	Name              string     `json:"name"`                // name of the menu item
	Description       string     `json:"description"`         // description of the menu item
	Price             float64    `json:"price"`               // price of the menu item
	CategoryID        uuid.UUID  `json:"category_id"`         // category of the menu item
	AvalabilityStatus ItemStatus `json:"availability_status"` // status of the menu item in stock (e.g., "in_stock", "out_of_stock")
	CreatedAt         time.Time  `json:"created_at"`          // when the menu item was created
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
