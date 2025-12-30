# Restaurant Management System - Service Layer Documentation

## Overview

This document provides comprehensive documentation for the Service layer of the Restaurant Management System. The Service layer contains business logic, input validation, and orchestrates repository operations.

## Architecture

### Layered Architecture
```
Handler (HTTP API) -> Service (Business Logic) -> Repository (Data Access) -> Database
```

### Service Layer Responsibilities
The Service layer acts as the orchestration layer:
- **Business Logic**: Implements validation and business rules
- **Cross-Domain Orchestration**: Coordinates multiple repository calls
- **Domain Integration**: Services can depend on other services for validation
- **Transaction Management**: Handles multi-step operations

**Example: OrderService depends on MenuService**
- Validates menu items exist before creating order items
- Checks menu item availability status
- Ensures data consistency across domains

### Key Technologies
- **Language**: Go
- **Database**: PostgreSQL
- **Driver**: `github.com/lib/pq`
- **UUID Generation**: `github.com/google/uuid`

## Service Layer

The Service layer implements business logic, input validation, and orchestrates repository operations.

### Common Patterns
- Services wrap repository calls with validation
- Implement comprehensive input validation
- Handle business rules and constraints
- Generate UUIDs for new entities
- Return user-friendly error messages
- **Cross-domain orchestration**: Services coordinate multiple repository calls
- **Domain validation**: Services call other services for cross-domain checks

### Session Service

#### Interface: `SessionService`
```go
type SessionService interface {
    CreateSession(tableNumber int) (string, error)
    GetSession(id string) (*models.Session, error)
    EndSession(id string) error
    ListActiveSessions() ([]*models.Session, error)
    ChangeTable(sessionID string, newTableNumber int) error
}
```

#### Methods Documentation

**CreateSession(tableNumber int) (string, error)**
- Validates table exists
- Creates new session with generated UUID
- Returns session ID

**GetSession(id string) (*models.Session, error)**
- Validates session ID format
- Retrieves session details

**EndSession(id string) error**
- Validates session exists and is active
- Updates status to completed

**ListActiveSessions() ([]*models.Session, error)**
- Retrieves all active sessions

**ChangeTable(sessionID string, newTableNumber int) error**
- Validates both session and new table exist
- Updates table assignment

### Order Service

#### Dependencies
- `OrderRepository`: Data access for orders and order items
- `MenuService`: Cross-domain validation for menu items

#### Interface: `OrderService`
```go
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
```

#### Methods Documentation

**CreateOrder(sessionID string) error**
- Validates session ID is provided
- Creates order with "cart" status
- Generates UUID and timestamp

**GetOrder(id string) (*models.Order, error)**
- Validates order ID
- Retrieves order details

**ListOrders(limit int, offset int) ([]*models.Order, error)**
- Validates pagination parameters (limit 1-100, offset >= 0)
- Retrieves paginated orders

**UpdateOrder(orderID string, status string) error**
- Validates order ID and status
- Checks status against allowed values: cart, pending, preparing, served, cancelled

**CreateOrderItem(itemID string, quantity int, orderID string) error**
- Validates all inputs (item ID, quantity > 0, order ID)
- **Cross-domain validation**: Checks menu item exists via MenuService
- **Availability check**: Ensures menu item status is "available"
- Creates order item with generated UUID

**GetOrderItems(orderID string) ([]*models.OrderItems, error)**
- Validates order ID
- Retrieves associated items

**GetOrdersBySession(sessionID string) ([]*models.Order, error)**
- Validates session ID
- Retrieves all orders for session

**GetOrderItemsBySessionID(sessionID string) ([]*models.OrderItems, error)**
- **Service-level orchestration**: Combines multiple repository calls
- Step 1: Calls `repo.GetOrdersBySession(sessionID)` to get all orders
- Step 2: Extracts order IDs from the orders
- Step 3: Calls `repo.GetOrderItemsByOrderIDs(orderIDs)` to get all items
- Returns empty slice if session has no orders
- Respects domain boundaries by avoiding cross-domain joins in repository

### Menu Service

#### Interface: `MenuService`
```go
type MenuService interface {
    CreateMenuItem(Name string, Description string, Price float64, Category string, AvalabilityStatus models.ItemStatus) error
    GetMenuItem(id string) (*models.MenuItem, error)
    ListMenuItems() ([]*models.MenuItem, error)
    UpdateMenuItem(id string, name string, desc string, category string, price float64, status models.ItemStatus) error
    DeleteMenuItem(id string) error
    GetMenuItemsByCategory(category string) ([]*models.MenuItem, error)
    ListCategories() ([]string, error)
    CreateCategory(name string) error
}
```

#### Methods Documentation

**CreateMenuItem(...) error**
- Validates all inputs (name, description, price > 0, category, status)
- Length limits: name <= 100 chars, description <= 500 chars
- Auto-creates category if needed
- Generates UUID and timestamp

**GetMenuItem(id string) (*models.MenuItem, error)**
- Validates ID
- Retrieves menu item

**ListMenuItems() ([]*models.MenuItem, error)**
- Retrieves all menu items

**UpdateMenuItem(...) error**
- Validates all inputs with same rules as create
- Ensures category exists
- Updates menu item details

**DeleteMenuItem(id string) error**
- Validates ID
- Removes menu item

**GetMenuItemsByCategory(category string) ([]*models.MenuItem, error)**
- Validates category
- Filters items by category

**ListCategories() ([]string, error)**
- Retrieves all category names

**CreateCategory(name string) error**
- Validates category name
- Prevents duplicates
- Creates new category

## Validation Rules

### Common Validations
- **IDs**: Must be non-empty strings
- **Strings**: Trimmed and checked for emptiness
- **Numbers**: Range validation where applicable

### Order Status Validation
Allowed values: "cart", "pending", "preparing", "served", "cancelled"

### Menu Item Validation
- **Name**: Required, <= 100 characters
- **Description**: Required, <= 500 characters
- **Price**: Required, > 0
- **Category**: Required, auto-created if missing
- **Status**: Required

### Pagination Validation
- **Limit**: 1-100
- **Offset**: >= 0

## Error Handling

### Service Layer
- Validates inputs before repository calls
- Returns user-friendly error messages
- Handles business logic errors

## Usage Examples

### Creating a Session
```go
service := NewSessionService(repo)
sessionID, err := service.CreateSession(5) // Table 5
```

### Creating an Order
```go
orderService := NewOrderService(orderRepo)
err := orderService.CreateOrder(sessionID)
```

### Adding Menu Items
```go
menuService := NewMenuService(menuRepo)
err := menuService.CreateMenuItem("Burger", "Delicious burger", 9.99, "Main Course", models.StatusAvailable)
```

## Future Enhancements

- Transaction management for complex operations
- Caching layer for frequently accessed data
- Advanced filtering and search capabilities
- Audit logging for changes
- Soft deletes for data preservation

## Dependencies

- `database/sql`: Standard SQL interface
- `github.com/lib/pq`: PostgreSQL driver
- `github.com/google/uuid`: UUID generation
- Custom models package for data structures

This documentation covers the complete Service layer implementation.</content>
<parameter name="filePath">/home/jntbatra/Resturant/documentation/service-docs.md