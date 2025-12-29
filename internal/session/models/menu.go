package models

import (
	"time"
)

type MenuItem struct {
	ID          string    // unique menu item ID
	Name        string    // name of the menu item
	Description string    // description of the menu item
	Price       float64   // price of the menu item
	Category    string    // category of the menu item
	AvalabilityStatus    	ItemStatus// status of the menu item in stock (e.g., "in_stock", "out_of_stock")
	CreatedAt   time.Time // when the menu item was created
}

type ItemStatus string

const (
	ItemStatusInStock    ItemStatus = "in_stock"
	ItemStatusOutOfStock ItemStatus = "out_of_stock"
)