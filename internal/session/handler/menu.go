package handler

import (
	"net/http"
	"restaurant/internal/session/service"
)

// MenuHandler handles HTTP requests for menu items
type MenuHandler struct {
	svc service.MenuService
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

// CreateMenuItem handles POST /menu
func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: Implement
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetMenuItem handles GET /menu/{id}
func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: Implement
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// ListMenuItems handles GET /menu
func (h *MenuHandler) ListMenuItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: Implement
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}