# Restaurant Management System - Models Documentation

## Overview

This document provides comprehensive documentation for the data models used in the Restaurant Management System. The models define the core entities and their relationships within the application.

## Architecture

### Model Structure
- All models are defined in the `internal/session/models` package
- Models use Go struct tags for JSON serialization where applicable
- Time fields use `time.Time` for timestamps
- String enums are used for status fields
- UUIDs are used for unique identifiers

## Session Models

### SessionStatus Enum
Defines the possible states of a dining session.

```go
type SessionStatus string
```

**Constants:**
- `StatusActive` - Session is currently active
- `StatusCompleted` - Session has been completed
- `StatusPending` - Session is pending activation
- `StatusCancelled` - Session has been cancelled

### Session Struct
Represents a dining session started by QR scan.

```go
type Session struct {
    ID          string        // unique session ID (UUID)
    TableID     int           // which table this session is for
    CreatedAt   time.Time     // when the session was created
    CompletedAt *time.Time    // when the session was completed, nil if not completed
    Status      SessionStatus // e.g., StatusActive, StatusCompleted, or StatusPending
}
```

**Fields:**
- **ID**: Unique identifier for the session (UUID format)
- **TableID**: Reference to the table number where the session is active
- **CreatedAt**: Timestamp when the session was initiated
- **CompletedAt**: Timestamp when the session ended (nullable)
- **Status**: Current status of the session

### Table Struct
Represents a physical table in the restaurant.

```go
type Table struct {
    ID        string    // unique table ID
    Number    int       // table number in the restaurant
    CreatedAt time.Time // when the table was created
}
```

**Fields:**
- **ID**: Unique identifier for the table
- **Number**: Human-readable table number displayed to customers
- **CreatedAt**: Timestamp when the table record was created

### Bill Struct
Represents a bill for a completed session.

```go
type Bill struct {
    ID        string
    SessionID string
    Total     float64
    Subtotal  float64
    Tax       float64
    CreatedAt time.Time
}
```

**Fields:**
- **ID**: Unique identifier for the bill
- **SessionID**: Reference to the session this bill belongs to
- **Total**: Final total amount including tax
- **Subtotal**: Amount before tax
- **Tax**: Tax amount
- **CreatedAt**: Timestamp when the bill was generated

## Order Models

### Order Struct
Represents an order placed within a session.

```go
type Order struct {
    ID        string      // unique order ID
    SessionID string      // associated session ID
    CreatedAt time.Time   // when the order was created
    Status    OrderStatus // e.g., OrderStatusPending, OrderStatusPreparing, etc.
}
```

**Fields:**
- **ID**: Unique identifier for the order (UUID format)
- **SessionID**: Reference to the session this order belongs to
- **CreatedAt**: Timestamp when the order was created
- **Status**: Current status of the order processing

### OrderItems Struct
Represents individual items within an order.

```go
type OrderItems struct {
    ID         string // unique order ID
    OrderID    string // associated order ID
    MenuItemID string // associated menu item ID
    Quantity   int    // quantity of the menu item in the order
}
```

**Fields:**
- **ID**: Unique identifier for the order item
- **OrderID**: Reference to the parent order
- **MenuItemID**: Reference to the menu item being ordered
- **Quantity**: Number of units of the menu item

### OrderStatus Enum
Defines the possible states of an order.

```go
type OrderStatus string
```

**Constants:**
- `OrderStatusCart` - Items added to cart, not yet submitted
- `OrderStatusPending` - Order submitted, waiting for preparation
- `OrderStatusPreparing` - Order is being prepared
- `OrderStatusServed` - Order has been served
- `OrderStatusCancelled` - Order has been cancelled

## Menu Models

### MenuItem Struct
Represents an item available on the restaurant menu.

```go
type MenuItem struct {
    ID                string     // unique menu item ID
    Name              string     // name of the menu item
    Description       string     // description of the menu item
    Price             float64    // price of the menu item
    Category          string     // category of the menu item
    AvalabilityStatus ItemStatus // status of the menu item in stock (e.g., "in_stock", "out_of_stock")
    CreatedAt         time.Time  // when the menu item was created
}
```

**Fields:**
- **ID**: Unique identifier for the menu item (UUID format)
- **Name**: Display name of the menu item
- **Description**: Detailed description of the menu item
- **Price**: Price in the restaurant's currency
- **Category**: Menu category (e.g., "Appetizers", "Main Course")
- **AvalabilityStatus**: Current availability status
- **CreatedAt**: Timestamp when the menu item was added

### ItemStatus Enum
Defines the availability status of menu items.

```go
type ItemStatus string
```

**Constants:**
- `ItemStatusInStock` - Item is available for ordering
- `ItemStatusOutOfStock` - Item is temporarily unavailable

### Category Struct (Private)
Internal representation of menu categories.

```go
type category struct {
    ID   string // unique category ID
    Name string // name of the category
}
```

**Note:** This struct is not exported (lowercase 'c'), suggesting it may be used internally or replaced by a simpler string-based approach in the public API.

## Relationships

### Entity Relationships
```
Session (1) -> (many) Orders
Session (1) -> (1) Table
Session (1) -> (1) Bill

Order (1) -> (many) OrderItems
Order (1) -> (1) Session

OrderItems (many) -> (1) MenuItem

MenuItem (many) -> (1) Category
```

### Database Mapping
- **Sessions Table**: Maps to `Session` struct
- **Orders Table**: Maps to `Order` struct
- **Order_Items Table**: Maps to `OrderItems` struct
- **Menu_Items Table**: Maps to `MenuItem` struct
- **Categories Table**: Maps to `category` struct
- **Tables Table**: Maps to `Table` struct
- **Bills Table**: Maps to `Bill` struct

## Validation Rules

### Common Field Validations
- **IDs**: UUID format strings, non-empty
- **Timestamps**: Valid `time.Time` values
- **Prices**: Positive float64 values
- **Quantities**: Positive integers
- **Names**: Non-empty strings with reasonable length limits

### Status Validations
- Status fields must match defined enum constants
- Invalid status values should be rejected at service layer

## Usage Examples

### Creating a Session
```go
session := &models.Session{
    ID:        uuid.New().String(),
    TableID:   5,
    CreatedAt: time.Now(),
    Status:    models.StatusActive,
}
```

### Creating an Order
```go
order := &models.Order{
    ID:        uuid.New().String(),
    SessionID: sessionID,
    CreatedAt: time.Now(),
    Status:    models.OrderStatusCart,
}
```

### Creating a Menu Item
```go
menuItem := &models.MenuItem{
    ID:                uuid.New().String(),
    Name:              "Grilled Salmon",
    Description:       "Fresh Atlantic salmon grilled to perfection",
    Price:             24.99,
    Category:          "Main Course",
    AvalabilityStatus: models.ItemStatusInStock,
    CreatedAt:         time.Now(),
}
```

## Future Enhancements

- Add image URLs for menu items
- Include nutritional information
- Add allergen warnings
- Support for multiple price tiers (regular, large, etc.)
- Localization support for multi-language menus
- Seasonal availability flags

## Dependencies

- `time`: For timestamp fields
- `github.com/google/uuid`: For ID generation (used in services)

This documentation covers all data models in the Restaurant Management System.</content>
<parameter name="filePath">/home/jntbatra/Resturant/documentation/models-docs.md