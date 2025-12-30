package service

import (
	"errors"
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"

	"github.com/google/uuid"
)

// Service defines business logic for sessions
type Service interface {
	CreateSession(tableID int) (*models.Session, error)
	GetSession(id string) (*models.Session, error)
	UpdateSession(id string, status models.SessionStatus) (*models.Session, error)
	ListSessions(offset, limit int) ([]*models.Session, error)
	ListActiveSessions() ([]*models.Session, error)
	ChangeTable(id string, tableNumber int) error
}

// sessionService implements Service
type sessionService struct {
	repo repository.Repository
}

// NewService creates a new session service
func NewService(repo repository.Repository) Service {
	return &sessionService{repo: repo}
}

// CreateSession creates a new session
func (s *sessionService) CreateSession(tableID int) (*models.Session, error) {
	if tableID <= 0 {
		return nil, errors.New("table ID must be greater than 0")
	}
	id := uuid.New().String()
	return s.repo.CreateSession(id, tableID)
}

// GetSession retrieves a session
func (s *sessionService) GetSession(id string) (*models.Session, error) {
	if id == "" {
		return nil, errors.New("session ID is required")
	}
	return s.repo.GetSession(id)
}

// UpdateSession updates the status of a session
func (s *sessionService) UpdateSession(id string, status models.SessionStatus) (*models.Session, error) {
	if id == "" {
		return nil, errors.New("session ID is required")
	}
	if status == "" {
		return nil, errors.New("status is required")
	}
	// Validate status is valid
	validStatuses := []models.SessionStatus{"active", "completed", "pending", "cancelled"}
	valid := false
	for _, s := range validStatuses {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("invalid status")
	}
	err := s.repo.UpdateSession(id, status)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSession(id)
}

// ListSessions lists sessions with pagination
func (s *sessionService) ListSessions(offset, limit int) ([]*models.Session, error) {
	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}
	if limit <= 0 || limit > 100 {
		return nil, errors.New("limit must be between 1 and 100")
	}
	return s.repo.ListSessions(offset, limit)
}

// ListActiveSessions lists active sessions
func (s *sessionService) ListActiveSessions() ([]*models.Session, error) {
	return s.repo.ListActiveSessions()
}

// ChangeTable changes the table of a session
func (s *sessionService) ChangeTable(id string, tableNumber int) error {
	if id == "" {
		return errors.New("session ID is required")
	}
	if tableNumber <= 0 {
		return errors.New("table number must be greater than 0")
	}
	return s.repo.ChangeSessionTable(id, tableNumber)
}