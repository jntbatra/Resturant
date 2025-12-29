package handler

import (
	"net/http"
	"restaurant/internal/session/service"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	svc service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: Implement
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetOrder handles GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: Implement
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: Implement
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}