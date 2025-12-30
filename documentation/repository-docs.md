# Restaurant Management System - Repository Layer Documentation

## Overview

This document provides comprehensive documentation for the Repository layer of the Restaurant Management System. The Repository layer handles data access and persistence operations using Go with PostgreSQL.

## Architecture

### Layered Architecture
```
Handler (HTTP API) -> Service (Business Logic) -> Repository (Data Access) -> Database
```

### Domain Separation Principle
The repository layer strictly adheres to domain boundaries:
- **Order Repository**: Only queries `orders` and `order_items` tables
- **Menu Repository**: Only queries `menu_items` and `categories` tables  
- **Session Repository**: Only queries `sessions` and `tables` tables

**Benefits:**
- Clear separation of concerns
- Easy to split into microservices later
- No cross-domain JOINs in repository layer
- Service layer orchestrates cross-domain operations

### Key Technologies
- **Language**: Go
- **Database**: PostgreSQL
- **Driver**: `github.com/lib/pq`
- **UUID Generation**: `github.com/google/uuid`

## Repository Layer

The Repository layer provides data access abstractions using interfaces and PostgreSQL implementations.

### Common Patterns
- All repositories follow the Repository pattern with interfaces
- Use `*sql.DB` for database connections
- Implement CRUD operations with proper error handling
- Use prepared statements with parameterized queries for security
- **Domain Separation**: Each repository only queries its own domain tables
- Cross-domain operations are orchestrated at the service layer

### Session Repository

#### Interface: `SessionRepository`
```go
type SessionRepository interface {
    CreateSession(session *models.Session) error
    GetSession(id string) (*models.Session, error)
    UpdateSessionStatus(id string, status models.SessionStatus) error
    ListActiveSessions() ([]*models.Session, error)
    ChangeTable(id string, tableNumber int) error
}
```

#### Methods Documentation

**CreateSession(session *models.Session) error**
- Inserts a new session into the `sessions` table
- Generates UUID for session ID
- Sets initial status to "active"
- Records creation timestamp

**GetSession(id string) (*models.Session, error)**
- Retrieves a session by ID
- Returns nil if session not found (no error)
- Includes table information via JOIN

**UpdateSessionStatus(id string, status models.SessionStatus) error**
- Updates session status and completion timestamp
- Sets `completed_at` when status changes from active

**ListActiveSessions() ([]*models.Session, error)**
- Retrieves all sessions with "active" status
- Ordered by creation date descending

**ChangeTable(id string, tableNumber int) error**
- Updates the table assignment for a session
- Validates table exists in `tables` table

### Order Repository

#### Interface: `OrderRepository`
```go
type OrderRepository interface {
    CreateOrder(order *models.Order) error
    GetOrder(id string) (*models.Order, error)
    ListOrders(limit int, offset int) ([]*models.Order, error)
    UpdateOrder(orderID string, status string) error
    CreateOrderItem(item *models.OrderItems) error
    GetOrderItems(orderID string) ([]*models.OrderItems, error)
    GetOrdersBySession(sessionID string) ([]*models.Order, error)
    GetOrderItemsByOrderIDs(orderIDs []string) ([]*models.OrderItems, error)
}
```

#### Methods Documentation

**CreateOrder(order *models.Order) error**
- Inserts a new order into the `orders` table
- Links order to a session

**GetOrder(id string) (*models.Order, error)**
- Retrieves order details by ID
- Returns nil if not found

**ListOrders(limit int, offset int) ([]*models.Order, error)**
- Retrieves orders with pagination
- Ordered by creation date descending

**UpdateOrder(orderID string, status string) error**
- Updates order status

**CreateOrderItem(item *models.OrderItems) error**
- Inserts order item into `order_items` table

**GetOrderItems(orderID string) ([]*models.OrderItems, error)**
- Retrieves all items for a specific order

**GetOrdersBySession(sessionID string) ([]*models.Order, error)**
- Retrieves all orders for a session

**GetOrderItemsByOrderIDs(orderIDs []string) ([]*models.OrderItems, error)**
- Bulk retrieval of order items for multiple orders
- Uses PostgreSQL `ANY` operator with array parameter
- Enables efficient batch retrieval without cross-domain queries

### Menu Repository

#### Interface: `MenuRepository`
```go
type MenuRepository interface {
    CreateMenuItem(item *models.MenuItem) error
    GetMenuItem(id string) (*models.MenuItem, error)
    ListMenuItems() ([]*models.MenuItem, error)
    UpdateMenuItem(item *models.MenuItem) error
    DeleteMenuItem(id string) error
    GetMenuItemsByCategory(category string) ([]*models.MenuItem, error)
    ListCategories() ([]string, error)
    CreateCategory(name string) error
}
```

#### Methods Documentation

**CreateMenuItem(item *models.MenuItem) error**
- Inserts menu item into `menu_items` table
- Auto-creates category if it doesn't exist

**GetMenuItem(id string) (*models.MenuItem, error)**
- Retrieves menu item by ID

**ListMenuItems() ([]*models.MenuItem, error)**
- Retrieves all menu items

**UpdateMenuItem(item *models.MenuItem) error**
- Updates menu item details

**DeleteMenuItem(id string) error**
- Removes menu item from database

**GetMenuItemsByCategory(category string) ([]*models.MenuItem, error)**
- Filters menu items by category

**ListCategories() ([]string, error)**
- Retrieves all category names

**CreateCategory(name string) error**
- Inserts new category into `categories` table

## Database Schema

### Tables
- `tables`: Table management
- `sessions`: Customer sessions
- `orders`: Order headers
- `order_items`: Order line items
- `menu_items`: Menu catalog
- `categories`: Menu categories

### Relationships
- Sessions reference tables
- Orders reference sessions
- Order items reference orders and menu items
- Menu items reference categories

## Error Handling

### Repository Layer
- Returns database errors directly
- Handles `sql.ErrNoRows` for not found cases
- Uses parameterized queries to prevent SQL injection

## Dependencies

- `database/sql`: Standard SQL interface
- `github.com/lib/pq`: PostgreSQL driver
- Custom models package for data structures

This documentation covers the complete Repository layer implementation.</content>
<parameter name="filePath">/home/jntbatra/Resturant/documentation/repository-docs.md