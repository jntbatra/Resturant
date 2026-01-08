package session

import (
	"fmt"
	"net/http"
	"strconv"

	"restaurant/internal/errors"
	"restaurant/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for sessions
type Handler struct {
	svc SessionService
}

// NewHandler creates a new handler
func NewHandler(svc SessionService) *Handler {
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

		// Table management routes under sessions
		sessionGroup.GET("/tables", h.ListTables)
		sessionGroup.POST("/tables", h.CreateTable)
		sessionGroup.POST("/tables/bulk", h.BulkCreateTables)
		sessionGroup.GET("/tables/:id", h.GetTable)
		sessionGroup.DELETE("/tables/:id", h.DeleteTable)
	}
}

// CreateSession handles POST /sessions
// @Summary Create a new session
// @Description Create a new dining session for a table
// @Tags Sessions
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "Session creation request"
// @Success 200 {object} Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions [post]
func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := ValidateCreateSession(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	session, err := h.svc.CreateSession(c.Request.Context(), req.TableID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(201, session)
}

// GetSession handles GET /sessions/:id
// @Summary Get session by ID
// @Description Retrieve a specific session by its ID
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} Session
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id} [get]
func (h *Handler) GetSession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := ValidateSessionID(id); err != nil {
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
// @Param request body UpdateSessionRequest true "Status update request"
// @Success 200 {object} Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id} [put]
func (h *Handler) UpdateSession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := ValidateSessionID(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateUpdateSession(req); err != nil {
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
// @Success 200 {array} Session
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions [get]
func (h *Handler) ListSessions(c *gin.Context) {
	var req ListSessionsRequest

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

	if err := ValidateListSessions(req); err != nil {
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
// @Param request body ChangeSessionTableRequest true "Table change request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id}/table [put]
func (h *Handler) ChangeSessionTable(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := ValidateSessionID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	var req ChangeSessionTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := ValidateChangeSessionTable(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	err := h.svc.ChangeTable(c.Request.Context(), id, req.TableID)
	if err != nil {
		middleware.HandleError(c, err)
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
// @Success 200 {array} Session
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
// @Success 204 {string} string "No Content"
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/{id} [delete]
func (h *Handler) DeleteSession(c *gin.Context) {
	id, ok := middleware.UUIDParam(c, "id")
	if !ok {
		return
	}

	if err := ValidateSessionID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	err := h.svc.DeleteSession(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.Status(204) // No Content
}

// CreateTable handles POST /sessions/tables
// @Summary Create a new table
// @Description Create a new restaurant table
// @Tags Tables
// @Accept json
// @Produce json
// @Param request body CreateTableRequest true "Table creation request"
// @Success 201 {object} Table
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/tables [post]
func (h *Handler) CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := ValidateCreateTable(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	table, err := h.svc.CreateTable(c.Request.Context(), &CreateTableRequest{ID: req.ID})
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, table)
}

// BulkCreateTables handles POST /sessions/tables/bulk
// @Summary Bulk create tables
// @Description Create multiple tables in a specified range
// @Tags Tables
// @Accept json
// @Produce json
// @Param request body BulkCreateTablesRequest true "Bulk create tables request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/tables/bulk [post]
func (h *Handler) BulkCreateTables(c *gin.Context) {
	var req BulkCreateTablesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := ValidateBulkCreateTables(req); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	err := h.svc.BulkCreateTables(c.Request.Context(), req.Start, req.End)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Tables created from %d to %d", req.Start, req.End)})
}

// @Summary Get table by ID
// @Description Retrieve a specific table by its ID
// @Tags Tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Success 200 {object} Table
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/tables/{id} [get]
func (h *Handler) GetTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(c, errors.NewValidationError("invalid table ID"))
		return
	}

	if err := ValidateTableID(id); err != nil {
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

// ListTables handles GET /sessions/tables
// @Summary List all tables
// @Description Retrieve all restaurant tables
// @Tags Tables
// @Accept json
// @Produce json
// @Success 200 {array} Table
// @Failure 500 {object} middleware.ErrorResponse
// @Router /sessions/tables [get]
func (h *Handler) ListTables(c *gin.Context) {
	tables, err := h.svc.ListTables(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tables)
}

// DeleteTable handles DELETE /sessions/tables/:id
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
// @Router /sessions/tables/{id} [delete]
func (h *Handler) DeleteTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(c, errors.NewValidationError("invalid table ID"))
		return
	}

	if err := ValidateTableID(id); err != nil {
		middleware.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	err = h.svc.DeleteTable(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
}
