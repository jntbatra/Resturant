package handler

import (
	"restaurant/internal/errors"
	"restaurant/internal/middleware"
	"restaurant/internal/session/service"
	"restaurant/internal/session/validation"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	svc service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// RegisterRoutes registers all order routes with the Gin router
func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.GET("", h.ListOrders)
		orderGroup.POST("", h.CreateOrder)
		orderGroup.GET("/:id", h.GetOrder)
		orderGroup.PUT("/:id", h.UpdateOrder)
		orderGroup.POST("/:id/items", h.CreateOrderItem)
		orderGroup.GET("/:id/items", h.GetOrderItems)
	}

	// Session-related order routes
	sessionGroup := router.Group("/sessions")
	{
		sessionGroup.GET("/:id/orders", h.GetOrdersBySession)
		sessionGroup.GET("/:id/order-items", h.GetOrderItemsBySessionIDs)
	}
}

// CreateOrder handles POST /orders
// @Summary Create order
// @Description Create a new order for a session
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body validation.CreateOrderRequest true "Order creation request"
// @Success 201 {object} models.Order
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req validation.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateCreateOrder(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	order, err := h.svc.CreateOrder(c.Request.Context(), req.SessionID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(201, order)
}

// GetOrder handles GET /orders/:id
// @Summary Get order
// @Description Retrieve a specific order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} models.Order
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, order)
}

// ListOrders handles GET /orders
// @Summary List orders
// @Description List all orders with pagination
// @Tags Orders
// @Accept json
// @Produce json
// @Param offset query int false "Offset (default 0)"
// @Param limit query int false "Limit (default 10, max 100)"
// @Param session_id query string false "Filter by session ID"
// @Success 200 {array} models.Order
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var req validation.ListOrdersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	// Set defaults only if not provided
	if req.Offset == 0 && c.Query("offset") == "" {
		req.Offset = 0
	}
	if req.Limit == 0 && c.Query("limit") == "" {
		req.Limit = 10
	}

	if err := validation.ValidateListOrders(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	orders, err := h.svc.ListOrders(c.Request.Context(), req.Limit, req.Offset)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, orders)
}

// UpdateOrder handles PUT /orders/:id
// @Summary Update order status
// @Description Update the status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Param request body validation.UpdateOrderRequest true "Status update request"
// @Success 200 {object} models.Order
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	var req validation.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateUpdateOrder(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	order, err := h.svc.UpdateOrder(c.Request.Context(), id, string(req.Status))
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, order)
}

// CreateOrderItem handles POST /orders/:id/items
// @Summary Add item to order
// @Description Add a menu item to an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Param request body validation.CreateOrderItemRequest true "Order item creation request"
// @Success 201 {object} models.OrderItems
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id}/items [post]
func (h *OrderHandler) CreateOrderItem(c *gin.Context) {
	orderID, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	var req validation.CreateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateCreateOrderItem(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	item, err := h.svc.CreateOrderItem(c.Request.Context(), req.MenuItemID, req.Quantity, orderID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(201, item)
}

// GetOrderItems handles GET /orders/:id/items
// @Summary Get order items
// @Description Retrieve all items in an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {array} models.OrderItems
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id}/items [get]
func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	items, err := h.svc.GetOrderItems(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, items)
}

// GetOrdersBySession handles GET /sessions/:id/orders
// @Summary Get orders for session
// @Description Retrieve all orders for a specific session
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {array} models.Order
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id}/orders [get]
func (h *OrderHandler) GetOrdersBySession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	orders, err := h.svc.GetOrdersBySession(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, orders)
}

// GetOrderItemsBySessionIDs handles GET /sessions/:id/order-items
// @Summary Get order items for session
// @Description Retrieve all order items for a specific session
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {array} models.OrderItems
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id}/order-items [get]
func (h *OrderHandler) GetOrderItemsBySessionIDs(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	items, err := h.svc.GetOrderItemsBySessionID(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, items)
}
