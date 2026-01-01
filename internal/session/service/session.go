package service

import (
	"context"
	"errors"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"

	"github.com/google/uuid"
)

// SessionService defines business logic for sessions
type SessionService interface {
	CreateSession(ctx context.Context, tableID int) (*models.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error)
	UpdateSession(ctx context.Context, id uuid.UUID, status models.SessionStatus) (*models.Session, error)
	ListSessions(ctx context.Context, offset, limit int) ([]*models.Session, error)
	ListActiveSessions(ctx context.Context) ([]*models.Session, error)
	ChangeTable(ctx context.Context, id uuid.UUID, tableNumber int) error
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
	// Shape validation (tableID > 0) already done by handler using ValidateStruct
	id := uuid.New()
	return s.repo.CreateSession(ctx, id, tableID)
}

// GetSession retrieves a session
func (s *sessionService) GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	return s.repo.GetSession(ctx, id)
}

// UpdateSession updates the status of a session
func (s *sessionService) UpdateSession(ctx context.Context, id uuid.UUID, status models.SessionStatus) (*models.Session, error) {
	// Get current session to validate state transition (BUSINESS LOGIC)
	currentSession, err := s.repo.GetSession(ctx, id)
	if err != nil {
		return nil, err
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
		return nil, errors.New("invalid current status")
	}

	valid := false
	for _, allowed := range allowedStatuses {
		if status == allowed {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("invalid status transition from " + string(currentSession.Status) + " to " + string(status))
	}

	// Shape validation (format, ranges) already done by handler using ValidateStruct
	err = s.repo.UpdateSession(ctx, id, status)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSession(ctx, id)
}

// ListSessions lists sessions with pagination
func (s *sessionService) ListSessions(ctx context.Context, offset, limit int) ([]*models.Session, error) {
	// Shape validation (offset, limit ranges) already done by handler using ValidateStruct
	return s.repo.ListSessions(ctx, offset, limit)
}

// ListActiveSessions lists active sessions
func (s *sessionService) ListActiveSessions(ctx context.Context) ([]*models.Session, error) {
	return s.repo.ListActiveSessions(ctx)
}

// ChangeTable changes the table of a session
func (s *sessionService) ChangeTable(ctx context.Context, id uuid.UUID, tableNumber int) error {
	// Shape validation (tableNumber > 0) already done by handler using ValidateStruct
	return s.repo.ChangeSessionTable(ctx, id, tableNumber)
}
