package handler

import (
	"fmt"
	"restaurant/internal/errors"
	"restaurant/internal/menu/models"
	"restaurant/internal/menu/service"
	"restaurant/internal/menu/validation"
	"restaurant/internal/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

// MenuHandler handles HTTP requests for menu items
type MenuHandler struct {
	svc service.MenuService
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

// RegisterRoutes registers all menu routes with the Gin router
func (h *MenuHandler) RegisterRoutes(router *gin.Engine) {
	menuGroup := router.Group("/menu")
	{
		menuGroup.GET("", h.ListMenuItems)
		menuGroup.POST("", h.CreateMenuItem)
		menuGroup.GET("/:id", h.GetMenuItem)
		menuGroup.GET("/category/:name", h.GetMenuItemsByCategory)
		menuGroup.PUT("/:id", h.UpdateMenuItem)
		menuGroup.DELETE("/:id", h.DeleteMenuItem)
	}
	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("", h.ListCategories)
		categoryGroup.POST("", h.CreateCategory)
		categoryGroup.GET("/:name", h.GetCategoryByName)
		categoryGroup.GET("/id/:id", h.GetCategoryByID)
		categoryGroup.PUT("/:name", h.UpdateCategory)
		categoryGroup.DELETE("/:name", h.DeleteCategory)
		categoryGroup.GET("/:name/id", h.CategoryIDByName)
	}
}

// CreateMenuItem handles POST /menu
// @Summary Create menu item
// @Description Create a new menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param request body validation.CreateMenuItemRequest true "Menu item creation request"
// @Success 201 {object} models.MenuItem
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu [post]
func (h *MenuHandler) CreateMenuItem(c *gin.Context) {
	var req validation.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	// Trim whitespace from string fields
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Category = strings.TrimSpace(req.Category)

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "in_stock"
	}

	// Debug log
	fmt.Printf("DEBUG: req = %+v\n", req)

	if err := validation.ValidateCreateMenuItem(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	item, err := h.svc.CreateMenuItem(c.Request.Context(), req.Name, req.Description, req.Price, req.Category, models.ItemStatus(req.Status))
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(201, item)
}

// GetMenuItem handles GET /menu/:id
// @Summary Get menu item
// @Description Retrieve a specific menu item by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID (UUID)"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/{id} [get]
func (h *MenuHandler) GetMenuItem(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.svc.GetMenuItem(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, item)
}

// GetMenuItemsByCategory handles GET /menu/category/:name
// @Summary Get menu items by category
// @Description Get all menu items for a specific category
// @Tags Menu
// @Accept json
// @Produce json
// @Param name path string true "Category name"
// @Success 200 {array} models.MenuItem
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/category/{name} [get]
func (h *MenuHandler) GetMenuItemsByCategory(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		middleware.HandleError(c, errors.ErrInvalidRequest)
		return
	}

	items, err := h.svc.GetMenuItemsByCategory(c.Request.Context(), name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, items)
}

// ListMenuItems handles GET /menu
// @Summary List menu items
// @Description List all menu items with optional filtering and pagination
// @Tags Menu
// @Accept json
// @Produce json
// @Param offset query int false "Offset (default 0)"
// @Param limit query int false "Limit (default 10, max 100)"
// @Param category query string false "Filter by category"
// @Success 200 {array} models.MenuItem
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu [get]
func (h *MenuHandler) ListMenuItems(c *gin.Context) {
	var req validation.ListMenuItemsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 10
	}

	if err := validation.ValidateListMenuItems(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	items, err := h.svc.ListMenuItems(c.Request.Context(), req.Offset, req.Limit)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, items)
}

// UpdateMenuItem handles PUT /menu/:id
// @Summary Update menu item
// @Description Update an existing menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID (UUID)"
// @Param request body validation.UpdateMenuItemRequest true "Menu item update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/{id} [put]
func (h *MenuHandler) UpdateMenuItem(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	var req validation.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	// Trim whitespace from string fields
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Category = strings.TrimSpace(req.Category)

	if err := validation.ValidateUpdateMenuItem(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	var err error
	err = h.svc.UpdateMenuItem(c.Request.Context(), id, req.Name, req.Description, req.Category, req.Price, models.ItemStatus(req.Status))
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Menu item updated successfully"})
}

// DeleteMenuItem handles DELETE /menu/:id
// @Summary Delete menu item
// @Description Delete a menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/{id} [delete]
func (h *MenuHandler) DeleteMenuItem(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	err := h.svc.DeleteMenuItem(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(204, gin.H{"message": "Menu item deleted successfully"})
}

// ListCategories handles GET /menu/categories
// @Summary List categories
// @Description List all menu categories
// @Tags Menu
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories [get]
func (h *MenuHandler) ListCategories(c *gin.Context) {
	categories, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, categories)
}

// GetCategoryByName handles GET /categories/:name
// @Summary Get category by name
// @Description Get category ID by name
// @Tags Menu
// @Accept json
// @Produce json
// @Param name path string true "Category name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories/{name} [get]
func (h *MenuHandler) GetCategoryByName(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		middleware.HandleError(c, errors.NewValidationError("category name is required"))
		return
	}

	id, err := h.svc.CategoryIDByName(c.Request.Context(), name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{"name": name, "id": id})
}

// GetCategoryByID handles GET /menu/categories/id/:id
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} models.Category
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories/id/{id} [get]
func (h *MenuHandler) GetCategoryByID(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	category, err := h.svc.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, category)
}

// CreateCategory handles POST /menu/categories
// @Summary Create a new category
// @Description Create a new menu category
// @Tags Menu
// @Accept json
// @Produce json
// @Param request body validation.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} models.Category
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories [post]
func (h *MenuHandler) CreateCategory(c *gin.Context) {
	var req validation.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateCreateCategory(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	category, err := h.svc.CreateCategory(c.Request.Context(), req.Name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(201, category)
}

// UpdateCategory handles PUT /menu/categories/:name
// @Summary Update category
// @Description Update an existing category
// @Tags Menu
// @Accept json
// @Produce json
// @Param name path string true "Current category name"
// @Param request body validation.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} models.Category
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories/{name} [put]
func (h *MenuHandler) UpdateCategory(c *gin.Context) {
	old_name := c.Param("name")

	var req validation.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateUpdateCategory(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	category, err := h.svc.UpdateCategory(c.Request.Context(), old_name, req.Name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, category)
}

// DeleteCategory handles DELETE /menu/categories/:name
// @Summary Delete category
// @Description Delete a menu category
// @Tags Menu
// @Accept json
// @Produce json
// @Param name path string true "Category name"
// @Success 200 {object} map[string]string
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories/{name} [delete]
func (h *MenuHandler) DeleteCategory(c *gin.Context) {
	name := c.Param("name")
	err := h.svc.DeleteCategory(c.Request.Context(), name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(204, gin.H{"message": "Category deleted successfully"})
}

// CategoryIDByName handles GET /menu/categories/:name/id
// @Summary Get category ID by name
// @Description Get category ID by name
// @Tags Menu
// @Accept json
// @Produce json
// @Param name path string true "Category name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /menu/categories/{name}/id [get]
func (h *MenuHandler) CategoryIDByName(c *gin.Context) {
	name := c.Param("name")
	name = strings.TrimSpace(name)

	if name == "" {
		middleware.HandleError(c, errors.NewValidationError("category name is required"))
		return
	}

	id, err := h.svc.CategoryIDByName(c.Request.Context(), name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}
	c.JSON(200, gin.H{"id": id})
}
