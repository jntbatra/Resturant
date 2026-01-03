package handler

import (
	"restaurant/internal/errors"
	"restaurant/internal/middleware"
	"restaurant/internal/session/models"
	"restaurant/internal/session/service"
	"restaurant/internal/session/validation"

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
		menuGroup.PUT("/:id", h.UpdateMenuItem)
		menuGroup.DELETE("/:id", h.DeleteMenuItem)

		categoryGroup := menuGroup.Group("/categories")
		{
			categoryGroup.GET("", h.ListCategories)
			categoryGroup.POST("", h.CreateCategory)
			categoryGroup.GET("/:name", h.GetCategoryByName)
			categoryGroup.PUT("/:name", h.UpdateCategory)
			categoryGroup.DELETE("/:name", h.DeleteCategory)
			categoryGroup.GET("/:name/id", h.CategoryIDByName)
			categoryGroup.GET("/:name/id_or_create", h.CategoryIDorCreate)
		}
	}
}

// CreateMenuItem handles POST /menu
func (h *MenuHandler) CreateMenuItem(c *gin.Context) {
	var req validation.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

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

// ListMenuItems handles GET /menu
func (h *MenuHandler) ListMenuItems(c *gin.Context) {
	var req validation.ListMenuItemsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
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

	if err := validation.ValidateUpdateMenuItem(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	var err error
	err = h.svc.UpdateMenuItem(c.Request.Context(), id, req.Name, req.Description, req.Category, req.Price, models.ItemStatus("in_stock"))
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Menu item updated successfully"})
}

// DeleteMenuItem handles DELETE /menu/:id
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

	c.JSON(200, gin.H{"message": "Menu item deleted successfully"})
}

// ListCategories handles GET /menu/categories
func (h *MenuHandler) ListCategories(c *gin.Context) {
	categories, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, categories)
}

// GetCategoryByName handles GET /categories/:name
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

// CreateCategory handles POST /menu/categories
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

	err := h.svc.UpdateCategory(c.Request.Context(), old_name, req.Name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory handles DELETE /menu/categories/:name
func (h *MenuHandler) DeleteCategory(c *gin.Context) {
	name := c.Param("name")
	err := h.svc.DeleteCategory(c.Request.Context(), name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Category deleted successfully"})
}

func (h *MenuHandler) CategoryIDByName(c *gin.Context) {
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
	c.JSON(200, gin.H{"id": id})
}

func (h *MenuHandler) CategoryIDorCreate(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		middleware.HandleError(c, errors.NewValidationError("category name is required"))
		return
	}

	id, err := h.svc.CategoryIDByNameCreateIfNotPresent(c.Request.Context(), name)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}
	c.JSON(200, gin.H{"id": id})
}
