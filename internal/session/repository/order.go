package repository

import (
	"database/sql"

	"restaurant/internal/session/models"

	"github.com/lib/pq"
)

// OrderRepository defines methods for order database operations
type OrderRepository interface {
	// CreateOrder creates a new order
	CreateOrder(order *models.Order) error

	// GetOrder retrieves an order by ID
	GetOrder(id string) (*models.Order, error)

	// ListOrders lists all orders
	ListOrders(limit int, offset int) ([]*models.Order, error)

	// UpdateOrder updates an order
	UpdateOrder(orderID string, status string) error

	// CreateOrderItem creates a new order item
	CreateOrderItem(item *models.OrderItems) error

	// GetOrderItems retrieves order items by order ID
	GetOrderItems(orderID string) ([]*models.OrderItems, error)

	// GetOrdersBySession retrieves orders by session ID
	GetOrdersBySession(sessionID string) ([]*models.Order, error)

	// GetOrderItemsByOrderIDs retrieves order items by multiple order IDs
	GetOrderItemsByOrderIDs(orderIDs []string) ([]*models.OrderItems, error)
}

// postgresOrderRepository implements OrderRepository using PostgreSQL
type postgresOrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new PostgreSQL-based order repository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

// Implementations (stubs for now)

// CreateOrder inserts a new order into the database
func (r *postgresOrderRepository) CreateOrder(order *models.Order) error {
	// Execute INSERT query with order details
	_, err := r.db.Exec("INSERT INTO orders (id, session_id, status, created_at) VALUES ($1, $2, $3, $4)", order.ID, order.SessionID, order.Status, order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetOrder retrieves an order by ID from the database
func (r *postgresOrderRepository) GetOrder(id string) (*models.Order, error) {
	var order models.Order
	// Execute SELECT query and scan result
	err := r.db.QueryRow("SELECT id, session_id, status, created_at FROM orders WHERE id = $1", id).Scan(&order.ID, &order.SessionID, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// ListOrders lists all orders with pagination
func (r *postgresOrderRepository) ListOrders(limit int, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	// Execute SELECT query with LIMIT and OFFSET
	rows, err := r.db.Query("SELECT id, session_id, status, created_at FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through rows and scan into order structs
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.SessionID, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// CreateOrderItem creates a new order item in the database
func (r *postgresOrderRepository) CreateOrderItem(item *models.OrderItems) error {
	// Execute INSERT query for order item
	_, err := r.db.Exec("INSERT INTO order_items (id, order_id, menu_item_id, quantity) VALUES ($1, $2, $3, $4)", item.ID, item.OrderID, item.MenuItemID, item.Quantity)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrder updates an order status in the database
func (r *postgresOrderRepository) UpdateOrder(orderID string, status string) error {
	// Execute UPDATE query
	_, err := r.db.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderItems retrieves order items by order ID
func (r *postgresOrderRepository) GetOrderItems(orderID string) ([]*models.OrderItems, error) {
	// Execute SELECT query
	rows, err := r.db.Query("SELECT id, order_id, menu_item_id, quantity FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.OrderItems
	// Iterate through rows and scan into item structs
	for rows.Next() {
		var item models.OrderItems
		err := rows.Scan(&item.ID, &item.OrderID, &item.MenuItemID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// GetOrdersBySession retrieves orders by session ID
func (r *postgresOrderRepository) GetOrdersBySession(sessionID string) ([]*models.Order, error) {
	// Execute SELECT query
	rows, err := r.db.Query("SELECT id, session_id, status, created_at FROM orders WHERE session_id = $1", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	// Iterate through rows and scan into order structs
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.SessionID, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderItemsByOrderIDs retrieves order items by multiple order IDs
func (r *postgresOrderRepository) GetOrderItemsByOrderIDs(orderIDs []string) ([]*models.OrderItems, error) {
	// Execute SELECT query using ANY with array parameter
	rows, err := r.db.Query("SELECT id, order_id, menu_item_id, quantity FROM order_items WHERE order_id = ANY($1)", pq.Array(orderIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.OrderItems
	// Iterate through rows and scan into item structs
	for rows.Next() {
		var item models.OrderItems
		err := rows.Scan(&item.ID, &item.OrderID, &item.MenuItemID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
