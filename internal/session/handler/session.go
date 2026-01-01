package handler

import (
	"restaurant/internal/session/service"
	"restaurant/internal/session/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	}
}

// CreateSession handles POST /sessions
func (h *Handler) CreateSession(c *gin.Context) {
	var req validation.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateCreateSession(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session, err := h.svc.CreateSession(c.Request.Context(), req.TableID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, session)
}

// GetSession handles GET /sessions/:id
func (h *Handler) GetSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := validation.ValidateSessionID(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session, err := h.svc.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if session == nil {
		c.JSON(404, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(200, session)
}

// UpdateSession handles PUT /sessions/:id/status
func (h *Handler) UpdateSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
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
func (h *Handler) ListSessions(c *gin.Context) {
	var req validation.ListSessionsRequest
	req.Offset = 0
	req.Limit = 10

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
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
func (h *Handler) ChangeSessionTable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
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

	err = h.svc.ChangeTable(c.Request.Context(), id, req.TableID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Session table changed successfully"})
}

func (h *Handler) GetSessionsByTable(c *gin.Context) {
	// TODO: Implement
	c.JSON(501, gin.H{"error": "Not implemented"})
}

func (h *Handler) DeleteSession(c *gin.Context) {
	// TODO: Implement
	c.JSON(501, gin.H{"error": "Not implemented"})
}
