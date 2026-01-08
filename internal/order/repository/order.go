package repository

import (
	"context"
	"database/sql"

	"restaurant/internal/errors"
	"restaurant/internal/order/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Transaction provides methods for executing operations within a transaction
type Transaction interface {
	// Rollback cancels the transaction
	Rollback() error
	// Commit commits the transaction
	Commit() error
}

// OrderRepository defines methods for order database operations
type OrderRepository interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, order *models.Order) error

	// GetOrder retrieves an order by ID
	GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error)

	// ListOrders lists all orders
	ListOrders(ctx context.Context, limit int, offset int) ([]*models.Order, error)

	// UpdateOrder updates an order
	UpdateOrder(ctx context.Context, orderID uuid.UUID, status string) error

	// CreateOrderItem creates a new order item
	CreateOrderItem(ctx context.Context, item *models.OrderItems) error

	// UpdateOrderItemQuantity updates the quantity of an existing order item
	UpdateOrderItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error

	// GetOrderItems retrieves order items by order ID
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItems, error)

	// GetOrdersBySession retrieves orders by session ID
	GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.Order, error)

	// GetOrderItemsByOrderIDs retrieves order items by multiple order IDs
	GetOrderItemsByOrderIDs(ctx context.Context, orderIDs []uuid.UUID) ([]*models.OrderItems, error)

	// BeginTx begins a new database transaction
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

// TxOrderRepository provides transaction-aware order operations
type TxOrderRepository interface {
	// CreateOrderWithItems atomically creates an order and its items in a transaction
	CreateOrderWithItems(ctx context.Context, order *models.Order, items []*models.OrderItems, tx *sql.Tx) error

	// UpdateOrderItemsInTx updates multiple order items within a transaction
	UpdateOrderItemsInTx(ctx context.Context, items []*models.OrderItems, tx *sql.Tx) error
}

// postgresOrderRepository implements OrderRepository and TxOrderRepository using PostgreSQL
type postgresOrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new PostgreSQL-based order repository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

// NewTxOrderRepository creates a new transaction-aware order repository
func NewTxOrderRepository(db *sql.DB) TxOrderRepository {
	return &postgresOrderRepository{db: db}
}

// BeginTx begins a new database transaction
func (r *postgresOrderRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WrapError(500, "failed to begin transaction", err)
	}
	return tx, nil
}

// CreateOrder inserts a new order into the database
func (r *postgresOrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	// Execute INSERT query with order details
	_, err := r.db.ExecContext(ctx, "INSERT INTO orders (id, session_id, status, created_at) VALUES ($1, $2, $3, $4)", order.ID, order.SessionID, order.Status, order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetOrder retrieves an order by ID from the database
func (r *postgresOrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	// Execute SELECT query and scan result
	err := r.db.QueryRowContext(ctx, "SELECT id, session_id, status, created_at FROM orders WHERE id = $1", id).Scan(&order.ID, &order.SessionID, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrOrderNotFound
		}
		return nil, errors.NewInternalError("failed to get order", err)
	}
	return &order, nil
}

// ListOrders lists all orders with pagination
func (r *postgresOrderRepository) ListOrders(ctx context.Context, limit int, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	// Execute SELECT query with LIMIT and OFFSET
	rows, err := r.db.QueryContext(ctx, "SELECT id, session_id, status, created_at FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
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
func (r *postgresOrderRepository) CreateOrderItem(ctx context.Context, item *models.OrderItems) error {
	// Execute INSERT query for order item
	_, err := r.db.ExecContext(ctx, "INSERT INTO order_items (id, order_id, menu_item_id, quantity) VALUES ($1, $2, $3, $4)", item.ID, item.OrderID, item.MenuItemID, item.Quantity)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrder updates an order status in the database
func (r *postgresOrderRepository) UpdateOrder(ctx context.Context, orderID uuid.UUID, status string) error {
	// Execute UPDATE query
	_, err := r.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderItems retrieves order items by order ID
func (r *postgresOrderRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItems, error) {
	// Execute SELECT query
	rows, err := r.db.QueryContext(ctx, "SELECT id, order_id, menu_item_id, quantity FROM order_items WHERE order_id = $1", orderID)
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

// UpdateOrderItemQuantity updates the quantity of an existing order item
func (r *postgresOrderRepository) UpdateOrderItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE order_items SET quantity = $1 WHERE id = $2", quantity, itemID)
	if err != nil {
		return err
	}
	return nil
}

// CreateOrderWithItems atomically creates an order and its items in a transaction
func (r *postgresOrderRepository) CreateOrderWithItems(ctx context.Context, order *models.Order, items []*models.OrderItems, tx *sql.Tx) error {
	// Insert order within transaction
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, session_id, status, created_at) VALUES ($1, $2, $3, $4)",
		order.ID, order.SessionID, order.Status, order.CreatedAt,
	)
	if err != nil {
		return errors.WrapError(500, "failed to create order in transaction", err)
	}

	// Insert order items within same transaction
	for _, item := range items {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO order_items (id, order_id, menu_item_id, quantity) VALUES ($1, $2, $3, $4)",
			item.ID, item.OrderID, item.MenuItemID, item.Quantity,
		)
		if err != nil {
			return errors.WrapError(500, "failed to create order item in transaction", err)
		}
	}

	return nil
}

// UpdateOrderItemsInTx updates multiple order items within a transaction
func (r *postgresOrderRepository) UpdateOrderItemsInTx(ctx context.Context, items []*models.OrderItems, tx *sql.Tx) error {
	for _, item := range items {
		_, err := tx.ExecContext(
			ctx,
			"UPDATE order_items SET quantity = $1 WHERE id = $2",
			item.Quantity, item.ID,
		)
		if err != nil {
			return errors.WrapError(500, "failed to update order item in transaction", err)
		}
	}
	return nil
}

// GetOrdersBySession retrieves orders by session ID
func (r *postgresOrderRepository) GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.Order, error) {
	// Execute SELECT query
	rows, err := r.db.QueryContext(ctx, "SELECT id, session_id, status, created_at FROM orders WHERE session_id = $1", sessionID)
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
func (r *postgresOrderRepository) GetOrderItemsByOrderIDs(ctx context.Context, orderIDs []uuid.UUID) ([]*models.OrderItems, error) {
	// Execute SELECT query using ANY with array parameter
	rows, err := r.db.QueryContext(ctx, "SELECT id, order_id, menu_item_id, quantity FROM order_items WHERE order_id = ANY($1)", pq.Array(orderIDs))
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
