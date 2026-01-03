package docs

// API Documentation
// @title Restaurant Management API
// @version 1.0
// @description Restaurant management system with sessions, menus, and orders
// @termsOfService http://example.com/terms

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @basePath /api/v1

// Session Endpoints

// @Summary Create a new session
// @Description Create a new dining session for a table
// @Tags Sessions
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "Session creation request"
// @Success 200 {object} models.Session
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions [post]

// @Summary Get session by ID
// @Description Retrieve a specific session by its ID
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} models.Session
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id} [get]

// @Summary List sessions
// @Description List all sessions with pagination
// @Tags Sessions
// @Accept json
// @Produce json
// @Param offset query int false "Offset (default 0)"
// @Param limit query int false "Limit (default 10, max 100)"
// @Success 200 {array} models.Session
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions [get]

// @Summary Update session status
// @Description Update the status of a session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Param request body UpdateSessionRequest true "Status update request"
// @Success 200 {object} models.Session
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id} [put]

// @Summary Change session table
// @Description Change the table assignment for a session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Param request body ChangeSessionTableRequest true "Table change request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/table [put]

// @Summary Get sessions for a table
// @Description Retrieve all sessions for a specific table
// @Tags Sessions
// @Accept json
// @Produce json
// @Param tableID path int true "Table ID"
// @Success 200 {array} models.Session
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/table/{tableID} [get]

// @Summary Delete a session
// @Description Delete a session by ID
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id} [delete]

// Menu Endpoints

// @Summary Create menu item
// @Description Create a new menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param request body CreateMenuItemRequest true "Menu item creation request"
// @Success 201 {object} models.MenuItem
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /menu [post]

// @Summary Get menu item
// @Description Retrieve a specific menu item by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID (UUID)"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /menu/{id} [get]

// @Summary List menu items
// @Description List all menu items with optional filtering and pagination
// @Tags Menu
// @Accept json
// @Produce json
// @Param offset query int false "Offset (default 0)"
// @Param limit query int false "Limit (default 10, max 100)"
// @Param category query string false "Filter by category"
// @Success 200 {array} models.MenuItem
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /menu [get]

// @Summary Update menu item
// @Description Update an existing menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID (UUID)"
// @Param request body UpdateMenuItemRequest true "Menu item update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /menu/{id} [put]

// @Summary Delete menu item
// @Description Delete a menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID (UUID)"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /menu/{id} [delete]

// Order Endpoints

// @Summary Create order
// @Description Create a new order for a session
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Order creation request"
// @Success 201 {object} models.Order
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]

// @Summary Get order
// @Description Retrieve a specific order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} models.Order
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [get]

// @Summary List orders
// @Description List all orders with pagination
// @Tags Orders
// @Accept json
// @Produce json
// @Param offset query int false "Offset (default 0)"
// @Param limit query int false "Limit (default 10, max 100)"
// @Param session_id query string false "Filter by session ID"
// @Success 200 {array} models.Order
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]

// @Summary Update order status
// @Description Update the status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Param request body UpdateOrderRequest true "Status update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [put]

// @Summary Add item to order
// @Description Add a menu item to an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Param request body CreateOrderItemRequest true "Order item creation request"
// @Success 201 {object} models.OrderItems
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/items [post]

// @Summary Get order items
// @Description Retrieve all items in an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {array} models.OrderItems
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/items [get]

// @Summary Get orders for session
// @Description Retrieve all orders for a specific session
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {array} models.Order
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sessions/{id}/orders [get]

// Response types for documentation
type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type CreateSessionRequest struct {
	TableID int `json:"table_id" validate:"required,gt=0"`
}

type UpdateSessionRequest struct {
	Status string `json:"status" validate:"required,oneof=active completed pending cancelled"`
}

type ChangeSessionTableRequest struct {
	TableID int `json:"table_id" validate:"required,gt=0"`
}

type CreateMenuItemRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required,min=1,max=100"`
	Status      string  `json:"status" validate:"required,oneof=available unavailable discontinued"`
}

type UpdateMenuItemRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=1,max=255"`
	Description string  `json:"description" validate:"omitempty,max=1000"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Category    string  `json:"category" validate:"omitempty,min=1,max=100"`
}

type CreateOrderRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" validate:"required,oneof=cart pending preparing served cancelled"`
}

type CreateOrderItemRequest struct {
	OrderID    string `json:"order_id" validate:"required"`
	MenuItemID string `json:"menu_item_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,gt=0"`
}
