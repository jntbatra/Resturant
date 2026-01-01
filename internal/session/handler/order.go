package handler

import (
	"restaurant/internal/session/service"
	"restaurant/internal/session/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req validation.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateCreateOrder(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	order, err := h.svc.CreateOrder(c.Request.Context(), req.SessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, order)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if order == nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(200, order)
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var req validation.ListOrdersRequest
	req.Offset = 0
	req.Limit = 10

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateListOrders(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.svc.ListOrders(c.Request.Context(), req.Limit, req.Offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

// UpdateOrder handles PUT /orders/:id
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	var req validation.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateUpdateOrder(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.UpdateOrder(c.Request.Context(), id, string(req.Status))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Order updated successfully"})
}

// CreateOrderItem handles POST /orders/:id/items
func (h *OrderHandler) CreateOrderItem(c *gin.Context) {
	var req validation.CreateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateCreateOrderItem(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.CreateOrderItem(c.Request.Context(), req.MenuItemID, req.Quantity, req.OrderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, item)
}

// GetOrderItems handles GET /orders/:id/items
func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	items, err := h.svc.GetOrderItems(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, items)
}

// GetOrdersBySession handles GET /sessions/:id/orders
func (h *OrderHandler) GetOrdersBySession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}

	orders, err := h.svc.GetOrdersBySession(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

// GetOrderItemsByOrderIDs handles batch retrieval of order items
func (h *OrderHandler) GetOrderItemsByOrderIDs(c *gin.Context) {
	// TODO: Implement batch order items retrieval
	c.JSON(501, gin.H{"error": "Not implemented"})
}

// GetOrderItemsBySessionIDs handles GET /sessions/:id/order-items
func (h *OrderHandler) GetOrderItemsBySessionIDs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}

	items, err := h.svc.GetOrderItemsBySessionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, items)
}

// UpdateOrderItem handles PUT /orders/:id/items/:itemId
func (h *OrderHandler) UpdateOrderItem(c *gin.Context) {
	var req validation.UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateUpdateOrderItem(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement order item update in service
	c.JSON(501, gin.H{"error": "Order item update not yet implemented"})
}

// DeleteOrderItem handles DELETE /orders/:id/items/:itemId
func (h *OrderHandler) DeleteOrderItem(c *gin.Context) {
	// TODO: Implement order item deletion in service
	c.JSON(501, gin.H{"error": "Order item deletion not yet implemented"})
}
