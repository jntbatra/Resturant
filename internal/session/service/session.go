package service

import (
	"context"
	apperrors "restaurant/internal/errors"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// SessionService defines business logic for sessions
type SessionService interface {
	CreateSession(ctx context.Context, tableID int) (*models.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error)
	UpdateSession(ctx context.Context, id uuid.UUID, status models.SessionStatus) (*models.Session, error)
	ListSessions(ctx context.Context, offset, limit int) ([]*models.Session, error)
	ListActiveSessions(ctx context.Context) ([]*models.Session, error)
	ChangeTable(ctx context.Context, id uuid.UUID, tableNumber int) error
	GetSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error)
	GetActiveSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error

	// Table operations
	CreateTable(ctx context.Context, req *models.CreateTableRequest) (*models.Table, error)
	GetTable(ctx context.Context, id int) (*models.Table, error)
	GetTableByNumber(ctx context.Context, number int) (*models.Table, error)
	ListTables(ctx context.Context) ([]*models.Table, error)
	DeleteTable(ctx context.Context, id int) error
	BulkCreateTables(ctx context.Context, start, end int) error
	IsTableAvailable(ctx context.Context, tableID int) (bool, error)
}

// sessionService implements Service
type sessionService struct {
	repo repository.Repository
}

// NewService creates a new session service
func NewService(repo repository.Repository) SessionService {
	return &sessionService{repo: repo}
}

// CreateSession creates a new session
func (s *sessionService) CreateSession(ctx context.Context, tableID int) (*models.Session, error) {
	// Check if table is available (no active or pending sessions)
	available, err := s.IsTableAvailable(ctx, tableID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to check table availability", err)
	}
	if !available {
		return nil, apperrors.WrapError(409, "table is not available (has active or pending session)", nil)
	}

	// Shape validation (tableID > 0) already done by handler using ValidateStruct
	id := uuid.New()
	session, err := s.repo.CreateSession(ctx, id, tableID)
	if err != nil {
		// Check for foreign key constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, apperrors.ErrForeignKeyViolation
		}
		return nil, apperrors.WrapError(500, "failed to create session", err)
	}

	return session, nil
}

// GetSession retrieves a session
func (s *sessionService) GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	session, err := s.repo.GetSession(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve session", err)
	}
	return session, nil
}

// UpdateSession updates the status of a session
func (s *sessionService) UpdateSession(ctx context.Context, id uuid.UUID, status models.SessionStatus) (*models.Session, error) {
	// Get current session to validate state transition (BUSINESS LOGIC)
	currentSession, err := s.repo.GetSession(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve session", err)
	}

	// Validate state transitions (BUSINESS LOGIC - cannot change completed/cancelled sessions)
	validTransitions := map[models.SessionStatus][]models.SessionStatus{
		"active":    {"pending", "cancelled"},
		"pending":   {"completed", "cancelled"},
		"completed": {},
		"cancelled": {},
	}

	allowedStatuses, exists := validTransitions[currentSession.Status]
	if !exists {
		return nil, apperrors.NewValidationError("invalid current session status")
	}

	valid := false
	for _, allowed := range allowedStatuses {
		if status == allowed {
			valid = true
			break
		}
	}
	if !valid {
		return nil, apperrors.NewValidationError("invalid status transition from " + string(currentSession.Status) + " to " + string(status))
	}

	// Shape validation (format, ranges) already done by handler using ValidateStruct
	err = s.repo.UpdateSession(ctx, id, status)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to update session", err)
	}

	return s.repo.GetSession(ctx, id)
}

// ListSessions lists sessions with pagination
func (s *sessionService) ListSessions(ctx context.Context, offset, limit int) ([]*models.Session, error) {
	// Shape validation (offset, limit ranges) already done by handler using ValidateStruct
	sessions, err := s.repo.ListSessions(ctx, offset, limit)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list sessions", err)
	}
	return sessions, nil
}

// ListActiveSessions lists active sessions
func (s *sessionService) ListActiveSessions(ctx context.Context) ([]*models.Session, error) {
	sessions, err := s.repo.ListActiveSessions(ctx)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list active sessions", err)
	}
	return sessions, nil
}

// ChangeTable changes the table of a session
func (s *sessionService) ChangeTable(ctx context.Context, id uuid.UUID, tableNumber int) error {
	// Shape validation (tableNumber > 0) already done by handler using ValidateStruct

	// Check if the new table is available (no active or pending sessions)
	available, err := s.IsTableAvailable(ctx, tableNumber)
	if err != nil {
		return apperrors.WrapError(500, "failed to check table availability", err)
	}
	if !available {
		return apperrors.NewValidationError("table is not available")
	}

	err = s.repo.ChangeSessionTable(ctx, id, tableNumber)
	if err != nil {
		// Check for foreign key constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return apperrors.ErrForeignKeyViolation
		}
		return apperrors.WrapError(500, "failed to change session table", err)
	}
	return nil
}

// GetSessionsByTable retrieves all sessions for a specific table
func (s *sessionService) GetSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error) {
	sessions, err := s.repo.GetSessionsByTable(ctx, tableID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve sessions for table", err)
	}
	return sessions, nil
}

// GetActiveSessionsByTable retrieves only active sessions for a specific table
func (s *sessionService) GetActiveSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error) {
	sessions, err := s.repo.GetActiveSessionsByTable(ctx, tableID)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve active sessions for table", err)
	}
	return sessions, nil
}

// DeleteSession deletes a session by ID
func (s *sessionService) DeleteSession(ctx context.Context, id uuid.UUID) error {
	// Check if session exists
	_, err := s.repo.GetSession(ctx, id)
	if err != nil {
		// If session not found, return not found error
		if strings.Contains(err.Error(), "Session not found") {
			return apperrors.ErrSessionNotFound
		}
		return apperrors.WrapError(500, "failed to check session existence", err)
	}

	// Session exists, proceed to delete
	err = s.repo.DeleteSession(ctx, id)
	if err != nil {
		return apperrors.WrapError(500, "failed to delete session", err)
	}
	return nil
}

// BulkCreateTables creates multiple tables in the specified range
func (s *sessionService) BulkCreateTables(ctx context.Context, start, end int) error {
	tableIDs := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		tableIDs = append(tableIDs, i)
	}
	err := s.repo.BulkCreateTables(ctx, tableIDs)
	if err != nil {
		return apperrors.WrapError(500, "failed to bulk create tables", err)
	}
	return nil
}

// CreateTable creates a new table
func (s *sessionService) CreateTable(ctx context.Context, req *models.CreateTableRequest) (*models.Table, error) {
	table, err := s.repo.CreateTable(ctx, req)
	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, apperrors.ErrConflict
		}
		return nil, apperrors.WrapError(500, "failed to create table", err)
	}
	return table, nil
}

// GetTable retrieves a table by ID
func (s *sessionService) GetTable(ctx context.Context, id int) (*models.Table, error) {
	table, err := s.repo.GetTable(ctx, id)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve table", err)
	}
	return table, nil
}

// GetTableByNumber retrieves a table by table number
func (s *sessionService) GetTableByNumber(ctx context.Context, number int) (*models.Table, error) {
	table, err := s.repo.GetTableByNumber(ctx, number)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to retrieve table by number", err)
	}
	return table, nil
}

// ListTables lists all tables
func (s *sessionService) ListTables(ctx context.Context) ([]*models.Table, error) {
	tables, err := s.repo.ListTables(ctx)
	if err != nil {
		return nil, apperrors.WrapError(500, "failed to list tables", err)
	}
	return tables, nil
}

// DeleteTable deletes a table
func (s *sessionService) DeleteTable(ctx context.Context, id int) error {
	err := s.repo.DeleteTable(ctx, id)
	if err != nil {
		return apperrors.WrapError(500, "failed to delete table", err)
	}
	return nil
}

// IsTableAvailable checks if a table has no active or pending sessions
func (s *sessionService) IsTableAvailable(ctx context.Context, tableID int) (bool, error) {
	available, err := s.repo.IsTableAvailable(ctx, tableID)
	if err != nil {
		return false, apperrors.WrapError(500, "failed to check table availability", err)
	}
	return available, nil
}
