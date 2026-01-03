package handler

import (
	"net/http"
	"strconv"

	"restaurant/internal/errors"
	"restaurant/internal/middleware"
	"restaurant/internal/session/models"
	"restaurant/internal/session/service"
	"restaurant/internal/session/validation"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for sessions
type Handler struct {
	svc service.SessionService
}

// NewHandler creates a new handler
func NewHandler(svc service.SessionService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all session routes with the Gin router
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	sessionGroup := router.Group("/sessions")
	{
		sessionGroup.GET("", h.ListSessions)
		sessionGroup.POST("", h.CreateSession)
		sessionGroup.GET("/active", h.ListActiveSessions)
		sessionGroup.GET("/:id", h.GetSession)
		sessionGroup.PUT("/:id", h.UpdateSession)
		sessionGroup.PUT("/:id/table", h.ChangeSessionTable)
		sessionGroup.GET("/table/:tableID", h.GetSessionsByTable)
		sessionGroup.GET("/table/:tableID/active", h.GetActiveSessionsByTable)
		sessionGroup.DELETE("/:id", h.DeleteSession)
	}

	tableGroup := router.Group("/tables")
	{
		tableGroup.GET("", h.ListTables)
		tableGroup.POST("", h.CreateTable)
		tableGroup.GET("/:id", h.GetTable)
		tableGroup.PUT("/:id", h.UpdateTable)
		tableGroup.DELETE("/:id", h.DeleteTable)
	}
}

// CreateSession handles POST /sessions
// @Summary Create a new session
// @Description Create a new dining session for a table
// @Tags Sessions
// @Accept json
// @Produce json
// @Param request body validation.CreateSessionRequest true "Session creation request"
// @Success 200 {object} models.Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions [post]
func (h *Handler) CreateSession(c *gin.Context) {
	var req validation.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateCreateSession(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	session, err := h.svc.CreateSession(c.Request.Context(), req.TableID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, session)
}

// GetSession handles GET /sessions/:id
// @Summary Get session by ID
// @Description Retrieve a specific session by its ID
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} models.Session
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id} [get]
func (h *Handler) GetSession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := validation.ValidateSessionID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	session, err := h.svc.GetSession(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(200, session)
}

// UpdateSession handles PUT /sessions/:id/status
// @Summary Update session status
// @Description Update the status of a session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Param request body validation.UpdateSessionRequest true "Status update request"
// @Success 200 {object} models.Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id} [put]
func (h *Handler) UpdateSession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := validation.ValidateSessionID(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var req validation.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateUpdateSession(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updatedSession, err := h.svc.UpdateSession(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, updatedSession)
}

// ListSessions handles GET /sessions
// @Summary List sessions
// @Description List all sessions with pagination
// @Tags Sessions
// @Accept json
// @Produce json
// @Param offset query int false "Offset (default 0)"
// @Param limit query int false "Limit (default 10, max 100)"
// @Success 200 {array} models.Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions [get]
func (h *Handler) ListSessions(c *gin.Context) {
	var req validation.ListSessionsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Set defaults only if not provided
	if req.Offset == 0 && c.Query("offset") == "" {
		req.Offset = 0
	}
	if req.Limit == 0 && c.Query("limit") == "" {
		req.Limit = 10
	}

	if err := validation.ValidateListSessions(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sessions, err := h.svc.ListSessions(c.Request.Context(), req.Offset, req.Limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, sessions)
}

// ListActiveSessions handles GET /sessions/active
func (h *Handler) ListActiveSessions(c *gin.Context) {
	sessions, err := h.svc.ListActiveSessions(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, sessions)
}

// ChangeSessionTable handles PUT /sessions/:id/table
// @Summary Change session table
// @Description Change the table assignment for a session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Param request body validation.ChangeSessionTableRequest true "Table change request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id}/table [put]
func (h *Handler) ChangeSessionTable(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := validation.ValidateSessionID(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var req validation.ChangeSessionTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateChangeSessionTable(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.ChangeTable(c.Request.Context(), id, req.TableID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Session table changed successfully"})
}

// GetSessionsByTable handles GET /sessions/table/:tableID
// @Summary Get sessions for a table
// @Description Retrieve all sessions for a specific table
// @Tags Sessions
// @Accept json
// @Produce json
// @Param tableID path int true "Table ID"
// @Success 200 {array} models.Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/table/{tableID} [get]
func (h *Handler) GetSessionsByTable(c *gin.Context) {
	tableIDStr := c.Param("tableID")
	tableID, err := strconv.Atoi(tableIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid table ID - must be an integer"})
		return
	}

	sessions, err := h.svc.GetSessionsByTable(c.Request.Context(), tableID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, sessions)
}

// GetActiveSessionsByTable handles GET /sessions/table/:tableID/active
func (h *Handler) GetActiveSessionsByTable(c *gin.Context) {
	tableIDStr := c.Param("tableID")
	tableID, err := strconv.Atoi(tableIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid table ID - must be an integer"})
		return
	}

	sessions, err := h.svc.GetActiveSessionsByTable(c.Request.Context(), tableID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, sessions)
}

// DeleteSession handles DELETE /sessions/:id
// @Summary Delete a session
// @Description Delete a session by ID
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id} [delete]
func (h *Handler) DeleteSession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := validation.ValidateSessionID(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.DeleteSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Session deleted successfully"})
}

// CreateTable handles POST /tables
// @Summary Create a new table
// @Description Create a new restaurant table
// @Tags Tables
// @Accept json
// @Produce json
// @Param request body validation.CreateTableRequest true "Table creation request"
// @Success 201 {object} models.Table
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables [post]
func (h *Handler) CreateTable(c *gin.Context) {
	var req validation.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateCreateTable(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	table, err := h.svc.CreateTable(c.Request.Context(), &models.CreateTableRequest{ID: req.ID})
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, table)
}

// GetTable handles GET /tables/:id
// @Summary Get table by ID
// @Description Retrieve a specific table by its ID
// @Tags Tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Success 200 {object} models.Table
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables/{id} [get]
func (h *Handler) GetTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(c, errors.NewValidationError("invalid table ID"))
		return
	}

	if err := validation.ValidateTableID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	table, err := h.svc.GetTable(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, table)
}

// ListTables handles GET /tables
// @Summary List all tables
// @Description Retrieve all restaurant tables
// @Tags Tables
// @Accept json
// @Produce json
// @Success 200 {array} models.Table
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables [get]
func (h *Handler) ListTables(c *gin.Context) {
	tables, err := h.svc.ListTables(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tables)
}

// UpdateTable handles PUT /tables/:id
// @Summary Update a table
// @Description Update an existing restaurant table
// @Tags Tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Param request body validation.UpdateTableRequest true "Table update request"
// @Success 200 {object} models.Table
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables/{id} [put]
func (h *Handler) UpdateTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(c, errors.NewValidationError("invalid table ID"))
		return
	}

	if err := validation.ValidateTableID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	var req validation.UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := validation.ValidateUpdateTable(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	table, err := h.svc.UpdateTable(c.Request.Context(), id, &models.UpdateTableRequest{ID: req.ID})
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, table)
}

// DeleteTable handles DELETE /tables/:id
// @Summary Delete a table
// @Description Delete a restaurant table (only if no active sessions)
// @Tags Tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Success 204 "No Content"
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables/{id} [delete]
func (h *Handler) DeleteTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(c, errors.NewValidationError("invalid table ID"))
		return
	}

	if err := validation.ValidateTableID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	err = h.svc.DeleteTable(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
